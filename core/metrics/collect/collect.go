package collect

import (
	"context"

	"github.com/chaos-io/chaos/core/metrics"
)

type Func func(ctx context.Context, r metrics.Registry, c metrics.CollectPolicy)
