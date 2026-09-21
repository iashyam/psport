package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: psport [-q|--quiet] [-f|--free] <port>")
		os.Exit(1)
	}

	quiet := false
	free := false
	var portArg string
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-q", "--quiet":
			quiet = true
		case "-f", "--free":
			free = true
		default:
			portArg = arg
		}
	}

	port, err := strconv.Atoi(portArg)
	if err != nil {
		fmt.Println("invalid port:", portArg)
		os.Exit(1)
	}

	entries, err := lookupPort(port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		if !quiet {
			fmt.Printf("nothing listening on port %d\n", port)
		}
		return
	}

	if quiet && !free {
		printPIDs(entries)
		return
	}

	printTable(entries)

	if free {
		freeProcesses(entries)
	}
}
