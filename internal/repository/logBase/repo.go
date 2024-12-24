package logBase

type RedisRepository interface {
	Add(key string) bool
	Free() bool
}
