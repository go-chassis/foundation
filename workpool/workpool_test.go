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
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func Test_PoolSubmit(t *testing.T) {
	assert := assert.New(t)
	grNum := runtime.NumGoroutine()
	pool := NewDefaultPool("test", 2, time.Second*5)
	// 1个dispatcher goroutine + N个 work goroutine
	assert.Equal(grNum+1, runtime.NumGoroutine())

	var c atomic.Int32

	finished := make(chan struct{})
	do := func(iterations int) {
		for i := 0; i < iterations; i++ {
			pool.Submit(func() {
				c.Inc()
			})
		}
		finished <- struct{}{}
	}
	go do(100)
	<-finished
	assert.True(grNum+2+1 <= runtime.NumGoroutine())
	pool.Stop()
	pool.Stop()
	// reject all task
	go do(100)
	<-finished
	assert.Equal(int32(100), c.Load())
}

func Test_PoolSubmitAndWait(t *testing.T) {
	assert := assert.New(t)
	pool := NewDefaultPool("test", 2, time.Second*5)
	var c atomic.Int32

	finished := make(chan struct{})
	do := func(iterations int) {
		for i := 0; i < iterations; i++ {
			pool.SubmitAndWait(func() {
				c.Inc()
			})
		}
		finished <- struct{}{}
	}
	go do(1000)
	<-finished
	assert.Equal(int32(1000), c.Load())
	pool.Stop()
}

func Test_PoolStoped(t *testing.T) {
	assert := assert.New(t)
	pool := NewDefaultPool("test", 2, time.Second*5)
	var c atomic.Int32

	finished := make(chan struct{})
	do := func(iterations int) {
		for i := 0; i < iterations; i++ {
			pool.SubmitAndWait(func() {
				c.Inc()
			})
		}
		finished <- struct{}{}
	}
	go do(1000)
	<-finished
	assert.Equal(int32(1000), c.Load())

	ret := pool.Stopped()
	assert.False(ret)
	time.Sleep(time.Second * 1)
	ret = pool.Stopped()
	assert.False(ret)
	pool.Stop()
	ret = pool.Stopped()
	assert.True(ret)
}
