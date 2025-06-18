package main

import (
	expander "github.com/vengaer/asmexpander/internal/expander"
	"fmt"
	"flag"
	"log"
	"path/filepath"
	"os"
	"strings"
)

var progname = filepath.Base(os.Args[0])

func main() {
	outfile := flag.String("outfile", "", "Desired output file")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Expand GNU-style assembly macros
Usage:
	%s [flags] code
`, progname)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Code must be provided as single positional argument")
	}

	args := strings.Split(flag.Args()[0], "\n")
	expanded, err := expander.Expand(args, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	code := strings.Join(expanded, "\n")
	if *outfile == "" {
		fmt.Println(code)
	} else {
		os.WriteFile(*outfile, []byte(code), 0644)
	}
}
