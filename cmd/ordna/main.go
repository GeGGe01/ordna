package main

import (
    "fmt"
    "os"

    "ordna/internal/apply"
    "ordna/internal/options"
    "ordna/internal/plan"
    "ordna/internal/scan"
)

func main() {
    cfg, err := options.Parse(os.Args[1:])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(2)
    }

    files, err := scan.Discover(cfg)
    if err != nil {
        fmt.Fprintf(os.Stderr, "discover error: %v\n", err)
        os.Exit(1)
    }

    p := plan.Build(cfg, files)
    if err := apply.Execute(cfg, p); err != nil {
        fmt.Fprintf(os.Stderr, "apply error: %v\n", err)
        os.Exit(1)
    }
}
