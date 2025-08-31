package scan

import (
    "io/fs"
    "path/filepath"

    "ordna/internal/options"
    "ordna/internal/types"
)

// Discover walks the provided sources and returns a list of files
// that fall within the optional date range filters.
func Discover(cfg options.Config) ([]types.FileRecord, error) {
    var out []types.FileRecord
    for _, src := range cfg.Sources {
        root := src
        err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
            if err != nil {
                return err
            }
            if d.IsDir() {
                return nil
            }
            info, err := d.Info()
            if err != nil {
                return err
            }
            t := info.ModTime()
            if !cfg.From.IsZero() && t.Before(cfg.From) {
                return nil
            }
            if !cfg.To.IsZero() && t.After(cfg.To.Add(24*time.Hour)) {
                return nil
            }
            // Compute relpath best-effort relative to the provided source.
            rel, _ := filepath.Rel(root, path)
            fr := types.FileRecord{
                SrcPath: path,
                RelPath: rel,
                Size:    info.Size(),
                ModTime: t,
                Ext:     filepath.Ext(path),
            }
            out = append(out, fr)
            return nil
        })
        if err != nil {
            // If a single source errors (e.g., permission), surface it.
            return nil, err
        }
    }
    return out, nil
}
