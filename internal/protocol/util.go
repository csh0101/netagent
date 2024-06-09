package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func PackMessage(route uint32, buf []byte) ([]byte, error) {

	r := new(bytes.Buffer)

	if err := binary.Write(r, binary.BigEndian, route); err != nil {
		return nil, fmt.Errorf("failed to encode route: %v", err)
	}

	rawBufLen := int32(len(buf))

	if err := binary.Write(r, binary.BigEndian, rawBufLen); err != nil {
		return nil, fmt.Errorf("failed to encode buffer length: %v", err)
	}

	rb := r.Bytes()
	rb = append(rb, buf...)

	return rb, nil
}

func UnPackMessage(r io.Reader) (uint32, []byte, error) {

	var route uint32
	if err := binary.Read(r, binary.BigEndian, &route); err != nil {
		if io.EOF == err {
			return 0, nil, err
		}
		return 0, nil, fmt.Errorf("failed to decode route: %v", err)
	}

	var rawBufLen int32
	if err := binary.Read(r, binary.BigEndian, &rawBufLen); err != nil {
		return 0, nil, fmt.Errorf("failed to decode buffer length: %v", err)
	}

	if rawBufLen < 0 {
		return 0, nil, fmt.Errorf("invalid buffer length: %v", rawBufLen)
	}

	buf := make([]byte, rawBufLen)
	if _, err := r.Read(buf); err != nil {
		return 0, nil, fmt.Errorf("failed to read buffer: %v", err)
	}

	return route, buf, nil
}
