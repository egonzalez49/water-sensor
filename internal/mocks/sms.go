package mocks

import "sync"

type Cache interface {
	Send(recipient string, body string) ([]byte, error)
}

type MockSms struct {
	mu sync.Mutex
	wg sync.WaitGroup

	SendCount int
}

func (s *MockSms) AddExpectedCalls(count int) {
	s.wg.Add(count)
}

func (s *MockSms) WaitForCalls() {
	s.wg.Wait()
}

func (s *MockSms) Send(recipient string, body string) ([]byte, error) {
	defer s.wg.Done()

	s.mu.Lock()
	s.SendCount += 1
	s.mu.Unlock()

	return nil, nil
}
