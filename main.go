package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/tabwriter"
)

type portEntry struct {
	PID      string
	Command  string
	User     string
	FD       string
	Protocol string
	Name     string
	State    string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: psport [-q|--quiet] <port>")
		os.Exit(1)
	}

	quiet := false
	var portArg string
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-q", "--quiet":
			quiet = true
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

	if quiet {
		printPIDs(entries)
		return
	}

	printTable(entries)
}

func lookupPort(port int) ([]portEntry, error) {
	cmd := exec.Command("lsof", "-nP", "-i", fmt.Sprintf(":%d", port), "-F", "pcuLPnT")
	out, err := cmd.Output()
	if err != nil {
		// lsof exits 1 when nothing matches; treat that as "no results", not an error.
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 && len(out) == 0 {
			return nil, nil
		}
		return nil, err
	}

	var entries []portEntry
	var pid, command, user string
	var cur *portEntry

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		tag, val := line[0], line[1:]

		switch tag {
		case 'p':
			pid = val
			command, user = "", ""
			cur = nil
		case 'c':
			command = val
		case 'L':
			user = val
		case 'f':
			entries = append(entries, portEntry{PID: pid, Command: command, User: user, FD: val})
			cur = &entries[len(entries)-1]
		case 'P':
			if cur != nil {
				cur.Protocol = val
			}
		case 'n':
			if cur != nil {
				cur.Name = val
			}
		case 'T':
			if cur != nil && strings.HasPrefix(val, "ST=") {
				cur.State = strings.TrimPrefix(val, "ST=")
			}
		}
	}

	return entries, scanner.Err()
}

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
