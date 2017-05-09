package protocol

import (
  capn "github.com/glycerine/go-capnproto"
  "io"
)




func (s *ConcatReply) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	ConcatReplyGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *ConcatReply) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootConcatReplyCapn(capMsg)
      ConcatReplyCapnToGo(z, s)
   return nil
}



func ConcatReplyCapnToGo(src ConcatReplyCapn, dest *ConcatReply) *ConcatReply {
  if dest == nil {
    dest = &ConcatReply{}
  }
  dest.RetCode = src.RetCode()
  dest.Val = src.Val()

  return dest
}



func ConcatReplyGoToCapn(seg *capn.Segment, src *ConcatReply) ConcatReplyCapn {
  dest := AutoNewConcatReplyCapn(seg)
  dest.SetRetCode(src.RetCode)
  dest.SetVal(src.Val)

  return dest
}



func (s *ConcatRequest) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	ConcatRequestGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *ConcatRequest) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootConcatRequestCapn(capMsg)
      ConcatRequestCapnToGo(z, s)
   return nil
}



func ConcatRequestCapnToGo(src ConcatRequestCapn, dest *ConcatRequest) *ConcatRequest {
  if dest == nil {
    dest = &ConcatRequest{}
  }
  dest.UserId = src.UserId()
  dest.Str1 = src.Str1()
  dest.Str2 = src.Str2()

  return dest
}



func ConcatRequestGoToCapn(seg *capn.Segment, src *ConcatRequest) ConcatRequestCapn {
  dest := AutoNewConcatRequestCapn(seg)
  dest.SetUserId(src.UserId)
  dest.SetStr1(src.Str1)
  dest.SetStr2(src.Str2)

  return dest
}



func (s *MatchReply) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	MatchReplyGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *MatchReply) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootMatchReplyCapn(capMsg)
      MatchReplyCapnToGo(z, s)
   return nil
}



func MatchReplyCapnToGo(src MatchReplyCapn, dest *MatchReply) *MatchReply {
  if dest == nil {
    dest = &MatchReply{}
  }
  dest.RetCode = src.RetCode()

  return dest
}



func MatchReplyGoToCapn(seg *capn.Segment, src *MatchReply) MatchReplyCapn {
  dest := AutoNewMatchReplyCapn(seg)
  dest.SetRetCode(src.RetCode)

  return dest
}



func (s *MatchRequest) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	MatchRequestGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *MatchRequest) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootMatchRequestCapn(capMsg)
      MatchRequestCapnToGo(z, s)
   return nil
}



func MatchRequestCapnToGo(src MatchRequestCapn, dest *MatchRequest) *MatchRequest {
  if dest == nil {
    dest = &MatchRequest{}
  }
  dest.GameGroup = src.GameGroup()
  dest.GameDb = src.GameDb()
  dest.ServerId = src.ServerId()
  dest.UserId = src.UserId()
  dest.Score = int(src.Score())

  return dest
}



func MatchRequestGoToCapn(seg *capn.Segment, src *MatchRequest) MatchRequestCapn {
  dest := AutoNewMatchRequestCapn(seg)
  dest.SetGameGroup(src.GameGroup)
  dest.SetGameDb(src.GameDb)
  dest.SetServerId(src.ServerId)
  dest.SetUserId(src.UserId)
  dest.SetScore(int64(src.Score))

  return dest
}



func (s *SumReply) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	SumReplyGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *SumReply) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootSumReplyCapn(capMsg)
      SumReplyCapnToGo(z, s)
   return nil
}



func SumReplyCapnToGo(src SumReplyCapn, dest *SumReply) *SumReply {
  if dest == nil {
    dest = &SumReply{}
  }
  dest.RetCode = src.RetCode()
  dest.Val = src.Val()

  return dest
}



func SumReplyGoToCapn(seg *capn.Segment, src *SumReply) SumReplyCapn {
  dest := AutoNewSumReplyCapn(seg)
  dest.SetRetCode(src.RetCode)
  dest.SetVal(src.Val)

  return dest
}



func (s *SumRequest) Save(w io.Writer) error {
  	seg := capn.NewBuffer(nil)
  	SumRequestGoToCapn(seg, s)
    _, err := seg.WriteTo(w)
    return err
}
 


func (s *SumRequest) Load(r io.Reader) error {
  	capMsg, err := capn.ReadFromStream(r, nil)
  	if err != nil {
  		//panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
        return err
  	}
  	z := ReadRootSumRequestCapn(capMsg)
      SumRequestCapnToGo(z, s)
   return nil
}



func SumRequestCapnToGo(src SumRequestCapn, dest *SumRequest) *SumRequest {
  if dest == nil {
    dest = &SumRequest{}
  }
  dest.UserId = src.UserId()
  dest.Num1 = src.Num1()
  dest.Num2 = src.Num2()

  return dest
}



func SumRequestGoToCapn(seg *capn.Segment, src *SumRequest) SumRequestCapn {
  dest := AutoNewSumRequestCapn(seg)
  dest.SetUserId(src.UserId)
  dest.SetNum1(src.Num1)
  dest.SetNum2(src.Num2)

  return dest
}
