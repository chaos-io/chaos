package main

import (
	"fmt"
	"reflect"
	"strings"
)

func main1() {
	tests := []struct {
		nums []int
		k    int
		want bool
	}{
		{[]int{1, 2, 3, 1}, 3, true},
		{[]int{1, 0, 1, 1}, 1, true},
		{[]int{1, 2, 3, 1, 2, 3}, 2, false},
	}

	for _, tt := range tests {
		got := containsNearbyDuplicate(tt.nums, tt.k)
		if got != tt.want {
			fmt.Printf("something wrong, nums=%v, k=%v, want=%v, got=%v\n", tt.nums, tt.k, tt.want, got)
		}
	}
	fmt.Println("all done!")
}

func isIsomorphic(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	s2t := make(map[byte]byte)
	t2s := make(map[byte]byte)

	for i := range s {
		_s, _t := s[i], t[i]
		if s2t[_s] > 0 && s2t[_s] != _t || t2s[_t] > 0 && t2s[_t] != _s {
			return false
		}
		s2t[_s] = _t
		t2s[_t] = _s
	}

	return true
}

func wordPattern(pattern string, s string) bool {
	ss := strings.Split(s, " ")
	if len(pattern) != len(ss) {
		return false
	}

	p2s := make(map[string]string)
	s2p := make(map[string]string)
	for i := range pattern {
		_p, _s := string(pattern[i]), ss[i]
		if p2s[_p] != "" && p2s[_p] != _s || s2p[_s] != "" && s2p[_s] != _p {
			return false
		}

		p2s[_p] = _s
		s2p[_s] = _p
	}

	return true
}

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	sMap := make(map[byte]int, len(s))
	tMap := make(map[byte]int, len(t))
	for i := range s {
		sMap[s[i]]++
	}
	for i := range t {
		tMap[t[i]]++
	}

	return reflect.DeepEqual(sMap, tMap)
}

func twoSum(nums []int, target int) []int {
	found := make(map[int]int, len(nums))
	for i, num := range nums {
		if idx, ok := found[target-num]; ok {
			return []int{idx, i}
		}
		found[num] = i
	}
	return nil
}

func isHappy(n int) bool {
	if n == 1 {
		return true
	}

	exist := make(map[int]bool)
	exist[n] = true
	return findHappy(n, exist)
}

func findHappy(n int, exist map[int]bool) bool {
	var nums []int
	for n > 0 {
		remainder := n % 10
		n /= 10
		nums = append(nums, remainder)
	}

	for _, num := range nums {
		n += num * num
	}

	if n == 1 {
		return true
	}

	if exist[n] {
		return false
	}

	exist[n] = true

	return findHappy(n, exist)
}

func containsNearbyDuplicate(nums []int, k int) bool {
	dup := make(map[int][]int, len(nums))

	for i, num := range nums {
		if _nums, ok := dup[num]; !ok {
			dup[num] = []int{i}
		} else {
			_nums = append(_nums, i)
			dup[num] = _nums
		}
	}

	for _, _nums := range dup {
		if len(_nums) > 1 {
			for i := len(_nums) - 1; i >= 1; i-- {
				if _nums[i]-_nums[i-1] <= k {
					return true
				}
			}
		}
	}

	return false
}
