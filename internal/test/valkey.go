package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	valkeyCnt "github.com/testcontainers/testcontainers-go/modules/valkey"
	"github.com/valkey-io/valkey-go"
)

var ValkeyAddr string
var ValkeyCl valkey.Client

func SetupValkey(t *testing.T) {
	valkeyContainer, err := valkeyCnt.Run(t.Context(), "valkey/valkey-bundle:8.1.3")
	require.NoError(t, err, "failed to start valkey")
	testcontainers.CleanupContainer(t, valkeyContainer)

	host, err := valkeyContainer.Host(t.Context())
	require.NoError(t, err)
	port, err := valkeyContainer.MappedPort(t.Context(), "6379/tcp")
	require.NoError(t, err)

	addr := host + ":" + port.Port()
	ValkeyAddr = addr
	opt := valkey.ClientOption{InitAddress: []string{addr}}
	client, err := valkey.NewClient(opt)
	require.NoError(t, err, "can not init valkey")
	ValkeyCl = client
}

func CleanCache(t *testing.T) {
	ValkeyCl.Do(t.Context(), ValkeyCl.B().Flushall().Build())
}
