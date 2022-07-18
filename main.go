package main

import (
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
		qns.Load()
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
