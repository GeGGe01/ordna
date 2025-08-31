package types

import "time"

// FileRecord captures file metadata discovered during indexing.
type FileRecord struct {
    SrcPath     string
    RelPath     string
    Size        int64
    ModTime     time.Time
    TakenTime   time.Time
    Ext         string
    SniffNeeded bool
    DeviceID    uint64
    Hash        [32]byte
}

// AnalysisResult is a summary from the analyze (dry-run) phase.
type AnalysisResult struct {
    Count       int
    Bytes       int64
    ByExt       map[string]int
    ByType      map[string]int
    HashBytes   int64
    CopyBytes   int64
    EstDuration time.Duration
    Warnings    []string
}

// PlanEntry is a single action to perform for a file.
type PlanEntry struct {
    SrcPath string
    DstPath string
    Action  string // move, copy, skip-duplicate, rename
}

