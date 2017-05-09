package protocol

type SumRequest struct {
	UserId int64 `capid:"0" amf:"userid"`
	Num1   int64 `capid:"1" amf:"num1"`
	Num2   int64 `capid:"2" amf:"num2"`
}

type SumReply struct {
	RetCode int8  `capid:"0" amf:"retcode"`
	Val     int64 `capid:"1" amf:"val"`
}

type ConcatRequest struct {
	UserId int64  `capid:"0" amf:"userid"`
	Str1   string `capid:"1" amf:"str1"`
	Str2   string `capid:"2" amf:"str2"`
}

type ConcatReply struct {
	RetCode int8   `capid:"0" amf:"retcode"`
	Val     string `capid:"1" amf:"val"`
}

type MatchRequest struct {
	GameGroup string `capid:"0" amf:"gameGroup"`
	GameDb    string `capid:"1" amf:"gameDb"`
	ServerId  int64  `capid:"2" amf:"serverId"`
	UserId    int64  `capid:"3" amf:"userid"`
	Score     int    `capid:"4" amf:"score"`
}

type MatchReply struct {
	RetCode int8 `capid:"0" amf:"retcode"`
}
