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

import "time"

// TODO Task 应该有输入和输出
type Task func()

// 开辟一个协程池：当有任务提交时，提交到协程池中运行；如果协程池都在工作，任务挂起
type Pool interface {
	// 提交任务
	Submit(task Task)        // 提交任务
	SubmitAndWait(task Task) // 提交任务并等待其执行
	Stopped() bool           // 如果协程停止，返回true
	Stop()                   // 停下来优雅地停止所有的勾当，所有挂起的任务将在退出前完成
	//metrics
}

func NewDefaultPool(name string, maxWorkers int, idleTimeout time.Duration) Pool {
	return NewWorkerPool(name, maxWorkers, idleTimeout)
}
