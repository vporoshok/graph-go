package graph

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGraph(t *testing.T) {
	const n = 5

	V := make(map[interface{}]interface{}, n)
	for i := 0; i < n; i++ {
		V[i] = i
	}
	E := [][2]interface{}{
		{0, 1}, {0, 2},
		{1, 2}, {1, 3},
		{2, 3}, {2, 0},
		{3, 2}, {3, 0},
		{4, 0}, {4, 1},
	}
	g := New(V, E)

	t.Run("New", func(t *testing.T) {
		ok, err := g.AddLink(0, 1)
		require.NoError(t, err)
		require.False(t, ok)
		ok, err = g.AddLink(0, 10)
		require.Error(t, err)
		require.False(t, ok)
		ok, err = g.AddLink(3, 0)
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("AddVertex", func(t *testing.T) {
		require.False(t, g.AddVertex(0, 0))
		require.True(t, g.AddVertex(10, 10))
	})

	t.Run("AddLink", func(t *testing.T) {
		ok, err := g.AddLink(0, 3)
		require.NoError(t, err)
		require.True(t, ok)
		ok, err = g.AddLink(0, 10)
		require.NoError(t, err)
		require.True(t, ok)
		ok, err = g.AddLink(0, 10)
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("Resolve", func(t *testing.T) {
		res := g.Resolve(Out, 0)
		require.Len(t, res, 5)
		require.Contains(t, res, 0)
		require.Contains(t, res, 1)
		require.Contains(t, res, 2)
		require.Contains(t, res, 3)
		require.Contains(t, res, 10)
	})
}

func BenchmarkResolve(b *testing.B) {
	n := 500000
	step := 100000

	V := make(map[interface{}]interface{}, n)
	for i := 0; i < n; i++ {
		V[i] = i
	}
	g := New(V, nil)

	for i := 0; i < n*(n-1); i += step {
		timer := time.NewTimer(10 * time.Second)
		for j := 0; j < step; {
			select {
			case <-timer.C:
				b.Log("too long")

				return

			default:
			}
			src := rand.Intn(n)
			dst := rand.Intn(n)
			if ok, _ := g.AddLink(src, dst); ok {
				j++
			}
		}
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			b.ReportAllocs()

			for j := 0; j < b.N; j++ {
				g.Resolve(Out, 0)
			}
		})
	}
}
