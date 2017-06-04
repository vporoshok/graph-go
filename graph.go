package graph

import (
	"sync"

	"github.com/pkg/errors"
)

// Graph is a main object of this package
type Graph struct {
	sync.RWMutex
	vertices map[interface{}]*vertex

	labPool sync.Pool
	verPool sync.Pool
}

// New Graph constructor
func New(V map[interface{}]interface{}, E [][2]interface{}) *Graph {
	g := &Graph{
		vertices: make(map[interface{}]*vertex, len(V)),
	}
	g.labPool = sync.Pool{
		New: func() interface{} {
			return newLabels(len(g.vertices))
		},
	}
	g.verPool = sync.Pool{
		New: func() interface{} {
			return make([]*vertex, len(g.vertices))[:0]
		},
	}
	i := 0
	for name, data := range V {
		g.vertices[name] = &vertex{
			id:   i,
			name: name,
			data: data,
		}
		i++
	}
	for _, e := range E {
		if u, ok := g.vertices[e[0]]; ok {
			if v, ok := g.vertices[e[1]]; ok {
				g.addLink(u, v)
			}
		}
	}

	return g
}

// AddVertex to graph
func (g *Graph) AddVertex(name, data interface{}) bool {
	g.Lock()

	defer g.Unlock()

	if _, ok := g.vertices[name]; ok {

		return false
	}
	g.vertices[name] = &vertex{
		id:   len(g.vertices),
		name: name,
		data: data,
	}

	return true
}

// AddLink from src to dst
//
// Return error if one of vertex not found. Return true if link created
// and false if link already exists. Concurrent safe.
func (g *Graph) AddLink(src, dst interface{}) (bool, error) {
	g.Lock()

	defer g.Unlock()

	if u, ok := g.vertices[src]; ok {
		if v, ok := g.vertices[dst]; ok {

			return g.addLink(u, v), nil
		}

		return false, errors.Errorf("vertex %q not found", src)
	}

	return false, errors.Errorf("vertex %q not found", src)

}

func (g *Graph) addLink(u, v *vertex) bool {
	if len(u.link[Out]) < len(v.link[In]) {
		if u.addLink(Out, v) {

			return v.addLink(In, u)
		}

		return false
	}
	if v.addLink(In, u) {

		return u.addLink(Out, v)
	}

	return false
}

// Resolve all achievable vertices from given srcs in given direction
func (g *Graph) Resolve(dir direction, srcs ...interface{}) []interface{} {
	lab := g.labPool.Get().(labels)
	res := g.verPool.Get().([]*vertex)

	for _, src := range srcs {
		if v, ok := g.vertices[src]; ok {
			lab.set(v.id)
			res = append(res, v)
		}
	}

	for i := 0; i < len(res) && len(res) < len(g.vertices); i++ {
		for _, v := range res[i].link[dir] {
			if !lab.check(v.id) {
				lab.set(v.id)
				res = append(res, v)
			}
		}
	}
	g.labPool.Put(lab)

	dd := make([]interface{}, len(res))
	for i, v := range res {
		dd[i] = v.data
	}
	g.verPool.Put(res)

	return dd
}
