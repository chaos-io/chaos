package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func parseFormData(body string) (string, error) {
	var boundary string
	pairs := strings.Split(body, "Content-Disposition: form-data;")
	if len(pairs) > 0 {
		boundary = strings.TrimPrefix(pairs[0], "--")
		for strings.HasSuffix(boundary, "\r") || strings.HasSuffix(boundary, "\n") {
			boundary = strings.TrimSuffix(boundary, "\r")
			boundary = strings.TrimSuffix(boundary, "\n")
		}
	}

	bodyBytes := []byte(body)

	// 创建一个新的multipart.Reader，使用指定的boundary
	mr := multipart.NewReader(bytes.NewReader(bodyBytes), boundary)

	// 初始化参数和文件存储
	formParams := make(map[string]string)
	fileParams := make(map[string][]byte)

	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}

		// 检查是否是文件字段
		if part.FileName() != "" {
			fileContent, err := io.ReadAll(part)
			if err != nil {
				return "", err
			}
			fileParams[part.FormName()] = fileContent
		} else {
			// 普通表单字段
			fieldValue, err := io.ReadAll(part)
			if err != nil {
				return "", err
			}
			formParams[part.FormName()] = string(fieldValue)
		}
	}

	jsonData, err := jsoniter.MarshalToString(formParams)
	if err != nil {
		return "", err
	}

	return jsonData, nil
}

func main() {
	// 示例form-data请求体字符串
	formDataString := `--a18e148bfc61b7e072e2c700581cf1ca062202957a48919d254ef1f96cab
Content-Disposition: form-data; name="password"

secretpassword
--a18e148bfc61b7e072e2c700581cf1ca062202957a48919d254ef1f96cab
Content-Disposition: form-data; name="age"

25
--a18e148bfc61b7e072e2c700581cf1ca062202957a48919d254ef1f96cab
Content-Disposition: form-data; name="username"

john_doe
--a18e148bfc61b7e072e2c700581cf1ca062202957a48919d254ef1f96cab--
	`

	// 解析form-data请求体字符串
	formParams, err := parseFormData(formDataString)
	if err != nil {
		panic(err)
	}

	// 打印解析结果
	fmt.Println("Form Parameters:", formParams)
}
