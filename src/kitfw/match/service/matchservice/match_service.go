package matchservice

import (
	"fmt"
	"kitfw/commom/store"
)

type MatchStrategy interface {
	start()
}

const (
	PREFIX_TIMEOUT_STRING = "match:timeout:player_"
)

const (
	MATCH_TYPE_ONE = 1
)

func StartMatchService() error {

	if !store.RedisConnPoolOk() {
		return fmt.Errorf("strart match err!redis not init!")
	}

	// 1v1 start
	strategyOne := NewMatchStrategyOne()
	go strategyOne.Start()

	// 2v2 start
	// xx:=Newxxxxxx()
	// xxx.Start()

	// 5v5 start
	// xx:=Newxxxxxx()
	// xxx.Start()

	return nil
}
