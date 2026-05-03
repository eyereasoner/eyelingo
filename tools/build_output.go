package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func stripGoTail(text string) string {
	for _, marker := range []string{"\n## Check", "\n## Go audit details"} {
		if pos := strings.Index(text, marker); pos >= 0 {
			text = text[:pos]
		}
	}
	return strings.TrimRight(text, "\r\n \t")
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: build_output EXAMPLE GO_STDOUT_FILE")
		os.Exit(2)
	}
	name := os.Args[1]
	rawPath := os.Args[2]
	raw, err := os.ReadFile(rawPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	prefix := stripGoTail(string(raw))
	prefixPath := rawPath + ".prefix"
	if err := os.WriteFile(prefixPath, []byte(prefix+"\n"), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	checker := os.Getenv("SEE_CHECKER")
	var cmd *exec.Cmd
	if checker != "" {
		cmd = exec.Command(checker, name, prefixPath)
	} else {
		cmd = exec.Command("go", "run", ".", name, prefixPath)
		cmd.Dir = filepath.Join(root, "examples", "checks")
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(), "SEE_ROOT="+root)
	if err := cmd.Run(); err != nil {
		if stdout.Len() > 0 {
			fmt.Print(stdout.String())
		}
		if stderr.Len() > 0 {
			fmt.Fprint(os.Stderr, stderr.String())
		}
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(prefix)
	fmt.Println()
	fmt.Print(strings.TrimRight(stdout.String(), "\r\n"))
	fmt.Println()
}
