// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cmd

import (
	"context"
	"net"
	"os"
	"testing"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEndpointPrecedenceWithHostShellFallback(t *testing.T) {
	const legacyKey = "AZURE_AI_PROJECT_ENDPOINT"
	const canonicalKey = "FOUNDRY_PROJECT_ENDPOINT"
	const legacyShell = "https://legacyshell.services.ai.azure.com/api/projects/project"
	const canonicalShell = "https://canonicalshell.services.ai.azure.com/api/projects/project"
	const persisted = "https://persisted.services.ai.azure.com/api/projects/project"
	const legacyPersisted = "https://legacypersisted.services.ai.azure.com/api/projects/project"
	const global = "https://global.services.ai.azure.com/api/projects/project"
	t.Setenv("AZD_SERVER", "")
	for _, tt := range []struct {
		name        string
		values      map[string]string
		legacyShell string
		foundry     string
		noEnv       bool
		noGlobal    bool
		want        string
	}{
		{
			name:    "persisted canonical key wins",
			values:  map[string]string{canonicalKey: persisted},
			foundry: canonicalShell, legacyShell: legacyShell,
			want: persisted,
		},
		{
			name:    "persisted canonical key wins over persisted legacy key",
			values:  map[string]string{canonicalKey: persisted, legacyKey: legacyPersisted},
			foundry: canonicalShell, legacyShell: legacyShell,
			want: persisted,
		},
		{
			name:   "persisted canonical key only (issue 10407)",
			values: map[string]string{canonicalKey: persisted},
			want:   persisted,
		},
		{
			name:   "persisted canonical key beats global config",
			values: map[string]string{canonicalKey: persisted},
			want:   persisted,
		},
		{
			name:   "persisted legacy key used when canonical key is explicitly empty",
			values: map[string]string{canonicalKey: "", legacyKey: legacyPersisted},
			want:   legacyPersisted,
		},
		{
			name:        "persisted legacy key beats shell values when canonical key is empty",
			values:      map[string]string{canonicalKey: "", legacyKey: legacyPersisted},
			legacyShell: legacyShell,
			want:        legacyPersisted,
		},
		{
			name:    "missing canonical key uses canonical shell before global",
			foundry: canonicalShell, legacyShell: legacyShell,
			want: canonicalShell,
		},
		{
			name:        "missing keys use legacy shell before global",
			legacyShell: legacyShell,
			want:        legacyShell,
		},
		{
			name:    "empty keys defer to global",
			values:  map[string]string{canonicalKey: "", legacyKey: ""},
			foundry: canonicalShell, legacyShell: legacyShell,
			want: global,
		},
		{
			name:    "no active environment uses global first",
			noEnv:   true,
			foundry: canonicalShell, legacyShell: legacyShell,
			want: global,
		},
		{
			name:     "no environment or global uses canonical shell not legacy",
			noEnv:    true,
			noGlobal: true,
			foundry:  canonicalShell, legacyShell: legacyShell,
			want: canonicalShell,
		},
		{
			name:     "empty keys without global reach canonical shell",
			values:   map[string]string{canonicalKey: "", legacyKey: ""},
			noGlobal: true,
			foundry:  canonicalShell, legacyShell: legacyShell,
			want: canonicalShell,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(legacyKey, tt.legacyShell)
			t.Setenv(canonicalKey, tt.foundry)
			host := &precedenceHost{values: tt.values, noEnv: tt.noEnv}
			server := grpc.NewServer()
			azdext.RegisterEnvironmentServiceServer(server, host)
			azdext.RegisterUserConfigServiceServer(server, &precedenceConfig{noGlobal: tt.noGlobal, global: global})
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			go func() { _ = server.Serve(listener) }()
			t.Cleanup(func() { server.Stop(); _ = listener.Close() })
			t.Setenv("AZD_SERVER", listener.Addr().String())
			got, err := resolveProjectEndpoint(t.Context(), "")
			require.NoError(t, err)
			require.Equal(t, tt.want, got.Endpoint)
		})
	}
}

type precedenceHost struct {
	azdext.UnimplementedEnvironmentServiceServer
	values map[string]string
	noEnv  bool
}

func (s *precedenceHost) GetCurrent(context.Context, *azdext.EmptyRequest) (*azdext.EnvironmentResponse, error) {
	if s.noEnv {
		return nil, status.Error(codes.NotFound, "no active environment")
	}
	return &azdext.EnvironmentResponse{Environment: &azdext.Environment{Name: "dev"}}, nil
}

func (s *precedenceHost) GetValue(_ context.Context, req *azdext.GetEnvRequest) (*azdext.KeyValueResponse, error) {
	// Mirror the host's Environment.Getenv contract: only absent keys fall back to the process environment.
	value, present := s.values[req.Key]
	if !present {
		value = os.Getenv(req.Key)
	}
	return &azdext.KeyValueResponse{Key: req.Key, Value: value}, nil
}

type precedenceConfig struct {
	azdext.UnimplementedUserConfigServiceServer
	noGlobal bool
	global   string
}

func (s *precedenceConfig) Get(
	_ context.Context, req *azdext.GetUserConfigRequest,
) (*azdext.GetUserConfigResponse, error) {
	if s.noGlobal || req.Path != "extensions.ai-projects.context" {
		return &azdext.GetUserConfigResponse{}, nil
	}
	return &azdext.GetUserConfigResponse{Found: true, Value: []byte(`{"endpoint":"` + s.global + `"}`)}, nil
}
