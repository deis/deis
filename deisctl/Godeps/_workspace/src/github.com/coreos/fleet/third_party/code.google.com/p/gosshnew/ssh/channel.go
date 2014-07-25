// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
)

const (
	minPacketLength = 9
	// As per RFC 4253, section 6.1, 32k is also the minimum.
	channelMaxPacket = 1 << 15
	// We follow OpenSSH here.
	channelWindowSize = 64 * channelMaxPacket
)

// NewChannel is an incoming request to a channel. It must either be
// accepted for use with the Accept method or rejected with Reject.
type NewChannel interface {
	// Accept accepts the channel creation request. It returns the
	// Channel and the side-channel for requests. The side-channel
	// must be serviced.
	Accept() (Channel, <-chan *Request, error)

	// Reject rejects the channel creation request. After calling
	// this, no other methods on the Channel may be called.
	Reject(reason RejectionReason, message string) error

	// ChannelType returns the type of the channel, as supplied by the
	// client.
	ChannelType() string

	// ExtraData returns the arbitrary payload for this channel, as supplied
	// by the client. This data is specific to the channel type.
	ExtraData() []byte
}

// A Channel is an ordered, reliable, flow-controlled, duplex stream
// that is multiplexed over an SSH connection.
type Channel interface {
	// Read reads up to len(data) bytes from the channel.
	Read(data []byte) (int, error)

	// Write writes len(data) bytes to the channel.
	Write(data []byte) (int, error)

	// Signals end of channel use. No data may be sent after this
	// call.
	Close() error

	// CloseWrite signals the end of sending in-band
	// data. Requests may still be sent, and the other side may
	// still send data
	CloseWrite() error

	// SendRequest sends a channel request.  If wantReply is true,
	// it will wait for a reply and return the result as a
	// boolean, otherwise the return value will be false. Channel
	// requests are out-of-band messages, so they may be sent even
	// if the data stream is closed or blocked by flow control.
	SendRequest(name string, wantReply bool, payload []byte) (bool, error)

	// Stderr returns an io.ReadWriter that writes to this channel with the
	// extended data type set to stderr.
	Stderr() io.ReadWriter
}

// Request is a request sent outside of the normal stream of
// data. Requests can either be specific to an SSH channel, or they
// can be global.
type Request struct {
	Type      string
	WantReply bool
	Payload   []byte

	ch  *channel
	mux *mux
}

// Reply sends a message back to the request.  It should be called for
// all requests that have WantReply set.  The payload argument is
// ignored for channel specific requests
func (r *Request) Reply(ok bool, payload []byte) error {
	if !r.WantReply {
		return errors.New("ssh: request did not have WantReply")
	}

	if r.ch == nil {
		return r.mux.ackRequest(ok, payload)
	}

	return r.ch.ackRequest(ok)
}

// RejectionReason is an enumeration used when rejecting channel creation
// requests. See RFC 4254, section 5.1.
type RejectionReason uint32

const (
	Prohibited RejectionReason = iota + 1
	ConnectionFailed
	UnknownChannelType
	ResourceShortage
)

// String converts the rejection reason to human readable form.
func (r RejectionReason) String() string {
	switch r {
	case Prohibited:
		return "administratively prohibited"
	case ConnectionFailed:
		return "connect failed"
	case UnknownChannelType:
		return "unknown channel type"
	case ResourceShortage:
		return "resource shortage"
	}
	return fmt.Sprintf("unknown reason %d", int(r))
}

func min(a uint32, b int) uint32 {
	if a < uint32(b) {
		return a
	}
	return uint32(b)
}

// channel is an implementation of the Channel interface that works
// with the mux class.
type channel struct {
	// R/O after creation
	chanType          string
	extraData         []byte
	localId, remoteId uint32

	// maxPayload is the maximum payload size of normal and
	// extended data packets. The wire packet will be 9 or 13 bytes
	// larger (excluding encryption overhead).
	maxPayload uint32
	mux        *mux

	// If set, we have called Accept or Reject on this channel
	decided bool

	// Pending internal channel messages.
	msg chan interface{}

	// Since requests have no ID, there can be only one request
	// with WantReply=true outstanding.  This lock is held by a
	// goroutine that has such an outgoing request pending.
	sentRequestMu sync.Mutex

	incomingRequests chan *Request

	sentEOF bool

	// thread-safe data
	remoteWin  window
	pending    *buffer
	extPending *buffer

	// windowMu protects myWindow, the flow-control window.
	windowMu sync.Mutex
	myWindow uint32

	// writeMu serializes calls to mux.conn.writePacket() and
	// protects sentClose. This mutex must be different from
	// windowMu, as writePacket can block if there is a key
	// exchange pending
	writeMu   sync.Mutex
	sentClose bool
}

// writePacket sends the packet over the wire. If the packet is a
// channel close, it updates sentClose. This method takes the lock
// c.mu.
func (c *channel) writePacket(packet []byte) error {
	c.writeMu.Lock()
	if c.sentClose {
		c.writeMu.Unlock()
		return io.EOF
	}
	c.sentClose = (packet[0] == msgChannelClose)
	err := c.mux.conn.writePacket(packet)
	c.writeMu.Unlock()
	return err
}

func (c *channel) sendMessage(msg interface{}) error {
	if debugMux {
		log.Printf("send %d: %#v", c.mux.chanList.offset, msg)
	}

	p := Marshal(msg)
	binary.BigEndian.PutUint32(p[1:], c.remoteId)
	return c.writePacket(p)
}

func (c *channel) WriteExtended(data []byte, extendedCode uint32) (n int, err error) {
	if c.sentEOF {
		return 0, io.EOF
	}
	// 1 byte message type, 4 bytes remoteId, 4 bytes data length
	opCode := byte(msgChannelData)
	headerLength := uint32(9)
	if extendedCode > 0 {
		headerLength += 4
		opCode = msgChannelExtendedData
	}

	for len(data) > 0 {
		space := min(c.maxPayload, len(data))
		if space, err = c.remoteWin.reserve(space); err != nil {
			return n, err
		}
		todo := data[:space]

		packet := make([]byte, headerLength+uint32(len(todo)))
		packet[0] = opCode
		marshalUint32(packet[1:], c.remoteId)
		if extendedCode > 0 {
			marshalUint32(packet[5:], uint32(extendedCode))
		}
		marshalUint32(packet[headerLength-4:], uint32(len(todo)))
		copy(packet[headerLength:], todo)
		if err = c.writePacket(packet); err != nil {
			return n, err
		}

		n += len(todo)
		data = data[len(todo):]
	}

	return n, err
}

func (c *channel) handleData(packet []byte) error {
	headerLen := 9
	if packet[0] == msgChannelExtendedData {
		headerLen = 13
	}
	if len(packet) < headerLen {
		// malformed data packet
		return parseError(packet[0])
	}

	var extended uint32
	if headerLen > 9 {
		extended = binary.BigEndian.Uint32(packet[5:])
	}

	length := binary.BigEndian.Uint32(packet[headerLen-4 : headerLen])
	if length == 0 {
		return nil
	}
	if uint32(length) > c.maxPayload {
		// TODO(hanwen): should send Disconnect?
		return errors.New("ssh: incoming packet exceeds maximum payload size")
	}

	data := packet[headerLen:]
	if length != uint32(len(data)) {
		return errors.New("ssh: wrong packet length")
	}

	c.windowMu.Lock()
	if c.myWindow < length {
		c.windowMu.Unlock()
		// TODO(hanwen): should send Disconnect with reason?
		return errors.New("ssh: remote side wrote too much")
	}
	c.myWindow -= length
	c.windowMu.Unlock()

	if extended == 1 {
		c.extPending.write(data)
	} else if extended > 0 {
		// discard other extended data.
	} else {
		c.pending.write(data)
	}
	return nil
}

func (c *channel) adjustWindow(n uint32) error {
	c.windowMu.Lock()
	// Since myWindow is managed on our side, and can never exceed
	// the initial window setting, we don't worry about overflow.
	c.myWindow += uint32(n)
	c.windowMu.Unlock()
	return c.sendMessage(windowAdjustMsg{
		AdditionalBytes: uint32(n),
	})
}

func (c *channel) ReadExtended(data []byte, extended uint32) (n int, err error) {
	switch extended {
	case 1:
		n, err = c.extPending.Read(data)
	case 0:
		n, err = c.pending.Read(data)
	default:
		return 0, fmt.Errorf("ssh: extended code %d unimplemented", extended)
	}

	if n > 0 {
		err = c.adjustWindow(uint32(n))
		// sendWindowAdjust can return io.EOF if the remote
		// peer has closed the connection, however we want to
		// defer forwarding io.EOF to the caller of Read until
		// the buffer has been drained.
		if n > 0 && err == io.EOF {
			err = nil
		}
	}

	return n, err
}

func (c *channel) close() {
	c.pending.eof()
	c.extPending.eof()
	close(c.msg)
	close(c.incomingRequests)
	c.writeMu.Lock()
	// This is not necesary for a normal channel teardown, but if
	// there was another error, it is.
	c.sentClose = true
	c.writeMu.Unlock()
	// Unblock writers.
	c.remoteWin.close()
}

func (c *channel) handlePacket(packet []byte) error {
	switch packet[0] {
	case msgChannelData, msgChannelExtendedData:
		return c.handleData(packet)
	case msgChannelClose:
		c.sendMessage(channelCloseMsg{
			PeersId: c.remoteId})
		c.mux.chanList.remove(c.localId)
		c.close()
		return nil
	case msgChannelEOF:
		// RFC 4254 is mute on how EOF affects dataExt messages but
		// it is logical to signal EOF at the same time.
		c.extPending.eof()
		c.pending.eof()
		return nil
	}

	decoded, err := decode(packet)
	if err != nil {
		return err
	}

	switch msg := decoded.(type) {
	case *channelOpenFailureMsg:
		c.mux.chanList.remove(msg.PeersId)
		c.msg <- msg
	case *channelOpenConfirmMsg:
		if msg.MaxPacketSize < minPacketLength || msg.MaxPacketSize > 1<<31 {
			return errors.New("ssh: invalid MaxPacketSize from peer")
		}
		// fixup remoteId field
		c.remoteId = msg.MyId
		c.maxPayload = msg.MaxPacketSize
		c.remoteWin.add(msg.MyWindow)
		c.decided = true
		c.msg <- msg
	case *windowAdjustMsg:
		if !c.remoteWin.add(msg.AdditionalBytes) {
			return fmt.Errorf("invalid window update %d", msg.AdditionalBytes)
		}
	case *channelRequestMsg:
		req := Request{
			Type:      msg.Request,
			WantReply: msg.WantReply,
			Payload:   msg.RequestSpecificData,
			ch:        c,
		}

		c.incomingRequests <- &req
	default:
		c.msg <- msg
	}
	return nil
}

func (m *mux) newChannel(chanType string, extraData []byte) *channel {
	ch := &channel{
		remoteWin:        window{Cond: newCond()},
		myWindow:         channelWindowSize,
		pending:          newBuffer(),
		extPending:       newBuffer(),
		incomingRequests: make(chan *Request, 16),
		msg:              make(chan interface{}, 16),
		chanType:         chanType,
		extraData:        extraData,
		mux:              m,
	}
	ch.localId = m.chanList.add(ch)
	return ch
}

var errUndecided = errors.New("ssh: must Accept or Reject channel")
var errDecidedAlready = errors.New("ssh: can call Accept or Reject only once")

type extChannel struct {
	code uint32
	ch   *channel
}

func (e *extChannel) Write(data []byte) (n int, err error) {
	return e.ch.WriteExtended(data, e.code)
}

func (e *extChannel) Read(data []byte) (n int, err error) {
	return e.ch.ReadExtended(data, e.code)
}

func (c *channel) Accept() (Channel, <-chan *Request, error) {
	if c.decided {
		return nil, nil, errDecidedAlready
	}
	confirm := channelOpenConfirmMsg{
		PeersId:       c.remoteId,
		MyId:          c.localId,
		MyWindow:      c.myWindow,
		MaxPacketSize: c.maxPayload,
	}
	c.decided = true
	if err := c.sendMessage(confirm); err != nil {
		return nil, nil, err
	}

	return c, c.incomingRequests, nil
}

func (ch *channel) Reject(reason RejectionReason, message string) error {
	if ch.decided {
		return errDecidedAlready
	}
	reject := channelOpenFailureMsg{
		PeersId:  ch.remoteId,
		Reason:   reason,
		Message:  message,
		Language: "en",
	}
	ch.decided = true
	return ch.sendMessage(reject)
}

func (ch *channel) Read(data []byte) (int, error) {
	if !ch.decided {
		return 0, errUndecided
	}
	return ch.ReadExtended(data, 0)
}

func (ch *channel) Write(data []byte) (int, error) {
	if !ch.decided {
		return 0, errUndecided
	}
	return ch.WriteExtended(data, 0)
}

func (ch *channel) CloseWrite() error {
	if !ch.decided {
		return errUndecided
	}
	ch.sentEOF = true
	return ch.sendMessage(channelEOFMsg{
		PeersId: ch.remoteId})
}

func (ch *channel) Close() error {
	if !ch.decided {
		return errUndecided
	}

	return ch.sendMessage(channelCloseMsg{
		PeersId: ch.remoteId})
}

func (ch *channel) Extended(code uint32) io.ReadWriter {
	if !ch.decided {
		return nil
	}
	return &extChannel{code, ch}
}

func (ch *channel) Stderr() io.ReadWriter {
	return ch.Extended(1)
}

func (ch *channel) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	if !ch.decided {
		return false, errUndecided
	}

	if wantReply {
		ch.sentRequestMu.Lock()
		defer ch.sentRequestMu.Unlock()
	}

	msg := channelRequestMsg{
		PeersId:             ch.remoteId,
		Request:             name,
		WantReply:           wantReply,
		RequestSpecificData: payload,
	}

	if err := ch.sendMessage(msg); err != nil {
		return false, err
	}

	if wantReply {
		m, ok := (<-ch.msg)
		if !ok {
			return false, io.EOF
		}
		switch m.(type) {
		case *channelRequestFailureMsg:
			return false, nil
		case *channelRequestSuccessMsg:
			return true, nil
		default:
			return false, fmt.Errorf("unexpected response %#v", m)
		}
	}

	return false, nil
}

// ackRequest either sends an ack or nack to the channel request.
func (ch *channel) ackRequest(ok bool) error {
	if !ch.decided {
		return errUndecided
	}

	var msg interface{}
	if !ok {
		msg = channelRequestFailureMsg{
			PeersId: ch.remoteId,
		}
	} else {
		msg = channelRequestSuccessMsg{
			PeersId: ch.remoteId,
		}
	}
	return ch.sendMessage(msg)
}

func (ch *channel) ChannelType() string {
	return ch.chanType
}

func (ch *channel) ExtraData() []byte {
	return ch.extraData
}
