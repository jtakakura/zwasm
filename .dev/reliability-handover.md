# Reliability — Session Handover

> Plan: `@./.dev/reliability-plan.md`. Rules: `@./.claude/rules/reliability-work.md`.

## Branch
`strictly-check/reliability-003` (from main at d55a72b)

## Progress

### ✅ Completed
- A-F: Environment, compilation, compat, E2E expansion, benchmarks, analysis, W34 fix
- G.1-G.3: Ubuntu spec 62,158/62,158 (100%). Real-world: all pass without JIT, 6/9 fail with JIT → Phase J
- I.0-I.7: E2E 792/792 (100%). FP precision fix (JIT getOrLoad dirty FP cache),
  funcref validation, import type checking, memory64 bulk ops,
  GC array alloc guard, externref encoding, thread/wait sequential simulation.
- J.1-J.3: x86_64 JIT bug fixes complete. All C/C++ real-world pass with JIT.
  Fixes: division safety (SIGFPE), ABI register clobbering (global.set, mem ops),
  SCRATCH2/vreg10 alias (R11 reserved), call liveness (rd as USE for return/store).

### Active / TODO

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

1. **Phase K: Performance optimization** — reduce benchmark gaps to ≤1.5x wasmtime.
2. After K: Phase H (documentation audit).

## x86_64 JIT status (Phase J complete)
All C/C++ real-world programs pass with JIT on Ubuntu x86_64:
- cpp_string_ops: FIXED (division safety + register clobbering)
- c_string_processing: FIXED (SCRATCH2/vreg10 alias, global.set clobbering)
- cpp_vector_sort: FIXED (SCRATCH2/vreg10 alias + call liveness analysis)
- c_math_compute, c_matrix_multiply: PASS
- c_hello_wasi: EXIT=71 (WASI issue, not JIT — same with --profile)
- go_*: EXIT=0 but no output (WASI compat issue, not JIT — same with --profile)

## Benchmark gaps (Phase K input)
rw_c_math: 5.9x, rw_c_string: 4.1x, gc_tree: 3.2x, st_matrix: 2.8x, rw_c_matrix: 2.7x.
Root cause: libm/libc inner functions stay on interpreter.
