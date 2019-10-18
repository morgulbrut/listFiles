package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/homedepot/flop"
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
	Date     time.Time
	Size     int64
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
	cp := parser.Flag("", "copy", &argparse.Options{Required: false, Help: "Copy the collected files to one directory"})
	cpd := parser.String("", "copydir", &argparse.Options{Required: false, Help: "Directory to copy the collected files to", Default: "."})

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

	if *cp {
		color.Green("Copying files")
		copyFiles(d, *cpd)
	}

}

func exportCSV(docs []doc, path string) {
	var sb strings.Builder
	sb.WriteString("File, Path, Size, Mod.time\n")
	for _, d := range docs {
		sb.WriteString(fmt.Sprintf("%s, %s, %d, %v\n", d.Filename, d.Path, d.Size, d.Date.Format("2006-01-02 15:04:05")))
	}
	exportFile(sb.String(), path+".csv")
}

func exportMD(docs []doc, path string) {
	var sb strings.Builder
	sb.WriteString("| File | Path | Size | Mod.time\n")
	sb.WriteString("| ----- | ---- |---- |---- |\n")
	for _, d := range docs {
		sb.WriteString(fmt.Sprintf("| %s | %s | %d | %v | \n", d.Filename, d.Path, d.Size, d.Date.Format("2006-01-02 15:04:05")))
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

func copyFiles(docs []doc, path string) {
	options := flop.Options{
		Backup:   "numbered",
		MkdirAll: true,
	}
	for _, d := range docs {
		color.Yellow("Copying %s -> %s", d.Path, path+"/"+d.Filename)
		flop.Copy(d.Path, path+"/"+d.Filename, options)
	}
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
		docs = append(docs, doc{Path: path, Filename: info.Name(), Date: info.ModTime(), Size: info.Size()})
		return nil
	})

	if err != nil {
		panic(err)
	}

	return docs
}
