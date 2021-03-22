package benchmark

import (
	"sync"
	"time"
)

type Simple struct {
	sampleSize int
	samples    int
	total      int64
	replace    int64
	lock       sync.Mutex
}

func NewSimple(n int) *Simple {
	newSimple := &Simple{
		sampleSize: n,
	}
	return newSimple
}

func (s *Simple) MeasureStart() func() {
	now := time.Now()
	return func() {
		diff := int64(time.Since(now).Seconds())
		s.addDiff(diff)
	}

}

func (s *Simple) AvgTime() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.total / int64(s.samples)
}

func (s *Simple) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.total = 0
	s.replace = 0
	s.samples = 0
}

func (s *Simple) addDiff(d int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.samples < s.sampleSize || s.sampleSize == 0 {
		s.samples++
		s.total = s.total + d
		if s.samples == 1 {
			s.replace = d
		}
	} else {
		s.total = s.total + d - s.replace
		s.replace = d
	}
}
