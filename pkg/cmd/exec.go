package cmd

import (
	"context"
	"io"
	"os"
	"os/exec"
)

func Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return nil
	}

	name, arg := args[0], args[1:]

	cmd := exec.CommandContext(ctx, name, arg...)

	op, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer op.Close()

	ep, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer ep.Close()

	if err := cmd.Start(); err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, op)
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stderr, ep)
	if err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
