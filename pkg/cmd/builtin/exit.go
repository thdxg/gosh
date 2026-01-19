package builtin

import (
	"context"
)

func Exit(ctx context.Context, arg []string) error {
	return ErrExit
}
