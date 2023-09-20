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

package gopool

import (
	"context"
	"log"
	"time"
)

const (
	DefaultWorkers     = 1000
	DefaultIdleTimeout = 60 * time.Second
)

type Config struct {
	Ctx         context.Context
	Concurrent  int
	IdleTimeout time.Duration
	// RecoverFunc execute after recover a panic
	RecoverFunc func(r interface{})
}

func (c *Config) Workers(max int) *Config {
	c.Concurrent = max
	return c
}

func (c *Config) Idle(time time.Duration) *Config {
	c.IdleTimeout = time
	return c
}

func (c *Config) WithRecoverFunc(f func(r interface{})) *Config {
	c.RecoverFunc = f
	return c
}

func (c *Config) WithContext(ctx context.Context) *Config {
	c.Ctx = ctx
	return c
}

func Configure() *Config {
	return &Config{
		Ctx:         context.Background(),
		Concurrent:  DefaultWorkers,
		IdleTimeout: DefaultIdleTimeout,
		RecoverFunc: func(r interface{}) {
			log.Println("gopool recover:", r)
		},
	}
}
