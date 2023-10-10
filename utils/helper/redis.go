package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/usagifm/product-service/config"
	"github.com/usagifm/product-service/utils/logger"
)

func GetRedisDocumentName(redisKey interface{}, redisConfig *config.Redis) string {

	var redisName string = redisConfig.RootDocument + redisKey.(string)

	return redisName
}

func SetObjectToRedis(ctx context.Context, redisClient *redis.Client, object interface{}, objectName string, rediPath string, ttl time.Duration) {

	marshaledObject, errMarshal := json.Marshal(&object)
	if errMarshal != nil {
		logger.GetLogger(ctx).Errorf("Failed to marshal "+objectName+" (for redis), data "+objectName+" : %v | error : %d ", &object, errMarshal)
	}

	errRedis := redisClient.Set(ctx, rediPath, marshaledObject, ttl).Err()
	if errRedis != nil {
		logger.GetLogger(ctx).Errorf("Failed to set redis "+objectName+" (for redis), data "+objectName+" : %s | error : %d ", &object, errRedis.Error())
	}

}

func DeleteAllRedisDocumentByPath(ctx context.Context, redisClient *redis.Client, redisConfig *config.Redis, path string) {

	iter := redisClient.Scan(ctx, 0, redisConfig.RootDocument+path+"*", 0).Iterator()
	var foundedRecordCount = 0
	for iter.Next(ctx) {
		fmt.Println(iter.Val())
		logger.GetLogger(ctx).Infof("Deleted= %s\n", iter.Val())
		redisClient.Del(ctx, iter.Val())
		foundedRecordCount++
	}
	if err := iter.Err(); err != nil {
		logger.GetLogger(ctx).Error(err)
	}

}
