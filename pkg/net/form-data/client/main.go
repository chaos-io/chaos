//go:build local
// +build local

package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
)

func main() {
	url := "http://127.0.0.1:8088/submit-formdata"

	// 准备表单数据
	formData := map[string]string{
		"username": "john_doe",
		"password": "secretpassword",
		"age":      "25",
	}

	// 创建一个buffer来存储请求体
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加表单字段
	for key, value := range formData {
		_ = writer.WriteField(key, value)
	}

	// 关闭multipart writer
	err := writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	// 发送请求
	fmt.Printf("---body=|%v|\n", requestBody.String())
	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer response.Body.Close()

	// 处理响应
	fmt.Println("Response Status:", response.Status)
	// 这里可以添加进一步处理响应的逻辑
}
