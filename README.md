# eyelingo

Small Go translations of selected EyeReasoner/Eyeling N3 examples.

## Run

Run one example directly:

```sh
go run examples/sudoku.go
```

Run the full regression test:

```sh
./test
```

The test prints `OK` or `FAIL` for each example plus elapsed time. It compares against `examples/output/*.out` while normalizing Go runtime/platform audit lines, since those vary by machine.

Regenerate expected outputs after intentional changes:

```sh
./test --update
```
