package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hacker4257/pet_charity/pkg/utils"
)

const (
	uploadDir = "uploads"
	MaxSize   = 5 << 20 //5mb
)

var AllowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

var allowedMIME = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

//验证文件
func Validate(header *multipart.FileHeader) error {
	//检查大小
	if header.Size > MaxSize {
		return fmt.Errorf("file size exceeds 5mb limit")
	}

	//检查扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !AllowedExts[ext] {
		return fmt.Errorf("file type %s not allowed", ext)
	}

	//检查 MIME 类型
	f, err := header.Open()
	if err != nil {
		return fmt.Errorf("cannot open file")
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	mime := http.DetectContentType(buf[:n])
	if !allowedMIME[mime] {
		return fmt.Errorf("file content type %s not allowed", mime)
	}

	return nil
}

func SaveFile(file *multipart.FileHeader, subDir string) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	//生成唯一文件
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixMilli(), utils.RandomCode(8), ext)

	dir := filepath.Join(uploadDir, subDir)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create directory failed: %w", err)
	}

	savePath := filepath.Join(dir, filename)

	urlPath := fmt.Sprintf("/%s/%s/%s", uploadDir, subDir, filename)

	return urlPath, saveFileTocal(file, savePath)
}

func saveFileTocal(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
