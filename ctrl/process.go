package ctrl

import (
	"bufio"
	"bytes"
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

func ScanLinesWithCR(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func Execute(c *Command) {
	prefix := fmt.Sprint(" \x1b[35m[", c.Label, "] \x1b[0m")

	// Get stdout pipe (for OnOutput)
	stdout, err := c.Cmd.StdoutPipe()
	if err != nil {
		log.Printf(prefix, "\x1b[31mError obtaining stdout pipe: %v\x1b[0m\n", err)
		c.OnError(err)
	}

	log.Println(prefix, "executing \x1b[33m", c.Cmd.String(), "\x1b[0m")
	if err := c.Cmd.Start(); err != nil {
		log.Println(prefix, "error: \x1b[31m", err.Error(), "\x1b[0m")
		c.OnError(err)
	}

	c.OnSuccess()
	log.Println(prefix, "\x1b[32msuccess.\x1b[0m")

	if c.OnOutput != nil {
		go func() {
			scanner := bufio.NewScanner(stdout)
			scanner.Split(ScanLinesWithCR)
			for scanner.Scan() {
				line := scanner.Text()
				log.Printf("%s\x1b[36m%s\n", prefix, line)
				c.OnOutput(line)
			}
			if err := c.Cmd.Wait(); err != nil {
				log.Printf("%s\x1b[31mError:%s\n", err.Error())
				c.OnError(err)
			}
			c.OnOutput("Success")
		}()
	}
}
