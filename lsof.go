package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
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
