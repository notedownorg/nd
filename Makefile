# Copyright 2025 Notedown Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Default target - run full hygiene check
hygiene: tidy generate format licenser dirty

# Available targets:
# build         - Build the nd server binary
# test-unit     - Run unit tests only (no race detection)
# test-unit-race - Run unit tests with race detection  
# test-functional - Build binary and run functional tests
# test-all      - Run unit and functional tests
# tidy          - Tidy go modules
# generate      - Generate protobuf code
# format        - Format Go code
# licenser      - Apply license headers
# dirty         - Check for uncommitted changes

licenser:
	nix develop --command licenser apply -r "Notedown Authors"

tidy:
	nix develop --command go mod tidy

format:
	nix develop --command gofmt -w .

generate:
	nix develop --command buf generate
	nix develop --command go generate ./...

dirty:
	nix develop --command git diff --exit-code

build:
	nix develop --command go build -o ./bin/nd ./cmd/nd

test-unit:
	nix develop --command go test -timeout=30s ./pkg/...

test-unit-race:
	nix develop --command go test -race -timeout=30s ./pkg/...

test-functional: build
	nix develop --command go test -timeout=60s ./tests/functional/...

test-all: test-unit test-functional

