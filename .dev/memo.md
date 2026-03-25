# zwasm Development Memo

Session handover document. Read at session start.

## Current State

- Stages 0-46 + Phase 1, 3, 5, 8, 10, 11, 13, 15, 19, **20** complete.
- Spec: 62,263/62,263 Mac+Ubuntu (100.0%, 0 skip).
- E2E: 792/792 (Mac+Ubuntu).
- Real-world: Mac 50/50, Ubuntu 50/50. go_math_big fixed 2026-03-25.
- JIT: Register IR + ARM64/x86_64 + SIMD (NEON 253/256, SSE 244/256).
- HOT_THRESHOLD=3 (lowered from 10 in W38).
- Binary: 1.29MB stripped. Memory: ~3.5MB RSS.
- Platforms: macOS ARM64, Linux x86_64/ARM64, Windows x86_64.
- **main = stable**. ClojureWasm updated to v1.5.0.

## Current Task

**Phase 20: JIT Correctness Sweep — remaining W41 bugs**

All Phase 20 fixes merged to main (2026-03-25). Merge Gate passed (Mac + Ubuntu).
**Remainder aliasing fix** added 2026-03-25 (this branch). Mac 50/50.

### Completed fixes (all Phase 20)

| Fix                                              | Impact                    |
|--------------------------------------------------|---------------------------|
| void-call reloadVreg (ARM64 + x86)               | +5 Mac programs           |
| ARM64 emitMemFill/emitMemCopy/emitMemGrow ABI    | ARM64 memory ops          |
| written_vregs pre-scan (ARM64 + x86)             | +2 Mac (tinygo_hello/json) |
| void self-call result clobber (ARM64 + x86)      | Preventive fix            |
| ARM64 fuel check x0 clobber                      | tinygo_sort FIXED         |
| Stale scratch cache in signed div                | rust_enum_match FIXED     |
| **Remainder rd==rs1 aliasing** (this branch)     | **go_math_big FIXED**     |

**Root cause — remainder register aliasing in emitRem32/emitRem64**: When
rd == rs1 (same physical register), UDIV/SDIV overwrites the original
dividend before MSUB can use it as the Xa operand. MSUB then computes
`quo - quo * divisor` instead of `dividend - quo * divisor`.
This caused Go's `math/bits.Div64` to return remainder=0 in the hi==0
fast path (where i64.rem_u has rd==rs1), making big integer decimal
conversion produce truncated output (e.g. 2^100 showed 19 digits
instead of 31). Fix: save rs1 to SCRATCH before division when d==rs1.

### Open Work Items

| Item       | Description                                       | Status         |
|------------|---------------------------------------------------|----------------|
| **W43**    | **SIMD v128 base address cache (D132 Phase A)**   | **Next**       |
| W44        | SIMD register class (D132 Phase B)                | Future         |
| Phase 18   | Lazy Compilation + CLI Extensions                 | Future         |
| Zig 0.16   | API breaking changes                              | When released  |

## Completed Phases (summary)

| Phase    | Name                                  | Date       |
|----------|---------------------------------------|------------|
| 1        | Guard Pages + Module Cache            | 2026-03    |
| 3        | CI Automation + Documentation         | 2026-03    |
| 5        | C API + Conditional Compilation       | 2026-03    |
| 8        | Real-World Coverage + WAT Parity      | 2026-03    |
| 10       | Quality / Stabilization               | 2026-03    |
| 11       | Allocator Injection + Embedding       | 2026-03    |
| 13       | SIMD JIT (NEON + SSE)                 | 2026-03-23 |
| 15       | Windows Port                          | 2026-03    |
| 19       | JIT Reliability                       | 2026-03    |
| 20       | JIT Correctness Sweep                 | 2026-03-25 |

## Next Session Reference Chain

1. **Orient**: `git log --oneline -5 && git status && git branch`
2. **This memo**: current task, open work items
3. **D132**: `@./.dev/decisions.md` → search `## D132` — SIMD two-phase plan
   - Phase A (W43): v128 base address cache — implementation plan, register choices
   - Phase B (W44): SIMD register class — design challenges, deferred
4. **JIT SIMD code**: `src/jit.zig` → search `emitSimdV128Addr`, `emitLoadV128`,
   `emitStoreV128`, `emitPrologue`, `has_simd`, `simd_v128_offset`
5. **x86 SIMD code**: `src/x86.zig` → same patterns
6. **SIMD benchmarks**: `bench/run_simd_bench.sh`, `bench/simd_comparison.yaml`
7. **Roadmap**: `@./.dev/roadmap.md` — SIMD register class listed as Medium priority
8. **Ubuntu testing**: `@./.dev/references/ubuntu-testing-guide.md` — OrbStack VM
9. **Merge gate checklist**: CLAUDE.md → "Merge Gate Checklist" section

### Key next task
- **W43**: SIMD v128 base address cache (D132 Phase A). 2-3 days.
  ARM64: cache `vm_ptr + simd_v128_offset` in x17 when `has_simd`.
  Cuts v128 addr computation from 3-4 insns to 1-2.
  Pattern identical to MEM_BASE/MEM_SIZE caching.

## References

- `@./.dev/roadmap.md` — Phase roadmap
- `@./.dev/checklist.md` — resolved work items
- `@./.dev/decisions.md` — D130 (SIMD arch), D132 (SIMD perf plan)
- `@./.dev/jit-debugging.md` — JIT debug techniques
- `@./.dev/references/w38-osr-research.md` — OSR research
- `bench/simd_comparison.yaml` — SIMD performance data
- `bench/history.yaml` — benchmark history (latest: phase20-rem-fix)
- External: wasmtime (`~/Documents/OSS/wasmtime/`), zware (`~/Documents/OSS/zware/`)
