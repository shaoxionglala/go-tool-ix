package filetools

var (
	// LocalFilePath 本地上传文件目录
	LocalFilePath = ""
)

func InitLocalFile(path string) {
	LocalFilePath = path
}
