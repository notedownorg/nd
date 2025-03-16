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

package workspace

import (
	"log/slog"

	. "github.com/notedownorg/nd/pkg/workspace/node"
	"github.com/notedownorg/nd/pkg/workspace/reader"
)

type Workspace struct {
	log     *slog.Logger
	closers []func()

	// Nodes stored by kind for faster lookup
	// ids are prefixed with the kind if we need to workout which map to look in see kindFromID
	documents map[string]*Document
}

func NewWorkspace(name string, r reader.Reader) (*Workspace, error) {
	ws := &Workspace{
		log:       slog.Default().With("workspace", name),
		documents: make(map[string]*Document),
		closers:   make([]func(), 0, 2),
	}

	readerSubscription := make(chan reader.Event)
	subId := r.Subscribe(readerSubscription, true)
	ws.closers = append(ws.closers, func() { r.Unsubscribe(subId); close(readerSubscription) })

	go ws.processDocuments(readerSubscription)

	return ws, nil
}

func (w *Workspace) Close() {
	for _, closer := range w.closers {
		closer()
	}
}

func (w *Workspace) deleteNode(id string) {
	kind := KindFromID(id)
	switch kind {
	case DocumentKind:
		delete(w.documents, id)
		// TODO: verify that we don't need to delete nodes not accessible from the workspace struct
		// I think Go garbage collection should take care of this for us?
	default:
		return
	}
	w.log.Debug("deleted node", "id", id)
}
