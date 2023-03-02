/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package quic

import (
	"context"
	"crypto/tls"
	"io"
	"net"

	"github.com/cloudwego/hertz/pkg/common/config"
	errs "github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network"
	"github.com/hertz-contrib/http3/network/quic-go/testdata"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

type transport struct {
	QuicConfig      *quic.Config
	TLSConfig       *tls.Config
	EnableDatagrams bool
	Addr            string

	tcpTransporter network.Transporter
	listener       io.Closer
	handler        network.OnData
}

func (t *transport) Close() error {
	return t.listener.Close()
}

func (t *transport) Shutdown(ctx context.Context) error {
	return t.Close()
}

func (t *transport) ListenAndServe(onData network.OnData) error {
	t.handler = onData
	return t.ListenAndServeTLS(testdata.GetCertificatePaths())
}

// ListenAndServeTLS listens on the UDP address s.Addr and calls s.Handler to handle HTTP/3 requests on incoming connections.
//
// If s.Addr is blank, ":https" is used.
func (t *transport) ListenAndServeTLS(certFile, keyFile string) error {
	var err error
	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	// We currently only use the cert-related stuff from tls.Config,
	// so we don't need to make a full copy.
	config := &tls.Config{
		Certificates: certs,
	}
	return t.serveConn(config, nil)
}

func (t *transport) serveConn(tlsConf *tls.Config, conn net.PacketConn) error {
	if tlsConf == nil {
		return errs.NewPublic("not support quic without tls")
	}

	baseConf := http3.ConfigureTLSConfig(tlsConf)
	quicConf := t.QuicConfig
	if quicConf == nil {
		quicConf = &quic.Config{}
	} else {
		quicConf = t.QuicConfig.Clone()
	}
	if t.EnableDatagrams {
		quicConf.EnableDatagrams = true
	}

	var ln quic.EarlyListener
	var err error
	if conn == nil {
		addr := t.Addr
		if addr == "" {
			addr = ":https"
		}
		ln, err = quic.ListenAddrEarly(addr, baseConf, quicConf)
	} else {
		ln, err = quic.ListenEarly(conn, baseConf, quicConf)
	}
	if err != nil {
		return err
	}
	t.listener = ln

	return t.serveListener(ln)
}

func (t *transport) serveListener(ln quic.EarlyListener) error {
	for {
		conn, err := ln.Accept(context.Background())
		if err != nil {
			return err
		}
		go func() {
			if err := t.handler(context.Background(), newStreamConn(conn)); err != nil {
				hlog.Debugf(err.Error())
			}
		}()
	}
}

// For transporter switch
func NewTransporter(options *config.Options) network.Transporter {
	return &transport{
		TLSConfig: options.TLS,
		Addr:      options.Addr,
	}
}
