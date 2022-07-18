package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikngtty/qns/pkg/qns"
)

func main() {
	var loadSettings qns.LoadSettings
	loadFlag := flag.NewFlagSet("load", flag.ExitOnError)
	loadFlag.IntVar(&loadSettings.Pages, "pages", 50, "the number of pages to load")

	if len(os.Args) < 2 {
		fmt.Println("expected a subcommand")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "load":
		loadFlag.Parse(os.Args[2:])
		qns.Load(loadSettings)
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
