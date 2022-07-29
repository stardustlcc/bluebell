package logic

import (
	"go-jichu/dao/redis"
	"go-jichu/models"
	"strconv"

	"go.uber.org/zap"
)

func VoteForPost(userId int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userId", userId), zap.String("postId", p.PostId), zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userId)), p.PostId, float64(p.Direction))
}
