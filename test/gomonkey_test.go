package test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

// 要进行monkey patch的函数
func hello() string {
	return "Hello, World!"
}

func TestFunc(t *testing.T) {
	output := "hello"
	patches := gomonkey.ApplyFunc(hello, func() string { return output })
	assert.Equal(t, output, hello())
	patches.Reset()

	output2 := "hello2"
	patches = gomonkey.ApplyFuncReturn(hello)
	assert.Equal(t, output2, hello())
	patches.Reset()
}

var num = 2

func TestVar(t *testing.T) {
	patches := gomonkey.ApplyGlobalVar(&num, 3)
	defer patches.Reset()

	assert.Equal(t, 3, num)
}

type Task struct {
	A, B int
}

func (t *Task) Add() int {
	return t.A + t.B
}

func (t *Task) sub() int {
	return t.A - t.B
}

func TestMethod(t *testing.T) {
	task := &Task{A: 3, B: 2}
	patches := gomonkey.ApplyMethod(task, "Add", func(task *Task) int { return task.A - task.B })
	result := task.Add()
	assert.Equal(t, 1, result)
	patches.Reset()

	patches = gomonkey.ApplyPrivateMethod(task, "sub", func(_ *Task) int { return 10 })
	sub := task.sub()
	assert.Equal(t, 10, sub)
	patches.Reset()
}
