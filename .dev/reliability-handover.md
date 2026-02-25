# Reliability — Session Handover

> Plan: `@./.dev/reliability-plan.md`. Rules: `@./.claude/rules/reliability-work.md`.

## Branch
`strictly-check/reliability-003` (from main at d55a72b)

## Progress

### ✅ Completed
- A-F: Environment, compilation, compat, E2E expansion, benchmarks, analysis, W34 fix
- G.1-G.3: Ubuntu spec 62,158/62,158 (100%). Real-world: all pass without JIT, 6/9 fail with JIT → Phase J

### Active / TODO

**Phase I: E2E 100% + FP correctness**
- [ ] I.0: FP precision root cause (c_math_compute diff — bug, not acceptable)
- [ ] I.1: Typed funcref validation (30 assert_invalid)
- [ ] I.2: Import type checking (7 assert_unlinkable)
- [ ] I.3: Memory64 bounds edge cases (9 failures)
- [ ] I.4: GC ref.test type combinations (2 failures)
- [ ] I.5: GC array-alloc-too-large (2 failures)
- [ ] I.6: Memory64 linking validation (3 failures)
- [ ] I.7: Threads SB_atomic ordering (1 failure)

**Phase J: x86_64 JIT bug fixes**
- [ ] J.1: Investigate x86_64 JIT codegen crash patterns
- [ ] J.2: Fix x86_64 JIT bugs
- [ ] J.3: Verify all real-world pass on Ubuntu with JIT

**Phase K: Performance optimization (target: all ≤1.5x wasmtime)**
- [ ] K.1: JIT call threshold tuning
- [ ] K.2: Library function JIT coverage (try 2-3 approaches)
- [ ] K.3: Register allocation for f64-heavy code
- [ ] K.4: GC allocation optimization
- [ ] K.5: Benchmark re-recording on BOTH platforms

**Phase H: Documentation (LAST — requires Phase H Gate pass, see plan)**
- [ ] H.0: Phase H Gate — all 9 conditions verified (see `@./.dev/reliability-plan.md`)
- [ ] H.1: Audit README claims
- [ ] H.2: Fix discrepancies
- [ ] H.3: Update benchmark table

## Next session: start here

1. **Commit** the uncommitted plan/doc/script changes on `strictly-check/reliability-003`
   (reliability-plan.md, reliability-handover.md, memo.md, bench/run_bench.sh,
   .claude/rules/reliability-work.md — plan revision, no code changes)
2. **G.4: Ubuntu benchmarks** — not yet completed (previous SSH timed out).
   Re-run: `ssh ubuntu ... bash bench/run_bench.sh --quick` in background.
   Do NOT wait — proceed to Phase I or J while it runs.
3. **Start Phase I or J** (whichever is easiest to unblock first).
   I and J can proceed in parallel. See plan for task details.

## x86_64 JIT failures (Phase J input)
All PASS with `--profile` (JIT disabled). Failures with JIT:
- cpp_string_ops: Arithmetic exception (signal 6)
- c_string_processing, cpp_vector_sort: OOB memory access
- go_hello_wasi, go_json_marshal, go_sort_benchmark: OOB memory access

## Benchmark gaps (Phase K input)
rw_c_math: 5.9x, rw_c_string: 4.1x, gc_tree: 3.2x, st_matrix: 2.8x, rw_c_matrix: 2.7x.
Root cause: libm/libc inner functions stay on interpreter.
