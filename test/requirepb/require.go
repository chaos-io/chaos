package requirepb

import "github.com/chaos-io/chaos/test/assertpb"

func Equal(t assertpb.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if !assertpb.Equal(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

func Equalf(t assertpb.TestingT, expected, actual interface{}, msg string, args ...interface{}) {
	t.Helper()

	Equal(t, expected, actual, append([]interface{}{msg}, args...)...)
}
