# zwasm Reliability Improvement — Plan

> Updated: 2026-02-26
> Principles & branch strategy: `@./.claude/rules/reliability-work.md`
> Progress: `@./.dev/reliability-handover.md`

## Goal

Make zwasm **undeniably correct and fast** on Mac (aarch64) and Ubuntu (x86_64).
Zero "known limitations". All tests 100%. All benchmarks ≤1.5x wasmtime.

## Phases

| # | Phase | Status |
|---|-------|--------|
| A-F | Environment, compilation, compat, E2E, benchmarks, analysis | ✅ complete |
| G | Ubuntu cross-platform verification | ✅ spec/unit pass, JIT bugs found → J |
| **I** | **E2E 100% + FP correctness** | active |
| **J** | **x86_64 JIT bug fixes** | active |
| **K** | **Performance deep optimization** | next |
| H | Documentation accuracy audit | LAST (after I, J, K) |

Execution: I + J (parallel) → K → H.

---

## Phase I: E2E 100% + FP Correctness

Target: 778/778 (100%) E2E + 13/13 (100%) real-world compat.

**For every sub-task**: read failing test → check wasmtime source → WebSearch spec → implement → test → no regressions.

### I.0: FP precision root cause fix
c_math_compute: zwasm 21304744.877962 vs wasmtime 21304744.878669.
Same wasm bytecode → IEEE 754 mandates identical results → zwasm has FP bug.
**Approach**: Binary search subsets of computation to isolate diverging f64 operation.
Check `wasmNearest()`, JIT FP codegen (x87 vs SSE, NEON rounding).

### I.1: Typed funcref validation (30 assert_invalid)
Implement type validation in validator for call_indirect with typed function references.
**Ref**: wasmtime `cranelift/wasm/src/code_translator.rs`.

### I.2: Import type checking (7 assert_unlinkable)
Implement import type compatibility validation during instantiation.
**Ref**: wasmtime `crates/wasmtime/src/runtime/instantiate.rs`.

### I.3: Memory64 bounds edge cases (9 failures)
Zero-length ops at OOB addresses should trap. Fix bounds checking.

### I.4: GC ref.test type combinations (2 failures)
### I.5: GC array-alloc-too-large (2 failures)
### I.6: Memory64 linking validation (3 failures)
### I.7: Threads SB_atomic ordering (1 failure)
Try multiple approaches: compiler fences, atomic ordering, memory barriers.

---

## Phase J: x86_64 JIT Bug Fixes

All real-world programs PASS with JIT disabled. x86_64 JIT-specific bugs:
- cpp_string_ops: Arithmetic exception (signal 6)
- c_string_processing, cpp_vector_sort: OOB memory access
- go_hello_wasi, go_json_marshal, go_sort_benchmark: OOB memory access

**Investigation** (do while Ubuntu SSH runs in background):
1. Read `src/x86.zig` — compare with `src/jit.zig` (aarch64)
2. `--dump-jit=N` to diff aarch64 vs x86_64 codegen
3. Study cranelift x86_64: `~/Documents/OSS/wasmtime/cranelift/codegen/src/isa/x64/`
4. WebSearch for x86_64 calling convention edge cases
5. Look for: register clobbering, stack alignment, addressing mode overflow, sign extension

Fix root cause → verify all real-world programs pass on Ubuntu with JIT.

---

## Phase K: Performance Deep Optimization

Target: all benchmarks ≤1.5x wasmtime (1.0x ideal). No ROI-based deferral.

**Methodology per benchmark**: profile → read cranelift codegen → list 2-3 approaches →
try each (implement, measure, keep or revert) → record in runtime_comparison.yaml.

### K.1: JIT call threshold tuning
Try: lower threshold (50→20), selective lowering for JIT callees.

### K.2: Library function JIT coverage (biggest gap: 4-6x)
Root cause: libm/libc inner loops stay on interpreter.
**Try in order**: (1) lower call-count threshold, (2) interprocedural hotness propagation,
(3) eager JIT for large modules, (4) profile-guided second-pass JIT.

### K.3: Register allocation for f64-heavy code (st_matrix 2.8x, rw_c_matrix 2.7x)
Study cranelift regalloc. Try: spill/reload heuristics, register coalescing for FP loads.

### K.4: GC allocation optimization (gc_tree 3.2x)
Profile heap ops. Try: inline bump allocator fast path. Compare wasmtime GC impl.

### K.5: Benchmark re-recording
After each optimization: re-run ALL benchmarks on BOTH platforms. No cross-regression.

---

## Phase H Gate — Entry Criteria

**Phase H may NOT begin until ALL of the following are satisfied.**
Every item requires a root-cause fix — no workarounds, no "close enough", no
ROI-based deferral. If a condition fails, go back and fix the underlying issue.

| # | Condition | Verification |
|---|-----------|-------------|
| 1 | E2E: **778/778 (100%)** | Mac: e2e runner reports 0 failures |
| 2 | Real-world Mac: **13/13 PASS** (FP diff = 0) | `bash test/realworld/run_compat.sh` exits 0 |
| 3 | Real-world Ubuntu: **all PASS with JIT enabled** | `ssh ubuntu ... run_compat.sh` exits 0 |
| 4 | Spec Mac: **62,158/62,158 (100%)** | `python3 test/spec/run_spec.py --build --summary` |
| 5 | Spec Ubuntu: **62,158/62,158 (100%)** | SSH same as above |
| 6 | Benchmarks Mac: **all ≤1.5x wasmtime** | `bash bench/compare_runtimes.sh` (5 runs / 3 warmup) |
| 7 | Benchmarks Ubuntu: **all ≤1.5x wasmtime** | SSH same as above |
| 8 | Unit tests: **Mac + Ubuntu PASS** | `zig build test` on both |
| 9 | Benchmark regression: **none** | `bash bench/run_bench.sh` (5 runs / 3 warmup, no `--quick`) |

**If any condition is not met, do not proceed to Phase H. Fix the root cause first.**

---

## Phase H: Documentation Accuracy (LAST)

Only begins after Phase H Gate passes. Audit README claims, fix discrepancies, update benchmark table.
