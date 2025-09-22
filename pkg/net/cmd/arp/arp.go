package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

// arp 192.168.0.255
func main() {
	wg := sync.WaitGroup{}
	for i := 0; i <= 255; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cmd := newCmd(fmt.Sprintf("arp %s", "192.168.0."+strconv.Itoa(i)))
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return
			}
			bytes, err := cmd.Output()
			if err == nil {
				fmt.Println(string(bytes))
			}
		}(i)
	}
	wg.Wait()
}

func newCmd(cmdStr string) *exec.Cmd {
	return exec.Command("/bin/bash", "-c", cmdStr)
}
