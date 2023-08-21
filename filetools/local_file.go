package filetools

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// UploadFileToLocal 上传文件到本地
func UploadFileToLocal(path string, r *http.Request) error {

	// 从请求中解析文件
	file, header, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("获取文件失败: %v", err)
	}
	defer file.Close()

	// 创建本地文件
	filePath := fmt.Sprintf("%s%s%s", LocalFilePath, path, header.Filename)
	localFile, err := os.Create(filePath)
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
