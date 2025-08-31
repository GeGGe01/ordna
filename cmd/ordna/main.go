package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	moveMode = flag.Bool("m", false, "move files")
	copyMode = flag.Bool("c", false, "copy files")
	fromStr  = flag.String("from", "", "start date (YYYY-mm-dd)")
	toStr    = flag.String("to", "", "end date (YYYY-mm-dd)")
	sortExt  = flag.Bool("ext", false, "group by extension")
	dryRun   = flag.Bool("dry-run", false, "show actions without changes")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: ordna SOURCE... DEST {-m|-c} [OPTIONS]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if (*moveMode && *copyMode) || (!*moveMode && !*copyMode) {
		fmt.Fprintln(os.Stderr, "must specify exactly one of -m or -c")
		os.Exit(2)
	}

	args := flag.Args()
	if len(args) < 2 {
		usage()
	}
	dest := args[len(args)-1]
	sources := args[:len(args)-1]

	var from, to time.Time
	var err error
	if *fromStr != "" {
		from, err = time.Parse("2006-01-02", *fromStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid --from: %v\n", err)
			os.Exit(2)
		}
	}
	if *toStr != "" {
		to, err = time.Parse("2006-01-02", *toStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid --to: %v\n", err)
			os.Exit(2)
		}
	}

	for _, src := range sources {
		err := filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
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
			if !to.IsZero() && t.After(to.Add(24*time.Hour)) {
				return nil
			}

			dstDir := buildDest(t, filepath.Ext(path))
			dstDir = filepath.Join(dest, dstDir)
			if err := os.MkdirAll(dstDir, 0o755); err != nil {
				return err
			}
			dstPath := filepath.Join(dstDir, filepath.Base(path))

			if *dryRun {
				fmt.Printf("%s -> %s\n", path, dstPath)
				return nil
			}

			if *moveMode {
				if err := os.Rename(path, dstPath); err != nil {
					if errors.Is(err, os.ErrInvalid) {
						// Fallback to copy+delete
						if err := copyFile(path, dstPath); err != nil {
							return err
						}
						return os.Remove(path)
					}
					return err
				}
			} else {
				if err := copyFile(path, dstPath); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error processing %s: %v\n", src, err)
		}
	}
}

func buildDest(t time.Time, ext string) string {
	year, month, _ := t.Date()
	monthName := month.String()
	dir := filepath.Join(fmt.Sprintf("%04d", year), fmt.Sprintf("%02d_%s", month, monthName))
	if *sortExt {
		ext = strings.TrimPrefix(strings.ToLower(ext), ".")
		if ext == "" {
			ext = "unknown"
		}
		dir = filepath.Join(dir, ext)
	}
	return dir
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
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
		var atime time.Time
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			// Use access time from syscall.Stat_t
			sec := stat.Atim.Sec
			nsec := stat.Atim.Nsec
			atime = time.Unix(sec, nsec)
		} else {
			// Fallback to current time if unable to get access time
			atime = time.Now()
		}
		os.Chtimes(dst, atime, info.ModTime())
		os.Chmod(dst, info.Mode())
	}
	return nil
}
