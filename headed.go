package headed

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrTruncatedHeader = errors.New("truncated header")
)

type Headed struct {
	io.ReadWriteCloser
	Header []byte
}

func NewHeaded(backing io.ReadWriteCloser) (*Headed, error) {
	h := Headed{ReadWriteCloser: backing}

	var headerLength uint32
	if err := binary.Read(h, binary.BigEndian, &headerLength); err != nil {
		return nil, err
	}

	h.Header = make([]byte, headerLength)

	for complete := uint32(0); complete < headerLength; {
		n, err := h.Read(h.Header[complete:])

		complete += uint32(n)

		if err == io.EOF {
			if complete < headerLength {
				return nil, ErrTruncatedHeader
			} else {
				break
			}
		} else if err != nil {
			return nil, err
		}
	}

	return &h, nil
}
