package idgen

import "context"

//go:generate mockgen -destination=mocks/idgen.go -package=mocks . IDGenerator
type IDGenerator interface {
	GenID(ctx context.Context) (int64, error)
	GenMultiIDs(ctx context.Context, counts int) ([]int64, error)
}
