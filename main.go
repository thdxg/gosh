package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"gosh/pkg/cmd"
	"gosh/pkg/hook"

	"github.com/chzyer/readline"
)

func main() {
	ctx := context.Background()

	prompt := "$ "
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      prompt,
		HistoryFile: "~/dev/gosh/history",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		if err := hook.PrePrompt.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run pre-run hook: %v\n", err)
		}

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
			fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
		}

		args = nil

		rl.SetPrompt(prompt)
		if err := rl.SaveHistory(strings.Join(args, " ")); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save history: %v\n", err)
		}

		if err := hook.PostPrompt.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run post-run hook: %v\n", err)
		}
	}
}
