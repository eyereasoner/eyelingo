# eyelingo

Small Go translations of selected EyeReasoner/Eyeling N3 examples. The suite includes route planning, air-route path discovery over the full airport/flight graph, a 100,000-level deep taxonomy benchmark, privacy policy checks, ranked risk reasoning, combinatorial search, an exact crisis-dispatch optimizer, an exact biology marker-panel selector, a tamper-evident cold-chain release allocator, explicit complex-number arithmetic, and exact Ackermann / hyperoperation arithmetic.

## Layout

```text
examples/             Go examples
examples/output/      expected output for each example
test                  run examples and compare with expected output
```

## Run

Run one example directly:

```sh
go run examples/sudoku.go
```

For the complex-number translation:

```sh
go run examples/complex_numbers.go
```

For the Ackermann / hyperoperation translation:

```sh
go run examples/ackermann.go
```

Run the full regression test:

```sh
./test
```

The test prints `OK` or `FAIL` for each example, per-example timing, and total time. It compares against `examples/output/*.txt` while normalizing Go runtime/platform audit lines, since those vary by machine.

Regenerate expected outputs after intentional changes:

```sh
./test --update
```
