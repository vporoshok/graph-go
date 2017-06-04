package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVertex_AddLink(t *testing.T) {
	v := &vertex{}
	require.True(t, v.addLink(In, &vertex{id: 5}))
	require.True(t, v.addLink(In, &vertex{id: 4}))
	require.True(t, v.addLink(In, &vertex{id: 6}))
	require.False(t, v.addLink(In, &vertex{id: 4}))
	require.False(t, v.addLink(In, &vertex{id: 5}))
	require.False(t, v.addLink(In, &vertex{id: 6}))
}
