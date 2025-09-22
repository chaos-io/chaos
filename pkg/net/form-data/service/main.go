//go:build local
// +build local

package main

import (
	"fmt"
	"net/http"
)

func handleFormData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 设置最大内存为10 MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	// 获取表单字段的值
	username := r.FormValue("username")
	password := r.FormValue("password")
	age := r.FormValue("age")

	// 处理表单数据，这里简单打印到控制台
	fmt.Printf("Received FormData: username=%s, password=%s, age=%s\n", username, password, age)

	// 可以根据需要进行进一步的业务逻辑处理

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("FormData received successfully"))
}

func main() {
	http.HandleFunc("/submit-formdata", handleFormData)

	port := 8088
	fmt.Printf("Server is running on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
