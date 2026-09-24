package main

import (
    "bufio"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strings"
)

// UnusedDef represents a single unused function definition.
type UnusedDef struct {
    File string
    Line int
    Type string
    Name string
}

// UnusedReport holds all unused definitions found.
type UnusedReport struct {
    Defs []UnusedDef
}

func (r *UnusedReport) ToMarkdown() string {
    var sb strings.Builder
    sb.WriteString("| File | Line | Type | Name |\n")
    sb.WriteString("|------|------|------|------|\n")
    for _, d := range r.Defs {
        sb.WriteString(fmt.Sprintf("| %s | %d | %s | %s |\n", d.File, d.Line, d.Type, d.Name))
    }
    return sb.String()
}

func (r *UnusedReport) ToJSON() string {
    var sb strings.Builder
    sb.WriteString("[\n")
    for i, d := range r.Defs {
        sb.WriteString(fmt.Sprintf("  {\"file\":\"%s\",\"line\":%d,\"type\":\"%s\",\"name\":\"%s\"}", d.File, d.Line, d.Type, d.Name))
        if i < len(r.Defs)-1 {
            sb.WriteString(",\n")
        } else {
            sb.WriteString("\n")
        }
    }
    sb.WriteString("]\n")
    return sb.String()
}

// Scanner walks the file tree and collects definitions and references.
type Scanner struct {
    root     string
    excludes []string
}

func NewScanner(root string, excludes []string) *Scanner {
    return &Scanner{root: root, excludes: excludes}
}

func (s *Scanner) Scan() (*UnusedReport, error) {
    defs := make(map[string]UnusedDef)
    refs := make(map[string]bool)

    err := filepath.WalkDir(s.root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            for _, ex := range s.excludes {
                if strings.HasPrefix(path, ex) {
                    return filepath.SkipDir
                }
            }
            return nil
        }
        if !isSourceFile(path) {
            return nil
        }
        fileDefs, fileRefs, err := parseFile(path)
        if err != nil {
            return err
        }
        for k, v := range fileDefs {
            defs[k] = v
        }
        for k := range fileRefs {
            refs[k] = true
        }
        return nil
    })
    if err != nil {
        return nil, err
    }

    report := &UnusedReport{}
    for name, def := range defs {
        if !refs[name] {
            report.Defs = append(report.Defs, def)
        }
    }
    return report, nil
}

func (s *Scanner) Fix(report *UnusedReport) error {
    // Group definitions by file for efficient processing
    defsByFile := make(map[string][]UnusedDef)
    for _, d := range report.Defs {
        defsByFile[d.File] = append(defsByFile[d.File], d)
    }

    for file, defs := range defsByFile {
        // Read all lines of the file
        data, err := os.ReadFile(file)
        if err != nil {
            return err
        }
        lines := strings.Split(string(data), "\n")

        // Sort definitions by line number descending to avoid index shift when deleting
        sort.Slice(defs, func(i, j int) bool {
            return defs[i].Line > defs[j].Line
        })

        for _, d := range defs {
            idx := d.Line - 1
            if idx >= 0 && idx < len(lines) {
                lines = append(lines[:idx], lines[idx+1:]...)
            }
        }

        // Write the modified content back to the file
        newData := strings.Join(lines, "\n")
        if err := os.WriteFile(file, []byte(newData), 0644); err != nil {
            return err
        }
    }
    return nil
}

func isSourceFile(path string) bool {
    ext := filepath.Ext(path)
    switch ext {
    case ".py", ".js", ".ts", ".go", ".rs":
        return true
    default:
        return false
    }
}

// Regular expressions for detecting function definitions and references.
var (
    pyDefRe   = regexp.MustCompile(`^\s*def\s+(\w+)\s*\(`)
    pyRefRe   = regexp.MustCompile(`\b(\w+)\b`)
    goDefRe   = regexp.MustCompile(`^\s*func\s+(\w+)\s*\(`)
    goRefRe   = regexp.MustCompile(`\b(\w+)\b`)
    jsDefRe   = regexp.MustCompile(`^\s*function\s+(\w+)\s*\(`)
    jsRefRe   = regexp.MustCompile(`\b(\w+)\b`)
    rsDefRe   = regexp.MustCompile(`^\s*fn\s+(\w+)\s*\(`)
    rsRefRe   = regexp.MustCompile(`\b(\w+)\b`)
)

func parseFile(path string) (map[string]UnusedDef, map[string]bool, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, nil, err
    }
    defer f.Close()

    defs := make(map[string]UnusedDef)
    refs := make(map[string]bool)

    scanner := bufio.NewScanner(f)
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        var defName string
        switch filepath.Ext(path) {
        case ".py":
            if m := pyDefRe.FindStringSubmatch(line); m != nil {
                defName = m[1]
                defs[defName] = UnusedDef{File: path, Line: lineNum, Type: "function", Name: defName}
            }
        case ".go":
            if m := goDefRe.FindStringSubmatch(line); m != nil {
                defName = m[1]
                defs[defName] = UnusedDef{File: path, Line: lineNum, Type: "function", Name: defName}
            }
        case ".js", ".ts":
            if m := jsDefRe.FindStringSubmatch(line); m != nil {
                defName = m[1]
                defs[defName] = UnusedDef{File: path, Line: lineNum, Type: "function", Name: defName}
            }
        case ".rs":
            if m := rsDefRe.FindStringSubmatch(line); m != nil {
                defName = m[1]
                defs[defName] = UnusedDef{File: path, Line: lineNum, Type: "function", Name: defName}
            }
        }
        // Collect references, skipping the function name defined on this line.
        for _, m := range pyRefRe.FindAllStringSubmatch(line, -1) {
            name := m[1]
            if name == defName {
                continue
            }
            refs[name] = true
        }
    }
    if err := scanner.Err(); err != nil {
        return nil, nil, err
    }
    return defs, refs, nil
}
