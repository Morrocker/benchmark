package benchmark

import (
	"container/list"
	"sync"
	"time"
)

type SingleRate interface {
	MeasureStart(int64) func(int64) // returns MeasureEnd function
	AvgRate() int64
	Reset()
	Values() (sampleSize uint, total int64, listLen int)
	SampleSize(...uint) uint
}

type singleRate struct {
	sampleSize uint
	total      int64
	list       *list.List

	lock sync.Mutex
}

// NewSingleRate returns a new SingleRate object with the given sample size
func NewSingleRate(n uint) SingleRate {
	newSingleRate := &singleRate{
		sampleSize: n,
		list:       list.New(),
	}
	return newSingleRate
}

// MeasureStart starts a measurement and returns a function to end and record the measurement
func (s *singleRate) MeasureStart(x int64) func(int64) {
	now := time.Now()
	return func(m int64) {
		diff := int64(time.Since(now).Seconds())
		diff2 := m - x
		rate := diff2 * 1000000000 / diff
		s.add(rate)
	}
}

// AvgRate returns the average rate from measurements taken
func (s *singleRate) AvgRate() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.list.Len() == 0 {
		return 0
	}
	return s.total / int64(s.list.Len())
}

// Reset sets all measurements values to their initial value
func (s *singleRate) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.total = 0
	s.list = list.New()
}

// Values returns the SingleRate relevant values
func (s *singleRate) Values() (sampleSize uint, total int64, listLen int) {
	return s.sampleSize, s.total, s.list.Len()
}

// SampleSize sets the sample size or returns it no value is given
func (s *singleRate) SampleSize(n ...uint) uint {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(n) != 0 {
		s.sampleSize = n[0]
		return s.sampleSize
	}
	return s.sampleSize
}

func (s *singleRate) add(r int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.list.Len() < int(s.sampleSize) || s.sampleSize == 0 {
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
