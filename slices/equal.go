package slices

// EqualUnordered checks if slices of type E are equal, order independent.
func EqualUnordered[E comparable](a []E, b []E) bool {
	if len(a) != len(b) {
		return false
	}

	ma := make(map[E]struct{})
	for _, v := range a {
		ma[v] = struct{}{}
	}

	l := len(ma)
	for _, v := range b {
		ma[v] = struct{}{}
	}

	return len(ma) == l
}

func Equal[E comparable](a, b []E) bool {
	if len(a) != len(b) {
		return false
	}

	if (len(a) == 0) != (len(b) == 0) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
