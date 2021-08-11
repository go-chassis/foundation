// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tlsutil

import (
	"crypto/tls"
)

type TLSOptions struct {
	VerifyPeer     bool
	VerifyHostName bool
	CipherSuites   []uint16
	MinVersion     uint16
	MaxVersion     uint16
	CACertFile     string
	CertFile       string
	KeyFile        string
	KeyPassphase   string
	Decrypt        Decrypt
}

type TLSOption func(*TLSOptions)

func WithVerifyPeer(b bool) TLSOption      { return func(c *TLSOptions) { c.VerifyPeer = b } }
func WithVerifyHostName(b bool) TLSOption  { return func(c *TLSOptions) { c.VerifyHostName = b } }
func WithCipherSuits(s []uint16) TLSOption { return func(c *TLSOptions) { c.CipherSuites = s } }
func WithVersion(min, max uint16) TLSOption {
	return func(c *TLSOptions) { c.MinVersion, c.MaxVersion = min, max }
}
func WithCert(f string) TLSOption     { return func(c *TLSOptions) { c.CertFile = f } }
func WithKey(k string) TLSOption      { return func(c *TLSOptions) { c.KeyFile = k } }
func WithKeyPass(p string) TLSOption  { return func(c *TLSOptions) { c.KeyPassphase = p } }
func WithCA(f string) TLSOption       { return func(c *TLSOptions) { c.CACertFile = f } }
func WithDecrypt(f Decrypt) TLSOption { return func(c *TLSOptions) { c.Decrypt = f } }

func toTLSOptions(opts ...TLSOption) (op TLSOptions) {
	for _, opt := range opts {
		opt(&op)
	}
	return
}

func DefaultClientTLSOptions() []TLSOption {
	return []TLSOption{
		WithVerifyPeer(true),
		WithVerifyHostName(true),
		WithVersion(tls.VersionTLS12, MaxSupportedTLSVersion),
	}
}

func DefaultServerTLSOptions() []TLSOption {
	return []TLSOption{
		WithVerifyPeer(true),
		WithVersion(tls.VersionTLS12, MaxSupportedTLSVersion),
		WithCipherSuits(TLSCipherSuits()),
	}
}
