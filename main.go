package main

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
)

type RedisInfo struct {
	OrderId string `json:"order_id"`
	Token   string `json:"token"`
}

func main() {
	redisInit := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", "127.0.0.1", "6379"),
		PoolSize:   10,
		DB:         0,
		MaxRetries: 10,
	})

	sabaRedisClient := NewSabaRedisManager(redisInit)
	repository := NewChargeRedisRepository(sabaRedisClient)

	var token = "1234123asdawe1232"

	err := repository.Insert(token, &RedisInfo{
		OrderId: "1",
		Token:   token,
	})
	if err != nil {
		panic(err)
	}

	wg := new(sync.WaitGroup)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go tryToGetInfo(repository, token, wg)
	}

	wg.Wait()

}

func tryToGetInfo(repo ChargeRedisRepository, token string, wg *sync.WaitGroup) {
	defer wg.Done()
	v := &RedisInfo{}
	err := repo.FetchAndRemoveWithLock(token, v)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("********** %+v \n", v)
	}
}

func (s *SabaRedisManager) fetchAndRemoveWithLock(key string, value interface{}, lockDurationAsSecond int) error {
	lock := NewRedisLock(s.client, "Lock"+key)
	lock.SetExpire(lockDurationAsSecond)
	acquire, err := lock.Acquire()
	if err != nil {
		return err
	}

	fmt.Println("acquire", acquire)

	if acquire {
		fmt.Println("********************* acquire", acquire)
		err = s.fetch(key, value)
		if err != nil {
			go lock.Release()
			return err
		}

		err = s.remove(key)
		if err != nil {
			go lock.Release()
			return err
		}

		go lock.Release()
	} else {
		return errors.New("can not take lock")
	}

	return nil
}
