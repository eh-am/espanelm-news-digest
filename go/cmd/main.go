package main

import (
	"bilingual-articles/download"
	"bilingual-articles/providers"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	var (
		rootFlagSet = flag.NewFlagSet("espanelm", flag.ExitOnError)

		mapFlagSet  = flag.NewFlagSet("map", flag.ExitOnError)
		mapProvider = mapFlagSet.String("provider", "elpais", "what provider to use. (Available: 'elpais')")

		downloadFlagSet = flag.NewFlagSet("download", flag.ExitOnError)
		downloadFrom    = downloadFlagSet.String("from", "", "what file generated by rss to use")
		downloadDest    = downloadFlagSet.String("dest", "", "where to download")
	)

	mapCmd := &ffcli.Command{
		Name:       "map",
		ShortUsage: "espanelm map [-p provider]",
		ShortHelp:  "generate the page map for a given provider",
		FlagSet:    mapFlagSet,
		Exec: func(_ context.Context, args []string) error {
			switch *mapProvider {
			case "elpais":
				elpais := providers.NewElPais(&providers.RSSGet{}, &http.Client{}, providers.Config{})

				pages, err := elpais.FetchPagesList()
				if err != nil {
					return err
				}

				marshalled, err := json.Marshal(pages)
				if err != nil {
					return err
				}

				fmt.Fprint(os.Stdout, string(marshalled))
			default:
				return errors.New("invalid provider")
			}
			return nil
		},
	}

	download := &ffcli.Command{
		Name:       "download",
		ShortUsage: "espanelm download [-f file] [-d dest]",
		ShortHelp:  "download articles",
		FlagSet:    downloadFlagSet,
		Exec: func(_ context.Context, args []string) error {
			if *downloadFrom == "" || *downloadDest == "" {
				return flag.ErrHelp
			}

			// read the file and unmarshal into pages
			var pages []providers.Page
			dat, err := ioutil.ReadFile(*downloadFrom)
			if err != nil {
				return err
			}
			err = json.Unmarshal(dat, &pages)
			if err != nil {
				return err
			}

			err = download.Download(pages, *downloadDest)
			if err != nil {
				return err
			}

			return nil
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "espanelm <subcommand>",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{mapCmd, download},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
