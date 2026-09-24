package main

import (
    "flag"
    "fmt"
    "os"
    "strings"
)

var (
    excludeDirs string
    dryRun      bool
    fixMode     bool
    outputJSON  bool
)

func init() {
    flag.StringVar(&excludeDirs, "exclude", "", "Comma-separated list of directories to exclude")
    flag.BoolVar(&dryRun, "dry-run", false, "Show what would be removed without modifying files")
    flag.BoolVar(&fixMode, "fix", false, "Automatically remove unused definitions")
    flag.BoolVar(&outputJSON, "json", false, "Output report in JSON format")
}

func main() {
    flag.Parse()
    root := "."
    if flag.NArg() > 0 {
        root = flag.Arg(0)
    }

    excludes := parseExcludes(excludeDirs)

    scanner := NewScanner(root, excludes)
    unused, err := scanner.Scan()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
        os.Exit(1)
    }

    if outputJSON {
        fmt.Println(unused.ToJSON())
    } else {
        fmt.Println(unused.ToMarkdown())
    }

    if fixMode {
        if dryRun {
            fmt.Println("\nDry run: the following definitions would be removed:")
            fmt.Println(unused.ToMarkdown())
        } else {
            err := scanner.Fix(unused)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error fixing: %v\n", err)
                os.Exit(1)
            }
            fmt.Println("\nUnused definitions removed.")
        }
    }
}

func parseExcludes(s string) []string {
    if s == "" {
        return nil
    }
    parts := strings.Split(s, ",")
    for i := range parts {
        parts[i] = strings.TrimSpace(parts[i])
    }
    return parts
}
