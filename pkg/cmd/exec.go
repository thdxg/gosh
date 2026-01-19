package cmd

import (
	"context"
	"gosh/pkg/cmd/builtin"
	"os"
	"os/exec"
)

func Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return nil
	}

	name, arg := args[0], args[1:]

	matched, err := builtin.Handle(ctx, name, arg)
	if err != nil {
		return err
	}
	if matched {
		return nil
	}

	// run independently from the parent context
	cmd := exec.CommandContext(context.Background(), name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
