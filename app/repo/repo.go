package repo

type RedisRepository interface {
	Add(key string) bool
	Free() bool
}
