package plan

import (
    "fmt"
    "path/filepath"
    "strings"
    "time"

    "ordna/internal/options"
    "ordna/internal/types"
)

func destDirFor(t time.Time, ext string, groupExt bool) string {
    year, month, _ := t.Date()
    monthName := month.String()
    dir := filepath.Join(fmt.Sprintf("%04d", year), fmt.Sprintf("%02d_%s", month, monthName))
    if groupExt {
        e := strings.TrimPrefix(strings.ToLower(ext), ".")
        if e == "" {
            e = "unknown"
        }
        dir = filepath.Join(dir, e)
    }
    return dir
}

// Build creates a simple plan based on discovered files and config.
func Build(cfg options.Config, files []types.FileRecord) []types.PlanEntry {
    action := "copy"
    if cfg.Move {
        action = "move"
    }
    var entries []types.PlanEntry
    for _, f := range files {
        dstDir := destDirFor(f.ModTime, f.Ext, cfg.GroupExt)
        dstPath := filepath.Join(cfg.Dest, dstDir, filepath.Base(f.SrcPath))
        entries = append(entries, types.PlanEntry{
            SrcPath: f.SrcPath,
            DstPath: dstPath,
            Action:  action,
        })
    }
    return entries
}

