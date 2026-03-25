# zwasm Development Memo

Session handover document. Read at session start.

## Current State

- Stages 0-46 + Phase 1, 3, 5, 8, 10, 11, 13, 15, 19, **20 (partial)** complete.
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

### Remaining

- W42: go_math_big was NOT env-dependent — it was this JIT bug! Now PASS.

### Open Work Items

| Item     | Description                                       | Status         |
|----------|---------------------------------------------------|----------------|
| W41      | JIT real-world: ALL FIXED (Mac 50/50)             | **Done**       |
| W42      | go_math_big: FIXED (remainder aliasing bug)       | **Done**       |
| Phase 18 | Lazy Compilation + CLI Extensions                 | Future         |

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
| 20 (wip) | JIT Correctness Sweep                 | 2026-03-25 |

## Next Session Reference Chain

1. **Orient**: `git log --oneline -5 && git status && git branch`
2. **This memo**: current task, root causes found, remaining bugs
3. **Checklist**: `@./.dev/checklist.md` — W41 updated with tinygo_sort details
4. **JIT debug techniques**: `@./.dev/jit-debugging.md` — dump, ELF wrap, objdump
5. **JIT code** (ARM64): `src/jit.zig` — emitBinop32/64, emitMemStore/Load, getOrLoad, spillCallerSavedLive
6. **JIT code** (x86): `src/x86.zig` — same patterns
7. **Ubuntu testing**: `@./.dev/references/ubuntu-testing-guide.md` — OrbStack VM
8. **Merge gate checklist**: CLAUDE.md → "Merge Gate Checklist" section

### Key next tasks
- **W41 rust_enum_match**: FP JIT bug (garbage f64 in Triangle coords). Mac only.
- **W42 go_math_big**: env-dependent, low priority.

## References

- `@./.dev/roadmap.md` — Phase roadmap
- `@./.dev/checklist.md` — W41/W42 details + next steps
- `@./.dev/references/w38-osr-research.md` — OSR research (4 approaches)
- `@./.dev/decisions.md` — architectural decisions (D131: epoch JIT timeout)
- `@./.dev/jit-debugging.md` — JIT debug techniques
- `bench/simd_comparison.yaml` — SIMD performance data
- External: wasmtime (`~/Documents/OSS/wasmtime/`), zware (`~/Documents/OSS/zware/`)
