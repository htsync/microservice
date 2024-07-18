package repository

type RepositoryMock struct {
	messages map[int]string
	counter  int
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{messages: make(map[int]string)}
}

func (r *RepositoryMock) SaveMessage(content string) (int, error) {
	r.counter++
	r.messages[r.counter] = content
	return r.counter, nil
}

func (r *RepositoryMock) MarkMessageAsProcessed(id int) error {
	// For the mock, we don't need to implement this
	return nil
}

func (r *RepositoryMock) GetProcessedMessageCount() (int, error) {
	return r.counter, nil
}
