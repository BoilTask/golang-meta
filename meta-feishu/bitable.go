package metafeishu

import (
	"context"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	metaerror "meta/meta-error"
	"meta/retry"
	"time"
)

func SearchBitableRemoteRecord(
	ctx context.Context,
	appKey string,
	appToken string,
	tableId string,
	findReqFilter *larkbitable.FilterInfo,
) ([]*larkbitable.AppTableRecord, error) {
	feishuClient := GetSubsystem().GetFeishuClient(appKey)
	var records []*larkbitable.AppTableRecord
	var finalErr error
	err := retry.TryRetryDynamicSleep(
		"SearchBitableRemoteRecord", 5, func(i int) *time.Duration {
			finalErr = nil
			findReqBody := larkbitable.NewSearchAppTableRecordReqBodyBuilder().
				Filter(findReqFilter).
				Build()
			findReq := larkbitable.NewSearchAppTableRecordReqBuilder().
				AppToken(appToken).
				TableId(tableId).
				Body(findReqBody).
				Build()
			findResp, err := feishuClient.Bitable.AppTableRecord.
				Search(ctx, findReq)
			if err != nil {
				finalErr = metaerror.Wrap(
					err,
					"search app table recordId error, appToken:%s tableId:%s",
					appToken,
					tableId,
				)
				return nil
			}
			if !findResp.Success() {
				finalErr = metaerror.WrapFeishu(
					findResp,
					"search app table recordId not success, appToken:%s tableId:%s",
					appToken,
					tableId,
				)
				return GetFeishuErrorCodeRetrySleep(findResp.Code)
			}
			if findResp.Data == nil || findResp.Data.Items == nil || len(findResp.Data.Items) == 0 {
				return nil
			}
			records = findResp.Data.Items
			return nil
		},
	)
	finalErr = metaerror.Join(finalErr, err)
	return records, finalErr
}

func GetBitableRemoteRecord(
	ctx context.Context,
	appKey string,
	appToken string,
	tableId string,
	recordId string,
) (*larkbitable.AppTableRecord, error) {
	feishuClient := GetSubsystem().GetFeishuClient(appKey)
	var record *larkbitable.AppTableRecord
	var finalErr error
	err := retry.TryRetryDynamicSleep(
		"GetBitableRemoteRecord", 5, func(i int) *time.Duration {
			finalErr = nil
			body := larkbitable.NewBatchGetAppTableRecordReqBodyBuilder().
				RecordIds([]string{recordId}).
				WithSharedUrl(true).
				AutomaticFields(true).
				Build()
			req := larkbitable.NewBatchGetAppTableRecordReqBuilder().
				AppToken(appToken).
				TableId(tableId).
				Body(body).
				Build()
			resp, err := feishuClient.Bitable.AppTableRecord.BatchGet(ctx, req)
			if err != nil {
				finalErr = metaerror.Wrap(err, "get app table recordId error, recordId:%s", recordId)
				return nil
			}
			if !resp.Success() {
				finalErr = metaerror.WrapFeishu(resp, "get app table recordId not success, recordId:%s", recordId)
				return GetFeishuErrorCodeRetrySleep(resp.Code)
			}
			if len(resp.Data.Records) == 0 {
				return nil
			}
			record = resp.Data.Records[0]
			return nil
		},
	)
	finalErr = metaerror.Join(finalErr, err)
	return record, finalErr
}

func UpdateBitableRemoteRecord(
	ctx context.Context,
	appKey string,
	appToken string,
	tableId string,
	recordId string,
	record *larkbitable.AppTableRecord,
) error {
	feishuClient := GetSubsystem().GetFeishuClient(appKey)
	var finalErr error
	err := retry.TryRetryDynamicSleep(
		"UpdateBitableRemoteRecord", 5, func(i int) *time.Duration {
			finalErr = nil
			reqBitable := larkbitable.NewUpdateAppTableRecordReqBuilder().
				AppToken(appToken).
				TableId(tableId).
				RecordId(recordId).
				AppTableRecord(record).
				Build()
			resp, err := feishuClient.Bitable.AppTableRecord.Update(ctx, reqBitable)
			if err != nil {
				finalErr = metaerror.Wrap(err, "update app table record error, recordId:%s", recordId)
				return nil
			}
			if !resp.Success() {
				finalErr = metaerror.WrapFeishu(resp, "update app table record not success, recordId:%s", recordId)
				return GetFeishuErrorCodeRetrySleep(resp.Code)
			}
			return nil
		},
	)
	finalErr = metaerror.Join(finalErr, err)
	return finalErr
}

func CreateBitableRemoteRecord(
	ctx context.Context,
	appKey string,
	appToken string,
	tableId string,
	record *larkbitable.AppTableRecord,
) (*string, error) {
	feishuClient := GetSubsystem().GetFeishuClient(appKey)
	var finalErr error
	var recordId *string
	err := retry.TryRetryDynamicSleep(
		"CreateBitableRemoteRecord", 5, func(i int) *time.Duration {
			finalErr = nil
			reqBitable := larkbitable.NewCreateAppTableRecordReqBuilder().
				AppToken(appToken).
				TableId(tableId).
				AppTableRecord(record).
				Build()
			resp, err := feishuClient.Bitable.AppTableRecord.Create(ctx, reqBitable)
			if err != nil {
				finalErr = metaerror.Wrap(err, "create app table record error")
				return nil
			}
			if !resp.Success() {
				finalErr = metaerror.WrapFeishu(resp, "create app table record not success")
				return GetFeishuErrorCodeRetrySleep(resp.Code)
			}
			recordId = resp.Data.Record.RecordId
			return nil
		},
	)
	finalErr = metaerror.Join(finalErr, err)
	return recordId, finalErr
}
