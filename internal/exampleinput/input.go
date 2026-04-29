package exampleinput

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Load reads examples/input/<exampleName>.json into the same concrete type as
// fallback. The fallback value is only used to infer the target type.
func Load[T any](exampleName string, fallback T) T {
	path, data, err := Read(exampleName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input/%s.json: %v\n", exampleName, err)
		os.Exit(1)
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse %s: %v\n", path, err)
		os.Exit(1)
	}
	return value
}

// Read returns the path and bytes for examples/input/<exampleName>.json.
// Multiple candidate paths keep examples runnable from the repository root,
// from the examples directory, and from compiled binaries.
func Read(exampleName string) (string, []byte, error) {
	candidates := []string{
		filepath.Join("examples", "input", exampleName+".json"),
		filepath.Join("input", exampleName+".json"),
	}
	if _, file, _, ok := runtime.Caller(0); ok {
		root := filepath.Dir(filepath.Dir(filepath.Dir(file)))
		candidates = append(candidates, filepath.Join(root, "examples", "input", exampleName+".json"))
	}

	var lastErr error
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err == nil {
			if !json.Valid(data) {
				return path, data, fmt.Errorf("file exists but is not valid JSON")
			}
			return path, data, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = os.ErrNotExist
	}
	return candidates[0], nil, lastErr
}
