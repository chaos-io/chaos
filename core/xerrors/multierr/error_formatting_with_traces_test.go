package multierr

import (
	"testing"

	"github.com/chaos-io/chaos/core/xerrors"
	"github.com/chaos-io/chaos/core/xerrors/assertxerrors"
)

func TestFormatWithTracesPerMode(t *testing.T) {
	t.Parallel()

	constructor := func(t *testing.T) error {
		t.Helper()

		return Append(
			xerrors.New("foo"),
			xerrors.New("bar"),
		)
	}

	expected := assertxerrors.Expectations{
		ExpectedS: "foo; bar",
		ExpectedV: "foo; bar",
		Frames: assertxerrors.NewStackTraceModeExpectation(`
foo
      github.com/chaos-io/chaos/core/xerrors/multierr.TestFormatWithTracesPerMode.func1
          library/go/core/xerrors/multierr/error_formatting_with_traces_test.go:17

bar
      github.com/chaos-io/chaos/core/xerrors/multierr.TestFormatWithTracesPerMode.func1
          library/go/core/xerrors/multierr/error_formatting_with_traces_test.go:18
`),
		Stacks: assertxerrors.NewStackTraceModeExpectation(`
foo
      github.com/chaos-io/chaos/core/xerrors/multierr.TestFormatWithTracesPerMode.func1
          library/go/core/xerrors/multierr/error_formatting_with_traces_test.go:17
      github.com/chaos-io/chaos/core/xerrors/assertxerrors.RunTestsPerMode.func1
          /Users/djerys/arcadia/library/go/core/xerrors/assertxerrors/assertxerrors.go:18
      testing.tRunner
          /Users/djerys/.ya/tools/v4/1217046835/src/testing/testing.go:1259

bar
      github.com/chaos-io/chaos/core/xerrors/multierr.TestFormatWithTracesPerMode.func1
          library/go/core/xerrors/multierr/error_formatting_with_traces_test.go:18
      github.com/chaos-io/chaos/core/xerrors/assertxerrors.RunTestsPerMode.func1
          /Users/djerys/arcadia/library/go/core/xerrors/assertxerrors/assertxerrors.go:18
      testing.tRunner
          /Users/djerys/.ya/tools/v4/1217046835/src/testing/testing.go:1259
`, 3, 4, 5, 6, 10, 11, 12, 13),
		StackThenFrames: assertxerrors.NewStackTraceModeExpectation(`
foo
      github.com/chaos-io/chaos/core/xerrors/multierr.TestFormatWithTracesPerMode.func1
          library/go/core/xerrors/multierr/error_formatting_with_traces_test.go:17
      github.com/chaos-io/chaos/core/xerrors/assertxerrors.RunTestsPerMode.func1
          /Users/djerys/arcadia/library/go/core/xerrors/assertxerrors/assertxerrors.go:18
      testing.tRunner
          /Users/djerys/.ya/tools/v4/1217046835/src/testing/testing.go:1259

bar
      github.com/chaos-io/chaos/core/xerrors/multierr.TestFormatWithTracesPerMode.func1
          library/go/core/xerrors/multierr/error_formatting_with_traces_test.go:18
      github.com/chaos-io/chaos/core/xerrors/assertxerrors.RunTestsPerMode.func1
          /Users/djerys/arcadia/library/go/core/xerrors/assertxerrors/assertxerrors.go:18
      testing.tRunner
          /Users/djerys/.ya/tools/v4/1217046835/src/testing/testing.go:1259
`, 3, 4, 5, 6, 10, 11, 12, 13),
		StackThenNothing: assertxerrors.NewStackTraceModeExpectation(`
foo
      github.com/chaos-io/chaos/core/xerrors/multierr.TestFormatWithTracesPerMode.func1
          library/go/core/xerrors/multierr/error_formatting_with_traces_test.go:17
      github.com/chaos-io/chaos/core/xerrors/assertxerrors.RunTestsPerMode.func1
          /Users/djerys/arcadia/library/go/core/xerrors/assertxerrors/assertxerrors.go:18
      testing.tRunner
          /Users/djerys/.ya/tools/v4/1217046835/src/testing/testing.go:1259

bar
      github.com/chaos-io/chaos/core/xerrors/multierr.TestFormatWithTracesPerMode.func1
          library/go/core/xerrors/multierr/error_formatting_with_traces_test.go:18
      github.com/chaos-io/chaos/core/xerrors/assertxerrors.RunTestsPerMode.func1
          /Users/djerys/arcadia/library/go/core/xerrors/assertxerrors/assertxerrors.go:18
      testing.tRunner
          /Users/djerys/.ya/tools/v4/1217046835/src/testing/testing.go:1259
`, 3, 4, 5, 6, 10, 11, 12, 13),
		Nothing: assertxerrors.NewStackTraceModeExpectation(`
foo
bar`),
	}

	assertxerrors.RunTestsPerMode(t, expected, constructor)
}
