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

func ListFunctions(file string, disableSyntax bool) {
	f, _ := os.ReadFile(file)
	c, _ := format.Source(f)
	var re = regexp.MustCompile(`(?m)^func.*?{$`)
	for _, match := range re.FindAllString(string(c), -1) {
		if !disableSyntax {
			quick.Highlight(os.Stdout, match+"}\n", "go", "terminal256", "monokai")
		} else {
			fmt.Println(match)
		}
	}
}

func CatFunction(funcName, file string, disableSyntax bool) {
	f, _ := os.ReadFile(file)
	c, _ := format.Source(f)
	//var re = regexp.MustCompile(`(?ms)^func.*` + Fname + `.*?{\n.*?^}(\s+?)$`)
	var re = regexp.MustCompile(`(?ms)^func (\([A-Za-z0-9_]+.*?\)[ |\t])?` + funcName + `.*?{\n.*?^}(\s+?)$`)
	for _, match := range re.FindAllString(string(c), -1) {
		if !disableSyntax {
			quick.Highlight(os.Stdout, match, "go", "terminal256", "monokai")
		} else {
			fmt.Println(match)
		}
	}
}

func ListTypes(file string, disableSyntax bool) {
	f, _ := os.ReadFile(file)
	c, _ := format.Source(f)
	var re = regexp.MustCompile(`(?ms)^( +|\t+)?type\s[A-Za-z0-9_]+.*?\sstruct\s{\n.*?^( +|\t+)?}$`)
	for _, match := range re.FindAllString(string(c), -1) {
		if !disableSyntax {
			quick.Highlight(os.Stdout, match+"\n", "go", "terminal256", "monokai")
		} else {
			fmt.Println(match)
		}
	}
}

func ListMethod(Type, file string, disableSyntax bool) {
	f, _ := os.ReadFile(file)
	c, _ := format.Source(f)
	var re = regexp.MustCompile(`(?ms)^func (\([A-Za-z0-9_]+ [\*]?` + Type + `\)[ |\t])[A-Za-z0-9_].*?{$`)
	for _, match := range re.FindAllString(string(c), -1) {
		if !disableSyntax {
			quick.Highlight(os.Stdout, match+"}\n", "go", "terminal256", "monokai")
		} else {
			fmt.Println(match)
		}
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

func main() {
	var syntax, list, listtype bool
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
	op.On("-d", "--disable-syntax", "Disable syntax highlighting", &syntax)
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

	if len(op.Extra) == 0 || op.Extra[0] == "." || op.Extra[0] == "*" {
		files = GoFiles()
	} else {
		files = op.Extra
	}

	for _, i := range files {
		if len(files) > 1 {
			fmt.Println(i + ":")
		}
		switch {
		case list:
			ListFunctions(i, syntax)
		case cat != "":
			CatFunction(cat, i, syntax)
		case method != "":
			ListMethod(method, i, syntax)
		case listtype:
			ListTypes(i, syntax)
		}
	}
}
