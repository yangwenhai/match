package matchservice

import (
	"fmt"
	logger "kitfw/commom/log"
	"kitfw/match/envconf"
	"time"
)

const (
	MATCH_ONE_PREFIX_STRING         = "match_one_list_"
	MATCH_ONE_PREFIX_TIMEOUT_STRING = "match_one:timeout:player_"
)

var superInterValMap map[int][]int64

//1v1 match
type matchStrategyOne struct {
	netMatchMap map[int]*matchOneService
	superMatch  *matchOneService
}

type matchOneService struct {
	scoreQueue []*struct{ score int }
}

func NewMatchStrategyOne() *matchStrategyOne {

	ret := &matchStrategyOne{
		netMatchMap: make(map[int]*matchOneService),
		superMatch: &matchOneService{
			scoreQueue: make([]*struct{ score int }, 0),
		},
	}

	for i := 0; i < envconf.EnvMatchCfg.MATCH_ONE_NET_THREAD; i++ {
		ret.netMatchMap[i] = &matchOneService{
			scoreQueue: make([]*struct{ score int }, 0),
		}
	}

	intervals := getAllScoreInterval()
	superInterValMap = make(map[int][]int64, len(intervals))

	for i := 0; i < len(intervals); i++ {
		a := ret.netMatchMap[i%envconf.EnvMatchCfg.MATCH_ONE_NET_THREAD]
		a.scoreQueue = append(a.scoreQueue, &struct{ score int }{score: int(intervals[i])})
		ret.superMatch.scoreQueue = append(ret.superMatch.scoreQueue, &struct{ score int }{score: int(intervals[i])})
	}

	return ret
}

func (s *matchStrategyOne) Start() {
	for k, v := range s.netMatchMap {
		go v.run(k)
	}
	go s.superMatch.runSuper()
}

func (s *matchOneService) run(index int) {
	if len(s.scoreQueue) == 0 {
		logger.Info("ignore", fmt.Sprintf("match one return!queue empty!index:%d"), index)
		return
	}
	for {

		need_sleep := true
		for i := 0; i < len(s.scoreQueue); i++ {

			//pop players
			listname := GetListKeyByScore(s.scoreQueue[i].score)
			matchplayerids, err := CheckAndPopPlayerIds(listname)
			if err != nil {
				fmt.Println("CheckAndPopPlayerIds err", err)
				logger.Error("err", fmt.Sprintf("doMatchingNormal error!pop user empty!listname:%s err:%v", listname, err))
				continue
			}
			if matchplayerids == nil || len(matchplayerids) == 0 {
				continue
			}
			need_sleep = false

			uniqlist := removeDuplicates(matchplayerids)
			for i := 0; i < len(uniqlist)-1; i = i + 2 {
				player1 := uniqlist[i]
				player2 := uniqlist[i+1]
				go processMatch(listname, player1, player2)
				go processMatch(listname, player2, player1)
				logger.Info("matched", fmt.Sprintf("listname:%s playerid1:%d playerid2:%d", listname, player1, player2))
			}
		}

		if need_sleep == true {
			time.Sleep(time.Millisecond * 500)
		}

	}

}

func (s *matchOneService) runSuper() {

	intervals := getAllScoreInterval()
	for {
		for i := 0; i < len(intervals); i++ {
			score := intervals[i]
			listname := GetListKeyByScore(score)
			playerid, err := CheckAndPopOneFromList(listname)
			if err != nil {
				logger.Error("err", fmt.Sprintf("super match error!pop user empty!listname:%s err:%v playerid:%d", listname, err, playerid))
				continue
			}
			if playerid == 0 {
				continue
			}
			logger.Info("supermatch", fmt.Sprintf("pop playerid:%d listname:%s", playerid, listname))
			superInterValMap[score] = append(superInterValMap[score], playerid)
		}
		s.deleteMatchedPlayerids(intervals)
		s.doSuperMatch(intervals)
		time.Sleep(time.Millisecond * 500)
	}
}

func (s *matchOneService) deleteMatchedPlayerids(intervals []int) {
	for i := 0; i < len(intervals); i++ {
		score := intervals[i]
		length := len(superInterValMap[score])
		if length == 0 {
			continue
		}
		deleteindexs := []int{}

		for j := 0; j < length; j++ {
			playerid := superInterValMap[score][j]
			m, _ := GetMatchPlayerInfo(GetTimeOutKey(playerid))
			if m == nil {
				deleteindexs = append(deleteindexs, j)
				logger.Info("super-match-ignore", fmt.Sprintf("playerid:%d listname:%s aleady matched or timeout!", playerid, GetListKeyByScore(score)))
			}
		}
		for k := len(deleteindexs) - 1; k >= 0; k-- {
			index := deleteindexs[k]
			superInterValMap[score] = append(superInterValMap[score][:index], superInterValMap[score][index+1:]...)
		}
	}
}

func (s *matchOneService) doSuperMatch(intervals []int) {

	deleteindexs := make(map[int][]int)
	for i := 0; i < len(intervals); i++ {

		score := intervals[i]
		if len(superInterValMap[score]) == 0 {
			continue
		}

		superInterValMap[score] = removeDuplicates(superInterValMap[score])

		if len(superInterValMap[score]) > 1 {

			for k := 0; k < len(superInterValMap[score])-1; k = k + 2 {

				playerid1 := superInterValMap[score][k]
				playerid2 := superInterValMap[score][k+1]
				listname := GetListKeyByScore(score)

				go processMatch(listname, playerid1, playerid2)
				go processMatch(listname, playerid2, playerid1)
				deleteindexs[score] = append(deleteindexs[score], k, k+1)

				logger.Info("match-super-1", fmt.Sprintf("score:%d playerid1:%d playerid2:%d", score, playerid1, playerid2))
			}
		}
	}

	for score, v := range deleteindexs {
		for k := len(v) - 1; k >= 0; k-- {
			index := v[k]
			superInterValMap[score] = append(superInterValMap[score][:index], superInterValMap[score][index+1:]...)
		}
	}

	for i := 0; i < len(intervals); i++ {
		curscore := intervals[i]
		if len(superInterValMap[curscore]) == 0 {
			continue
		}
		curplayerid := superInterValMap[curscore][0]
		targetscore := int(0)
		targetplayerid := int64(0)
		for k := i + 1; k < len(intervals); k++ {
			targetscore := intervals[k]
			if targetplayerid == 0 && len(superInterValMap[targetscore]) > 0 && targetscore-curscore <= envconf.EnvMatchCfg.MATCH_ONE_SCORE_LIMIT {
				targetplayerid = superInterValMap[targetscore][0]
				listname := GetListKeyByScore(curscore)
				go processMatch(listname, curplayerid, targetplayerid)
				go processMatch(listname, targetplayerid, curplayerid)
				logger.Info("match-super-2", fmt.Sprintf("curscore:%d targetscore:%d playerid1:%d playerid2:%d", curscore, targetscore, curplayerid, targetplayerid))
			}
		}
		if targetplayerid > 0 {
			superInterValMap[curscore] = []int64{}
			superInterValMap[targetscore] = []int64{}
		}
	}

}
