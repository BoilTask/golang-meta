package metazip

import (
	"archive/zip"
	"io"
	metaerror "meta/meta-error"
	metapanic "meta/meta-panic"
	"os"
	"path"
	"path/filepath"
)

// PackagePath 将某个路径下的所有文件压缩成一个压缩包，并删除原文件
func PackagePath(srcPath string, fileName string) error {
	var targetPath string
	if path.Base(fileName) == fileName {
		targetPath = filepath.Join(srcPath, fileName)
	} else {
		targetPath = fileName
	}
	zipFile, err := os.Create(targetPath)
	if err != nil {
		return metaerror.Wrap(err)
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			metapanic.ProcessError(err)
			return
		}
	}(zipFile)
	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			metapanic.ProcessError(err)
			return
		}
	}(zipWriter)
	err = filepath.Walk(
		srcPath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Clean(filePath) == filepath.Clean(zipFile.Name()) {
				return nil
			}
			relPath, err := filepath.Rel(srcPath, filePath)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			zipFileWriter, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}
			srcFile, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer func(srcFile *os.File) {
				err := srcFile.Close()
				if err != nil {
					metapanic.ProcessError(err)
					return
				}
			}(srcFile)
			_, err = io.Copy(zipFileWriter, srcFile)
			return err
		},
	)
	if err != nil {
		return metaerror.Wrap(err)
	}
	err = filepath.Walk(
		srcPath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Clean(filePath) == filepath.Clean(srcPath) {
				return nil
			}
			if filepath.Clean(filePath) == filepath.Clean(zipFile.Name()) {
				return nil
			}
			if err := os.RemoveAll(filePath); err != nil {
				return err
			}
			return nil
		},
	)
	return metaerror.Wrap(err)
}

func UzipFile(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			metapanic.ProcessError(err)
			return
		}
	}(r)
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return metaerror.Wrap(err)
	}
	for _, f := range r.File {
		fPath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fPath, os.ModePerm)
			if err != nil {
				return metaerror.Wrap(err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.Create(fPath)
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			err := outFile.Close()
			if err != nil {
				return metaerror.Wrap(err)
			}
			return err
		}
		_, err = io.Copy(outFile, rc)
		err = outFile.Close()
		if err != nil {
			return metaerror.Wrap(err)
		}
		err = rc.Close()
		if err != nil {
			return metaerror.Wrap(err)
		}
	}
	return nil
}
