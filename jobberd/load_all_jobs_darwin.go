// Mac OS X does not store users in /etc/passwd.  Instead we shell out
// to OpenDirectory.
// http://apple.stackexchange.com/questions/29874/how-can-i-list-all-user-accounts-in-the-terminal

package main

import (
	"os/exec"
	"strings"
)

func (m *JobManager) LoadAllJobs() (int, error) {

	cmdName := "/usr/bin/dscl"
	cmdArgs := []string{".", "list", "/Users"}

	var (
		cmdOut []byte
		err    error
	)

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		ErrLogger.Printf("There was an error running dscl to get users: %v\n", err)
		return 0, err
	}

	users := strings.Split(string(cmdOut), "\n")
	totalJobs := 0

	for _, user := range users {
		user = strings.Trim(user, "\n")
		user = strings.Trim(user, " ")
		if len(user) > 0 && user[0:1] != "_" {
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
