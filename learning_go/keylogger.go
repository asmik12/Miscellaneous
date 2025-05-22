package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

const (
	logFile       = "login_events.txt"
	checkInterval = time.Minute
)

var seen = make(map[string]bool)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}
	username := currentUser.Username
	fmt.Println("Watching for user login events...")

	for {
		checkLoginEvents(username)
		time.Sleep(checkInterval)
	}
}

func checkLoginEvents(username string) {
	cmd := exec.Command("log", "show",
		"--predicate", `eventMessage contains "login" && process == "loginwindow"`,
		"--style", "compact",
		"--last", "1m",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		fmt.Println("log show error:", err)
		return
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), "login") &&
			(strings.Contains(line, username) || strings.Contains(line, "Login Window")) {

			// Use line hash to avoid duplicates
			lineKey := strings.TrimSpace(line)
			if seen[lineKey] {
				continue
			}
			seen[lineKey] = true

			ts := time.Now().Format("2006-01-02 15:04:05")
			entry := fmt.Sprintf("[%s] Login detected for user %s\n", ts, username)

			f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(entry)
				f.Close()
				fmt.Print(entry)
			}
		}
	}
}
