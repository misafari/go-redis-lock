package main

import (
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type SabaRedisManager struct {
	client      redis.Cmdable
	expiredTime time.Time
}

func NewSabaRedisManager(client redis.Cmdable) *SabaRedisManager {
	return &SabaRedisManager{
		client: client,
	}
}

func (s *SabaRedisManager) remove(key string) error {
	return s.client.Del(key).Err()
}

func (s *SabaRedisManager) insert(key string, value interface{}) error {
	val, _ := jsoniter.Marshal(value)
	return s.client.Set(key, string(val), time.Minute*15).Err()
}

func (s *SabaRedisManager) insertPermanent(key string, value interface{}) error {
	val, _ := jsoniter.Marshal(value)
	return s.client.Set(key, string(val), 0).Err()
}

func (s *SabaRedisManager) insertString(key string, value string) error {
	return s.client.Set(key, value, time.Minute*15).Err()
}

//func (s *SabaRedisManager) scanKeys(keyPat string) ([]string, error) {
//	keys, _, err := s.client.Scan(0, keyPat, 0).Result()
//	return keys, err
//}

func (s *SabaRedisManager) fetch(key string, value interface{}) error {
	val, err := s.client.Get(key).Result()
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal([]byte(val), value)
}

func (s *SabaRedisManager) fetchString(key string) (string, error) {
	return s.client.Get(key).Result()
}
