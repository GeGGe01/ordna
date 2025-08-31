package options

import (
    "errors"
    "flag"
    "fmt"
    "time"
)

// Config holds user-selected options for a run.
type Config struct {
    Move       bool
    Copy       bool
    DryRun     bool
    GroupExt   bool
    From       time.Time
    To         time.Time
    Sources    []string
    Dest       string
}

func (c Config) Validate() error {
    if (c.Move && c.Copy) || (!c.Move && !c.Copy) {
        return errors.New("must specify exactly one of -m (move) or -c (copy)")
    }
    if len(c.Sources) == 0 || c.Dest == "" {
        return errors.New("need at least one source and a destination")
    }
    return nil
}

// Parse parses CLI args (like os.Args[1:]) into a Config.
// It mirrors the existing short flags for compatibility.
func Parse(args []string) (Config, error) {
    fs := flag.NewFlagSet("ordna", flag.ContinueOnError)
    fs.Usage = func() {
        fmt.Fprintf(fs.Output(), "Usage: ordna SOURCE... DEST {-m|-c} [OPTIONS]\n")
        fs.PrintDefaults()
    }

    var (
        moveMode = fs.Bool("m", false, "move files")
        copyMode = fs.Bool("c", false, "copy files")
        fromStr  = fs.String("from", "", "start date (YYYY-mm-dd)")
        toStr    = fs.String("to", "", "end date (YYYY-mm-dd)")
        sortExt  = fs.Bool("ext", false, "group by extension")
        dryRun   = fs.Bool("dry-run", false, "show actions without changes")
    )

    if err := fs.Parse(args); err != nil {
        return Config{}, err
    }

    var from, to time.Time
    var err error
    if *fromStr != "" {
        from, err = time.Parse("2006-01-02", *fromStr)
        if err != nil {
            return Config{}, fmt.Errorf("invalid --from: %w", err)
        }
    }
    if *toStr != "" {
        to, err = time.Parse("2006-01-02", *toStr)
        if err != nil {
            return Config{}, fmt.Errorf("invalid --to: %w", err)
        }
    }

    argsLeft := fs.Args()
    if len(argsLeft) < 2 {
        return Config{}, errors.New("not enough arguments: need SOURCE... DEST")
    }
    dest := argsLeft[len(argsLeft)-1]
    sources := argsLeft[:len(argsLeft)-1]

    cfg := Config{
        Move:     *moveMode,
        Copy:     *copyMode,
        DryRun:   *dryRun,
        GroupExt: *sortExt,
        From:     from,
        To:       to,
        Sources:  sources,
        Dest:     dest,
    }
    return cfg, cfg.Validate()
}

