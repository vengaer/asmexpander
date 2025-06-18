package expander

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

type lineRange struct {
	first, last int
}

type macro struct {
	name string
	args []string
	lineRange
}

func collectMacros(code []string, verbose bool) ([]macro, error) {
	regex, err := regexp.Compile(`\.macro\s*(\w+)\s*(.*)\s*$`)
	if err != nil {
		return nil, err
	}

	macros := []macro{}
	inMacro := false

	for i, line := range code {
		if strings.Contains(line, ".endm") {
			if !inMacro {
				return nil, fmt.Errorf("Unbalanced .endm on line %d\n", i)
			}
			inMacro = false
			macros[len(macros) - 1].lineRange.last = i - 1
		}
		matches := regex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		var args []string = nil
		if len(matches[2]) > 0 {
			args = strings.Split(matches[2], ",")
		}

		if verbose {
			log.Printf("Found macro '%s' with %d arguments at line %d (%s) %v %v\n",
				matches[1], len(args), i, line, reflect.TypeOf(args), len(args))
		}

		for i := 0; i < len(args); i++ {
			args[i] = strings.TrimSpace(args[i])
		}
		macros = append(macros, macro{matches[1], args, lineRange{ i + 1, -1 }})
		inMacro = true
	}

	return macros, nil
}

func expandMacros(code []string, macros []macro, verbose bool) ([]string, bool, error) {
	names := []string{}
	macroLookup := map[string]macro{}
	for _, macro := range macros {
		names = append(names, macro.name)
		macroLookup[macro.name] = macro
	}
	pattern := `^\s*(` + strings.Join(names, "|") + ")"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, false, err
	}

	found := false
	expanded := []string{}
	for i, line := range code {
		matches := regex.FindStringSubmatch(line)
		if matches == nil {
			expanded = append(expanded, line)
			continue
		}

		line = strings.ReplaceAll(line, ",", "")
		line = strings.Split(line, "#")[0]
		line = strings.TrimSpace(line)
		words := strings.Split(line, " ")
		name := words[0]
		repl := map[string]string{}
		macro := macroLookup[name]

		if len(words) - 1 != len(macro.args) {
			return nil, false,
				fmt.Errorf("Invocation of '%s' at line %d passed %d parameters rather than %d (%s)\n",
					name, i, len(words) - 1, len(macro.args), strings.Join(words, " "))
		}

		for j, arg := range macro.args {
			repl[arg] = words[j + 1]
		}

		if verbose {
			log.Printf("Found invocation of '%s' at line %d\n", macro, i)
		}

		lineRange := macro.lineRange
		for j := lineRange.first; j <= lineRange.last; j++ {
			modified := code[j]
			for k, v := range repl {
				modified = strings.ReplaceAll(modified, `\` + k, v)
			}
			expanded = append(expanded, modified)
		}

		found = true
	}

	return expanded, found, nil
}

func Expand(code []string, verbose bool) ([]string, error) {
	macros, err := collectMacros(code, verbose)
	if err != nil {
		return nil, err
	}

	found := true
	for found {
	 	code, found, err = expandMacros(code, macros, verbose)
		if err != nil {
			return nil, err
		}
	}

	return code, nil
}
