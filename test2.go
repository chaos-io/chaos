package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func GetFreePort() (int, error) {
	net.LookupAddr("")
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func main1() {
	var wg sync.WaitGroup
	n := 1000

	for i := 0; i < n; i++ {
		wg.Add(1)

		i := i

		time.Sleep(time.Millisecond)
		go func() {
			if i%100 == 0 {
				fmt.Printf("\ni=%d\n", i)
			}
			port, err := GetFreePort()
			if err != nil {
				fmt.Println("错误:", err)
				return
			}
			fmt.Printf("%v ", port)
			if 30000 < port && port < 40000 {
				fmt.Printf("(%d) ", port)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func listen(port int) {
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Printf("listen %d\n", port)
	}
}
