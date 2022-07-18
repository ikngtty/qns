package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikngtty/qns/pkg/qns"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected a subcommand")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "load":
		var settings qns.LoadSettings
		flagSet := flag.NewFlagSet("load", flag.ExitOnError)
		flagSet.IntVar(&settings.Pages, "pages", 50, "the number of pages to load")
		flagSet.Parse(os.Args[2:])

		qns.Load(settings)
	case "view":
		if len(os.Args) < 3 {
			fmt.Println("expected a kind of notifications")
			os.Exit(1)
		}
		qns.View(os.Args[2])
	default:
		fmt.Println("subcommand invalid")
		os.Exit(1)
	}
}
