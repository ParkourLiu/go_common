package credis

import (
	"github.com/gomodule/redigo/redis"
	"go_common/clogs"
)

var log *clogs.Log

func init() {
	log = clogs.NewLog("7")
}

type RedisClient struct {
	pool *redis.Pool
}
type RedisInfo struct {
	Addr      string
	Password  string
	MaxIdle   int //最大空闲数,最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
	MaxActive int //最大连接数，表示同时最多有N个连接 0为不限制连接
}

func NewRedisClient(redisInfo *RedisInfo) *RedisClient {
	p := &redis.Pool{
		MaxIdle:   redisInfo.MaxIdle,
		MaxActive: redisInfo.MaxActive,
		Dial: func() (redis.Conn, error) {
			if redisInfo.Password == "" {
				return redis.Dial("tcp", redisInfo.Addr)
			} else {
				return redis.Dial("tcp", redisInfo.Addr, redis.DialPassword(redisInfo.Password))
			}

		},
	}
	r := &RedisClient{pool: p}
	return r
}

//给key设置过期时间
func (c *RedisClient) Expire(key string, expireSecond int) error {
	return c.returnError("EXPIRE", key, expireSecond)
}

//判断key是否存在
func (c *RedisClient) Exists(key string) (bool, error) {
	return c.returnBoolError("EXISTS", key)
}

//删除key
func (c *RedisClient) Del(keys ...interface{}) error {
	return c.returnError("DEL", keys...)
}

//匹配key
func (c *RedisClient) Keys(key string) ([]string, error) {
	return c.returnStringsError("KEYS", key)
}

//string===============================================================begin
func (c *RedisClient) Get(key string) (string, error) {
	return c.returnStringError("GET", key)
}

func (c *RedisClient) Set(key, value string) error {
	return c.returnError("SET", key, value)
}

//key 中储存的数字加上指定的增量值，如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行增量操作。count可为正负数
func (c *RedisClient) Incrby(key string, count int) (int64, error) {
	return c.returnInt64Error("INCRBY", key, count)
}

//在原本值基础上加上一个value值
func (c *RedisClient) Append(key, value string) error {
	return c.returnError("APPEND", key, value)
}

//set时加上过期秒数
func (c *RedisClient) Setex(key string, expireSecond int, value string, ) error {
	return c.returnError("SETEX", key, expireSecond, value)
}

////如果没有此key则插入并返回1，存在则不插入返回0,一般用于分布式锁
func (c *RedisClient) Setnx(key, value string) (int64, error) {
	return c.returnInt64Error("SETNX", key, value)
}

//string===============================================================end

//* hash */============================================================begin
func (c *RedisClient) Hmset(key string, value interface{}) error {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return err
	}
	defer conn.Close()
	_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...)
	return err
}

func (c *RedisClient) Hmget(key string, subkey ...interface{}) (values []string, err error) {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return nil, err
	}
	defer conn.Close()
	return redis.Strings(conn.Do("HMGET", redis.Args{}.Add(key).AddFlat(subkey)...))
}
func (c *RedisClient) Hset(key string, subkey string, value string) error {
	return c.returnError("HSET", key, subkey, value)
}

func (c *RedisClient) Hget(key, subkey string) (values string, err error) {
	return c.returnStringError("HGET", key, subkey)
}
func (c *RedisClient) HgetAllMap(key string) (value map[string]string, err error) {
	return c.returnStringMapError("HGETALL", key)
}
func (c *RedisClient) HgetAllStruct(key string, value interface{}) (err error) {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return err
	}
	defer conn.Close()
	var v []interface{}
	v, err = redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		return
	}
	err = redis.ScanStruct(v, value)
	if err == redis.ErrNil {
		err = nil
	}
	return
}
func (c RedisClient) Hincrby(key, subkey string, count int) (int64, error) {
	return c.returnInt64Error("HINCRBY", key, subkey, count)
}

//* hash */============================================================end
//* list */============================================================begin
func (c *RedisClient) Lpush(key string, value string) error {
	return c.returnError("LPUSH", key, value)
}
func (c *RedisClient) Rpush(key string, value string) error {
	return c.returnError("RPUSH", key, value)
}
func (c *RedisClient) Lpop(key string) (string, error) {
	return c.returnStringError("LPOP", key)
}
func (c *RedisClient) Rpop(key string) (string, error) {
	return c.returnStringError("RPOP", key)
}
func (c *RedisClient) Lrange(key string, startIndex int64, endIndex int64) ([]string, error) {
	return c.returnStringsError("LRANGE", key, startIndex, endIndex)
}
func (c *RedisClient) Llen(key string) (int64, error) {
	return c.returnInt64Error("LLEN", key)
}

//count > 0 : 从表头开始向表尾搜索，移除与 VALUE 相等的元素，数量为 COUNT 。
//count < 0 : 从表尾开始向表头搜索，移除与 VALUE 相等的元素，数量为 COUNT 的绝对值。
//count = 0 : 移除表中所有与 VALUE 相等的值。
func (c *RedisClient) Lrem(key string, count int, value string) error {
	return c.returnError("LREM", key, count, value)
}

//* list */============================================================end

//* set */============================================================begin
func (c *RedisClient) Sadd(key string, value string) error {
	return c.returnError("SADD", key, value)
}

//返回成员数量
func (c *RedisClient) Scard(key string) (int64, error) {
	return c.returnInt64Error("SCARD", key)
}

//判断是否存在此元素
func (c *RedisClient) Sismember(key string, value string) (bool, error) {
	return c.returnBoolError("SISMEMBER", key, value)
}

//返回所有元素
func (c *RedisClient) Smembers(key string) ([]string, error) {
	return c.returnStringsError("SMEMBERS", key)
}

//* set */============================================================end

//* zset */============================================================begin

//* zset */============================================================end
//返回类型封装------------------------------------------------------------------------------------------------------------------------------------
func (c *RedisClient) returnError(commandName string, in ...interface{}) error {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return err
	}
	defer conn.Close()

	_, err := conn.Do(commandName, in...)
	return err
}
func (c *RedisClient) returnBoolError(commandName string, in ...interface{}) (bool, error) {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return false, err
	}
	defer conn.Close()

	return redis.Bool(conn.Do(commandName, in...))
}
func (c *RedisClient) returnStringError(commandName string, in ...interface{}) (string, error) {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return "", err
	}
	defer conn.Close()

	return redis.String(conn.Do(commandName, in...))
}
func (c *RedisClient) returnStringsError(commandName string, in ...interface{}) ([]string, error) {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return nil, err
	}
	defer conn.Close()

	return redis.Strings(conn.Do(commandName, in...))
}
func (c *RedisClient) returnInt64Error(commandName string, in ...interface{}) (int64, error) {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return 0, err
	}
	defer conn.Close()

	return redis.Int64(conn.Do(commandName, in...))
}

func (c *RedisClient) returnStringMapError(commandName string, in ...interface{}) (map[string]string, error) {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return nil, err
	}
	defer conn.Close()

	return redis.StringMap(conn.Do(commandName, in...))
}

func (c *RedisClient) returnArgsError(commandName string, key string, value ...interface{}) error {
	conn := c.pool.Get()
	if err := conn.Err(); err != nil {
		log.Error(err)
		return err
	}
	defer conn.Close()

	_, err := conn.Do(commandName, redis.Args{}.Add(key).AddFlat(value)...)
	return err
}
