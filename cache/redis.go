package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/usercoredev/usercore/utils/cipher"
	"os"
	"time"
)

var (
	RDB                        *redis.Client
	UserCacheExpiration        = time.Hour * 48
	UserProfileCacheExpiration = time.Hour * 48
	UserCachePrefix            = "user:"
	UserProfileCachePrefix     = "user:profile:"
	UserCacheKey               = UserCachePrefix + "%s"
	UserProfileCacheKey        = UserProfileCachePrefix + "%s"

	CollectionCacheExpiration = time.Hour * 48
	CollectionCachePrefix     = "collection:"
	CollectionCacheKey        = CollectionCachePrefix + "%d"

	WordCacheExpiration = time.Hour * 48
	WordCachePrefix     = "word:"
	WordCacheKey        = WordCachePrefix + "%d"

	WordListCacheExpiration = time.Hour * 48
	WordListCachePrefix     = "collection:"
	WordListCacheKey        = WordListCachePrefix + "%d" + ":words"

	BlogCacheExpiration = time.Hour * 48
	BlogCachePrefix     = "blog:"
	BlogCacheKey        = BlogCachePrefix + "%s"
)

type CacheNotEnabled string

func (e CacheNotEnabled) Error() string { return string(e) }

func (CacheNotEnabled) CacheNotEnabledError() {}

const CacheNotEnabledErr = CacheNotEnabled("cache not enabled")

func Redis() error {
	if os.Getenv("CACHE_ENABLED") != "true" {
		return nil
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		return fmt.Errorf("REDIS_PORT not set")
	}
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		return fmt.Errorf("REDIS_HOST not set")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		return fmt.Errorf("REDIS_PASSWORD not set")
	}

	redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

func Set(key string, value interface{}, expire time.Duration) error {
	if os.Getenv("CACHE_ENABLED") != "true" {
		return nil
	}
	encryptionKey := os.Getenv("REDIS_ENCRYPTION_KEY")
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}
	encryptedValue, err := cipher.EncryptWithKey(jsonVal, encryptionKey)
	if err != nil {
		return err
	}
	ctx := context.Background()
	result := RDB.Set(ctx, key, encryptedValue, expire)
	return result.Err()
}

func SetList(key string, value interface{}, expire time.Duration) error {
	if os.Getenv("CACHE_ENABLED") != "true" {
		return CacheNotEnabledErr
	}
	encryptionKey := os.Getenv("REDIS_ENCRYPTION_KEY")
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}
	encryptedValue, err := cipher.EncryptWithKey(jsonVal, encryptionKey)
	if err != nil {
		return err
	}
	ctx := context.Background()
	result := RDB.LPush(ctx, key, encryptedValue)
	return result.Err()
}

func Get(key string, value interface{}) error {
	if os.Getenv("CACHE_ENABLED") != "true" {
		return CacheNotEnabledErr
	}

	encryptionKey := os.Getenv("REDIS_ENCRYPTION_KEY")
	ctx := context.Background()
	result := RDB.Get(ctx, key)

	val, err := result.Result()
	if err != nil {
		return err
	}

	decryptedValue, err := cipher.DecryptWithKey(val, encryptionKey)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decryptedValue, value)
	if err != nil {
		return err
	}

	return nil
}
