package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gosh/pkg/cmd"
	"gosh/pkg/cmd/builtin"
	"gosh/pkg/hook"

	"github.com/ergochat/readline"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configDir, ok := os.LookupEnv("GOSH_CONFIG_DIR")
	if !ok {
		panic("GOSH_CONFIG_DIR not set")
	}

	historyFile := filepath.Join(configDir, "history")

	rl, err := readline.NewEx(&readline.Config{
		HistoryFile: historyFile,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	hook.PostCommand.Register(func() error {
		fmt.Println()
		return nil
	})

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			if err := hook.PrePrompt.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to run pre-run hook: %v\n", err)
			}

			pwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get working directory")
			}
			prompt := fmt.Sprintf("\x1b[36m%s\n> \x1b[0m", pwd)
			rl.SetPrompt(prompt)

			args := make([]string, 0)
			multiline := true

			for multiline {
				line, err := rl.Readline()
				if err != nil {
					break
				}

				line = strings.TrimSpace(line)
				multiline = strings.HasSuffix(line, "\\")
				line = strings.TrimSuffix(line, "\\")
				line = strings.TrimSpace(line)

				if len(line) > 0 {
					args = slices.Concat(
						args,
						strings.Split(line, " "),
					)
				}

				if multiline {
					rl.SetPrompt(">>> ")
				}
			}

			if err := cmd.Exec(ctx, args); err != nil {
				if errors.Is(err, builtin.ErrExit) {
					break loop
				}
				fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
			}

			if err := rl.SaveToHistory(strings.Join(args, " ")); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save history: %v\n", err)
			}

			if err := hook.PostCommand.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to run post-run hook: %v\n", err)
			}
		}
	}
}
