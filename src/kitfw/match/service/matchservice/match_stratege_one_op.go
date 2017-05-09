package matchservice

import (
	"encoding/json"
	"fmt"
	logger "kitfw/commom/log"
	"kitfw/commom/phpproxy"
	"kitfw/commom/store"
	"kitfw/match/envconf"

	"github.com/garyburd/redigo/redis"
)

type MatchPlayerInfo struct {
	Group    string
	Db       string
	Serverid int64
	Token    string
}

func NewMatchPlayerInfo(group string, db string, serverid int64, token string) *MatchPlayerInfo {
	return &MatchPlayerInfo{
		Group:    group,
		Db:       db,
		Serverid: serverid,
		Token:    token,
	}
}
func GetMatchPlayerInfo(key string) (*MatchPlayerInfo, error) {
	jsonbyts, err := redis.Bytes(store.Do(key, "GET", key))
	if redis.ErrNil == err {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if jsonbyts == nil {
		return nil, nil
	}
	ret := &MatchPlayerInfo{}
	err = json.Unmarshal(jsonbyts, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func SetMatchPlayerInfo(key string, info *MatchPlayerInfo, timeout int) error {
	bs, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = store.Do(key, "SET", key, bs)

	// expire
	_, err = store.Do(key, "EXPIRE", key, timeout)
	if err != nil {
		store.Do(key, "DEL", key)
		return err
	}
	return err
}

func DelMatchPlayerInfo(key string) error {
	_, err := store.Do(key, "DEL", key)
	return err
}

func CheckAndPopPlayerIds(listname string) ([]int64, error) {
	key := listname
	script := fmt.Sprintf("local len=redis.call('LLEN','%s'); len=len-math.mod(len,2)-1; if len > 0 then local val=redis.call('LRANGE','%s',0,len); redis.call('LTRIM','%s',len+1, -1); return val; end ;return {};", listname, listname, listname)
	cmd := make([]interface{}, 1)
	cmd[0] = 0
	strs, err := redis.Strings(store.DoScript(key, script, cmd))
	if redis.ErrNil == err {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	ret := []int64{}
	for i := 0; i < len(strs); i++ {
		if strs[i] != "" {
			r, err := redis.Int64([]byte(strs[i]), nil)
			if err != nil {
				return nil, err
			}
			ret = append(ret, r)
		}
	}
	return ret, nil
}

func CheckAndPopOneFromList(key string) (int64, error) {
	script := fmt.Sprintf("if redis.call('LLEN','%s') == 1 then return redis.call('LPOP','%s'); end ;return 0;", key, key)
	cmd := make([]interface{}, 1)
	cmd[0] = 0
	return redis.Int64(store.DoScript(key, script, cmd))
}

func PushBackMatchList(listname string, playerid int64) error {
	strplayerid := fmt.Sprintf("%d", playerid)
	_, err := store.Do(listname, "RPUSH", listname, []byte(strplayerid))
	return err
}

func processMatch(listname string, playerid int64, matched_playerid int64) {

	if playerid == 0 || matched_playerid == 0 {
		logger.Error("err", fmt.Sprintf("processMatch err!playerid zero!listname:%s playerid1:%d playerid2:%s", listname, playerid, matched_playerid))
		return
	}

	key := GetTimeOutKey(playerid)
	defer DelMatchPlayerInfo(key)

	info, err := GetMatchPlayerInfo(key)
	if err != nil {
		logger.Error("match_err", err)
		return
	}
	if info == nil {
		logger.Error("match_err", fmt.Sprintf("processMatch error!info nil!listname:%s playerid:%d matched_playerid:%d", listname, playerid, matched_playerid))
		return
	}
	err = phpproxy.ExecuteTask(
		playerid,
		info.Group,
		info.Db,
		info.Serverid,
		info.Token,
		"top.updateBePraiseInfo",
		[]string{fmt.Sprintf("%d matched player %d", playerid, matched_playerid)},
	)
	if err != nil {
		logger.Error("logid", info.Token, "match_err", err, "playerid", playerid, "matchedid", matched_playerid)
		return
	}
	logger.Info("logid", info.Token, "executeTask", "match-ok", "playerid", playerid, "matchedid", matched_playerid)
}

func removeDuplicates(playerids []int64) []int64 {
	// Use map to record duplicates as we find them.
	encountered := map[int64]bool{}
	result := []int64{}

	for v := range playerids {
		if encountered[playerids[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[playerids[v]] = true
			// Append to result slice.
			result = append(result, playerids[v])
		}
	}
	// Return the new slice.
	return result
}

func getAllScoreInterval() []int {
	ret := []int{}
	for i := envconf.EnvMatchCfg.MATCH_ONE_MIN_SCORE; i <= envconf.EnvMatchCfg.MATCH_ONE_MAX_SCORE; i = i + envconf.EnvMatchCfg.MATCH_ONE_SCORE_INTERVAL {
		ret = append(ret, i)
	}
	return ret
}

func GetListNameByScore(score int) string {
	itervals := getAllScoreInterval()
	for i := 0; i < len(itervals); i++ {
		if score >= itervals[i] && score <= itervals[i]+envconf.EnvMatchCfg.MATCH_ONE_SCORE_INTERVAL-1 {
			return GetListKeyByScore(itervals[i])
		}
	}
	return ""
}

func GetListKeyByScore(score int) string {
	str := fmt.Sprintf("%v%d_%d", MATCH_ONE_PREFIX_STRING, score, score+int(envconf.EnvMatchCfg.MATCH_ONE_SCORE_INTERVAL-1))
	return str
}

func GetTimeOutKey(playerid int64) string {
	return fmt.Sprintf("%v%d", MATCH_ONE_PREFIX_TIMEOUT_STRING, playerid)
}

/*type matchOneConfig struct {
	minScore      int
	maxScore      int
	scoreInterval int
	scoreLimit    int
	netThread     int
}
config := &matchOneConfig{
		minScore:      envconf.EnvMatchCfg.MATCH_ONE_MIN_SCORE,
		maxScore:      envconf.EnvMatchCfg.MATCH_ONE_MAX_SCORE,
		scoreInterval: envconf.EnvMatchCfg.MATCH_ONE_SCORE_INTERVAL,
		scoreLimit:    envconf.EnvMatchCfg.MATCH_ONE_SCORE_LIMIT,
		netThread:     envconf.EnvMatchCfg.MATCH_ONE_NET_THREAD,
	}
*/
