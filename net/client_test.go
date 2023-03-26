/*
Copyright © 2021 Thomas Meitz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package net_test

import (
	"testing"

	"github.com/oscarzhou/ksqldb-go/net"
	"github.com/stretchr/testify/require"
)

func TestClient_NotNil(t *testing.T) {
	client, err := net.NewHTTPClient(net.Options{}, nil)

	// Ensures that the HTTPClient interface is implemented.
	// Aborts the compiler if it does not.
	var _ net.HTTPClient = &client

	require.NotNil(t, client)
	require.Nil(t, err)
}

func TestClient_Nil(t *testing.T) {
	client, err := net.NewHTTPClient(net.Options{BaseUrl: "sf"}, nil)
	require.NotNil(t, client)
	require.NotNil(t, err)
	require.Equal(t, "invalid host name given", err.Error())
}
