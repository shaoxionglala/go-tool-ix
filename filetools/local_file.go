package filetools

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// 上传文件到本地
func UploadFileToLocal(path string, r *http.Request) error {
	// 检查请求方法
	if r.Method != http.MethodPost {
		return fmt.Errorf("方法不允许: 只允许POST方法")
	}

	// 从请求中解析文件
	file, header, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("获取文件失败: %v", err)
	}
	defer file.Close()

	// 创建本地文件
	localFile, err := os.Create(header.Filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer localFile.Close()

	// 将上传的文件内容写入本地文件
	_, err = io.Copy(localFile, file)
	if err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}
	return nil
}
