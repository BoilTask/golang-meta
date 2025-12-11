package metafeishu

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	metaerror "meta/meta-error"
	"net/http"
	"time"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	pkgerrors "github.com/pkg/errors"
)

func GetFeishuMessageTextByFeishuOpenId(feishuOpenId string) string {
	return "<at user_id=\"" + feishuOpenId + "\"></at>"
}

func GetFeishuMessageByFeishuOpenId(feishuOpenId string) string {
	return "<at id=\"" + feishuOpenId + "\"></at>"
}

func GetFeishuMessageText(message string) string {
	data := map[string]interface{}{
		"text": message,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData)
}

func GetFeishuCardTemplate(template string, templateVariables interface{}) string {
	data := map[string]interface{}{
		"type": "template",
		"data": map[string]interface{}{
			"template_id":       template,
			"template_variable": templateVariables,
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData)
}

func IsMessageInGroup(event *larkim.P2MessageReceiveV1) bool {
	return event.Event.Message.ChatType != nil && *event.Event.Message.ChatType == "group"
}

func IsFeishuErrorNoAvailabilityToUser(err error) bool {
	if err == nil {
		return false
	}
	type CodeStruct struct {
		Code int `json:"code"`
	}
	err = pkgerrors.Cause(err)
	var code CodeStruct
	if err = json.Unmarshal([]byte(err.Error()), &code); err != nil {
		return false
	}
	return code.Code == 230013
}

func GetEmployeeTypeName(employeeType int) string {
	switch employeeType {
	case 1:
		return "正式员工"
	case 2:
		return "实习生"
	case 3:
		return "外包"
	case 4:
		return "劳务"
	case 5:
		return "顾问"
	default:
		return "未知"
	}
}

func GetFeishuErrorCodeRetrySleep(code int) *time.Duration {
	switch code {
	case int(ErrorCodeCheckAppTenantFail):
		sleep := time.Second * 1
		return &sleep
	case int(ErrorCodeServerInternalError):
		sleep := time.Second * 1
		return &sleep
	case int(ErrorCodeBitableDataNotReady):
		sleep := time.Second * 1
		return &sleep
	case int(ErrorCodeBitableLockNotObtainedError):
		sleep := time.Second * 1
		return &sleep
	case int(ErrorCodeBitableTooManyRequest):
		sleep := time.Second * 15
		return &sleep
	case int(ErrorCodeFrequencyLimit):
		sleep := time.Second * 30
		return &sleep
	default:
		return nil
	}
}

func SendMessageTextToCustomRobot(
	robotUrl string,
	message string,
) error {
	if robotUrl == "" {
		return metaerror.New("robotUrl is empty")
	}
	//发起post请求
	body := struct {
		MsgType string `json:"msg_type"`
		Content struct {
			Text string `json:"text"`
		} `json:"content"`
	}{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: message,
		},
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return metaerror.Wrap(err, "failed to marshal request body")
	}
	resp, err := http.Post(robotUrl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return metaerror.New("send message failed, status code: %d", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return metaerror.Wrap(err, "failed to read response body")
	}
	var respData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return metaerror.Wrap(err, "failed to unmarshal response body")
	}
	if respData.Code != 0 {
		return metaerror.New("send message failed, code: %d, msg: %s", respData.Code, respData.Msg)
	}
	return nil
}
