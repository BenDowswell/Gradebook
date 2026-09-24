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
  fine for a single-threaded CLI. Add synchronisation (`sync.RWMutex`) at the HTTP
  stage (step 4). Raised again by an outside reviewer on 2026-09-24; still deferred.
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
- **Private storage, public data types.** `GradeBook`'s fields (`students`,
  `subjects`, `grades`, both counters) are lowercase, so another package can't
  write `gb.grades[99][88] = -500` and skip validation. Every write goes through
  a method. The value structs handed back (`Student`, `Subject`, `StudentGrade`)
  keep **capitalised** fields: callers must be able to read `s.Name` to display
  it, and since they only ever get copies, writing to one changes nothing inside
  the gradebook.
- **Sanitise → validate → store.** `AddStudent` / `AddSubject` run
  `strings.TrimSpace` on every string input *first*, then the empty check, then
  store the trimmed value. So `"   "` is `ErrEmptyField` (no separate sentinel:
  to the caller, blanks are just empty) and `"  Alice  "` is stored as `"Alice"`.
  This lives in the core, not the console, so the future HTTP API gets it too.
  The console may also check input early, but only to re-prompt the user sooner.
- **Validate inputs before mutating state** (e.g. empty-field check runs before
  `counter++`, so a rejected call doesn't consume an ID).
- **An ID is not a slice index.** IDs start at 1 and have gaps after a delete, so
  never use `ListStudents()[id-1]`. To find a student by ID, loop and compare, or
  add a `GetStudent(id)` method.
- **IDs** come from monotonic counters (`studentCounter` / `subjectCounter`),
  starting at 1, only ever incremented — never `len()+1`.
- **Duplicate names are allowed** — students are distinguished by their unique ID,
  mirroring a real class with two "Alice Smith"s.
- **`gofmt` on save.** Turn on `editor.formatOnSave` so formatting never drifts.
- Tests live in `*_test.go` only (the runner ignores `Test*` funcs elsewhere, and
  importing `testing` in a normal file breaks the build).
- **Two kinds of test file.** `gradebook_test.go` is `package gradebook` (internal):
  it can read private fields like `g.students`. `external_test.go` is
  `package gradebook_test` (external): it imports `"gradebook"`, needs the prefix
  (`gradebook.NewGradeBook()`, `gradebook.ErrEmptyField`), and sees only
  capitalised names, the same as the console will. Use it to check the public API
  is enough on its own.
- **Table rows state what goes in *and* what comes out.** When the output differs
  from the input (e.g. trimming), add an `expected…` column (`expectedName`) and
  compare against that, never against the input. Give every row a distinct label,
  because a duplicate shows up as `#01` / `#02` in failures.
- **Prove a new test catches the bug.** Temporarily switch off the fix, check the
  new rows go red for the right reason, then switch it back on.
- **List / read operations don't error on "empty".** "Zero students" is a valid,
  complete answer, so `ListStudents` / `ListSubjects` return an empty slice, not a
  sentinel. They take no input, so there's nothing to fail on — no `error` return
  at all. A read that takes an ID (`ListGradesForStudent`) *can* fail on a bad ID
  and returns `error` for that case only; "student exists but has no grades" is
  still an empty slice + `nil`.
- **List order is explicit.** Map iteration is randomised, so every `List*` builds
  a slice and `slices.SortFunc(... cmp.Compare(a.ID, b.ID))` before returning.
  Tests add entries out of ID order and assert the returned order.
- **Returned slices are snapshots.** `List*` append copies of the value structs
  (`Student` / `Subject` / `StudentGrade` — all plain value types), so callers
  can't mutate `GradeBook` through them. Never return an internal map directly.

## Core API (as of 2026-09-24)

| Create | Read | Delete |
|---|---|---|
| `AddStudent(name, school) (int, error)` | `ListStudents() []Student` | `DeleteStudent(id) error` |
| `AddSubject(name) (int, error)` | `ListSubjects() []Subject` | `DeleteSubject(id) error` |
| `AddGrade(sID, subID, grade) error` | `GetGrade(sID, subID) (float64, error)` | — |
| | `ListGradesForStudent(sID) ([]StudentGrade, error)` | |

- **`StudentGrade{ SubjectID int; SubjectName string; Grade float64 }`** is a view
  struct for report-card rows — it carries the subject *name* so the console never
  has to look it up again.
- **No `EditGrade`.** `AddGrade` already overwrites an existing entry (a test pins
  this), so "edit" is just `AddGrade` from the console. Caveat: a mistyped subject
  silently creates a new grade rather than erroring.
- **Invariant: every grade references an existing student and an existing subject.**
  `DeleteStudent` drops the student's whole `grades[id]` entry (one `delete`,
  because grades are keyed by studentID at the top level). `DeleteSubject` has to
  loop every student's inner map and `delete` the subjectID from each — the cost
  of the nested map. It also removes any inner map left empty by that, so a
  student whose only grade was for the deleted subject then reads as
  `ErrNoGradesForStudent`. The defensive `continue` in `ListGradesForStudent`
  (skip a subjectID with no `subjects` entry) is insurance for a state this
  invariant says can't occur — not normal-path logic.
- **Delete on an unknown ID returns the sentinel** (`ErrStudentID` /
  `ErrSubjectID`), not a silent no-op. Counters are never touched by delete.
- Deleting from a map while ranging *that same map* is allowed in Go (unreached
  entries just aren't yielded); adding during iteration is the unpredictable case.
  `DeleteSubject` relies on this.

## Next session

- Optional tidy-up: on rows that expect an error, `expectedName` is never
  checked, so use `""` there instead of leftover input values.
- Paper exercise: sketch the console menu, map each action to a core method, find
  the gaps. Two known candidates:
  - **Single-item getters** (`GetStudent(id)` / `GetSubject(id)`) so the console
    can echo "Student: Alice, Springfield High" after an ID is typed. Not built yet.
  - **Rename / edit** a student's name or school (and subject name) — undecided
    whether it's a feature or a delete-and-re-add for now.
- If every menu line maps to a method, start step 2 (console app in its own
  package).

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

Added 2026-09-06 (`List*` / `Delete*`): tests for list ordering (added out of ID
order, asserted sorted), empty-gradebook lists, `ListGradesForStudent` (bad ID /
no grades / sorted rows with names), `DeleteSubject` cascade across *two*
students with unrelated grades kept, empty-inner-map cleanup, surgical
`DeleteStudent`, and both delete sentinels. All green; `gofmt -l` and `go vet`
clean.

Added 2026-09-24 (outside review by a friend):
- `GradeBook` fields unexported; data-struct fields stay exported (see
  "Private storage, public data types").
- `strings.TrimSpace` on names and school, with test rows for blank and padded
  values in both `TestAddStudent` and `TestAddSubject` (`expectedName` /
  `expectedSchoolName` columns). Each trim line was switched off in turn to
  confirm its rows go red.
- New `external_test.go` (`package gradebook_test`) checks the public API from
  outside: add a student, read `.Name` back through `ListStudents()`.
- Mutex point acknowledged; still a deliberate deferral.

Remaining open items are the two deliberate deferrals above.

## Commands

- `go test ./... -v` — run tests
- `gofmt -l .` — list unformatted files (silent = clean)
- `gofmt -w gradebook.go` — format in place
- `go vet ./...` — static checks (silent = clean)
