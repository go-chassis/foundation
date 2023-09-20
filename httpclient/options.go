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

package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

//DefaultOptions is a struct object which has default client option
var DefaultOptions = Options{
	Compressed:            true,
	HandshakeTimeout:      30 * time.Second,
	ResponseHeaderTimeout: 60 * time.Second,
	RequestTimeout:        60 * time.Second,
	ConnsPerHost:          5,
	TLSVersion:            tls.VersionTLS13,
}

//Options is a struct which provides options for client
type Options struct {
	TLSVersion            int
	TLSConfig             *tls.Config
	Compressed            bool
	HandshakeTimeout      time.Duration
	ResponseHeaderTimeout time.Duration
	RequestTimeout        time.Duration
	ConnsPerHost          int
	SignRequest           func(*http.Request) error
}
