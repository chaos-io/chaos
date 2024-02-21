package main

import (
	"fmt"
	"strings"
)

func main() {
	s := []struct {
		roman   string
		integer int
	}{
		// {roman: "III", integer: 3},
		// {roman: "LVIII ", integer: 5},
		// {roman: "MCMXCIV ", integer: 7},
		// {roman: "Hello World", integer: 5},
		{roman: " a", integer: 1},
	}

	for _, _s := range s {
		toInt := lengthOfLastWord2(_s.roman)
		if toInt != _s.integer {
			fmt.Printf("err, _s.roman=%s, toInt=%d, _s.integer=%d\n", _s.roman, toInt, _s.integer)
		}
		fmt.Println("------------")
	}
	// wordLen := lengthOfLastWord("Today is a nice day")
	// fmt.Println(wordLen)

	// var my mystruct
	// my.CoveragesSnapshot = []string{"1"}
	// my.ma["a"] = "a"
	// fmt.Printf("my=%v", my)

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
