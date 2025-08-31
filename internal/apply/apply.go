package apply

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"

    "ordna/internal/options"
    "ordna/internal/types"
)

// Execute applies the plan. If DryRun is set, it only prints planned actions.
func Execute(cfg options.Config, plan []types.PlanEntry) error {
    for _, p := range plan {
        if cfg.DryRun {
            fmt.Printf("%s -> %s\n", p.SrcPath, p.DstPath)
            continue
        }
        if err := os.MkdirAll(filepath.Dir(p.DstPath), 0o755); err != nil {
            return err
        }
        switch p.Action {
        case "move":
            if err := os.Rename(p.SrcPath, p.DstPath); err != nil {
                // Fallback to copy+delete (e.g., cross-device)
                if err := copyFile(p.SrcPath, p.DstPath); err != nil {
                    return err
                }
                if err := os.Remove(p.SrcPath); err != nil {
                    return err
                }
            }
        case "copy":
            if err := copyFile(p.SrcPath, p.DstPath); err != nil {
                return err
            }
        default:
            // For now, ignore unknown actions.
        }
    }
    return nil
}

func copyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()
    info, err := in.Stat()
    if err != nil {
        return err
    }
    out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode().Perm())
    if err != nil {
        return err
    }
    if _, err := io.Copy(out, in); err != nil {
        out.Close()
        return err
    }
    if err := out.Close(); err != nil {
        return err
    }
    _ = os.Chmod(dst, info.Mode())
    _ = os.Chtimes(dst, time.Now(), info.ModTime())
    return nil
}

