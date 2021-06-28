package gears

import (
	"context"
	"io/ioutil"

	"github.com/go-redis/redis/v8"
)

func LoadGears(rdb *redis.Client, path string) error {
	gearScript, err := ioutil.ReadFile(path)
	var ctx = context.Background()
	_, err = rdb.Do(ctx, "RG.PYEXECUTE", gearScript, "REQUIREMENTS", "elasticsearch").Result()
	return err
}
