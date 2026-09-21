package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func printPIDs(entries []portEntry) {
	seen := make(map[string]bool)
	for _, e := range entries {
		if seen[e.PID] {
			continue
		}
		seen[e.PID] = true
		fmt.Println(e.PID)
	}
}

func printTable(entries []portEntry) {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(w, "COMMAND\tPID\tUSER\tPROTO\tADDRESS\tSTATE")
	for _, e := range entries {
		state := e.State
		if state == "" {
			state = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", e.Command, e.PID, e.User, e.Protocol, e.Name, state)
	}
	w.Flush()
}
