package graph

import (
	"fmt"
)

type direction int8

const (
	// In is a incoming direction
	In direction = iota
	// Out is a outgoing direction
	Out
)

type vertex struct {
	id   int
	name interface{}
	data interface{}
	link [2][]*vertex
}

func (v *vertex) String() string {
	return fmt.Sprintf("%d:%s", v.id, v.name)
}

func (v *vertex) addLink(dir direction, u *vertex) bool {
	if len(v.link[dir]) == 0 {
		v.link[dir] = append(v.link[dir], u)

		return true
	}

	var l, r, c int

	l = 0
	r = len(v.link[dir]) - 1
	for l < r {
		c = l + (r-l)/2
		if v.link[dir][c].id == u.id {

			return false
		}
		if v.link[dir][c].id < u.id {
			l = c + 1
		} else {
			r = c
		}
	}
	if v.link[dir][r].id == u.id {

		return false
	}
	if v.link[dir][r].id < u.id {
		v.link[dir] = append(v.link[dir], u)
	} else {
		v.link[dir] = append(v.link[dir], nil)
		copy(v.link[dir][r+1:], v.link[dir][r:])
		v.link[dir][r] = u
	}

	return true
}
