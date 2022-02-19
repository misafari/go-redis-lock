package main

type ChargeRedisRepository interface {
	Insert(string, *RedisInfo) error
	Fetch(string, *RedisInfo) error
	FetchAndRemoveWithLock(string, *RedisInfo) error
	Remove(string) error
}

type chargeRedisRepository struct {
	sabaRedisManager *SabaRedisManager
}

func NewChargeRedisRepository(sabaRedisManager *SabaRedisManager) ChargeRedisRepository {
	return &chargeRedisRepository{
		sabaRedisManager: sabaRedisManager,
	}
}

func (s *chargeRedisRepository) Insert(key string, value *RedisInfo) error {
	return s.sabaRedisManager.insertPermanent(key, value)
}

func (s *chargeRedisRepository) FetchAndRemoveWithLock(key string, result *RedisInfo) error {
	return s.sabaRedisManager.fetchAndRemoveWithLock(key, result, 10)
}

func (s *chargeRedisRepository) Fetch(key string, result *RedisInfo) error {
	return s.sabaRedisManager.fetch(key, result)
}

func (s *chargeRedisRepository) Remove(key string) error {
	return s.sabaRedisManager.remove(key)
}
