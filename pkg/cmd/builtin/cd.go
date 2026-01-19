package builtin

import (
	"context"
	"fmt"
	"os"
)

func Cd(ctx context.Context, arg []string) error {
	if len(arg) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}
	dir := arg[0]
	return os.Chdir(dir)
}
