// +build !darwin

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (m *JobManager) LoadAllJobs() (int, error) {
	// get all users by reading passwd
	f, err := os.Open("/etc/passwd")
	if err != nil {
		ErrLogger.Printf("Failed to open /etc/passwd: %v\n", err)
		return 0, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	totalJobs := 0
	for scanner.Scan() {
		line := scanner.Text()
		commentStart := strings.Index(line, "#")
		fmt.Printf("%s commentStart: %v\n", line, commentStart)
		if commentStart >= 0 {
			line = line[0:commentStart]
		}

		fmt.Printf("Line: %s\n", line)
		parts := strings.Split(line, ":")
		if len(line) > 0 && len(parts) > 0 {
			user := parts[0]
			nbr, err := m.loadJobsForUser(user)
			totalJobs += nbr
			if err != nil {
				ErrLogger.Printf("Failed to load jobs for %s: %v\n", err, user)
			}
		}
	}

	ErrLogger.Printf("totalJobs: %v; len(m.jobs): %v", totalJobs, len(m.jobs))

	return len(m.jobs), nil
}
