package rediskit

import (
	"github.com/gomodule/redigo/redis"
	"testing"
)

func show(t *testing.T, format string, data interface{}, err error) {
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(format, data)
}

func TestConn(t *testing.T) {
	p, err := NewRedis(RedisConfig{Db: 3})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Status:%+v", p.Stats())

	_conn := p.Get()
	defer _conn.Close()

	var _str string
	var _bool bool

	_str, err = redis.String(_conn.Do("set", "Key001", "你好，中文"))
	show(t, "set:%+v", _str, err)

	_bool, err = redis.Bool(_conn.Do("EXISTS", "Key001"))
	show(t, "exists:%+v", _bool, err)

	_str, err = redis.String(_conn.Do("GET", "Key001"))
	show(t, "get1:%+v", _str, err)

	_str, err = redis.String(_conn.Do("GET", "Key002"))
	if err != redis.ErrNil {
		show(t, "get2:%+v", _str, err)
	}

	_str, err = redis.String(_conn.Do("INFO"))
	show(t, "info:%+v", _str, err)

	t.Logf("The End Status:%+v", p.Stats())
	_conn.Close()
	_ = p.Close()
	t.Logf("Status:%+v", p.Stats())

}
