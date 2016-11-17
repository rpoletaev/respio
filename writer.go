package respio

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type RESPWriter struct {
	writer *bufio.Writer
}

func NewWriter(w io.Writer) *RESPWriter {
	return &RESPWriter{
		writer: bufio.NewWriter(w),
	}
}

func (w *RESPWriter) SendArray(length int64) error {
	blkLen := []byte{'*'}
	w.writer.Write(strconv.AppendInt(blkLen, length, 10))
	_, err := w.writer.WriteString("\r\n")
	return err
}

//SendBulkString convert val to Bulk string format
func (w *RESPWriter) SendBulkString(val string) error {
	blkLen := []byte{'$'}
	blkLen = strconv.AppendInt(blkLen, int64(len(val)), 10)
	w.writer.Write(blkLen)
	w.writer.WriteString("\r\n")
	w.writer.WriteString(val)
	_, err := w.writer.WriteString("\r\n")
	if err != nil {
		println("error on sending string: ", err.Error())
	}
	return err
}

//SendBulk convert val to Bulk format and send
func (w *RESPWriter) SendBulk(val []byte) error {
	blkLen := []byte{'$'}
	blkLen = strconv.AppendInt(blkLen, int64(len(val)), 10)
	w.writer.Write(blkLen)
	w.writer.WriteString("\r\n")
	w.writer.Write(val)
	_, err := w.writer.WriteString("\r\n")
	return err
}

//SendRESPInt format and send RESP Int val as :val
//http://redis.io/topics/protocol#resp-integers
func (w *RESPWriter) SendRESPInt(val int64) error {
	bts := []byte{':'}
	bts = strconv.AppendInt(bts, val, 10)
	w.writer.Write(bts)
	_, err := w.writer.WriteString("\r\n")
	return err
}

//SendBulkInt format integer value as bulk string to send command to server
func (w *RESPWriter) SendBulkInt(val int64) error {
	var valBuf []byte
	valBuf = strconv.AppendInt(valBuf, val, 10)
	return w.SendBulk(valBuf)
}

//SendCmd format cmd and parameters into resp array
func (w *RESPWriter) SendCmd(cmd string, prs []interface{}) error {
	bts := []byte{'*'}
	bts = strconv.AppendInt(bts, int64(len(prs)+1), 10)
	w.writer.Write(bts)
	w.writer.WriteString("\r\n")
	err := w.SendBulkString(cmd)

	if prs == nil {
		return nil
	}

	for _, prm := range prs {
		if err != nil {
			return err
		}
		switch param := prm.(type) {
		case []byte:
			err = w.SendBulk(param)
			break
		case string:
			err = w.SendBulkString(param)
			break
		case int:
			err = w.SendBulkInt(int64(param))
			break
		case int64:
			err = w.SendBulkInt(param)
			break
		default:
			err = errors.New("WRONGTYPE")
			break
		}
	}
	return err
}

//SendNil send nil value as $-1\r\n
func (w *RESPWriter) SendNil() error {
	w.writer.Write([]byte{'$', '-', '1'})
	_, err := w.writer.WriteString("\r\n")
	w.Flush()
	return err
}

// SendError write string in format -err\r\n
func (w *RESPWriter) SendError(err string) error {
	w.writer.Write([]byte{'-'})
	w.writer.WriteString(err)
	_, werr := w.writer.WriteString("\r\n")
	w.Flush()
	return werr
}

//SendSimpleString  write string in format +str\r\n
func (w *RESPWriter) SendSimpleString(str string) error {
	w.writer.Write([]byte{'+'})
	w.writer.WriteString(str)
	_, werr := w.writer.WriteString("\r\n")
	w.Flush()
	return werr
}

//Flush just flush
func (w *RESPWriter) Flush() error {
	return w.writer.Flush()
}
