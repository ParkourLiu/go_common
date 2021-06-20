package credis_test

import (
	"fmt"
	"go_common/credis"
	"testing"
	"time"
)

var redisClient *credis.RedisClient

func init() {
	redisClient = credis.NewRedisClient(&credis.RedisInfo{
		Addr:     "127.0.0.1:6379",
		Password: "",
		MaxIdle:  200,
	})
}
func TestRedisClient_Expire_AND_Exists(t *testing.T) {
	redisClient.Set("a", "aaa")
	err := redisClient.Expire("a", 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(redisClient.Exists("a"))
	time.Sleep(6 * time.Second)
	fmt.Println(redisClient.Exists("a"))
}

func TestRedisClient_Setex(t *testing.T) {
	redisClient.Setex("a", 3, "aaa")
	fmt.Println(redisClient.Exists("a"))
	time.Sleep(4 * time.Second)
	fmt.Println(redisClient.Exists("a"))
}

func TestRedisClient_Get(t *testing.T) {
	for i := 0; i < 400; i++ {
		time.Sleep(5 * time.Second)
		str, err := redisClient.Get("a")
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(i, str)
	}
}
func TestRedisClient_Set(t *testing.T) {
	err := redisClient.Set("a", "aaa")
	if err != nil {
		t.Fatal(err)
	}
}
func TestRedisClient_Incrby(t *testing.T) {
	for i := 0; i < 400; i++ {
		str, err := redisClient.Incrby("aa", 2)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(i, str)
	}
}
func TestRedisClient_Exists(t *testing.T) {
	str, err := redisClient.Exists("aa")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}
func TestRedisClient_Del(t *testing.T) {
	err := redisClient.Del("aa", "a", "aaa")
	if err != nil {
		t.Fatal(err)
	}
}
func TestRedisClient_Keys(t *testing.T) {
	strs, err := redisClient.Keys("list*")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(strs)
}
func TestRedisClient_Append(t *testing.T) {
	err := redisClient.Append("a", "a")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisClient_Setnx(t *testing.T) {
	fmt.Println(redisClient.Setnx("a", "aa"))
	fmt.Println(redisClient.Setnx("a", "aaa"))
}
func TestRedisClient_Hmset(t *testing.T) {
	a := map[string]string{"aaaaa": "aaaaa"}
	err := redisClient.Hmset("aa", a)
	if err != nil {
		t.Fatal(err)
	}
}
func TestRedisClient_Hmget(t *testing.T) {
	values, err := redisClient.Hmget("aa", "aa", "aaa", "aaaa")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}
func TestRedisClient_Hset(t *testing.T) {
	err := redisClient.Hset("aa", "a", "a")
	if err != nil {
		t.Fatal(err)
	}
}
func TestRedisClient_Hget(t *testing.T) {
	values, err := redisClient.Hget("hashtest", "aaa")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}
func TestRedisClient_HgetAllMap(t *testing.T) {
	values, err := redisClient.HgetAllMap("aa")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}
func TestRedisClient_Hincrby(t *testing.T) {
	values, err := redisClient.Hincrby("aa", "1", 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}
func TestRedisClient_Lpush(t *testing.T) {
	err := redisClient.Lpush("bb", "a")
	if err != nil {
		t.Fatal(err)
	}
}
func TestRedisClient_Rpush(t *testing.T) {
	err := redisClient.Rpush("listtest", "a")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisClient_Lrange(t *testing.T) {
	values, err := redisClient.Lrange("listtest", 0, -1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}
func TestRedisClient_Llen(t *testing.T) {
	len, err := redisClient.Llen("bb")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len)
}
func TestRedisClient_Lrem(t *testing.T) {
	err := redisClient.Lrem("bb", 0, "a")
	if err != nil {
		t.Fatal(err)
	}
}
func TestRedisClient_Lpop(t *testing.T) {
	values, err := redisClient.Lpop("bb")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}

func TestRedisClient_Sadd(t *testing.T) {
	err := redisClient.Sadd("ss", "ss3 ss4 aaaa")
	if err != nil {
		t.Fatal(err)
	}
}
func TestRedisClient_Scard(t *testing.T) {
	values, err := redisClient.Scard("ss")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}
func TestRedisClient_Sismember(t *testing.T) {
	values, err := redisClient.Sismember("ss", "ss3 ss4 aaaa")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}
func TestRedisClient_Smembers(t *testing.T) {
	values, err := redisClient.Smembers("ss")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}

func TestRedisClient_ZrangeWithscores(t *testing.T) {
	values, err := redisClient.ZrangeWithscores("zsettest", 0, -1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(values)
}

func TestHdel(t *testing.T) {
	err := redisClient.Hdel("hashtest", "bbbb")
	if err != nil {
		t.Fatal(err)
	}
}

func TestHlen(t *testing.T) {
	len, err := redisClient.Hlen("IPa")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len)
}
