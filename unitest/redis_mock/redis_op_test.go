package redis_mock

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func TestDoSomethingWithRedis(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.Set("lvsoso", "lvsoso.com")
	s.SAdd(KeyValidWebsite, "lvsoso")

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	ok := DoSomethingWithRedis(rdb, "lvsoso")
	if !ok {
		t.Fatal()
	}

	if got, err := s.Get("blog"); err != nil || got != "https://lvsoso.com" {
		t.Fatalf("'blog' has the wrong value")
	}
	s.CheckGet(t, "blog", "https://lvsoso.com")

	s.FastForward(5 * time.Second) // 快进5秒
	if s.Exists("blog") {
		t.Fatal("'blog' should not have existed anymore")
	}
}
