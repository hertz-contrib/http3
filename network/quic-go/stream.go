/*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package quic

import (
	"github.com/cloudwego/hertz/pkg/network"
	"github.com/lucas-clemente/quic-go"
)

var _ network.Stream = &stream{}

type stream struct {
	network.ReceiveStream
	network.SendStream
}

type readStream struct {
	quic.ReceiveStream
}

func (r *readStream) CancelRead(err network.ApplicationError) {
	r.ReceiveStream.CancelRead(quic.StreamErrorCode(err.ErrCode()))
}

func (r *readStream) StreamID() int64 {
	return int64(r.ReceiveStream.StreamID())
}

type writeStream struct {
	quic.SendStream
}

func (w *writeStream) CancelWrite(err network.ApplicationError) {
	w.SendStream.CancelWrite(quic.StreamErrorCode(err.ErrCode()))
}

func (s *stream) StreamID() int64 {
	// the result is same for receiveStream and sendStream
	return s.SendStream.StreamID()
}

func (w *writeStream) StreamID() int64 {
	return int64(w.SendStream.StreamID())
}

func newReadStream(s quic.ReceiveStream) network.ReceiveStream {
	return &readStream{s}
}

func newWriteStream(s quic.SendStream) network.SendStream {
	return &writeStream{s}
}

func newStream(s quic.Stream) network.Stream {
	return &stream{
		ReceiveStream: newReadStream(s),
		SendStream:    newWriteStream(s),
	}
}
