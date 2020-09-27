package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		path = "."
	}

	types := map[string]int{}

	must(filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && len(info.Name()) > 1 {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(path)
		types[ext]++
		return nil
	}))

	keys := []string{}
	for k := range types {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s\t%d\n", k, types[k])
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
