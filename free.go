package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func freeProcesses(entries []portEntry) {
	type proc struct{ pid, command string }
	var procs []proc
	seen := make(map[string]bool)
	for _, e := range entries {
		if seen[e.PID] {
			continue
		}
		seen[e.PID] = true
		procs = append(procs, proc{pid: e.PID, command: e.Command})
	}

	noun := "process"
	if len(procs) > 1 {
		noun = "processes"
	}
	fmt.Printf("\nKill %d %s and free the port? [y/N]: ", len(procs), noun)

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer != "y" && answer != "yes" {
		fmt.Println("aborted, nothing killed")
		return
	}

	for _, p := range procs {
		if err := exec.Command("kill", p.pid).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to kill PID %s (%s): %v\n", p.pid, p.command, err)
			continue
		}
		fmt.Printf("killed PID %s (%s)\n", p.pid, p.command)
	}
}
