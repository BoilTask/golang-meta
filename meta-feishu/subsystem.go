package metafeishu

import (
	"context"
	"encoding/json"
	"log/slog"
	"meta/engine"
	metaerror "meta/meta-error"
	"meta/meta-feishu/variable"
	metaflag "meta/meta-flag"
	metalog "meta/meta-log"
	"meta/retry"
	"meta/subsystem"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type Subsystem struct {
	subsystem.Subsystem

	GetFeishuConfigs func() map[string]AppConfig

	feishuClients map[string]*lark.Client
}

func GetSubsystem() *Subsystem {
	if thisSubsystem := engine.GetSubsystem[*Subsystem](); thisSubsystem != nil {
		return thisSubsystem.(*Subsystem)
	}
	return nil
}

func (s *Subsystem) GetName() string {
	return "Feishu"
}

func (s *Subsystem) Start() error {
	configs := s.GetFeishuConfigs()
	if len(configs) == 0 {
		return metaerror.New("feishu configs is empty")
	}
	s.feishuClients = make(map[string]*lark.Client)
	for appKey, config := range configs {
		s.feishuClients[appKey] = startFeishuClient(&config)
	}
	return nil
}

func IsLogLevelDebug() bool {
	if IsSwitchDebugToInfo() {
		return true
	}
	return metalog.IsLogDebug()
}

func IsLogReqAtDebug() bool {
	//return metaflag.IsDebug()
	return true
}

func startFeishuClient(config *AppConfig) *lark.Client {
	if metaflag.IsDebug() {
		slog.Info("startFeishuClient", "config", config)
	}

	logLevel := larkcore.LogLevelInfo
	if IsLogLevelDebug() {
		logLevel = larkcore.LogLevelDebug
	}

	cli := lark.NewClient(
		config.AppId, config.AppSecret,
		lark.WithLogLevel(logLevel),
		lark.WithLogReqAtDebug(IsLogReqAtDebug()),
		lark.WithLogger(NewLogger()),
	)
	return cli
}

func (s *Subsystem) GetFeishuClient(appKey string) *lark.Client {
	if client, ok := s.feishuClients[appKey]; ok {
		return client
	}
	return nil
}

// IsUserInChat 用户是否在群里
// appKey: 应用类型
// chatId: 群ID
// openId: 用户openId
// 返回值:
// bool: 是否在群里
// bool: 机器人是否有权限
// error: 错误
func (s *Subsystem) IsUserInChat(
	ctx context.Context,
	appKey string,
	chatId string,
	openId string,
) (bool, bool, error) {
	pageToken := ""
	for {
		req := larkim.NewGetChatMembersReqBuilder().
			ChatId(chatId).
			PageToken(pageToken).
			Build()
		resp, err := s.GetFeishuClient(appKey).
			Im.
			ChatMembers.Get(ctx, req)
		if err != nil {
			return false, false, err
		}
		if !resp.Success() {
			return false, false, nil
		}
		for _, member := range resp.Data.Items {
			if *member.MemberId == openId {
				return true, true, nil
			}
		}
		if !*resp.Data.HasMore {
			break
		}
		pageToken = *resp.Data.PageToken
	}
	return false, true, nil
}

func (s *Subsystem) GetCreateMessageReqLog(req *larkim.CreateMessageReq) string {
	if req == nil {
		return "nil"
	}
	if req.Body == nil {
		return "nil body"
	}
	if req.Body.ReceiveId == nil {
		return "nil body.receive_id"
	}
	return *req.Body.ReceiveId
}

func (s *Subsystem) SendMessageByReq(
	ctx context.Context,
	appKey string,
	req *larkim.CreateMessageReq,
) (*string, error) {
	if req == nil {
		return nil, metaerror.New("request body is nil")
	}
	client := s.GetFeishuClient(appKey)
	if client == nil {
		return nil, metaerror.New("feishu client is nil, appKey:%s", appKey)
	}
	var err error
	var resp *larkim.CreateMessageResp
	retryErr := retry.TryRetryDynamicSleep(
		"Feishu Message Send", 6, func(i int) *time.Duration {
			resp, err = client.Im.Message.Create(ctx, req)
			if err != nil {
				err = metaerror.Wrap(err, "failed to send message, appKey:%s", appKey)
				return nil
			}
			if !resp.Success() {
				err = metaerror.WrapFeishu(
					resp,
					"send message failed, appKey:%s, req:%s",
					appKey,
					s.GetCreateMessageReqLog(req),
				)
				return GetFeishuErrorCodeRetrySleep(resp.Code)
			}
			return nil
		},
	)
	err = metaerror.Join(err, retryErr)
	if err != nil {
		return nil, err
	}
	return resp.Data.MessageId, nil
}

func (s *Subsystem) SendMessageTextToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	message string,
) (*string, error) {
	if openId == "" {
		return nil, metaerror.New("openId is empty")
	}
	content := GetFeishuMessageText(message)
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(
			larkim.NewCreateMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeText).
				ReceiveId(openId).
				Content(content).
				Build(),
		).
		Build()
	return s.SendMessageByReq(ctx, appKey, req)
}

func (s *Subsystem) SendMessageTextToChat(
	ctx context.Context,
	appKey string,
	chatId string,
	message string,
) (*string, error) {
	content := GetFeishuMessageText(message)
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(
			larkim.NewCreateMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeText).
				ReceiveId(chatId).
				Content(content).
				Build(),
		).
		Build()
	return s.SendMessageByReq(ctx, appKey, req)
}

func (s *Subsystem) SendMessageTextToEvent(
	ctx context.Context,
	appKey string,
	event *larkim.P2MessageReceiveV1,
	message string,
) (*string, error) {
	chatId := event.Event.Message.ChatId
	if IsMessageInGroup(event) {
		message = GetFeishuMessageTextByFeishuOpenId(*event.Event.Sender.SenderId.OpenId) + " " + message
	}
	return s.SendMessageTextToChat(ctx, appKey, *chatId, message)
}

func (s *Subsystem) SendMessageCardContentToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	cardContent string,
) (*string, error) {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(
			larkim.NewCreateMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeInteractive).
				ReceiveId(openId).
				Content(cardContent).
				Build(),
		).
		Build()
	return s.SendMessageByReq(ctx, appKey, req)
}

func (s *Subsystem) SendMessageCardContentToChat(
	ctx context.Context,
	appKey string,
	chatId string,
	cardContent string,
) (*string, error) {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(
			larkim.NewCreateMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeInteractive).
				ReceiveId(chatId).
				Content(cardContent).
				Build(),
		).
		Build()
	return s.SendMessageByReq(ctx, appKey, req)
}

func (s *Subsystem) SendMessageCardToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	card *larkcard.MessageCard,
) (*string, error) {
	content, err := card.String()
	if err != nil {
		return nil, metaerror.Wrap(err, "failed to get card content")
	}
	slog.Info("SendMessageCardToOpenId", "content", content)
	return s.SendMessageCardContentToOpenId(ctx, appKey, openId, content)
}

func (s *Subsystem) SendMessageCardToChat(
	ctx context.Context,
	appKey string,
	chatId string,
	card *larkcard.MessageCard,
) (*string, error) {
	content, err := card.String()
	if err != nil {
		return nil, metaerror.Wrap(err, "failed to get card content")
	}
	return s.SendMessageCardContentToChat(ctx, appKey, chatId, content)
}

func (s *Subsystem) SendMessageCardTemplateToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	templateId string,
	templateVariable interface{},
) (*string, error) {
	content := GetFeishuCardTemplate(templateId, templateVariable)
	return s.SendMessageCardContentToOpenId(ctx, appKey, openId, content)
}

func (s *Subsystem) SendMessageCardTemplateToChat(
	ctx context.Context,
	appKey string,
	chatId string,
	templateId string,
	templateVariable interface{},
) (*string, error) {
	content := GetFeishuCardTemplate(templateId, templateVariable)
	return s.SendMessageCardContentToChat(ctx, appKey, chatId, content)
}

func (s *Subsystem) SendBatchMessageCardContent(
	ctx context.Context,
	appKey string,
	openIds []string,
	departmentIds []string,
	cardContent string,
) (*string, error) {
	client := s.GetFeishuClient(appKey)
	var cardData interface{}
	err := json.Unmarshal([]byte(cardContent), &cardData)
	if err != nil {
		return nil, metaerror.Wrap(err, "failed to unmarshal card content")
	}
	body := struct {
		OpenIds       []string    `json:"open_ids"`
		DepartmentIds []string    `json:"department_ids"`
		MsgType       string      `json:"msg_type"`
		Card          interface{} `json:"card"`
	}{
		MsgType:       larkim.MsgTypeInteractive,
		Card:          cardData,
		OpenIds:       openIds,
		DepartmentIds: departmentIds,
	}
	url := "https://open.feishu.cn/open-apis/message/v4/batch_send"
	resp, err := client.Post(ctx, url, body, larkcore.AccessTokenTypeTenant)
	if err != nil {
		return nil, err
	}
	var respData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			MessageId string `json:"message_id"`
		} `json:"data"`
	}
	err = json.Unmarshal(resp.RawBody, &respData)
	if err != nil {
		return nil, err
	}
	if respData.Code != 0 {
		return nil, metaerror.New("send batch message failed, code: %d, msg: %s", respData.Code, respData.Msg)
	}
	return &respData.Data.MessageId, nil
}

func (s *Subsystem) SendMessageSuccessToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	templateVariable *variable.CommonContent,
) (*string, error) {
	content := GetFeishuCardTemplate("ctp_AAmaVxKT57DB", templateVariable)
	return s.SendMessageCardContentToOpenId(ctx, appKey, openId, content)
}

func (s *Subsystem) SendMessageSuccessToChat(
	ctx context.Context,
	appKey string,
	chatId string,
	templateVariable *variable.CommonContent,
) (*string, error) {
	content := GetFeishuCardTemplate("ctp_AAmaVxKT57DB", templateVariable)
	return s.SendMessageCardContentToChat(ctx, appKey, chatId, content)
}

func (s *Subsystem) SendMessageWarningToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	templateVariable *variable.CommonContent,
) (*string, error) {
	content := GetFeishuCardTemplate("ctp_AAruDiBOyZpO", templateVariable)
	return s.SendMessageCardContentToOpenId(ctx, appKey, openId, content)
}

func (s *Subsystem) SendMessageErrorToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	templateVariable *variable.CommonContent,
) (*string, error) {
	content := GetFeishuCardTemplate("ctp_AArYnsnCRkFu", templateVariable)
	return s.SendMessageCardContentToOpenId(ctx, appKey, openId, content)
}

func (s *Subsystem) ReplyMessageByReq(
	ctx context.Context,
	appKey string,
	req *larkim.ReplyMessageReq,
) (*string, error) {
	client := s.GetFeishuClient(appKey)
	if client == nil {
		return nil, metaerror.New("feishu client is nil")
	}
	resp, err := client.Im.Message.Reply(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.Success() {
		return nil, metaerror.WrapFeishu(resp, "send message failed")
	}
	return resp.Data.MessageId, nil
}

func (s *Subsystem) ReplyMessageText(
	ctx context.Context,
	appKey string,
	messageId string,
	message string,
) (*string, error) {
	content := GetFeishuMessageText(message)
	req := larkim.NewReplyMessageReqBuilder().
		MessageId(messageId).
		Body(
			larkim.NewReplyMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeText).
				Content(content).
				Build(),
		).
		Build()
	return s.ReplyMessageByReq(ctx, appKey, req)
}

func (s *Subsystem) ReplyMessageTextToEvent(
	ctx context.Context,
	appKey string,
	event *larkim.P2MessageReceiveV1,
	message string,
) (*string, error) {
	if IsMessageInGroup(event) {
		message = GetFeishuMessageTextByFeishuOpenId(*event.Event.Sender.SenderId.OpenId) + " " + message
	}
	return s.ReplyMessageText(ctx, appKey, *event.Event.Message.MessageId, message)
}

func (s *Subsystem) ReplyMessageCardContent(
	ctx context.Context,
	appKey string,
	messageId string,
	cardContent string,
) (*string, error) {
	req := larkim.NewReplyMessageReqBuilder().
		MessageId(messageId).
		Body(
			larkim.NewReplyMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeInteractive).
				Content(cardContent).
				Build(),
		).
		Build()
	return s.ReplyMessageByReq(ctx, appKey, req)
}

func (s *Subsystem) ReplyMessageCard(
	ctx context.Context,
	appKey string,
	messageId string,
	card *larkcard.MessageCard,
) (*string, error) {
	content, err := card.String()
	if err != nil {
		return nil, metaerror.Wrap(err, "failed to get card content")
	}
	return s.ReplyMessageCardContent(ctx, appKey, messageId, content)
}

func (s *Subsystem) UpdateMessageCardTemplate(
	ctx context.Context,
	appKey string,
	messageId string,
	templateId string,
	templateVariable interface{},
) error {
	client := s.GetFeishuClient(appKey)
	if client == nil {
		return metaerror.New("feishu client is nil")
	}
	content := GetFeishuCardTemplate(templateId, templateVariable)
	body := larkim.NewPatchMessageReqBodyBuilder().Content(content).Build()
	req := larkim.NewPatchMessageReqBuilder().
		MessageId(messageId).
		Body(body).
		Build()

	var err error
	var resp *larkim.PatchMessageResp
	retryErr := retry.TryRetryDynamicSleep(
		"Feishu Message Patch", 6, func(i int) *time.Duration {
			resp, err = client.Im.Message.Patch(ctx, req)
			if err != nil {
				err = metaerror.Wrap(err, "failed to update message, messageId:%s templateId:%s", messageId, templateId)
				sleep := time.Second * 30
				return &sleep
			}
			if !resp.Success() {
				return GetFeishuErrorCodeRetrySleep(resp.Code)
			}
			return nil
		},
	)
	err = metaerror.Join(err, retryErr)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return metaerror.WrapFeishu(resp, "update message failed")
	}
	return nil
}

func (s *Subsystem) SendMessageChatShareToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	chatId string,
) (*string, error) {
	messageText := larkim.MessageShareChat{ChatId: chatId}
	content, err := messageText.String()
	if err != nil {
		return nil, err
	}
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(
			larkim.NewCreateMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeShareChat).
				ReceiveId(openId).
				Content(content).
				Build(),
		).
		Build()
	return s.SendMessageByReq(ctx, appKey, req)
}

func (s *Subsystem) MergeForwardMessagesToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	messageIds []string,
) (*string, error) {
	client := s.GetFeishuClient(appKey)
	if client == nil {
		return nil, metaerror.New("feishu client is nil")
	}
	body := larkim.NewMergeForwardMessageReqBodyBuilder().
		ReceiveId(openId).
		MessageIdList(messageIds).
		Build()
	req := larkim.NewMergeForwardMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(body).
		Build()
	resp, err := client.Im.Message.MergeForward(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.Success() {
		return nil, metaerror.WrapFeishu(resp, "merge forward message failed")
	}
	return resp.Data.Message.MessageId, nil
}

func (s *Subsystem) MergeForwardMessageToOpenId(
	ctx context.Context,
	appKey string,
	openId string,
	messageId string,
) (*string, error) {
	return s.MergeForwardMessagesToOpenId(ctx, appKey, openId, []string{messageId})
}

func (s *Subsystem) InviteUserToChat(
	ctx context.Context,
	appKey string,
	chatId string,
	openId string,
) error {
	return s.InviteUsersToChat(ctx, appKey, chatId, []string{openId})
}

func (s *Subsystem) InviteUsersToChat(
	ctx context.Context,
	appKey string,
	chatId string,
	openIds []string,
) error {
	client := s.GetFeishuClient(appKey)
	if client == nil {
		return metaerror.New("feishu client is nil")
	}
	body := larkim.NewCreateChatMembersReqBodyBuilder().IdList(openIds).Build()
	req := larkim.NewCreateChatMembersReqBuilder().ChatId(chatId).Body(body).Build()
	resp, err := client.Im.ChatMembers.Create(ctx, req)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return metaerror.WrapFeishu(resp, "merge forward message failed")
	}
	return nil
}
