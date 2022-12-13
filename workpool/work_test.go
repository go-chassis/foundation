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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func Test_NewWorker(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	w := NewWorker(
		&workerPool{
			name:                "test",
			maxWorkers:          2,
			tasks:               make(chan Task, tasksCapacity),
			readyWorkers:        make(chan *worker, readyWorkerQueueSize),
			idleTimeout:         time.Second * 5,
			onDispatcherStopped: make(chan struct{}),
			stopped:             *atomic.NewBool(false),
			workersAlive:        *atomic.NewInt32(0),
			workersCreated:      *atomic.NewInt32(0),
			workersKilled:       *atomic.NewInt32(0),
			tasksConsumed:       *atomic.NewInt32(0),
			ctx:                 ctx,
			cancel:              cancel,
		},
	)

	w.execute(func() { fmt.Println("dongjiang1") })
	w.execute(func() { fmt.Println("dongjiang2") })

	w.stop(func() { fmt.Println("finished") })
	assert.True(true)
}
