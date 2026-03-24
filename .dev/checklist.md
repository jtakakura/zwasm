# Deferred Work Checklist

Open items only. Resolved items in git history.
Prefix: W## (to distinguish from CW's F## items).

## Invariants (always enforce)

- [ ] D100: No CW-specific types in public API (Value, GC, Env)
- [ ] D102: All types take allocator parameter, no global allocator

## Open Items

- [x] W38: SIMD JIT — compiler-generated code performance (2026-03-24)
  **Resolved via Lazy AOT approach** on branch `perf/w38-lazy-aot`:
  - HOT_THRESHOLD 10 → 3: earlier JIT compilation
  - `back_edge_bailed` flag: reentry guard bail no longer poisons call-path JIT
  - ARM64 extract_lane encoding fix (imm5 shift + upper-half memory load)
  - JIT memory_grow64 u32 → u64 (memory64 `-1` return fix)
  - JIT trampoline cross-module instance fix (callee's instance, not caller's)
  - Spec: 62,263/62,263 (100%). Benchmarks: no regression.

- [ ] W41: JIT real-world correctness — OOB/wrong results at HOT_THRESHOLD=3
  6 real-world programs produce wrong JIT output (correct with `--interp`):
  - `rust_compression`: OOB at test 2 (**T=10 でも再現、back-edge JIT バグ**)
  - `rust_enum_match`: float formatting garbage (T=3 で露出)
  - `rust_serde_json`: OOB (T=3 で露出)
  - `tinygo_hello`: OOB (T=3 で露出)
  - `tinygo_json`: OOB (T=3 で露出)
  - `tinygo_sort`: 出力差異 (T=3 で露出)
  Spec tests (62,263) all pass — these are code patterns spec tests don't cover.
  Likely root cause: JIT codegen for complex Rust/Go/TinyGo patterns (large
  functions, complex control flow, specific memory access patterns).

- [ ] W42: wasmtime 互換性差異 (JIT 無関係)
  3 Go programs produce different output from wasmtime (same in interp and JIT):
  - `go_crypto_sha256`, `go_math_big`, `go_regex`
  Likely Go runtime behavior differences (env, args, or WASI capability gaps).

## Resolved (summary)

W37: Contiguous v128 storage. W39: Multi-value return JIT (guard removed).
W40: Epoch-based JIT timeout (D131).
W38: SIMD JIT C-compiled perf — Lazy AOT (HOT_THRESHOLD 10→3, back_edge_bailed,
     extract_lane fix, memory_grow64 fix, cross-module instance fix).

W2-W36: See git history. All resolved through Stages 0-47 and Phases 1-19.
