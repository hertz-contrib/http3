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

package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/protocol/suite"
	"github.com/hertz-contrib/http3/network/quic-go"
	"github.com/hertz-contrib/http3/network/quic-go/testdata"
	http3 "github.com/hertz-contrib/http3/server/quic-go"
	"github.com/hertz-contrib/http3/server/quic-go/factory"
)

type Test struct {
	A string
	B string
}

func main() {
	run()
}

func run() {
	h := server.New(server.WithALPN(true), server.WithTLS(testdata.GetTLSConfig()), server.WithTransport(quic.NewTransporter), server.WithAltTransport(netpoll.NewTransporter), server.WithHostPorts("127.0.0.1:8080"))
	h.AddProtocol(suite.HTTP3, factory.NewServerFactory(&http3.Option{}))

	h.GET("/demo/tile", func(c context.Context, ctx *app.RequestContext) {
		// Small 40x40 png
		ctx.Write([]byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
			0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x28,
			0x01, 0x03, 0x00, 0x00, 0x00, 0xb6, 0x30, 0x2a, 0x2e, 0x00, 0x00, 0x00,
			0x03, 0x50, 0x4c, 0x54, 0x45, 0x5a, 0xc3, 0x5a, 0xad, 0x38, 0xaa, 0xdb,
			0x00, 0x00, 0x00, 0x0b, 0x49, 0x44, 0x41, 0x54, 0x78, 0x01, 0x63, 0x18,
			0x61, 0x00, 0x00, 0x00, 0xf0, 0x00, 0x01, 0xe2, 0xb8, 0x75, 0x22, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
		})
	})

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	h.GET("/struct", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, &Test{
			A: "aaa",
			B: "bbb",
		})
	})

	v1 := h.Group("/v1")
	{
		v1.GET("/hello/:name", func(c context.Context, ctx *app.RequestContext) {
			fmt.Fprintf(ctx, "Hi %s, this is the response from Hertz.\n", ctx.Param("name"))
		})
	}

	h.Spin()
}
