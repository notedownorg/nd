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
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	pb "github.com/notedownorg/nd/api/go/nodes/v1alpha1"
	"github.com/notedownorg/nd/pkg/workspace"
	"github.com/notedownorg/nd/pkg/workspace/node"
	"github.com/notedownorg/nd/pkg/workspace/reader"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type activeSubscription struct {
	subscriptionID string
	workspaceName  string
	eventChan      chan workspace.Event
	workspaceSubID int
	cancelFunc     context.CancelFunc
}

type Server struct {
	pb.UnimplementedNodeServiceServer
	log *slog.Logger

	workspaces    map[string]*workspace.Workspace
	subscriptions map[string]*activeSubscription
	subsMu        sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		log:           slog.Default(),
		workspaces:    make(map[string]*workspace.Workspace),
		subscriptions: make(map[string]*activeSubscription),
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
			go s.handleSubscriptionRequest(stream, sub)
		}
	}
}


func (s *Server) convertNodeToPB(n node.Node) (*pb.Node, error) {
	switch typedNode := n.(type) {
	case *node.Document:
		pbDoc := &pb.Document{
			Id:        typedNode.ID(),
			Workspace: "", // TODO: Get workspace name from context
		}

		// Convert children IDs
		typedNode.DepthFirstSearch(func(child node.Node) {
			if child != typedNode { // Don't include self
				pbDoc.Children = append(pbDoc.Children, child.ID())
			}
		})

		// Convert metadata from document frontmatter
		pbDoc.Metadata = s.convertMetadataToPB(typedNode)

		return &pb.Node{
			Node: &pb.Node_Document{
				Document: pbDoc,
			},
		}, nil

	case *node.Section:
		pbSection := &pb.Section{
			Id:     typedNode.ID(),
			Parent: "", // TODO: Get parent ID
			Title:  "", // TODO: Get title from section
			Level:  0,  // TODO: Get level from section
		}

		return &pb.Node{
			Node: &pb.Node_Section{
				Section: pbSection,
			},
		}, nil

	default:
		return nil, fmt.Errorf("unknown node type: %T", n)
	}
}

func (s *Server) convertMetadataToPB(doc *node.Document) []*pb.Document_MetadataEntry {
	// Get metadata directly from the document
	metadata := doc.GetMetadata()
	if metadata == nil {
		return nil
	}
	
	// Convert to protobuf metadata entries
	var entries []*pb.Document_MetadataEntry
	for key, value := range metadata {
		entry := &pb.Document_MetadataEntry{
			Key: key,
		}
		
		// Convert value to protobuf Any
		if anyValue, err := s.convertValueToAny(value); err == nil {
			entry.Value = anyValue
		} else {
			s.log.Error("failed to convert metadata value", "key", key, "error", err)
		}
		
		entries = append(entries, entry)
	}
	
	return entries
}

func (s *Server) convertValueToAny(value interface{}) (*anypb.Any, error) {
	// Convert common types to protobuf Any
	switch v := value.(type) {
	case string:
		return anypb.New(&wrapperspb.StringValue{Value: v})
	case int:
		return anypb.New(&wrapperspb.Int64Value{Value: int64(v)})
	case int64:
		return anypb.New(&wrapperspb.Int64Value{Value: v})
	case float64:
		return anypb.New(&wrapperspb.DoubleValue{Value: v})
	case bool:
		return anypb.New(&wrapperspb.BoolValue{Value: v})
	default:
		// For complex types, convert to JSON string
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		return anypb.New(&wrapperspb.StringValue{Value: string(jsonBytes)})
	}
}

func (s *Server) sendError(stream pb.NodeService_StreamServer, requestID string, errorCode pb.ErrorCode, message string) {
	errorEvent := &pb.StreamEvent{
		Event: &pb.StreamEvent_Error{
			Error: &pb.Error{
				RequestId:    requestID,
				ErrorCode:    errorCode,
				ErrorMessage: message,
			},
		},
	}

	if err := stream.Send(errorEvent); err != nil {
		s.log.Error("failed to send error event", "error", err, "request_id", requestID)
	}
}
