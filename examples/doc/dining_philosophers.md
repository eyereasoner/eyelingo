# Dining Philosophers

`dining_philosophers` is a Go translation/adaptation of Eyeling's `dining-philosophers.n3`.

## Background

The dining philosophers problem is a classic concurrency example: several processes compete for shared resources and can deadlock if each holds one resource while waiting for another. Chandy-Misra style solutions avoid deadlock by giving resources direction, ownership, or cleanliness rules. The important idea is to make every resource transfer explicit so conflicts and progress can be checked from the trace.

## What it demonstrates

**Category:** Mathematics. Chandy-Misra dining-philosophers trace with concurrency conflict checks.

## How the Go implementation works

The implementation replays a deterministic Chandy-Misra style schedule. Goroutines model philosophers making concurrent requests within a round, while phase boundaries keep fork transfers and dirty/clean state updates deterministic.

The report checks that adjacent philosophers never eat with the same fork, that dirty-fork transfer rules are respected, and that the repeated rounds remain fair.

## Files

Input JSON: [../input/dining_philosophers.json](../input/dining_philosophers.json)

Go implementation: [../dining_philosophers.go](../dining_philosophers.go)

Python check: [../checks/dining_philosophers.py](../checks/dining_philosophers.py)

Expected Markdown output: [../output/dining_philosophers.md](../output/dining_philosophers.md)
