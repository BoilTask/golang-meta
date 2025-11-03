package svnbranch

import (
	"fmt"
	"strings"
)

func GetSvnBranchNameByPath(filePath string, rootPath ...string) string {
	filePath = strings.TrimPrefix(filePath, "/")
	if strings.HasPrefix(filePath, "trunk") {
		return "trunk"
	}
	filePath = strings.TrimPrefix(filePath, "/")
	filePath = strings.TrimPrefix(filePath, "branches")
	filePath = strings.TrimPrefix(filePath, "/")
	for _, path := range rootPath {
		index := strings.Index(filePath, path)
		if index >= 0 {
			// 获取path之前的路径
			filePath = filePath[:index]
			break
		}
	}
	filePath = strings.TrimSuffix(filePath, "/")
	index := strings.Index(filePath, "/")
	if index >= 0 {
		branchFirstName := filePath[:index]
		indexSecond := strings.Index(filePath[index+1:], "/")
		var branchSecondName string
		if indexSecond >= 0 {
			branchSecondName = filePath[index+1 : index+1+indexSecond]
		} else {
			branchSecondName = filePath[index+1:]
		}
		if strings.Contains(branchSecondName, branchFirstName) {
			return branchSecondName
		}
		return branchFirstName
	}
	return filePath
}

func GetSvnBranchNameByUrl(url string, repository string, rootPath ...string) string {
	// 获取url中的路径部分
	prefix := fmt.Sprintf("svn/%s/", repository)
	index := strings.Index(url, prefix)
	if index < 0 {
		return ""
	}
	filePath := url[index+len(prefix):]
	return GetSvnBranchNameByPath(filePath, rootPath...)
}
