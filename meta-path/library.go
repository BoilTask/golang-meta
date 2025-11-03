package metapath

import (
	"io"
	"io/fs"
	metaerror "meta/meta-error"
	metapanic "meta/meta-panic"
	"os"
	"path/filepath"
	"strings"
)

// GetBaseName 根据路径的获取文件名（不包含扩展名）
func GetBaseName(path string) string {
	// 获取文件名
	base := filepath.Base(path)
	// 获取扩展名
	ext := filepath.Ext(path)
	// 去掉扩展名
	name := strings.TrimSuffix(base, ext)
	return name
}

// CheckCreateFile 确保创建文件的目录
func CheckCreateFile(filePath string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return metaerror.New("failed to create directory: %v", filePath)
	}
	return err
}

func WriteStringToFile(path string, content string) error {
	err := CheckCreateFile(path)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return metaerror.Wrap(err, "failed to create file: %s", path)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			metapanic.ProcessError(metaerror.Wrap(err, "failed to close file: %s", path))
		}
	}(file)
	_, err = file.WriteString(content)
	if err != nil {
		return metaerror.Wrap(err, "failed to write string to file: %s", path)
	}
	return nil
}

func CopyFile(srcPath, targetPath string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return metaerror.Wrap(err, "source path error")
	}

	isTargetDir := strings.HasSuffix(targetPath, string(os.PathSeparator)) ||
		strings.HasSuffix(targetPath, "/") || strings.HasSuffix(targetPath, "\\")

	if srcInfo.IsDir() {
		// src 是目录
		if !isTargetDir {
			// targetPath 是文件路径，但 src 是目录 —— 不合法
			return metaerror.New("cannot copy directory '%s' to file '%s'", srcPath, targetPath)
		}
		return copyDir(srcPath, targetPath)
	} else {
		// src 是文件
		if isTargetDir {
			// target 是目录 -> 保持文件名
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return metaerror.New("failed to create target directory '%s': %w", targetPath, err)
			}
			targetPath = filepath.Join(targetPath, filepath.Base(srcPath))
		} else {
			// target 是文件路径 -> 确保其父目录存在
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return metaerror.New("failed to create parent directory for '%s': %w", targetPath, err)
			}
		}
		return copySingleFile(srcPath, targetPath)
	}
}

func copySingleFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return metaerror.Wrap(err, "failed to open source file '%s'", srcFile)
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			metapanic.ProcessError(err)
		}
	}(src)

	dst, err := os.Create(dstFile)
	if err != nil {
		return metaerror.Wrap(err, "failed to create target file '%s'", dstFile)
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			metapanic.ProcessError(err)
		}
	}(dst)

	if _, err := io.Copy(dst, src); err != nil {
		return metaerror.Wrap(err, "copy failed")
	}

	// 可选：复制权限
	if info, err := os.Stat(srcFile); err == nil {
		err := os.Chmod(dstFile, info.Mode())
		if err != nil {
			return metaerror.Wrap(err)
		}
	}
	return nil
}

func copyDir(srcDir, dstDir string) error {
	return filepath.WalkDir(
		srcDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(srcDir, path)
			if err != nil {
				return err
			}
			targetPath := filepath.Join(dstDir, relPath)
			if d.IsDir() {
				return os.MkdirAll(targetPath, 0755)
			}
			return copySingleFile(path, targetPath)
		},
	)
}
