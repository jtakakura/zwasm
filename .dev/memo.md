# zwasm Development Memo

Session handover document. Read at session start.

## Current State

- Stages 0-46 complete. v1.2.0 released. ~50K LOC, 521 unit tests.
- Spec: 62,263/62,263 Mac+Ubuntu (100.0%, 0 skip). E2E: 792/792 (100.0%, 0 leak).
- Wasm 3.0: all 9 proposals. WASI: 46/46 (100%). WAT parser complete.
- JIT: Register IR + ARM64/x86_64. Size: 1.19MB / 1.52MB RSS.
- **main = stable**: ClojureWasm depends on main (v1.2.0 tag).

## Current Task

Phase 1.2: Module Cache / AOT Serialize (see `roadmap.md` Phase 1).

Phase 1.1 (Guard Pages) already complete — implemented in guard.zig, memory.zig,
store.zig, jit.zig, x86.zig, cli.zig. JIT bounds check elimination active.

### Design: Module Cache

Save predecoded IR to disk for fast startup on repeated execution.
Decision: D124 (to be written).

**What to cache**: Predecoded instruction stream (from `predecode.zig`).
RegIR is per-function JIT output — may be too large/complex to cache initially.
Start with predecoded IR only.

**Cache location**: `~/.cache/zwasm/<hash>.bin`
**Cache key**: SHA-256 of wasm binary + zwasm version string.
**CLI**: `zwasm run --cache file.wasm` (auto-cache), `zwasm compile file.wasm` (AOT).
**Invalidation**: version mismatch or hash mismatch → recompile.

**Files to create/modify**:
- New: `src/cache.zig` — serialize/deserialize predecoded IR
- Modify: `src/module.zig` — load from cache if available
- Modify: `src/cli.zig` — `--cache` flag + `compile` subcommand

Expected: 10-100x startup improvement for large modules.

Previous: v1.2.0 released (tagged 5d54ae9, CW updated).

## Known Bugs

None.

## References

- `@./.dev/roadmap.md` (future phases), `@./.dev/roadmap-archive.md` (completed stages)
- `@./private/future/03_zwasm_clojurewasm_roadmap_ja.md` (integrated roadmap)
- `@./.dev/decisions.md`, `@./.dev/checklist.md`, `@./.dev/spec-support.md`
- `@./.dev/jit-debugging.md`, `@./.dev/bench-strategy.md`
- External: wasmtime (`~/Documents/OSS/wasmtime/`), zware (`~/Documents/OSS/zware/`)
