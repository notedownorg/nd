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

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	pb "github.com/notedownorg/nd/api/go/nodes/v1alpha1"
	"github.com/notedownorg/nd/pkg/server"
	"github.com/notedownorg/nd/pkg/workspace/reader/filesystem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		port          = flag.String("port", "8080", "Port to listen on")
		workspaceDir  = flag.String("workspace", ".", "Workspace directory path")
		workspaceName = flag.String("name", "default", "Workspace name")
		logLevel      = flag.String("log", "info", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	// Set up structured logging
	var level slog.Level
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	// Validate workspace directory
	absWorkspaceDir, err := filepath.Abs(*workspaceDir)
	if err != nil {
		logger.Error("Invalid workspace directory", "error", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absWorkspaceDir); os.IsNotExist(err) {
		logger.Error("Workspace directory does not exist", "path", absWorkspaceDir)
		os.Exit(1)
	}

	logger.Info("Starting nd streaming server",
		"port", *port,
		"workspace", absWorkspaceDir,
		"workspace_name", *workspaceName,
		"log_level", *logLevel)

	// Create filesystem reader
	reader, err := filesystem.NewReader(absWorkspaceDir)
	if err != nil {
		logger.Error("Failed to create filesystem reader", "error", err)
		os.Exit(1)
	}
	defer reader.Close()

	// Create server
	srv := server.NewServer()

	// Register workspace with server
	err = srv.RegisterWorkspace(*workspaceName, reader)
	if err != nil {
		logger.Error("Failed to register workspace", "error", err)
		os.Exit(1)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterNodeServiceServer(grpcServer, srv)

	// Enable reflection for easier debugging
	reflection.Register(grpcServer)

	// Listen on specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		logger.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	logger.Info("Server listening", "address", lis.Addr().String())

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, stopping server...")
		grpcServer.GracefulStop()
	}()

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Failed to serve", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped")
}
