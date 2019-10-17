package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/morgulbrut/color"
	"github.com/morgulbrut/listFiles/version"
	"github.com/morgulbrut/toml"
)

type filter struct {
	Rgxp  string
	Match bool
}

type config struct {
	Paths   []string
	Filters []filter
}

type doc struct {
	Path     string
	Filename string
}

func readConfig(fn string) config {
	var c config
	if _, err := toml.DecodeFile(fn, &c); err != nil {
		fmt.Println(err)
	}
	return c
}

func main() {

	fmt.Println(version.DrawLogo())

	var d []doc

	parser := argparse.NewParser("listFiles", "generates a list of files based on filters (regex) in various file formats")
	cf := parser.String("c", "config", &argparse.Options{Required: false, Help: "Path to config file", Default: "config.toml"})
	out := parser.Flag("m", "markdown", &argparse.Options{Required: false, Help: "Export as markdown (default csv)"})
	fn := parser.String("f", "filename", &argparse.Options{Required: false, Help: "Output file name without ending", Default: "files"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	color.Green("Reading: %s", *cf)
	conf := readConfig(*cf)
	for _, p := range conf.Paths {
		color.Green("Collecting: %s", p)
		d = append(d, collect(p)...)
	}

	for _, f := range conf.Filters {
		color.Green("Filtering: %v -> %v", f.Rgxp, f.Match)
		d = useFilter(d, f)
	}
	if *out {
		color.Green("Writing: %s.md", *fn)
		exportMD(d, *fn)
	} else {
		color.Green("Writing: %s.csv", *fn)
		exportCSV(d, *fn)
	}
}

func exportCSV(docs []doc, path string) {
	var sb strings.Builder
	sb.WriteString("File, Path\n")
	for _, d := range docs {
		sb.WriteString(fmt.Sprintf("%s, %s\n", d.Filename, d.Path))
	}
	exportFile(sb.String(), path+".csv")
}

func exportMD(docs []doc, path string) {
	var sb strings.Builder
	sb.WriteString("| File | Path |\n")
	sb.WriteString("| ----- | ---- |\n")
	for _, d := range docs {
		sb.WriteString(fmt.Sprintf("| %s | %s | \n", d.Filename, d.Path))
	}
	exportFile(sb.String(), path+".md")
}

func exportFile(data string, path string) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(data)
	f.Close()
}

func useFilter(docs []doc, f filter) []doc {

	rgxp := regexp.MustCompile(f.Rgxp)

	var ret []doc
	for _, d := range docs {
		if rgxp.MatchString(d.Path) == f.Match {
			ret = append(ret, d)
		}
	}
	return ret
}

func collect(root string) []doc {
	var docs []doc

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		docs = append(docs, doc{Path: path, Filename: info.Name()})
		return nil
	})

	if err != nil {
		panic(err)
	}

	return docs
}
