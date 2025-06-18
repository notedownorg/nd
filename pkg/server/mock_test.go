// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"

	pb "github.com/notedownorg/nd/api/go/nodes/v1alpha1"
	"google.golang.org/grpc/metadata"
)

// Mock implementation for testing
type mockNodeServiceStreamServer struct {
	events chan *pb.StreamEvent
	ctx    context.Context
}

func (m *mockNodeServiceStreamServer) Send(event *pb.StreamEvent) error {
	select {
	case m.events <- event:
		return nil
	default:
		return nil // Drop if channel is full
	}
}

func (m *mockNodeServiceStreamServer) Recv() (*pb.StreamRequest, error) {
	// Not used in these tests
	return nil, nil
}

func (m *mockNodeServiceStreamServer) Context() context.Context {
	return m.ctx
}

func (m *mockNodeServiceStreamServer) SendMsg(interface{}) error {
	return nil
}

func (m *mockNodeServiceStreamServer) RecvMsg(interface{}) error {
	return nil
}

func (m *mockNodeServiceStreamServer) SendHeader(metadata.MD) error {
	return nil
}

func (m *mockNodeServiceStreamServer) SetHeader(metadata.MD) error {
	return nil
}

func (m *mockNodeServiceStreamServer) SetTrailer(metadata.MD) {
}
