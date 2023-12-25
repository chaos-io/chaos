package docker

import (
	"context"
	"testing"

	"github.com/chaos-io/chaos/logs"
)

func TestRun(t *testing.T) {
	// docker run -it --rm -v /opt/abfuzz/data/2cbbe09e45d3189965f0144f7cb31fe8/fuzzings/harness-cov-2ZtDl0iIAzdVcw71O0ZnaTBmqvE:/src abfuzz-abfast-linux-arm64-ubuntu-20_04.2ztdcmsgqbpozd5wk5bnxgvubrd /bin/bash -c genjson.sh /src/harness /src/fresh_files.txt
	cmd := []string{"/bin/bash", "-c", "genjson.sh /src/harness /src/fresh_files.txt"}
	_, _, err := Run(context.Background(), "abfuzz-abfast-linux-arm64-ubuntu-20_04.2ztdcmsgqbpozd5wk5bnxgvubrd", "", cmd, nil, "/opt/abfuzz/data/2cbbe09e45d3189965f0144f7cb31fe8/fuzzings/harness-cov-2ZtDl0iIAzdVcw71O0ZnaTBmqvE", "/src")
	if err != nil {
		logs.Errorw("----", "error", err)
	}
}
