GO_LIBRARY()

OWNER(g:go-library)

SRCS(singleinflight.go)

GO_TEST_SRCS(singleinflight_test.go)

END()

RECURSE(gotest)
