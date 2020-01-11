// Package streamer contains everything to stream an animation to a client
// It will write the protocol data to an io.Writer and maybe also read for bidirectional communication.
package streamer

import (
	"io"

	"github.com/sirupsen/logrus"

	"github.com/omniskop/firefly/pkg/scanner"
)

type Streamer struct {
	destination io.Writer
}

func New(dst io.Writer) Streamer {
	return Streamer{
		destination: dst,
	}
}

func (s Streamer) Stream(frame scanner.Frame) {
	if s.destination == nil {
		return
	}
	var data = make([]byte, 1+3*len(frame.Pixel))
	data[0] = 0
	for i, pixel := range frame.Pixel {
		r, g, b, _ := pixel.RGBA()
		data[i*3+1] = byte(r / 257)
		data[i*3+2] = byte(g / 257)
		data[i*3+3] = byte(b / 257)
	}
	_, err := s.destination.Write(data)
	if err != nil {
		logrus.Error("streaming error: %v", err)
	}
}
