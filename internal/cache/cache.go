package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/usercoredev/usercore/internal/cipher"
	"net/url"
	"os"
	"strings"
	"time"
)

type notEnabled string

func (e notEnabled) Error() string { return string(e) }

const NotEnabled = notEnabled("cache_not_enabled")

type Settings struct {
	Enabled                    string
	Host                       string
	Port                       string
	Password                   string
	PasswordFile               string
	EncryptionKey              string
	UserCacheExpiration        string
	UserProfileCacheExpiration string
	UserCachePrefix            string
	UserProfileCachePrefix     string
}
type redisCache struct {
	redis                      *redis.Client
	encryptionKey              string
	UserPrefix                 string
	UserProfilePrefix          string
	UserCacheExpiration        time.Duration
	UserProfileCacheExpiration time.Duration
}

var Client *redisCache

func (s *Settings) SetupCache() error {
	if s.Enabled != "true" {
		return nil
	}
	if s.Host == "" {
		panic("cache host not set")
	} else if s.Port == "" {
		panic("cache port not set")
	} else if s.Password == "" && s.PasswordFile == "" {
		panic("CACHE_PASSWORD or CACHE_PASSWORD_FILE is required")
	}
	if s.PasswordFile != "" {
		bin, err := os.ReadFile(s.PasswordFile)
		if err != nil {
			return err
		}
		s.Password = string(bin)
	}
	s.Password = url.QueryEscape(strings.TrimSpace(s.Password))

	userCacheExpiration, err := time.ParseDuration(s.UserCacheExpiration)
	if err != nil {
		return err
	}
	userProfileCacheExpiration, err := time.ParseDuration(s.UserProfileCacheExpiration)
	if err != nil {
		return err
	}

	Client = &redisCache{
		encryptionKey:              s.EncryptionKey,
		UserPrefix:                 s.UserCachePrefix,
		UserProfilePrefix:          s.UserProfileCachePrefix,
		UserCacheExpiration:        userCacheExpiration,
		UserProfileCacheExpiration: userProfileCacheExpiration,
	}
	Client.redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", s.Host, s.Port),
		Password: s.Password,
		DB:       0,
	})
	_, err = Client.redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return nil
}

func Set(key string, value interface{}, duration time.Duration) error {
	if Client == nil {
		return NotEnabled
	}
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}
	encryptedValue, err := cipher.EncryptWithKey(jsonVal, Client.encryptionKey)
	if err != nil {
		return err
	}
	ctx := context.Background()
	result := Client.redis.Set(ctx, key, encryptedValue, duration)
	return result.Err()
}

func Get(key string, value interface{}) error {
	if Client == nil {
		return NotEnabled
	}
	ctx := context.Background()
	result := Client.redis.Get(ctx, key)

	val, err := result.Result()
	if err != nil {
		return err
	}

	decryptedValue, err := cipher.DecryptWithKey(val, Client.encryptionKey)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decryptedValue, value)
	if err != nil {
		return err
	}

	return nil
}
