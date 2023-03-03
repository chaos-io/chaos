package slices

import (
	"sort"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Dedup removes duplicate values from slice.
// It will alter original non-empty slice, consider copy it beforehand.
func Dedup[E constraints.Ordered](s []E) []E {
	if len(s) < 2 {
		return s
	}
	slices.Sort(s)
	tmp := s[:1]
	cur := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != cur {
			tmp = append(tmp, s[i])
			cur = s[i]
		}
	}
	return tmp
}

// DedupBools removes duplicate values from bool slice.
// It will alter original non-empty slice, consider copy it beforehand.
func DedupBools(a []bool) []bool {
	if len(a) < 2 {
		return a
	}
	sort.Slice(a, func(i, j int) bool { return a[i] != a[j] })
	tmp := a[:1]
	cur := a[0]
	for i := 1; i < len(a); i++ {
		if a[i] != cur {
			tmp = append(tmp, a[i])
			cur = a[i]
		}
	}
	return tmp
}
