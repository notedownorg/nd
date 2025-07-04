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
	"fmt"

	pb "github.com/notedownorg/nd/api/go/nodes/v1alpha1"
	"github.com/notedownorg/nd/pkg/workspace"
	"github.com/notedownorg/nd/pkg/workspace/node"
)

func (s *Server) handleSubscriptionRequest(stream pb.NodeService_StreamServer, sub *pb.SubscriptionRequest) {
	ctx := stream.Context()
	subscriptionID := sub.GetSubscriptionId()

	s.log.Debug("handling subscription request", "subscription_id", subscriptionID)

	switch msg := sub.GetMsg().(type) {
	case *pb.SubscriptionRequest_Unsubscribe:
		s.handleUnsubscribe(subscriptionID)

	case *pb.SubscriptionRequest_DocumentSubscription:
		s.handleDocumentSubscription(ctx, stream, subscriptionID, msg.DocumentSubscription)

	default:
		s.log.Warn("unknown subscription message type", "subscription_id", subscriptionID)
	}
}

func (s *Server) handleUnsubscribe(subscriptionID string) {
	s.subsMu.Lock()
	defer s.subsMu.Unlock()

	sub, exists := s.subscriptions[subscriptionID]
	if !exists {
		s.log.Warn("attempted to unsubscribe from non-existent subscription", "subscription_id", subscriptionID)
		return
	}

	s.log.Debug("unsubscribing", "subscription_id", subscriptionID, "workspace", sub.workspaceName)

	if ws, exists := s.workspaces[sub.workspaceName]; exists {
		ws.Unsubscribe(sub.workspaceSubID)
	}

	sub.cancelFunc()
	delete(s.subscriptions, subscriptionID)
}

func (s *Server) handleDocumentSubscription(ctx context.Context, stream pb.NodeService_StreamServer, subscriptionID string, docSub *pb.DocumentSubscription) {
	workspaceName := docSub.GetWorkspaceName()

	// Use default workspace if none specified
	if workspaceName == "" {
		// Find the first available workspace as default
		for name := range s.workspaces {
			workspaceName = name
			break
		}
	}

	s.log.Debug("handling document subscription", "subscription_id", subscriptionID, "workspace", workspaceName)

	// Check if subscription ID already exists
	s.subsMu.RLock()
	if _, exists := s.subscriptions[subscriptionID]; exists {
		s.subsMu.RUnlock()
		s.log.Error("subscription ID already exists", "subscription_id", subscriptionID)
		s.sendError(stream, subscriptionID, pb.ErrorCode_SUBSCRIPTION_ID_CONFLICT, "subscription ID already exists")
		return
	}
	s.subsMu.RUnlock()

	ws, exists := s.workspaces[workspaceName]
	if !exists {
		s.log.Error("workspace not found", "workspace", workspaceName, "subscription_id", subscriptionID)
		s.sendError(stream, subscriptionID, pb.ErrorCode_WORKSPACE_NOT_FOUND, "workspace not found")
		return
	}

	// Send subscription confirmation
	confirmEvent := &pb.StreamEvent{
		Event: &pb.StreamEvent_SubscriptionConfirmation{
			SubscriptionConfirmation: &pb.SubscriptionConfirmation{
				SubscriptionId: subscriptionID,
			},
		},
	}
	if err := stream.Send(confirmEvent); err != nil {
		s.log.Error("failed to send subscription confirmation", "error", err, "subscription_id", subscriptionID)
		return
	}

	eventChan := make(chan workspace.Event, 1000)
	workspaceSubID := ws.Subscribe(eventChan, node.DocumentKind, true)

	ctx, cancel := context.WithCancel(ctx)

	activeSub := &activeSubscription{
		subscriptionID: subscriptionID,
		workspaceName:  workspaceName,
		eventChan:      eventChan,
		workspaceSubID: workspaceSubID,
		cancelFunc:     cancel,
	}

	s.subsMu.Lock()
	s.subscriptions[subscriptionID] = activeSub
	s.subsMu.Unlock()

	go s.forwardEvents(ctx, stream, activeSub)
}

func (s *Server) forwardEvents(ctx context.Context, stream pb.NodeService_StreamServer, sub *activeSubscription) {
	defer func() {
		s.subsMu.Lock()
		delete(s.subscriptions, sub.subscriptionID)
		s.subsMu.Unlock()

		if ws, exists := s.workspaces[sub.workspaceName]; exists {
			ws.Unsubscribe(sub.workspaceSubID)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			s.log.Debug("subscription context cancelled", "subscription_id", sub.subscriptionID)
			return

		case event, ok := <-sub.eventChan:
			if !ok {
				s.log.Debug("event channel closed", "subscription_id", sub.subscriptionID)
				return
			}

			streamEvent, err := s.convertWorkspaceEventToStreamEvent(sub.subscriptionID, event)
			if err != nil {
				s.log.Error("failed to convert workspace event", "error", err, "subscription_id", sub.subscriptionID)
				continue
			}

			if err := stream.Send(streamEvent); err != nil {
				s.log.Error("failed to send stream event", "error", err, "subscription_id", sub.subscriptionID)
				return
			}
		}
	}
}

func (s *Server) convertWorkspaceEventToStreamEvent(subscriptionID string, event workspace.Event) (*pb.StreamEvent, error) {
	switch event.Op {
	case workspace.Load:
		pbNode, err := s.convertNodeToPB(event.Node)
		if err != nil {
			return nil, fmt.Errorf("failed to convert node to protobuf: %w", err)
		}
		return &pb.StreamEvent{
			Event: &pb.StreamEvent_Load{
				Load: &pb.Load{
					SubscriptionId: subscriptionID,
					Node:           pbNode,
				},
			},
		}, nil

	case workspace.Change:
		pbNode, err := s.convertNodeToPB(event.Node)
		if err != nil {
			return nil, fmt.Errorf("failed to convert node to protobuf: %w", err)
		}
		return &pb.StreamEvent{
			Event: &pb.StreamEvent_Change{
				Change: &pb.Change{
					SubscriptionId: subscriptionID,
					Node:           pbNode,
				},
			},
		}, nil

	case workspace.Delete:
		// For delete events, create a minimal node with just the ID
		// The event.Node is nil for delete events, so we construct a minimal protobuf node
		pbNode := &pb.Node{
			Node: &pb.Node_Document{
				Document: &pb.Document{
					Id: event.Id,
				},
			},
		}
		return &pb.StreamEvent{
			Event: &pb.StreamEvent_Delete{
				Delete: &pb.Delete{
					SubscriptionId: subscriptionID,
					Node:           pbNode,
				},
			},
		}, nil

	case workspace.SubscriberLoadComplete:
		return &pb.StreamEvent{
			Event: &pb.StreamEvent_InitializationComplete{
				InitializationComplete: &pb.InitializationComplete{
					SubscriptionId: subscriptionID,
				},
			},
		}, nil

	default:
		return nil, fmt.Errorf("unknown workspace operation: %d", event.Op)
	}
}
