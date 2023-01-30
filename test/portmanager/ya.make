GO_LIBRARY()

OWNER(
    prime
    g:go-library
)

SRCS(manager.go)

GO_XTEST_SRCS(manager_test.go)

IF (OS_LINUX)
    SRCS(
        manager_linux.go
        manager_unix.go
    )
ENDIF()

IF (OS_DARWIN)
    SRCS(
        manager_darwin.go
        manager_unix.go
    )
ENDIF()

IF (OS_WINDOWS)
    SRCS(manager_other.go)
ENDIF()

END()

RECURSE(
    burn_ports
    gotest
)
