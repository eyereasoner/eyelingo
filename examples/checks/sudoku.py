# Independent Python checks for the sudoku example.
from .common import run_fragment_checks

CHECKS = [
    ('every given clue is preserved in the final grid', ['The solver starts from 23 clues and fills the remaining 58 cells by combining constraint propagation with depth-first search. At each step it chooses the empty cell with the fewest']),
    ('the final grid contains only digits 1 through 9, with no blanks left', ['1 6 2 | 8 5 7 | 4 9 3']),
    ('each row contains every digit exactly once', ['The puzzle is solved, and the completed grid is the unique valid Sudoku solution.']),
    ('each column contains every digit exactly once', ['5 3 4 | 1 2 9 | 6 7 8']),
    ('each 3×3 box contains every digit exactly once', ['7 8 9 | 6 4 3 | 5 2 1']),
    ('replaying the recorded placements from the original puzzle remains legal at every step', ['default puzzle : classic']),
    ('the search statistics and the successful proof path are internally consistent', ['4 7 5 | 3 1 2 | 9 8 6']),
    ('a second search found no alternative solution, so the solution is unique', ['9 1 3 | 5 8 6 | 7 4 2']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
