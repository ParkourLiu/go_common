package cthreadpool

type ThreadPllClient struct {
	queue         chan bool
	runtineNumber int
	threadCount   int
}

func NewThreadPllClient(threadCount int) *ThreadPllClient {
	queue := make(chan bool, threadCount)
	return &ThreadPllClient{queue: queue, threadCount: threadCount}
}

func (c *ThreadPllClient) Run(method func()) {
	c.queue <- true
	go func() {
		method()
	}()
	i := len(c.queue)
}
