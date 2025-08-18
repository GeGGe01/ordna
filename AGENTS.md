# Ordna (Go) — Functional Specification

This document defines a concise, implementable specification for the Ordna file organizer rewritten in Go with a native TUI. It converts the prior narrative design into actionable requirements, data flows, and milestones.

## Objectives

- Single portable binary with no runtime external tools.
- Native TUI (Bubbletea) with pre-run analysis and live progress.
- Fast pipeline: batch indexing, minimal redundant I/O, safe concurrency.
- Exact duplicate detection (content hash) with configurable policy.
- Robust cross-platform support (Linux, macOS, Windows).

## Non-Goals

- Fuzzy/visual duplicate detection (perceptual hashing).
- Cloud sync or remote destinations.
- Lossy media transforms (e.g., resizing). 

## User Stories

- As a user, I select sources and a destination, preview what will happen, then run with live progress.
- As a user, I filter by date range and optionally group output by extension.
- As a user, I choose move or copy and understand if move will degrade to copy+delete.
- As a user, I enable duplicate detection and avoid storing duplicate content.
- As a user, I see an ETA during analysis and while processing.

## TUI Flow (Bubbletea)

1. Welcome: create or pick a recent configuration.
2. Source selection: multi-select files/directories (recursive) and set destination.
3. Options: move/copy, date range, group-by-extension, strict type detection, duplicate policy, concurrency.
4. Analyze (dry run): scan+index, compute metrics and ETA, list warnings/heavy-cost drivers.
5. Review: summary of actions; allow editing options.
6. Execute: live progress, current step/file, ETA, togglable log.
7. Summary: totals, skipped/duplicates, errors, duration; option to export report.

Keyboard: arrows/tab to navigate, space to select, enter to confirm, `a` analyze, `r` run, `l` toggle log, `q` quit/back.

## CLI Mode (Non-Interactive)

`ordna --src <path...> --dst <path> [--copy|--move] [--from YYYY-MM-DD] [--to YYYY-MM-DD] [--group-ext] [--strict-ext] [--dupes policy] [--threads N] [--analyze]`

- `--dupes`: `skip` (default), `merge`, `rename` (append short hash).
- `--analyze`: print metrics and ETA then exit (no changes).

## Core Behavior

- Discovery: recursively list files from all sources; follow symlinks only if enabled.
- Date extraction: EXIF/metadata first; fallback to filesystem times.
- Destination scheme: `YYYY/MM_Month/` with optional `/ext/` grouping.
- Move/Copy: `os.Rename` if same device, else copy+delete; preserve perms and times.
- Conflict handling: apply duplicate policy; if `rename`, append first 12 hex of SHA-256.
- Strict type: sniff magic header when enabled; otherwise trust extension.

## Pre-Run Analysis (Dry Run)

Metrics collected:
- File count, total bytes, by-type breakdown, by-extension counts.
- Share of items requiring EXIF reads, magic sniffing, and hashing.

ETA model (configurable constants with sane defaults):
- Hashing throughput (MB/s) × bytes-to-hash (only same-size candidate groups if dupes on).
- Copy throughput (MB/s) × bytes-to-copy (consider cross-device move).
- Metadata overhead: fixed-per-file micro-costs (EXIF/stat/sniff) × item counts.
- Add 5–15% overhead buffer for directory creation and bookkeeping.

Outputs:
- “N files (~X GiB). Estimated time: ~T.”
- Top cost drivers: dupes hashing, cross-device copies, strict sniffing.
- Warnings: very large set thresholds and suggested option tweaks.

## Execution Progress

- Overall progress bar + percent; current phase: Indexing → De-dup → Organize.
- Current file + action (e.g., Hashing, Copying, Moving, Skipping).
- Live ETA based on moving average throughput.
- Toggle log for warnings/errors and decisions (e.g., duplicate skipped).

## Performance Design

- Single-pass indexing to collect: size, mtime, candidate EXIF time, extension, sniff-needed flag.
- Stage pipeline: Index → (optional) Duplicate grouping → Plan → Apply.
- Worker pools with bounded concurrency for I/O-bound steps; default `threads = min(4, GOMAXPROCS)`.
- Duplicate detection: group by size → selectively hash within groups; parallel hashing; cache results for apply stage.
- Avoid re-reading bytes: if hashed once, reuse for rename suffix and dedupe.
- Throttle UI updates (e.g., at 30–60 Hz max) to avoid terminal churn.

## Libraries (Pure Go)

- TUI: Charmbracelet Bubbletea (+ Bubbles/Lipgloss).
- EXIF/metadata: `github.com/dsoprea/go-exif` (images); MP4/MOV: suitable Go parser (e.g., QuickTime/ISO-BMFF reader) or fallback to file times.
- File type sniffing: `github.com/h2non/filetype`.
- Hashing: `crypto/sha256` (streaming via `io.Copy` to hasher).
- Filesystem ops: `os`, `io`, `fs`; times via `os.Chtimes`, perms via `os.Chmod`.
- Time parsing: `time.Parse` with layouts like `2006-01-02`.

## Data Structures

```go
type FileRecord struct {
    SrcPath      string
    RelPath      string // relative to source root chosen
    Size         int64
    ModTime      time.Time
    TakenTime    time.Time // from EXIF/metadata; zero if unknown
    Ext          string
    SniffNeeded  bool
    DeviceID     uint64 // for cross-device move detection (platform-specific)
    Hash         [32]byte // set if computed
}

type AnalysisResult struct {
    Count        int
    Bytes        int64
    ByExt        map[string]int
    ByType       map[string]int
    HashBytes    int64
    CopyBytes    int64
    EstDuration  time.Duration
    Warnings     []string
}

type PlanEntry struct {
    SrcPath  string
    DstPath  string
    Action   string // move, copy, skip-duplicate, rename
}
```

## Duplicate Policy

- `skip`: if content-equal exists at destination, do not copy/move; log once.
- `merge`: replace destination with source when equal content; preserve newest timestamp.
- `rename`: if same name but different content, append 12-hex of SHA-256 to stem.

Detection rules:
- Build groups by `Size`; hash only groups with size collisions.
- For conflicts at destination, compare content hash to avoid false positives.

## Error Handling & Logging

- Non-fatal per-file errors: log and continue; summarize counts.
- Fatal errors (e.g., destination unwritable): surface and abort gracefully.
- Provide `--report <path>` to emit JSON summary of actions, errors, and duplicates.

## Cross-Platform Notes

- Use `os.Rename` and detect EXDEV errors to fall back to copy+delete.
- Avoid platform-specific external commands; rely on Go stdlib.
- Windows: long paths (`\\\\?\\`), permissions semantics; test TUI behavior in default terminals.

## Testing Strategy

- Unit: date parsing, path templating, sniffing vs extension logic, hash compare.
- Integration: temp dirs with mixed media; simulate cross-device via distinct temp roots.
- Golden tests for planning: given inputs + options, expected set of `PlanEntry`.

## Milestones

1. Core discovery + date extraction + path templating.
2. Duplicate grouping and hashing pipeline.
3. Apply stage (move/copy) with metadata preservation.
4. TUI shell (navigation) + analyze view.
5. Live progress + ETA + logging pane.
6. Cross-platform polish and error handling.
7. CLI parity and JSON report.

## Configuration Options

- Concurrency: `--threads N` (default based on CPU/IO heuristics).
- Strict extension: `--strict-ext` to enable magic sniffing for unknown/mismatched extensions.
- Group by extension: `--group-ext` to add `/ext/` layer under date.
- Preserve timestamps: always; permissions preserved on copy; note platform differences.
- Symlink handling: `--follow-symlinks` off by default.

## Destination Path Template

Default: `YYYY/MM_Month/[ext/]OriginalName[+hash].Ext`

Example: `2025/07_July/jpg/IMG_1234.jpg` or `2025/07_July/IMG_1234+ab12cd34ef56.jpg` (rename policy).

## ETA Details

- Maintain exponential moving average for observed copy and hash throughput.
- Update ETA each tick using remaining bytes per stage and current EMA.
- During analysis, derive ETA from constants and metrics; during execution, refine with observed rates.

## Acceptance Criteria

- Runs without external tools; handles >100k files with stable memory usage.
- TUI supports analyze→review→execute with responsive updates and accurate ETA (±20%).
- Duplicate policy behaves deterministically and avoids redundant content.
- Cross-device moves gracefully degrade; metadata preserved on copies.
- CLI produces identical outcomes to TUI with equivalent options.

---

This spec is the implementation contract for the Go rewrite. Prefer simplicity and safety over premature micro-optimizations; measure and iterate on throughput-related constants.

