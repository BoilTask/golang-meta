package metamd5

import (
	"crypto/md5"
	"fmt"
	"io"
	metaerror "meta/meta-error"
	metapanic "meta/meta-panic"
	"os"
)

// MultiFileMD5 计算多个文件的综合MD5
func MultiFileMD5(paths []string) (string, error) {
	hasher := md5.New()
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("open file %s: %w", path, err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				metapanic.ProcessError(err)
			}
		}(file)
		// 把当前文件的内容写入到 hasher
		if _, err := io.Copy(hasher, file); err != nil {
			return "", metaerror.Wrap(err, "read file %s: %w", path)
		}
	}
	sum := hasher.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}
