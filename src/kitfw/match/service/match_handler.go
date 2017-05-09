package service

import (
	"context"
	"fmt"
	logger "kitfw/commom/log"
	"kitfw/commom/phpproxy"
	"kitfw/match/envconf"
	"kitfw/match/service/matchservice"
	"time"
)

func (handler *MatchHandler) doProcess(ctx context.Context) {
	handler.reply.RetCode = 0

	rlogger := ctx.Value("logger").(*logger.Logger)

	//正在匹配中
	key := matchservice.GetTimeOutKey(handler.request.UserId)
	info, err := matchservice.GetMatchPlayerInfo(key)
	if err != nil {
		handler.reply.RetCode = -1
		rlogger.Error("err", err)
		return
	}
	if info != nil {
		handler.reply.RetCode = -2
		rlogger.Error("err", err)
		return
	}

	//加入匹配队列
	listname := matchservice.GetListNameByScore(handler.request.Score)
	err = matchservice.PushBackMatchList(listname, handler.request.UserId)
	if err != nil {
		handler.reply.RetCode = -3
		rlogger.Error("err", err)
		return
	}

	//设置超时信息
	token := ctx.Value("logid").(string)
	newinfo := matchservice.NewMatchPlayerInfo(handler.request.GameGroup, handler.request.GameDb, handler.request.ServerId, token)
	if err := matchservice.SetMatchPlayerInfo(key, newinfo, envconf.EnvMatchCfg.MATCH_TIME_OUT+2); err != nil {
		handler.reply.RetCode = -4
		rlogger.Error("err", err)
		return
	}

	go handler.matchTimeOut(listname, token, rlogger)

	rlogger.Info("group", handler.request.GameGroup, "db", handler.request.GameDb, "serverid", handler.request.ServerId, "userid", handler.request.UserId, "score", handler.request.Score, "listname", listname)
}

func (handler *MatchHandler) matchTimeOut(listname string, token string, rloger *logger.Logger) {

	t := time.NewTimer(time.Duration(envconf.EnvMatchCfg.MATCH_TIME_OUT) * time.Second)
	key := matchservice.GetTimeOutKey(handler.request.UserId)
	defer matchservice.DelMatchPlayerInfo(key)
	for {
		select {
		case <-t.C:
			{
				m, _ := matchservice.GetMatchPlayerInfo(key)
				if m == nil {
					rloger.Info("ok", fmt.Sprintf("match timeout!player:%d alread matched", handler.request.UserId))
					return
				}
				err := phpproxy.ExecuteTask(
					handler.request.UserId,
					handler.request.GameGroup,
					handler.request.GameDb,
					handler.request.ServerId,
					token,
					"top.updateBePraiseInfo",
					[]string{fmt.Sprintf("playerid:%d_match_timeout", handler.request.UserId)},
				)
				if err != nil {
					rloger.Error("error", err)
					return
				}
				rloger.Info("match_timeout", fmt.Sprintf("playerid:%d", handler.request.UserId))
				return
			}
		}
	}
}
