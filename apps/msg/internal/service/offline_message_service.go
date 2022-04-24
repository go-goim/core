package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	redisv8 "github.com/go-redis/redis/v8"
	"github.com/yusank/goim/pkg/log"

	messagev1 "github.com/yusank/goim/api/message/v1"
	"github.com/yusank/goim/apps/msg/internal/app"
	"github.com/yusank/goim/apps/msg/internal/data"
)

type OfflineMessageService struct {
	messagev1.UnimplementedOfflineMessageServer
}

func (o *OfflineMessageService) QueryOfflineMessage(ctx context.Context, req *messagev1.QueryOfflineMessageReq) (*messagev1.QueryOfflineMessageResp, error) {
	var rsp = &messagev1.QueryOfflineMessageResp{
		Status:   -1,
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
	}

	log.Info("req=", req.String())
	msgID, err := primitive.UnmarshalMsgID([]byte(req.GetLastMsgSeq()))
	if err != nil {
		log.Info("unmarshal msg id err=", err)
		rsp.Reason = err.Error()
		return rsp, nil
	}

	log.Info("unmarshal msg", "host", msgID.Addr, "port", msgID.Port, "offset", msgID.Offset)

	cnt, err := app.GetApplication().Redis.ZCount(ctx,
		data.GetUserOfflineQueueKey(req.GetUserId()),
		// offset add 1 to skip the message user last online msg
		strconv.FormatInt(msgID.Offset+1, 10),
		"+inf").Result()
	if err != nil {
		rsp.Reason = err.Error()
		return rsp, nil
	}

	rsp.Status = 0
	rsp.TotalCount = int32(cnt)
	if req.GetOnlyCount() {
		return rsp, nil
	}

	results, err := app.GetApplication().Redis.ZRangeByScoreWithScores(ctx,
		data.GetUserOfflineQueueKey(req.GetUserId()), &redisv8.ZRangeBy{
			// offset add 1 to skip the message user last online msg
			Min:    strconv.FormatInt(msgID.Offset+1, 10),
			Max:    "+inf",
			Offset: int64((req.GetPage() - 1) * req.GetPageSize()),
			Count:  int64(req.GetPageSize()),
		}).Result()
	if err != nil {
		rsp.Reason = err.Error()
		return rsp, nil
	}

	rsp.Messages = make([]*messagev1.BriefMessage, len(results))
	for i, result := range results {
		str := result.Member.(string)
		msg := new(messagev1.BriefMessage)
		if err = json.Unmarshal([]byte(str), msg); err != nil {
			rsp.Reason = err.Error()
			return rsp, nil
		}

		rsp.Messages[i] = msg
	}

	rsp.Status = 0
	return rsp, nil
}
