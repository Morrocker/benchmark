package benchmark

import (
	"container/list"
	"sync"
	"time"
)

type SRate struct {
	sampleSize int
	total      int64
	list       *list.List

	lock sync.Mutex
}

func NewSRate(n int) *SRate {
	newSRate := &SRate{
		sampleSize: n,
		list:       list.New(),
	}
	return newSRate
}

func (s *SRate) MeasureStart(x int64) func(int64) {
	now := time.Now()
	return func(m int64) {
		diff := int64(time.Since(now).Nanoseconds())
		diff2 := m - x
		rate := diff2 * 1000000000 / diff
		s.add(rate)
	}
}

func (s *SRate) AvgRate() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.list.Len() == 0 {
		return 0
	}
	return s.total / int64(s.list.Len())
}

func (s *SRate) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()

}

func (s *SRate) add(r int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.list.Len() < s.sampleSize || s.sampleSize == 0 {
		s.list.PushFront(r)
		s.total = s.total + r
	} else {
		s.list.PushFront(r)
		e := s.list.Back()
		val := e.Value.(int64)
		s.total = s.total + r - val
		s.list.Remove(e)
	}
}