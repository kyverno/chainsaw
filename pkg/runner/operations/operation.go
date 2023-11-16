package operations

import (
	"context"
)

type Operation interface {
	Exec(ctx context.Context) error
	Cleanup()
}
