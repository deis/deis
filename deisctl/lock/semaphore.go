package lock

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

var (
	ErrExist    = errors.New("holder exists")
	ErrNotExist = errors.New("holder does not exist")
)

type Semaphore struct {
	Index     uint64   `json:"-"`
	Semaphore int      `json:"semaphore"`
	Max       int      `json:"max"`
	Holders   []string `json:"holders"`
}

func (s *Semaphore) SetMax(max int) error {
	diff := s.Max - max

	s.Semaphore = s.Semaphore - diff
	s.Max = s.Max - diff

	return nil
}

func (s *Semaphore) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Semaphore) addHolder(h string) error {
	loc := sort.SearchStrings(s.Holders, h)
	switch {
	case loc == len(s.Holders):
		s.Holders = append(s.Holders, h)
	case s.Holders[loc] == h:
		return ErrExist
	default:
		s.Holders = append(s.Holders[:loc], append([]string{h}, s.Holders[loc:]...)...)
	}

	return nil
}

func (s *Semaphore) removeHolder(h string) error {
	loc := sort.SearchStrings(s.Holders, h)
	if loc < len(s.Holders) && s.Holders[loc] == h {
		s.Holders = append(s.Holders[:loc], s.Holders[loc+1:]...)
	} else {
		return ErrNotExist
	}

	return nil
}

func (s *Semaphore) Lock(h string) error {
	if s.Semaphore <= 0 {
		return fmt.Errorf("semaphore is at %v", s.Semaphore)
	}

	if err := s.addHolder(h); err != nil {
		return err
	}

	s.Semaphore = s.Semaphore - 1

	return nil
}

func (s *Semaphore) Unlock(h string) error {
	if err := s.removeHolder(h); err != nil {
		return err
	}

	s.Semaphore = s.Semaphore + 1

	return nil
}

func newSemaphore() (sem *Semaphore) {
	return &Semaphore{0, 1, 1, nil}
}

type holder struct {
	ID        string `json:"-"`
	StartTime int64  `json:"startTime"`
}
