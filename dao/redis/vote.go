package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInseconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postId, communityId int64) error {
	//事务
	pipeline := rdb.TxPipeline()

	//贴子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	//贴子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	//把帖子id加到社区的set
	ckey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityId)))
	pipeline.SAdd(ckey, postId)
	_, err := pipeline.Exec()

	return err
}

func VoteForPost(userId, postId string, value float64) error {
	//判断投票限制
	//发布时间已过一个星期，则不让投票了
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postId).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInseconds {
		return ErrVoteTimeExpire
	}

	//先查当前用户给当前贴子的投票记录
	ovalue := rdb.ZScore(getRedisKey(KeyPostVotedZSetPF+postId), userId).Val()

	//如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ovalue {
		return ErrVoteRepeated
	}

	var op float64
	if value > ovalue {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ovalue - value) //计算两次投票的差值

	// //更新帖子的分数
	// _, err := rdb.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postId).Result()
	// if ErrVoteTimeExpire != nil {
	// 	return err
	// }
	// //记录用户为该贴子投票的数据
	// if value == 0 {
	// 	_, err = rdb.ZRem(getRedisKey(KeyPostVotedZSetPF+postId), postId).Result()
	// } else {
	// 	_, err = rdb.ZAdd(getRedisKey(KeyPostVotedZSetPF+postId), redis.Z{
	// 		Score:  value, //赞成票还是反对票
	// 		Member: userId,
	// 	}).Result()
	// }

	//事务
	pipeline := rdb.TxPipeline()
	//更新帖子的分数
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postId)

	//记录用户为该贴子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postId), postId)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postId), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userId,
		})
	}
	_, err := pipeline.Exec()
	return err

}
