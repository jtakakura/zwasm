# zwasm Development Memo

Session handover document. Read at session start.

## Current State

- Stages 0-46 + Phase 1, 3, 5 complete. **v1.3.0 released** (tagged 7570170).
- Spec: 62,263/62,263 Mac+Ubuntu (100.0%, 0 skip). E2E: 792/792 (100.0%, 0 leak).
- Wasm 3.0: all 9 proposals. WASI: 46/46 (100%). WAT parser complete.
- JIT: Register IR + ARM64/x86_64. Size: 1.20MB stripped. RSS: 4.48MB.
- Module cache: `zwasm run --cache`, `zwasm compile` (D124).
- **C API**: `libzwasm.so`/`.dylib`/`.a` — 25 exported `zwasm_*` functions (D126).
- **Conditional compilation**: `-Djit=false`, `-Dcomponent=false`, `-Dwat=false` (D127).
  Minimal build: ~940KB stripped (24% reduction).
- **main = stable**: v1.3.0 tagged. ClojureWasm updated to v1.3.0.

## Current Task

**Phase 5 complete** — branch `phase5/c-api`, ready for merge.

12 commits: D126 decision, c_api.zig (lifecycle/invoke/memory/exports/WASI/host-fns),
include/zwasm.h, build targets, C tests + Python example, feature flags,
conditional compilation guards, CI size-matrix, D127 decision + docs.

**Next**: Merge Gate (Mac + Ubuntu), then merge to main.

## Known Bugs

None.

## References

- `@./.dev/roadmap.md` (future phases), `@./.dev/roadmap-archive.md` (completed stages)
- `@./private/future/03_zwasm_clojurewasm_roadmap_ja.md` (integrated roadmap)
- `@./.dev/decisions.md`, `@./.dev/checklist.md`, `@./.dev/spec-support.md`
- `@./.dev/jit-debugging.md`, `@./.dev/bench-strategy.md`
- External: wasmtime (`~/Documents/OSS/wasmtime/`), zware (`~/Documents/OSS/zware/`)
