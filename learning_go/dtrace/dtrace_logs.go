package main

import (
	"bufio"
	"log"
	"os/exec"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetBorder(true).SetTitle("DTrace: Security-Related Syscalls")

	// Define DTrace script inline
	dtraceScript := `
    syscall::open*:entry
    {
        printf("%Y %s[%d] opened a file\n", walltimestamp, execname, pid);
    }`

	cmd := exec.Command("sudo", "dtrace", "-n", dtraceScript)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Failed to get stderr: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start DTrace: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			app.QueueUpdateDraw(func() {
				textView.Write([]byte(line + "\n"))
			})
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			app.QueueUpdateDraw(func() {
				textView.Write([]byte("[red]" + line + "[white]\n"))
			})
		}
	}()

	if err := app.SetRoot(textView, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("UI error: %v", err)
	}

	cmd.Wait()
}
