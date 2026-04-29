# eyelingo

Small Go translations of selected EyeReasoner/Eyeling N3 examples.

## Layout

```text
go.mod                         local module so examples can share input loading
internal/exampleinput/          shared JSON input loader
examples/                       Go examples
examples/input/                 example-specific JSON data and parameters
examples/output/                expected text output for each example
test                            run examples and compare with expected text output
```

## Example structure

Each example now has three pieces:

```text
examples/example_xyz.go
examples/input/example_xyz.json
examples/output/example_xyz.txt
```

The Go file contains the example logic and prints the original ARC-style text output. The matching input JSON file contains the example-specific facts, data, or parameters that are feasible to externalize.

The output remains text, with the existing structured sections such as:

```text
=== Answer ===
...

=== Reason Why ===
...

=== Check ===
...

=== Go audit details ===
...
```

Most examples load their domain fixture directly from `examples/input/<name>.json` through `exampleinput.Load`. A few examples still keep complex internal relation structures in Go, but they also have a matching JSON input file documenting the corresponding data or parameters.

## Run

Run one example from the repository root:

```sh
go run examples/bmi.go
```

The program writes the original text output to stdout.

Run the full regression test:

```sh
./test
```

The test prints `OK` or `FAIL` for each example, per-example timing, and total time. It compares against `examples/output/*.txt` while normalizing the Go platform audit value, since that varies by machine.

Regenerate expected outputs after intentional changes:

```sh
./test --update
```
