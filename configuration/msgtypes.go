package configuration

import (
	"encoding/gob"
	"fmt"
	"bytes"
	"time"
)

const (
	MsgTypeConfig = "cfgx"
	MsgTypeWorkerStats = "wkst"
	MsgTypeDebugMessage = "debg"
	MsgTypeHit = "hitt"
)

type WorkerStats struct {
	Id int
	Success, Fail uint64
}

type DebugMessage struct {
	Id int
	Message string
}

type Hit struct {
	Id int
	Started, Finished time.Time
	ErrorType         string
	BodySize          uint
}


func Decode(m *MsgItem) interface{} {
	buf := bytes.NewBuffer(m.Data)
	dec := gob.NewDecoder(buf)
	switch m.Type {
		case MsgTypeConfig:
			x := new(Configuration)
			return tryDecode(dec, x)
		
		case MsgTypeWorkerStats:
			x := new(WorkerStats)
			return tryDecode(dec, x)

		case MsgTypeDebugMessage:
			x := new(DebugMessage)
			return tryDecode(dec, x)

		case MsgTypeHit:
			x := new(Hit)
			return tryDecode(dec, x)
	}
	return fmt.Errorf("Unable to decode unknown message type %q", m.Type)
}

func tryDecode(d *gob.Decoder, x interface{}) interface{} {
	if err := d.Decode(x); err != nil {
		return fmt.Errorf("decoding data: %s", err)
	}
	return x
}
