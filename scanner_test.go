package main

import (
    "io/ioutil"
    "os"
    "path/filepath"
    "testing"
)

func TestScanner(t *testing.T) {
    tmpDir, err := ioutil.TempDir("", "deadcodefinder")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tmpDir)

    // create a simple python file
    pyContent := `def used_func():
    pass

def unused_func():
    pass

used_func()
`
    pyPath := filepath.Join(tmpDir, "test.py")
    if err := ioutil.WriteFile(pyPath, []byte(pyContent), 0644); err != nil {
        t.Fatal(err)
    }

    scanner := NewScanner(tmpDir, nil)
    report, err := scanner.Scan()
    if err != nil {
        t.Fatal(err)
    }

    if len(report.Defs) != 1 {
        t.Fatalf("expected 1 unused def, got %d", len(report.Defs))
    }
    if report.Defs[0].Name != "unused_func" {
        t.Fatalf("expected unused_func, got %s", report.Defs[0].Name)
    }
}
