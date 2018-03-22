package configuration

import (
	"bytes"
	"io"
	"encoding/gob"
)

type Writer struct {
	closed bool
	write io.Writer
	enc gob.Encoder
	item chan MsgItem
}


type MsgItem struct {
	Type string
	Data []byte
}

type MsgHeader struct {
	Type string
	Length uint32
}


func NewWriter(cw io.Writer) (w *Writer) {
	w = new(Writer)
	w.writer = cw
	w.enc = gob.NewEncoder(w.writer)
	w.item = make(chan MsgItem)
	go func() {
		for i := range w.item {
			w.enc.Encode(MsgHeader{ Type:i.Type, Length:uint32(len(i.Data)) })
			w.write.Write(i.Data)
		}
	}()
	return
}

func (w *Writer) Write(mi MsgItem) error {
	w.item <- mi
}
