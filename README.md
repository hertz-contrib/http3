# Hertz-HTTP3

This repo is the collection of Hertz HTTP3 implementations. Includes: Network layer & Protocol layer.
Detailed information can be found in
the [Hertz-HTTP3](https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/protocol/http3/).

## Network Layer

Currently, we provide 1 implementation of network layer which is based on: [quic-go](https://github.com/quic-go/quic-go)
.

### quic-go

#### Usage

```go
package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
	"github.com/hertz-contrib/http3/network/quic-go"
	"github.com/hertz-contrib/http3/network/quic-go/testdata"
)

func main() {
	h := server.New(server.WithALPN(true), server.WithTLS(testdata.GetTLSConfig()), server.WithTransport(quic.NewTransporter), server.WithAltTransport(netpoll.NewTransporter), server.WithHostPorts("127.0.0.1:8080"))
	...

	h.Spin()
}
```

QUIC is forced to depend on TLS, so you need to provide a TLS configuration.
For there is only Server side ready, we embed a testdata package from quic-go, which means the example server can
directly communicate with the example client
from [quic-go](https://github.com/quic-go/quic-go/blob/master/example/client/main.go).

#### Options

``server.WithTransport()``
Use it to set the network layer implementation.

``server.WithAltTransport()``
Use it to set the alternative network layer implementation. The AltTransporter will be used for parallel listening -
both in TCP and QUIC.

``server.WithALPN()``
Whether to enable ALPN.

``server.WithTLS()``
Which TLS configuration to use.

``server.WithHostPorts()``
Which host and port to listen on.

## Protocol Layer

Currently, we provide 1 implementation of protocol layer which is also based
on: [quic-go](https://github.com/quic-go/quic-go)
.

### quic-go

#### Usage

```go
package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
	"github.com/cloudwego/hertz/pkg/protocol/suite"
	"github.com/hertz-contrib/http3/network/quic-go"
	"github.com/hertz-contrib/http3/network/quic-go/testdata"
	http3 "github.com/hertz-contrib/http3/server/quic-go"
	"github.com/hertz-contrib/http3/server/quic-go/factory"
)

func main() {
	h := server.New(server.WithALPN(true), server.WithTLS(testdata.GetTLSConfig()), server.WithTransport(quic.NewTransporter), server.WithAltTransport(netpoll.NewTransporter), server.WithHostPorts("127.0.0.1:8080"))
	h.AddProtocol(suite.HTTP3, factory.NewServerFactory(&http3.Option{}))
    ...
	
	h.Spin()
}
```

### Example

For battery-included example, please refer
to [hertz-example](https://github.com/hertz-contrib/http3/blob/main/examples/quic-go/main.go).

Try using [quic-go client](https://github.com/quic-go/quic-go/blob/master/example/client/main.go) to say hello to the
server.
