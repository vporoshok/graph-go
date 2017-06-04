package graph

type labels []uint8

func newLabels(n int) labels {
	l := make(labels, n/8+1)
	for i := range l {
		l[i] &= 0
	}

	return l
}

func (l labels) check(n int) bool {
	return l[n/8]&(uint8(1)<<uint(n%8)) > 0
}

func (l labels) set(n int) {
	l[n/8] |= uint8(1) << uint(n%8)
}
