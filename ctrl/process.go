package ctrl

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

const LoadingText = "Loading..."

type Command struct {
	Label     string
	Cmd       *exec.Cmd
	OnSuccess func()
	OnError   func(error)
	// Optional
	OnOutput func(string)
}

func Execute(c *Command) {
	prefix := fmt.Sprint(" \x1b[35m[", c.Label, "] \x1b[0m")

	// Get stdout pipe (for OnOutput)
	stdout, err := c.Cmd.StdoutPipe()
	if err != nil {
		log.Printf(prefix, "\x1b[31mError obtaining stdout pipe: %v\x1b[0m\n", err)
		c.OnError(err)
	}

	log.Println(prefix, "executing \x1b[36m", c.Cmd.String(), "\x1b[0m")
	if err := c.Cmd.Start(); err != nil {
		log.Println(prefix, "error: \x1b[31m", err.Error(), "\x1b[0m")
		c.OnError(err)
	}

	c.OnSuccess()
	log.Println(prefix, "\x1b[32msuccess.\x1b[0m")

	if c.OnOutput != nil {
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				log.Printf("%s\x1b[36mOutput: %s\n", prefix, line)
				c.OnOutput(line)
			}

			// Check for scanner errors
			if err := scanner.Err(); err != nil {
				log.Printf("%s\x1b[32mError reading stdout: %v\n", prefix, err)
				c.OnError(err)
			}

			// Wait for the command to finish
			if err := c.Cmd.Wait(); err != nil {
				log.Printf("%s\x1b[32mCommand finished with error: %v\n", prefix, err)
				c.OnError(err)
			}
		}()
	}
}
