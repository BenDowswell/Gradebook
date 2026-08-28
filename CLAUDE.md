# Gradebook — project context

A learning exercise: an in-memory student gradebook in Go, built up in stages.
This file is the cross-machine memory (PC + laptop). Keep it updated as decisions
change.

## Working style

The user is learning by doing. On this project they often run sessions in
**coaching mode** (`/coaching`): Claude guides with questions and hints and
reviews the user's own code, but does not write the solution. If you want that
mode on a fresh session, start it with `/coaching`.

## Roadmap (all later stages "way off" as of 2026-08-28)

1. **Now:** in-memory core — `GradeBook` with maps, sentinel errors, table-driven tests.
2. Wrap the core in a simple console application.
3. Replace in-memory maps with a database.
4. Expose operations as HTTP API endpoints.
5. Add a React front end.

## Deliberate deferrals (don't "fix" these without asking)

- **No mutex / concurrency safety.** Maps aren't safe for concurrent use; this is
  fine for a single-threaded CLI. Add synchronisation at the HTTP stage (step 4).
- **Nested `map[int]map[int]float64` for grades.** A flat
  `map[struct{ studentID, subjectID int }]float64` (composite key) is the intended
  refactor — expected to fall out naturally at the DB step, since it maps to a
  compound primary key. Nested map + lazy inner-map init is fine until then; a
  test pins its overwrite behaviour.

## Conventions established in review

- **Sentinel errors:** package-level `var ErrX = errors.New("lowercase, no punctuation")`,
  matched by callers with `errors.Is`. Name them for the *absence*
  (`ErrNoGradesForStudent`), not the presence.
- **Error strings:** lowercase, no trailing punctuation (`ST1005`).
- **Constructors / mutators that can fail** return `(value, error)`; on the error
  path return the zero value. Contract: caller checks `err` before trusting the value.
- **Validate inputs before mutating state** (e.g. empty-field check runs before
  `counter++`, so a rejected call doesn't consume an ID).
- **IDs** come from monotonic counters (`studentCounter` / `subjectCounter`),
  starting at 1, only ever incremented — never `len()+1`.
- **Duplicate names are allowed** — students are distinguished by their unique ID,
  mirroring a real class with two "Alice Smith"s.
- **`gofmt` on save.** Turn on `editor.formatOnSave` so formatting never drifts.
- Tests live in `*_test.go` only (the runner ignores `Test*` funcs elsewhere, and
  importing `testing` in a normal file breaks the build).

## Review status

The findings in `REVIEW.md` are all addressed as of 2026-08-28:

- Monotonic ID counters; `AddStudent` / `AddSubject` return the new ID.
- `AddGrade` rejects NaN (`math.IsNaN`) alongside out-of-range — plain range
  checks miss NaN because every ordered comparison with NaN is `false`.
- Empty `name` / `school` (and empty subject name) rejected with `ErrEmptyField`.
- Sentinel errors + `errors.Is` throughout; strings lowercased, typos fixed.
- `package gradebook`, file renamed from `main.go` to `gradebook.go`.
- `NewGradeBook` tidied (no throwaway locals).
- `gradebook_test.go`: table-driven tests covering happy paths, every sentinel
  path, grade bounds + NaN, boundary values (0 and 100), grade overwrite, and
  ID sequencing across multiple adds.

Remaining open items are only the two deliberate deferrals above.

## Commands

- `go test ./... -v` — run tests
- `gofmt -l .` — list unformatted files (silent = clean)
- `gofmt -w gradebook.go` — format in place
- `go vet ./...` — static checks (silent = clean)
