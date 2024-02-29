package main

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

func main1() {
	// tests := []struct {
	// 	s    string
	// 	t    string
	// 	want bool
	// }{
	// 	// {s: "acb", t: "ahbgdc", want: false},
	// 	{s: "acb", t: "ahbgdcb", want: true},
	// }
	//
	// for _, test := range tests {
	// 	has := isSubsequence(test.s, test.t)
	// 	fmt.Printf("got=%v, want=%v\n", has, test.want)
	// }

	// value := "NaN"
	// v, err := strconv.ParseFloat(value, 32)
	// if err == nil && v > 0 {
	// 	if v == math.Inf(1) || v == math.Inf(-1) {
	// 		fmt.Printf("---1%T/%v", v, v)
	// 	}
	// 	fmt.Printf("---2%T/%v", v, v)
	// }
	// fmt.Printf("---3%T/%v\n", v, v)
	// fmt.Printf("---err=%v\n", err)
	// fmt.Printf("---math.IsNaN(v)=%v\n", math.IsNaN(v))

}

func GetFreePort() (int, error) {
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

func main() {
	port, err := GetFreePort()
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	fmt.Printf("找到未使用的端口：%d\n", port)
}

func main2() {
	// 假设文件中有一个包含非法浮点数的字段
	fileContent := "1.23\n4.56\ninvalid\n7.89"

	// 将文件内容按行拆分
	lines := strings.Split(fileContent, "\n")

	// 遍历每一行并解析为浮点数
	for _, line := range lines {
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			fmt.Printf("Error parsing float: %v %v\n", value, err)
			continue
		}

		// 打印解析后的值
		fmt.Printf("Parsed value: %v\n", value)

		// 检查是否为NaN
		if math.IsNaN(value) {
			fmt.Println("Encountered NaN!")
		}
	}
}

func isSubsequence(s string, t string) bool {
	for i, j := 0, 0; i < len(s); i++ {
		has := false
		for ; j < len(t); j++ {
			if s[i] == t[j] {
				has = true
				j++
				break
			}
		}
		if !has {
			return false
		}
	}
	return true
}

type mystruct struct {
	CoveragesSnapshot []string
	ma                map[string]string
}

func lengthOfLastWord(s string) int {
	s = strings.TrimSpace(s)
	split := strings.Split(s, " ")
	return len(split[len(split)-1])
}

func lengthOfLastWord2(s string) int {
	end := len(s) - 1
	for ; s[end] == ' '; end-- {
	}
	start := len(s[:end])
	for ; s[start] != ' '; start-- {
		if start >= 1 && s[start-1] == ' ' {
			break
		}

		if start == 0 {
			break
		}
	}
	return end - start + 1
}

func romanToInt(s string) int {
	total := 0
	for i := 0; i <= len(s)-1; i++ {
		var next byte
		if i < len(s)-1 {
			next = s[i+1]
		} else if i == len(s)-1 {
			next = 0
		}

		if n := seekNext(s[i], next); n > 0 {
			fmt.Println("---n1=", n)
			total += n
			fmt.Printf("c(%s%s)/i(%d), total=%d\n", string(s[i]), string(next), i, total)
			i++
		} else {
			n = seek(s[i])
			fmt.Println("---n2=", n)
			total += n
			fmt.Printf("c(%s)/i(%d), total=%d\n", string(s[i]), i, total)
		}
	}
	return total
}

func seek(c byte) int {
	switch c {
	case 'I':
		return 1
	case 'V':
		return 5
	case 'X':
		return 10
	case 'L':
		return 50
	case 'C':
		return 100
	case 'D':
		return 500
	case 'M':
		return 1000
	}
	return 0
}

func seekNext(c, next byte) int {
	switch c {
	case 'I':
		if next == 'V' {
			return 4
		} else if next == 'X' {
			return 9
		}
	case 'X':
		if next == 'L' {
			return 40
		} else if next == 'C' {
			return 90
		}
	case 'C':
		if next == 'D' {
			return 400
		} else if next == 'M' {
			return 900
		}
	}
	return 0
}
