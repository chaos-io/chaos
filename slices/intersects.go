package slices

// Intersection returns intersection for slices of various built-in types
func Intersection[E comparable](a, b []E) []E {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}

	p, s := a, b
	if len(b) > len(a) {
		p, s = b, a
	}

	m := make(map[E]struct{})
	for _, i := range p {
		m[i] = struct{}{}
	}

	var res []E
	for _, v := range s {
		if _, exists := m[v]; exists {
			res = append(res, v)
		}
	}

	return res
}
