package jun_util

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

const (
	//FileUploadFormat 文件上传按时间划分目录的时间格式
	FileUploadFormat = "200601"
	//FileUploadRoot 文件上传根目录名
	FileUploadRoot = "upload"
	//FileUploadThumb 文件上传图片缩略图存储目录名
	FileUploadThumb = "thumbnails"
)

// FileExt 从文件地址获取文件的后缀
func FileExt(name string) string {
	names := strings.Split(name, ".")
	return "." + names[len(names)-1]
}

// IsExist 检查文件目录是否已存在,不存在就创建
func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			errs := os.MkdirAll(path, 0755)
			if errs != nil {
				return false, errs
			}
		} else {
			return false, err
		}
	}
	return true, nil
}

// IsPathExist 检查文件目录是否已存在
func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// FileURL 根据传入的目录与文件名生成文件路径和相对URL地址
func FileURL(access string, filename string) (string, string, error) {
	getwd, _ := os.Getwd()
	month := time.Now().Format(FileUploadFormat)
	path := filepath.Join(getwd, FileUploadRoot, access, month)
	urlArr := []string{"/", FileUploadRoot, access, month, filename}
	url := strings.Join(urlArr, "/")
	_, err := IsExist(path)
	return filepath.Join(path, filename), url, err
}

// ThumbURL 根据传入的文件路径，文件名，宽度,高度生成缩略图，缩略图存储到常量设置的目录中，返回缩略图的URL地址
func ThumbURL(filepath string, filename string, width int, height int) (string, error) {
	var parseURL string
	if width == 0 {
		width = 300
	}
	if height == 0 {
		height = width
	}
	img, err := imaging.Open(filepath)
	if err != nil {
		return parseURL, err
	}
	thumb := imaging.Thumbnail(img, width, height, imaging.CatmullRom)
	dest := imaging.New(width, height, color.NRGBA{0, 0, 0, 0})
	dest = imaging.Paste(dest, thumb, image.Pt(0, 0))
	ThumbURL, parseURL, err := FileURL(FileUploadThumb, filename)
	if err != nil {
		return parseURL, err
	}
	err = imaging.Save(dest, ThumbURL)
	return parseURL, err

}

// ThumbStream 生成缩略图流
func ThumbStream(reader io.Reader, width int, height int) (*image.NRGBA, error) {
	if width == 0 {
		width = 300
	}
	if height == 0 {
		height = width
	}
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}
	thumb := imaging.Thumbnail(img, width, height, imaging.CatmullRom)
	dest := imaging.New(width, height, color.NRGBA{0, 0, 0, 0})
	dest = imaging.Paste(dest, thumb, image.Pt(0, 0))
	return dest, err

}

// ThumbStreamByPath 生成缩略图流
func ThumbStreamByPath(filepath string, width int, height int) (*image.NRGBA, error) {
	if width == 0 {
		width = 300
	}
	if height == 0 {
		height = width
	}
	img, err := imaging.Open(filepath)
	if err != nil {
		return nil, err
	}
	thumb := imaging.Thumbnail(img, width, height, imaging.CatmullRom)
	dest := imaging.New(width, height, color.NRGBA{0, 0, 0, 0})
	dest = imaging.Paste(dest, thumb, image.Pt(0, 0))
	return dest, err

}

// ChunkURL 根据传入的目录名，文件md5标识，分片序号返回分片的路径，分片存储目录地址
func ChunkURL(access string, identifier string, chunkNumber int) (string, string, error) {
	getwd, _ := os.Getwd()
	path := filepath.Join(getwd, FileUploadRoot, access, identifier)
	_, err := IsExist(path)
	return filepath.Join(path, strconv.Itoa(chunkNumber)), path, err
}

// WriteFile 往文件中写入数据
func WriteFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// CreateIfNotExist 创建文件
func CreateIfNotExist(file string) (*os.File, error) {
	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exist", file)
	}

	return os.Create(file)
}
