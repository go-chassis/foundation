/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package workpool

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"
)

const (
	// Queue 默认长度
	readyWorkerQueueSize = 32
	// Task 数据
	tasksCapacity = 8
	// 队列为空时，sleep 5ms
	sleepInterval = time.Millisecond * 5
)

type workerPool struct {
	name                string        // 工作协程池名称
	maxWorkers          int           // 最大工程协程池数据
	tasks               chan *Task    // Task channel
	readyWorkers        chan *worker  // 当前活跃工作协程
	idleTimeout         time.Duration // 空闲goroutine回收时间
	onDispatcherStopped chan struct{} // stop信号
	stopped             atomic.Bool   // 标记 协程池是否关闭
	workersAlive        atomic.Int32  // 当前协程使用数
	workersCreated      atomic.Int32  // 当前协程创建数
	workersKilled       atomic.Int32  // 当前协程完成数： 包括被kill
	tasksConsumed       atomic.Int32  // 处理的任务数
	ctx                 context.Context
	cancel              context.CancelFunc
}

func NewWorkerPool(name string, maxWorkers int, idleTimeout time.Duration) Pool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	pool := &workerPool{
		name:                name,
		maxWorkers:          maxWorkers,
		tasks:               make(chan *Task, tasksCapacity),
		readyWorkers:        make(chan *worker, readyWorkerQueueSize),
		idleTimeout:         idleTimeout,
		onDispatcherStopped: make(chan struct{}),
		stopped:             *atomic.NewBool(false),
		workersAlive:        *atomic.NewInt32(0),
		workersCreated:      *atomic.NewInt32(0),
		workersKilled:       *atomic.NewInt32(0),
		tasksConsumed:       *atomic.NewInt32(0),
		ctx:                 ctx,
		cancel:              cancel,
	}
	go pool.dispatch()
	return pool
}

func (p *workerPool) Submit(task *Task) {
	if task == nil || p.Stopped() {
		return
	}
	p.tasks <- task
}

func (p *workerPool) SubmitAndWait(task *Task) {
	if task == nil || p.Stopped() {
		return
	}
	worker := p.mustGetWorker()
	doneChan := make(chan struct{})
	worker.execute(&Task{
		ID: task.ID,
		F: func() {
			task.F()
			close(doneChan)
		},
	})
	<-doneChan
}

// 返回可用worker
func (p *workerPool) mustGetWorker() *worker {
	var worker *worker
	for {
		select {
		// 获得一个worker
		case worker = <-p.readyWorkers:
			return worker
		default:
			if int(p.workersAlive.Load()) >= p.maxWorkers {
				// 没有可用worker
				time.Sleep(sleepInterval)
				continue
			}
			w := NewWorker(p)
			return w
		}
	}
}

func (p *workerPool) dispatch() {
	defer func() {
		p.onDispatcherStopped <- struct{}{}
	}()

	idleTimeoutTimer := time.NewTimer(p.idleTimeout)
	defer idleTimeoutTimer.Stop()
	var (
		worker *worker
		task   *Task
	)

	for {
		idleTimeoutTimer.Reset(p.idleTimeout)
		select {
		case <-p.ctx.Done():
			return
		case task = <-p.tasks:
			worker := p.mustGetWorker()
			worker.execute(task)
		case <-idleTimeoutTimer.C:
			// 超时, kill掉worker
			if p.workersAlive.Load() > 0 {
				select {
				case worker = <-p.readyWorkers:
					worker.stop(func(chan *Task) {})
				default:
					// 所有worker都忙, continue
				}
			}
		}
	}
}

func (p *workerPool) Stopped() bool {
	return p.stopped.Load()
}

// stopWorkers stops all workers
func (p *workerPool) stopWorkers() {
	var wg sync.WaitGroup
	for p.workersAlive.Load() > 0 {
		wg.Add(1)
		worker := <-p.readyWorkers
		worker.stop(func(chan *Task) {
			wg.Done()
		})
	}
	wg.Wait()
}

// consumedRemainingTasks consumes all buffered tasks in the channel
func (p *workerPool) consumedRemainingTasks() {
	for {
		select {
		case task := <-p.tasks:
			task.F()
			p.tasksConsumed.Inc()
		default:
			return
		}
	}
}

// Stop tells the dispatcher to exit with pending tasks done.
func (p *workerPool) Stop() {
	if p.stopped.Swap(true) {
		return
	}
	// close dispatcher
	p.cancel()
	// wait dispatcher's exit
	<-p.onDispatcherStopped
	// close all workers
	p.stopWorkers()
	// consume remaining tasks
	p.consumedRemainingTasks()
}
