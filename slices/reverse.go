package slices

// Reverse reverses given slice.
// It will alter original non-empty slice, consider copy it beforehand.
func Reverse[E any](s []E) []E {
	if len(s) < 2 {
		return s
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
