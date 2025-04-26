package code_runner

import (
	"sync"
)

type Run struct {
	TimeLimitMs   int    `json:"time_limit"`
	MemoryLimitMb int    `json:"memory_limit"`
	PproblemID    int    `json:"problem_id"`
	SubmissionID  int    `json:"submission_id"`
	CallbackToken string `json:"calback_token"`
}

type QueueManager struct {
	tasks     chan Run
	semaphore chan struct{}
	wg        sync.WaitGroup
}

func NewQueueManager(maxConcurrent int) *QueueManager {
	qm := &QueueManager{
		tasks:     make(chan Run, 1000),
		semaphore: make(chan struct{}, maxConcurrent),
	}
	go qm.startWorkers()
	return qm
}

func (qm *QueueManager) startWorkers() {
	for task := range qm.tasks {
		qm.wg.Add(1)
		qm.semaphore <- struct{}{}
		go func(r Run) {
			defer qm.wg.Done()
			defer func() { <-qm.semaphore }()
			sendRunCallBack(runCodeInsideContainer(r), r)
		}(task)
	}
}

func (qm *QueueManager) Enqueue(r Run) {
	qm.tasks <- r
}

func (qm *QueueManager) Close() {
	close(qm.tasks)
	qm.wg.Wait()
}
