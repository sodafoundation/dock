// Copyright 2019 The OpenSDS Authors.
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

package constants

const (
	// It's RFC 8601 format that decodes and encodes with
	// exactly precision to seconds.
	TimeFormat = `2006-01-02T15:04:05`

	// OpenSDS current api version
	APIVersion = "v1beta"

	// OpensdsConfigPath indicates the absolute path of opensds global
	// configuration file.
	OpensdsConfigPath = "/etc/opensds/opensds.conf"

	// OpensdsDockBindEndpoint indicates the bind endpoint which the opensds
	// dock grpc server would listen to.
	OpensdsDockBindEndpoint = "0.0.0.0:50050"

	//Storage type for profile
	Block = "block"
	File  = "file"

	//StorageAccessCApability enum constants for profile
	Read    = "Read"
	Write   = "Write"
	Execute = "Execute"

	// Default value for pagination and sorting
	DefaultSortDir = "desc"
	DefaultLimit   = 50
	DefaultOffset  = 0
)
