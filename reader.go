package respio

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

//RESPReader allow read RESP
type RESPReader struct {
	reader *bufio.Reader
}

// NewReader creates new RESPReader
func NewReader(reader io.Reader) *RESPReader {
	return &RESPReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *RESPReader) getLine() ([]byte, error) {
	line, _, err := r.reader.ReadLine()
	if err != nil {
		return nil, err
	}

	if len(line) < 2 {
		return nil, fmt.Errorf("Wrong RESP line")
	}

	return line, nil
}

func (r *RESPReader) Read() (interface{}, error) {
	line, err := r.getLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case '-':
		return errors.New(string(line[1:])), nil
	case '+':
		return string(line[1:]), nil
	case ':':
		return strconv.ParseInt(string(line[1:]), 10, 64)
	case '$':
		return r.parseBulk(line[1:])
	case '*':
		length, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if length == -1 && err == nil {
			return nil, nil
		}

		response := make([]interface{}, length)
		for i := range response {
			response[i], err = r.Read()
			if err != nil {
				return nil, err
			}
		}
		return response, err
	default:
		return nil, errors.New("WRONGTYPE")
	}
}

//ReadCommand returns incomming command name, parameters or error
//RESP command format must be array
func (r *RESPReader) ReadCommand() (string, []interface{}, error) {
	rawCmd, err := r.Read()
	if err != nil {
		return "", nil, fmt.Errorf("Wrong command format: %s", err)
	}
	switch rawSlice := rawCmd.(type) {
	case []interface{}:
		if len(rawSlice) == 0 {
			return "", nil, fmt.Errorf("Wrong command format")
		}
		if len(rawSlice) == 1 {
			return string(rawSlice[0].([]byte)), nil, nil
		}

		return string(rawSlice[0].([]byte)), rawSlice[1:], nil
	default:
		return "", nil, fmt.Errorf("Wrong command format. Command must be an array")
	}
}

func (r *RESPReader) parseBulk(src []byte) (interface{}, error) {
	length, err := strconv.ParseInt(string(src), 10, 32)
	if length == -1 && err == nil {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Wrong BulkString format: %s", err)
	}

	bulk := make([]byte, length)
	_, err = io.ReadFull(r.reader, bulk)
	if err != nil {
		return nil, err
	}
	// must read emty line to delete "\r\n"
	line, _, err := r.reader.ReadLine()
	if err != nil {
		return nil, err
	}

	if len(line) > 0 {
		return nil, fmt.Errorf("Wrong BulkString format: Body of bulk greater then bulk length")
	}

	return bulk, nil
}
