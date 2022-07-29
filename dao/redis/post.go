package redis

import (
	"go-jichu/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIdsFormKey(key string, pageNum, pageLimit int64) ([]string, error) {
	//确定查询的索引起始点
	start := (pageNum - 1) * pageLimit
	end := start + pageLimit - 1
	//zrevarnge 按分数从大到小的顺序查询值指定数量
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIdsInSort(p *models.ParamPostList) ([]string, error) {
	//从redis中获取id
	//根据用户请求中携带的order参数确定cha'x
	key := getRedisKey(KeyPostTimeZSet)
	if p.Sort == models.SortScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIdsFormKey(key, p.PageNum, p.PageLimit)
}

//根据ids查询每篇贴子的投票数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVotedZSetPF + id)
	// 	//查询key中分数是1（赞成）的元素数量 -> 统计每篇贴子的赞成票数量
	// 	v := rdb.ZCount(key, "1", "1").Val()
	// 	data = append(data, v)
	// }

	//使用pipeline 一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

//按照社区查找id
func GetCommunityPostIdsInSort(p *models.ParamPostList) ([]string, error) {

	sortKey := getRedisKey(KeyPostTimeZSet)
	if p.Sort == models.SortScore {
		sortKey = getRedisKey(KeyPostScoreZSet)
	}

	//使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	//针对新的zset 按之前的逻辑取数据

	//社区的key
	ckey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityId)))

	//利用缓存key 减少zinterstore 执行的次数
	key := sortKey + strconv.Itoa(int(p.CommunityId))

	if rdb.Exists(key).Val() < 1 {
		//不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, ckey, sortKey) //zinterstore 计算
		pipeline.Expire(key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在，则直接根据key查询ids
	return getIdsFormKey(key, p.PageNum, p.PageLimit)

}
