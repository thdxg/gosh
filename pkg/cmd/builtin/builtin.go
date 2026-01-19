package builtin

import "context"

const (
	CmdCd   = "cd"
	CmdExit = "exit"
)

func Handle(ctx context.Context, name string, arg []string) (matched bool, err error) {
	switch name {
	case CmdCd:
		return true, Cd(ctx, arg)
	case CmdExit:
		return true, Exit(ctx, arg)
	}

	return false, nil
}
