package ctrl

import (
	"bufio"
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2/data/binding"
)

const LoadingText = "Loading..."

// BindExecOutput binds the cmd output to a binding.String that you used to bind on a UI element.
// Usage ctrl.BindExecOutput(exec.Command("bash", "-c", "for i in {1..5}; do echo $i; sleep 1; done"), loadStatusString)
func BindExecOutput(cmd *exec.Cmd, bindStr binding.String) {
	// Get the output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error obtaining stdout pipe: %v\n", err)
		bindStr.Set(err.Error())
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		bindStr.Set(err.Error())
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("Output: %s\n", line)
			acc, err := bindStr.Get()
			if err != nil {
				bindStr.Set(err.Error())
				break
			}

			// First overwriting of the loading text
			if acc == LoadingText {
				bindStr.Set(line)
			} else {
				bindStr.Set(acc + line)
			}
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading stdout: %v\n", err)
			bindStr.Set(err.Error())
		}

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			fmt.Printf("Command finished with error: %v\n", err)
			bindStr.Set(err.Error())
		}
	}()
}
