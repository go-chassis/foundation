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

// 执行任务的worker
type worker struct {
	pool   *workerPool
	tasks  chan *Task
	stopCh chan struct{}
}

func NewWorker(pool *workerPool) *worker {
	w := &worker{
		pool:   pool,
		tasks:  make(chan *Task),
		stopCh: make(chan struct{}),
	}
	w.pool.workersAlive.Inc()
	w.pool.workersCreated.Inc()
	go w.process()
	return w
}

func (w *worker) execute(task *Task) {
	w.tasks <- task
}

func (w *worker) stop(callable func(chan *Task)) {
	defer callable(w.tasks)
	w.stopCh <- struct{}{}
	w.pool.workersKilled.Inc()
	w.pool.workersAlive.Dec()
}

func (w *worker) process() {
	var task *Task
	for {
		select {
		case <-w.stopCh:
			return
		case task = <-w.tasks:
			task.F()
			w.pool.tasksConsumed.Inc()
			// 将w注册到readyWorkers
			w.pool.readyWorkers <- w
		}
	}
}
