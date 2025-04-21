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
	"fmt"
	"log/slog"

	pb "github.com/notedownorg/nd/api/go/nodes/v1alpha1"
	"github.com/notedownorg/nd/pkg/workspace"
	"github.com/notedownorg/nd/pkg/workspace/reader"
)

type Server struct {
	pb.UnimplementedNodeServiceServer
	log *slog.Logger

	workspaces    map[string]*workspace.Workspace
	subscriptions map[string]chan reader.Event
}

func NewServer() *Server {
	return &Server{
		log:        slog.Default(),
		workspaces: make(map[string]*workspace.Workspace),
	}
}

func (s *Server) RegisterWorkspace(name string, r reader.Reader) error {
	ws, err := workspace.NewWorkspace(name, r)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	s.workspaces[name] = ws
	return nil
}

func (s *Server) Stream(stream pb.NodeService_StreamServer) error {
	s.log.Debug("new client connected")
	for {
		req, err := stream.Recv()
		if err != nil {
			s.log.Debug("client disconnected", "error", err)
			return err
		}

		switch req.GetRequest().(type) {
		case *pb.StreamRequest_Subscription:
			sub := req.GetSubscription()
			go s.handleSubscriptionMessage(stream, sub)
		}
	}
}

func (s *Server) handleSubscriptionMessage(stream pb.NodeService_StreamServer, sub *pb.Subscription) {
	panic("not implemented")
}
