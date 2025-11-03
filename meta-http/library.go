package metahttp

import (
	"io"
	metaerror "meta/meta-error"
	metapanic "meta/meta-panic"
	"meta/retry"
	"net/http"
	"path"
	"strings"
	"time"
)

func UrlJoin(url string, paths ...string) string {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	joinedPath := path.Join(paths...)
	if strings.HasPrefix(joinedPath, "/") {
		joinedPath = joinedPath[1:]
	}
	return url + joinedPath
}

func SendRequest(client *http.Client, method, url string, body io.Reader) (int, []byte, error) {
	fileListRequest, err := http.NewRequest(method, url, body)
	if err != nil {
		return -1, nil, metaerror.Wrap(err, "failed to create file list request")
	}
	fileListResp, err := client.Do(fileListRequest)
	if err != nil {
		return -1, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			metapanic.ProcessError(err)
		}
	}(fileListResp.Body)
	responseBody, err := io.ReadAll(fileListResp.Body)
	return fileListResp.StatusCode, responseBody, err
}

func SendRequestRetry(
	client *http.Client,
	name string,
	maxCount int,
	sleep time.Duration,
	method, url string,
	headers map[string]string,
	body io.Reader,
	checkStatus bool,
) (int, []byte, error) {
	var status int
	var responseBody []byte
	var finalErr error
	_ = retry.TryRetrySleep(
		name, maxCount, sleep, func(i int) bool {
			fileListRequest, err := http.NewRequest(method, url, body)
			if err != nil {
				finalErr = err
				return true
			}
			if len(headers) > 0 {
				for key, value := range headers {
					fileListRequest.Header.Set(key, value)
				}
			}
			fileListResp, err := client.Do(fileListRequest)
			if err != nil {
				finalErr = err
				return false
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					metapanic.ProcessError(err)
				}
			}(fileListResp.Body)
			status = fileListResp.StatusCode
			if checkStatus {
				if status != http.StatusOK {
					finalErr = metaerror.New(
						"unexpected status code:%d, method:%s, url:%s",
						fileListResp.StatusCode,
						method,
						url,
					)
					return false
				}
			}
			responseBody, finalErr = io.ReadAll(fileListResp.Body)
			return true
		},
	)
	return status, responseBody, finalErr
}
