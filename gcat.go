package main

import (
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alecthomas/chroma/quick"
	"github.com/mmpx12/optionparser"
)

const version = "v1.0.1"

var DisableSyntax bool

func ListFunctions(file string) {
	var re = regexp.MustCompile(`(?m)^func.*?{$`)
	for _, match := range re.FindAllString(file, -1) {
		PrintRes(match + "}\n")
	}
}

func CatFunction(funcName, file string) {
	var re = regexp.MustCompile(`(?ms)^func (\([A-Za-z0-9_]+.*?\)[ |\t])?` + funcName + `.*?{\n.*?^}(\s+?)$`)
	for _, match := range re.FindAllString(file, -1) {
		PrintRes(match)
	}
}

func ListTypes(file string) {
	var re = regexp.MustCompile(`(?ms)^( +|\t+)?type\s[A-Za-z0-9_]+\sstruct\s{\n.*?^( +|\t+)?}$`)
	for _, match := range re.FindAllString(file, -1) {
		PrintRes(match + "\n")
	}
}

func ListMethod(Type, file string) {
	var re = regexp.MustCompile(`(?ms)^func (\([A-Za-z0-9_]+ [\*]?` + Type + `\)[ |\t])[A-Za-z0-9_].*?{$`)
	for _, match := range re.FindAllString(file, -1) {
		PrintRes(match + "}\n")
	}
}

func GoFiles() []string {
	var files []string
	filepath.WalkDir(".", func(s string, d fs.DirEntry, e error) error {
		if filepath.Ext(d.Name()) == ".go" {
			files = append(files, s)
		}
		return nil
	})
	return files
}

func PrintVersion() {
	fmt.Println("gcat version:", version)
	os.Exit(1)
}

func PrintRes(res string) {
	if !DisableSyntax {
		quick.Highlight(os.Stdout, res, "go", "terminal256", "monokai")
	} else {
		fmt.Printf(res)
	}
}

func main() {
	var list, listtype bool
	var cat, method string
	op := optionparser.NewOptionParser()
	op.Banner = `List or print functions, type, method in go files

Usage: gcat [OPTION] FILE1 FILE2 ...
Check recursivly on current directory if no files are passed

options:`
	op.On("-l", "--list-functions", "List all functions", &list)
	op.On("-p", "--print-function FUNC", "Cat FUNC function", &cat)
	op.On("-m", "--method TYPE", "List TYPE method", &method)
	op.On("-t", "--list-types", "List types", &listtype)
	op.On("-d", "--disable-syntax", "Disable syntax highlighting", &DisableSyntax)
	op.On("-v", "--version", "Print version and exit", PrintVersion)
	op.Exemple("List all function in all go files in current directory:")
	op.Exemple("gcat -l\n")
	op.Exemple("Print main function in file.go")
	op.Exemple("gcat -p main file.go")
	op.Parse()

	var files []string

	if !list && !listtype && cat == "" && method == "" {
		op.Help()
		os.Exit(1)
	}

	if len(op.Extra) == 0 || op.Extra[0] == "." {
		files = GoFiles()
	} else {
		files = op.Extra
	}

	for _, i := range files {
		f, err := os.ReadFile(i)
		if err != nil {
			panic(err)
		}
		formated, err := format.Source(f)
		if err != nil {
			panic(err)
		}
		if len(files) > 1 {
			fmt.Println(i + ":")
		}
		switch {
		case list:
			ListFunctions(string(formated))
		case cat != "":
			CatFunction(cat, string(formated))
		case method != "":
			ListMethod(method, string(formated))
		case listtype:
			ListTypes(string(formated))
		}
	}
}
