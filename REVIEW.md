# Code Review — main.go

Reviewed [main.go](main.go) — a single-file in-memory gradebook. Overall it's
clean and readable, the lazy init of the nested grades map is done correctly, and
factoring out `validateIDs` is good. No code was changed.

## Bugs / correctness

- **main.go:48, main.go:61 — ID generation via `len()+1` is fragile.** The
  comments acknowledge deletion will break it, but it's also just latent: any
  future delete produces duplicate IDs and silent overwrites. Safest fix is a
  monotonic counter field on `GradeBook` (`nextStudentID int`) that only ever
  increments.
- **main.go:46, main.go:59 — `AddStudent`/`AddSubject` don't return the new ID.**
  A caller has no reliable way to learn the ID that was assigned, so they can't
  then call `AddGrade` without independently reproducing the `len()+1` logic.
  These should return `int` (or `(int, error)`).
- **main.go:77 — `NaN` grade passes validation.** `NaN > 100` and `NaN < 0` are
  both false, so `AddGrade(s, sub, math.NaN())` succeeds. Add a `math.IsNaN`
  check if that input is reachable.
- **No duplicate / empty-value guarding.** Empty `name`/`school` are accepted
  (main.go:46, main.go:59), and adding "Alice" twice creates two records. May be
  intentional, but worth a decision.

## Go conventions

- **Error strings (main.go:38, 41, 78, 99, 104) violate `ST1005`:** they should
  not be capitalized and should not end with punctuation. main.go:78 also has a
  trailing space. Typos throughout: "doesnt", "availble".
- **Sentinel errors.** Every error is a fresh `errors.New`, so callers can't
  branch on them. Consider `var ErrStudentNotFound = errors.New("student not
  found")` etc. and let callers use `errors.Is`. `fmt.Errorf("student %d not
  found", studentID)` would also make them more debuggable.
- **main.go:1 — `package student` in a file named `main.go`.** Mismatched on both
  counts: the domain is "gradebook", not "student", and `main.go` conventionally
  implies `package main` with a `main()`. Rename the file (e.g. `gradebook.go`).
- **main.go:46 — parameter named `student` is actually the name.** Inconsistent
  with `AddSubject(name string)`; call it `name`.
- **main.go:22-34 — `NewGradeBook` is verbose.** The three local vars and the
  `g :=` / `return &g` dance can collapse to a single
  `return &GradeBook{Students: make(...), ...}`.

## Design considerations (for later)

- **Concurrency:** maps aren't safe for concurrent access and there's no mutex.
  Fine for a single-threaded CLI; a problem the moment this sits behind HTTP
  handlers.
- **Grades storage:** `map[int]map[int]float64` with lazy inner-map init works,
  but a flat `map[struct{ studentID, subjectID int }]float64` removes the
  nested-init branch in main.go:82-86 entirely.
- **No tests.** A `gradebook_test.go` covering the not-found paths, the grade
  bounds, and the overwrite behavior would lock in current intent.
