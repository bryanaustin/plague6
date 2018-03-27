package configuration

import (
	"encoding/binary"
	"encoding/gob"
	"io"
	"fmt"
	"bytes"
)

type Writer struct {
	// closed bool
	write  io.Writer
	item   chan MsgItem
}

type Reader struct {
	Item chan interface{}
}

type MsgItem struct {
	Type string
	Data []byte
}

type MsgHeader struct {
	Type   string
	Length uint32
}

func NewWriter(cw io.Writer) (w *Writer) {
	w = new(Writer)
	w.write = cw
	w.item = make(chan MsgItem)
	go func() {
		for i := range w.item {
			var header [12]byte
			copy(header[:4], []byte(i.Type))
			binary.PutUvarint(header[4:], uint64(len(i.Data)))
			w.write.Write(header[:])
			w.write.Write(i.Data)
		}
	}()
	return
}

func (w *Writer) Write(mi MsgItem) error {
	w.item <- mi
	return nil
}

func (w *Writer) WriteObj(x interface{}, objtype string) error {
	sbconf := new(bytes.Buffer)
	sconf := gob.NewEncoder(sbconf)
	if err := sconf.Encode(x); err != nil {
		return fmt.Errorf("converting %v to output format: %s", x, err)
	}

	// Output configuration
	w.Write(MsgItem{objtype, sbconf.Bytes()})
	return nil
}

func NewReader(cr io.Reader) (r *Reader) {
	r = new(Reader)
	r.Item = make(chan interface{})
	go func() {
		for {
			var input [12]byte
			_, err := io.ReadFull(cr, input[:])
			if err != nil {
				if err != io.EOF {
					r.Item <- fmt.Errorf("reading input header: %s", err)
				}
				close(r.Item)
				return
			}
			size, _ := binary.Uvarint(input[4:])
			data := make([]byte, int(size))
			_, err = io.ReadFull(cr, data)
			if err != nil {
				r.Item <- fmt.Errorf("reading input body: %s", err)
				close(r.Item)
				return
			}
			msg := &MsgItem{ Type:string(input[:4]), Data:data }
			r.Item <- Decode(msg)
		}
	}()
	return
}


func (r Reader) WaitRead() interface{} {
	return <-r.Item
}

