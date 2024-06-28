package git

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const (
	// defaultEditor is the name of the default editor program
	defaultEditor = "vi"
	// defaultPager is the name of the default pager program
	defaultPager = "less"
)

// Exec executes a git command and returns the contents of stdout and stderr.
func Exec(ctx context.Context, stdin io.Reader, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{name}, args...)...)
	cmd.Stdin = stdin

	out, err := cmd.CombinedOutput()
	switch {
	case err == nil:
		return bytes.TrimSuffix(out, []byte("\n")), nil

	case len(out) > 0:
		return nil, fmt.Errorf("%s", out)

	default:
		return nil, err
	}
}

// LaunchEditor launches a text editor opened to the file at the given path.
func LaunchEditor(ctx context.Context, path string) error {
	editor := os.Getenv("GIT_EDITOR")
	editorProgram, _ := Exec(ctx, nil, "config", "core.editor")

	term := os.Getenv("TERM")
	isTermDumb := term == "" || term == "dumb"

	if editor == "" && len(editorProgram) != 0 {
		editor = string(editorProgram)
	}
	if editor == "" && !isTermDumb {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" && isTermDumb {
		return fmt.Errorf("terminal is dumb, but EDITOR unset")
	}
	if editor == "" {
		editor = defaultEditor
	}

	cmd := exec.CommandContext(ctx, editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Pager attempts to launch a pager program with the given input and
// returns a bool indicating if the program was launched successfully.
func Pager(ctx context.Context, stdin io.Reader) (bool, error) {
	pager := os.Getenv("GIT_PAGER")
	pagerProgram, _ := Exec(ctx, nil, "config", "core.pager")

	if pager == "" && len(pagerProgram) != 0 {
		pager = string(pagerProgram)
	}
	if pager == "" {
		pager = os.Getenv("PAGER")
	}
	if pager == "" {
		pager = defaultPager
	}
	if pager == "" || pager == "cat" {
		return false, nil
	}

	cmd := exec.CommandContext(ctx, pager)
	cmd.Stdin = stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return true, cmd.Run()
}
