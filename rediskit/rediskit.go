package rediskit

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisConfig struct {
	Network  string // tcp
	Address  string //'127.0.0.1:6379'
	Password string // =''
	Db       int    // ='7'
}

type Conn = redis.Conn

var ErrNil = redis.ErrNil
var String = redis.String

func NewRedis(cfg RedisConfig) (*redis.Pool, error) {
	_network, _address := cfg.Network, cfg.Address
	_passwd, _db := cfg.Password, cfg.Db
	if _network == "" {
		_network = "tcp"
	}
	if _address == "" {
		_address = ":6379"
	}
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(_network, _address)
			if err != nil {
				return nil, err
			}
			if _passwd != "" {
				if _, err := c.Do("AUTH", _passwd); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			if _db > 0 {
				if _, err := c.Do("SELECT", _db); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     8,
		IdleTimeout: 60 * time.Second,
	}, nil
}
