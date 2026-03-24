# zwasm Development Memo

Session handover document. Read at session start.

## Current State

- Stages 0-46 + Phase 1, 3, 5, 8, 10, 11, 13, 15, 19 all complete.
- Spec: 62,263/62,263 Mac+Ubuntu+Windows (100.0%, 0 skip).
- E2E: 792/792 (Mac+Ubuntu). Real-world: 50/50.
- JIT: Register IR + ARM64/x86_64 + SIMD (NEON 253/256, SSE 244/256).
- Binary: 1.29MB stripped. Memory: ~3.5MB RSS.
- Platforms: macOS ARM64, Linux x86_64/ARM64, Windows x86_64.
- **main = stable**. ClojureWasm updated to v1.5.0.

## Current Task

**W38 Lazy AOT — merge gate** on branch `perf/w38-lazy-aot`.

### Changes

- HOT_THRESHOLD 10 → 3: earlier JIT compilation
- `back_edge_bailed` flag: reentry guard/br_table bail separated from `jit_failed`
- ARM64 extract_lane: imm5 shift fix + upper-half lane memory load
- JIT `jitMemGrow`: u32 → u64 for memory64 `-1` return
- JIT trampoline: cross-module call passes callee's instance (not caller's)
- Tests: force_interpreter for fuel/deadline/memory resource tests

### Gate Status

- `zig build test`: all pass
- Spec: **62,263/62,263 (100%)**
- E2E: 792/792
- Real-world: 41/50 (6 JIT bugs → W41, 3 wasmtime diffs → W42)
- Benchmarks: no regression

### Open Work Items

| Item     | Description                                       | Status         |
|----------|---------------------------------------------------|----------------|
| W41      | JIT real-world correctness (6 programs, T=3 露出)  | New            |
| W42      | wasmtime 互換性差異 (3 Go programs, JIT 無関係)     | New            |
| Phase 18 | Lazy Compilation + CLI Extensions                 | Future         |

## Completed Phases (summary)

| Phase | Name                                  | Date       |
|-------|---------------------------------------|------------|
| 1     | Guard Pages + Module Cache            | 2026-03    |
| 3     | CI Automation + Documentation         | 2026-03    |
| 5     | C API + Conditional Compilation       | 2026-03    |
| 8     | Real-World Coverage + WAT Parity      | 2026-03    |
| 10    | Quality / Stabilization               | 2026-03    |
| 11    | Allocator Injection + Embedding       | 2026-03    |
| 13    | SIMD JIT (NEON + SSE)                 | 2026-03-23 |
| 15    | Windows Port                          | 2026-03    |
| 19    | JIT Reliability                       | 2026-03    |

## References

- `@./.dev/roadmap.md` — Phase roadmap
- `@./.dev/checklist.md` — W38 details (investigation steps, benchmarks, sources)
- `@./.dev/references/w38-osr-research.md` — **W38 next steps: 4 approaches compared**
- `@./.dev/decisions.md` — architectural decisions (D131: epoch JIT timeout)
- `@./.dev/jit-debugging.md` — JIT debug techniques
- `bench/simd_comparison.yaml` — SIMD performance data (3 layers: baseline → post-opt → JIT)
- `bench/simd/src/` — C source for compiler-generated SIMD benchmarks
- `bench/run_simd_bench.sh` — SIMD microbenchmark runner
- External: wasmtime (`~/Documents/OSS/wasmtime/`), zware (`~/Documents/OSS/zware/`)
