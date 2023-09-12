package main

import (
	"flag"
	"fmt"
	"golang.org/x/mod/modfile"
	"os"
)

var root = os.Getenv("GOPATH") + "/pkg/mod"

func main() {
	var p = flag.String("path", "./", "main mod path")
	flag.Parse()
	iterator(*p)
}

func iterator(path string) {
	goModPath := fmt.Sprintf("%s/go.mod", path)
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to read go.mod file: %s\n", err.Error())
		return
	}
	modFile, err := modfile.Parse("go.mod", goModContent, nil)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to parse go.mod file: %s\n", err.Error())
		return
	}
	for _, require := range modFile.Require {
		if require.Indirect {
			continue
		}
		next := require.Mod.String()
		isReplaced := false
		for _, r := range modFile.Replace {
			if r.Old.Path == path {
				isReplaced = true
				next = r.New.String()
				break
			}
		}

		if isReplaced {
			next = fmt.Sprintf("%s", next)
		} else {
			next = fmt.Sprintf("%s/%s", root, next)
		}
		_, _ = fmt.Fprintf(os.Stdout, "%s %s\n", modFile.Module.Mod.String(), require.Mod.Path)
		iterator(next)
	}
}
