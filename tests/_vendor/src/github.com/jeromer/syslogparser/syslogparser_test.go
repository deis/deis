package syslogparser

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type CommonTestSuite struct {
}

var _ = Suite(&CommonTestSuite{})

func (s *CommonTestSuite) TestParsePriority_Empty(c *C) {
	pri := newPriority(0)
	buff := []byte("")
	start := 0

	s.assertPriority(c, pri, buff, start, start, ErrPriorityEmpty)
}

func (s *CommonTestSuite) TestParsePriority_NoStart(c *C) {
	pri := newPriority(0)
	buff := []byte("7>")
	start := 0

	s.assertPriority(c, pri, buff, start, start, ErrPriorityNoStart)
}

func (s *CommonTestSuite) TestParsePriority_NoEnd(c *C) {
	pri := newPriority(0)
	buff := []byte("<77")
	start := 0

	s.assertPriority(c, pri, buff, start, start, ErrPriorityNoEnd)
}

func (s *CommonTestSuite) TestParsePriority_TooShort(c *C) {
	pri := newPriority(0)
	buff := []byte("<>")
	start := 0

	s.assertPriority(c, pri, buff, start, start, ErrPriorityTooShort)
}

func (s *CommonTestSuite) TestParsePriority_TooLong(c *C) {
	pri := newPriority(0)
	buff := []byte("<1233>")
	start := 0

	s.assertPriority(c, pri, buff, start, start, ErrPriorityTooLong)
}

func (s *CommonTestSuite) TestParsePriority_NoDigits(c *C) {
	pri := newPriority(0)
	buff := []byte("<7a8>")
	start := 0

	s.assertPriority(c, pri, buff, start, start, ErrPriorityNonDigit)
}

func (s *CommonTestSuite) TestParsePriority_Ok(c *C) {
	pri := newPriority(190)
	buff := []byte("<190>")
	start := 0

	s.assertPriority(c, pri, buff, start, start+5, nil)
}

func (s *CommonTestSuite) TestNewPriority(c *C) {
	obtained := newPriority(165)

	expected := Priority{
		P: 165,
		F: Facility{Value: 20},
		S: Severity{Value: 5},
	}

	c.Assert(obtained, DeepEquals, expected)
}

func (s *CommonTestSuite) TestParseVersion_NotFound(c *C) {
	buff := []byte("<123>")
	start := 5

	s.assertVersion(c, NO_VERSION, buff, start, start, ErrVersionNotFound)
}

func (s *CommonTestSuite) TestParseVersion_NonDigit(c *C) {
	buff := []byte("<123>a")
	start := 5

	s.assertVersion(c, NO_VERSION, buff, start, start+1, nil)
}

func (s *CommonTestSuite) TestParseVersion_Ok(c *C) {
	buff := []byte("<123>1")
	start := 5

	s.assertVersion(c, 1, buff, start, start+1, nil)
}

func (s *CommonTestSuite) TestParseHostname_Invalid(c *C) {
	// XXX : no year specified. Assumed current year
	// XXX : no timezone specified. Assume UTC
	buff := []byte("foo name")
	start := 0
	hostname := "foo"

	s.assertHostname(c, hostname, buff, start, 3, nil)
}

func (s *CommonTestSuite) TestParseHostname_Valid(c *C) {
	// XXX : no year specified. Assumed current year
	// XXX : no timezone specified. Assume UTC
	hostname := "ubuntu11.somehost.com"
	buff := []byte(hostname + " ")
	start := 0

	s.assertHostname(c, hostname, buff, start, len(hostname), nil)
}

func (s *CommonTestSuite) BenchmarkParsePriority(c *C) {
	buff := []byte("<190>")
	var start int
	l := len(buff)

	for i := 0; i < c.N; i++ {
		start = 0
		_, err := ParsePriority(buff, &start, l)
		if err != nil {
			panic(err)
		}
	}
}

func (s *CommonTestSuite) BenchmarkParseVersion(c *C) {
	buff := []byte("<123>1")
	start := 5
	l := len(buff)

	for i := 0; i < c.N; i++ {
		start = 0
		_, err := ParseVersion(buff, &start, l)
		if err != nil {
			panic(err)
		}
	}
}

func (s *CommonTestSuite) assertPriority(c *C, p Priority, b []byte, cursor int, expC int, e error) {
	obtained, err := ParsePriority(b, &cursor, len(b))
	c.Assert(obtained, DeepEquals, p)
	c.Assert(cursor, Equals, expC)
	c.Assert(err, Equals, e)
}

func (s *CommonTestSuite) assertVersion(c *C, version int, b []byte, cursor int, expC int, e error) {
	obtained, err := ParseVersion(b, &cursor, len(b))
	c.Assert(obtained, Equals, version)
	c.Assert(cursor, Equals, expC)
	c.Assert(err, Equals, e)
}

func (s *CommonTestSuite) assertHostname(c *C, h string, b []byte, cursor int, expC int, e error) {
	obtained, err := ParseHostname(b, &cursor, len(b))
	c.Assert(obtained, Equals, h)
	c.Assert(cursor, Equals, expC)
	c.Assert(err, Equals, e)
}
