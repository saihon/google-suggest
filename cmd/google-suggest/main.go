package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	flag "github.com/saihon/flags"
	suggest "github.com/saihon/google-suggest"
)

const (
	URL = "https://www.google.com/complete/search"
)

type Flags struct {
	Language string
	Query    []string
}

var (
	Name    string
	Version string
	flags   Flags
)

func init() {
	flag.CommandLine.Init(Name, flag.ExitOnError, false)

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\nUsage: %s [options] [arguments]\n\n", Name)

		flag.VisitAll(func(f *flag.Flag) {
			s := ""
			if f.Alias > 0 {
				s = fmt.Sprintf("  -%c, --%s", f.Alias, f.Name)
			} else {
				s = fmt.Sprintf("  --%s", f.Name)
			}
			_, usage := flag.UnquoteUsage(f)
			if len(s) <= 4 {
				s += "\t"
			} else {
				s += "\n    \t"
			}
			s += strings.ReplaceAll(usage, "\n", "\n    \t")
			fmt.Fprint(w, s, "\n")
		})
	}

	flag.Bool("version", 'v', false,
		"Output version information\n",
		func(_ flag.Getter) error {
			fmt.Fprintf(flag.CommandLine.Output(), "%s: %s\n", Name, Version)
			return flag.ErrHelp
		})

	flag.String("query", 'q', "",
		"Specify search keywords. Also\ncan specify without this flag\n",
		func(g flag.Getter) error {
			flags.Query = append(flags.Query, g.String())
			return nil
		})

	flag.StringVar(&flags.Language, "lang", 'l', "",
		"Specify language. by default\nprobably your local language\n", nil)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Recover: %v", err)
			os.Exit(1)
		}
	}()

	os.Exit(_main())
}

func _main() int {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return 2
	}

	return 0
}

func Run() error {
	flag.Parse()

	if len(flags.Query) == 0 {
		flags.Query = flag.Args()
	}

	if len(flags.Query) == 0 {
		return errors.New("search query not specified")
	}

	rawurl := URL
	headers := map[string]string{
		`User-Agent`:      `Mozilla/5.0`,
		`Accept`:          `text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8`,
		`Accept-Language`: `en-US, en;q=0.5`,
	}
	timeout := 120 * time.Second

	f := suggest.NewFetcher(headers, timeout)
	a := make([]suggest.GoogleSuggestion, len(flags.Query), len(flags.Query))
	for i, query := range flags.Query {
		gs, err := f.Fetch(rawurl, query, flags.Language)
		if err != nil {
			return err
		}
		a[i] = gs
	}

	// struct to JSON
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(&a)
}
