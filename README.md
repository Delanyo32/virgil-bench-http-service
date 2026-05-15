# virgil-bench-http-service

A Go benchmark codebase for the [virgil-cli](https://github.com/Delanyo32/virgil-cli) flaw detector.

This repository contains an HTTP service with concurrency, context, and error-handling flaws. Every intentional flaw is cataloged in `DEBT_MANIFEST.toml` with its file, line, category, pattern, severity, and description. Running `virgil-cli audit` against this codebase produces findings that are scored against the manifest to compute precision and recall.

It is one of 10 program-specific benchmarks pulled in as submodules by the [Virgil Benchmark Suite](https://github.com/Delanyo32/virgil-skills).

## Manifest summary

- **Total cataloged flaws:** 2,116
- **Synced against:** virgil 0.4.3
- **Generated:** 2026-04-17
- **Severity breakdown:** warning (1,410), info (703), error (3)

## Top detection patterns

The codebase is seeded with the following patterns (top 15 by count):

| Pattern | Category | Count |
|---|---|---|
| `mutex_misuse` | mutex_misuse | 601 |
| `context_not_propagated` | context_not_propagated | 440 |
| `unvalidated_path_join` | go_path_traversal | 440 |
| `narrowing_conversion` | go_integer_overflow | 154 |
| `swallowed_error` | error_swallowing | 113 |
| `init_function_abuse` | init_abuse | 70 |
| `god_struct` | god_struct | 53 |
| `naked_interface` | naked_interface | 53 |
| `high_coupling` | coupling | 47 |
| `stringly_typed_config` | stringly_typed_config | 36 |
| `concrete_return_type` | concrete_return_type | 30 |
| `duplicate_symbol` | duplicate_symbols | 24 |
| `goroutine_leak_risk` | goroutine_leak | 11 |
| `unsafe_pointer_cast` | go_type_confusion | 9 |
| `loop_var_capture` | race_conditions | 6 |

## All categories covered (23)

- `mutex_misuse` (601)
- `context_not_propagated` (440)
- `go_path_traversal` (440)
- `go_integer_overflow` (154)
- `error_swallowing` (113)
- `init_abuse` (70)
- `god_struct` (53)
- `naked_interface` (53)
- `coupling` (47)
- `stringly_typed_config` (36)
- `concrete_return_type` (30)
- `duplicate_symbols` (24)
- `goroutine_leak` (11)
- `go_type_confusion` (9)
- `memory_leak_indicators` (6)
- `race_conditions` (6)
- `resource_exhaustion` (6)
- `deep_nesting` (5)
- `api_surface_area_go` (3)
- `cognitive_complexity` (3)
- `function_length` (3)
- `cyclomatic_complexity` (2)
- `sql_injection_go` (1)

## Repository layout

```
DEBT_MANIFEST.toml
cmd/
go.mod
internal/
pkg/
```

## Usage

### Standalone

```bash
git clone https://github.com/Delanyo32/virgil-bench-http-service.git
cd virgil-bench-http-service
virgil-cli audit .
```

Compare the resulting findings to `DEBT_MANIFEST.toml` to score precision and recall.

### As part of the full benchmark suite

```bash
git clone --recurse-submodules https://github.com/Delanyo32/virgil-skills.git
cd virgil-skills/benchmarks/tests
cargo test
```

The Rust test harness discovers this benchmark automatically, runs `virgil-cli audit` against it, parses the JSON output, compares to the manifest, and emits per-pattern recall/precision tracked over time in `benchmarks/baseline.json`.

## Manifest format

Each entry in `DEBT_MANIFEST.toml`:

```toml
[[debt]]
file = "src/example.rs"
line = 42
category = "panic_dos"
pattern = "unwrap_untrusted"
description = "unwrap/expect call detected -- may panic on untrusted input causing denial of service"
severity = "warning"
```

## Disclaimer

This repository **deliberately contains flawed code** (insecure patterns, hardcoded credentials that look real, memory issues, race conditions, etc.). It exists solely to measure flaw-detection tooling. **Do not use any code from this repository in production.**
