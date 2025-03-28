// Copyright 2024 The PipeCD Authors.
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

package pipedapi

import (
	"context"
	"slices"

	"google.golang.org/grpc"

	service "github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
)

type PluginServiceClient struct {
	service.PluginServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, address string, opts ...rpcclient.DialOption) (*PluginServiceClient, error) {
	// Clone the opts to avoid modifying the original opts slice.
	opts = slices.Clone(opts)

	// Append the required options.
	// The WithBlock option is required to make the client wait until the connection is up.
	// The WithInsecure option is required to disable the transport security.
	// The piped service does not require transport security because it is only used in localhost.
	opts = append(opts, rpcclient.WithBlock(), rpcclient.WithInsecure())

	conn, err := rpcclient.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, err
	}

	return &PluginServiceClient{
		PluginServiceClient: service.NewPluginServiceClient(conn),
		conn:                conn,
	}, nil
}

func (c *PluginServiceClient) Close() error {
	return c.conn.Close()
}
