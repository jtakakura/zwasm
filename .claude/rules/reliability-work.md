---
paths:
  - "src/**"
  - "test/**"
  - "bench/**"
  - ".dev/reliability-*"
---

# Reliability Work Rules

Active when on `strictly-check/reliability-*` branches.
Plan: `@./.dev/reliability-plan.md`. Progress: `@./.dev/reliability-handover.md`.

## Principles

1. **Zero tolerance** — every test failure is a bug. No "known limitations".
   FP precision diffs are bugs (same wasm bytecode = IEEE 754 mandates identical results).
2. **Performance: 1.0x target, 1.5x ceiling** — wasmtime is the reference.
   No ROI-based deferral. If first approach fails, try second. Exhaust options.
3. **Cross-platform** — Mac aarch64 + Ubuntu x86_64 must both pass.
   Ubuntu/WSL users matter. JIT bugs on x86_64 are blocking.
4. **Fair benchmarks** — all scripts: 5 runs / 3 warmup. No legacy defaults.
5. **Non-blocking Ubuntu workflow** — SSH in background, work on other tasks
   (code investigation, reading references) while waiting. Never block on SSH output.

## Investigation — Go Wide

- **wasmtime/cranelift**: `~/Documents/OSS/wasmtime/` — always check first.
  Key paths: `cranelift/codegen/src/isa/aarch64/` (ARM64),
  `cranelift/codegen/src/isa/x64/` (x86_64), `cranelift/codegen/src/opts/`.
- **Clone more if needed**: `~/Documents/OSS/` — wasm3, wasmer, wazero, etc.
- **Web search**: Use WebFetch/WebSearch for specs, blog posts, papers.
- **zware**: `~/Documents/OSS/zware/` (Zig idioms).
- Check references for **every** task: debugging, optimization, AND correctness.

## Experiment-First

- **Try boldly, revert cleanly.** Tests + benchmarks are the safety net.
  Implement and measure. No effect? `git checkout -- .` and try next approach.
- **Multiple approaches in sequence.** List 2-3 candidates, try lightest first.
  Don't over-analyze upfront. Iterate rapidly.
- **Regressions are the only hard stop.** `zig build test` pass + spec pass +
  no benchmark regression = safe to keep. Otherwise revert completely.
- **No partial fixes.** Every fix must be clean and spec-compliant.
  The codebase must be simpler after reliability work, not more complex.

## Branch Strategy

Incremental merge to main — users actively depend on main.

- Branches: `strictly-check/reliability-001`, `-002`, `-003`, ...
- Each branch = one **regression-free improvement unit**
- Workflow (continuous — do NOT pause):
  1. Work on `strictly-check/reliability-NNN`
  2. Passes Merge Gate → merge to main, push
  3. Create next: `git checkout -b strictly-check/reliability-NNN+1 main`
  4. Update `@./.dev/reliability-handover.md`, continue immediately
