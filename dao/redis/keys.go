package redis

const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "bluebell:post:time"  //zset 贴子以及发帖时间
	KeyPostScoreZSet   = "bluebell:post:score" //zset 贴子以及投票的分数
	KeyPostVotedZSetPF = "post:voted:"         //zset 记录及投票的类型
	KeyCommunitySetPF  = "community:"          //set 保存每一个分区下帖子的id
)

func getRedisKey(key string) string {
	return Prefix + key
}
