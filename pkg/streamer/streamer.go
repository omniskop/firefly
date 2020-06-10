// Package streamer contains everything to stream an animation to a client
// It will write the protocol data to an io.Writer and maybe also read for bidirectional communication.
package streamer

import (
	"io"
	"math"

	"github.com/sirupsen/logrus"

	"github.com/omniskop/firefly/pkg/scanner"
)

type Streamer struct {
	destination io.Writer
	gamma       float64
}

func New(dst io.Writer) Streamer {
	return Streamer{
		destination: dst,
		gamma:       2.2,
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
		// map from 0xffff to 0xff and apply gamma correction
		data[i*3+1] = byte(math.Pow(float64(r)/0xffff, s.gamma) * 0xff)
		data[i*3+2] = byte(math.Pow(float64(g)/0xffff, s.gamma) * 0xff)
		data[i*3+3] = byte(math.Pow(float64(b)/0xffff, s.gamma) * 0xff)
	}
	_, err := s.destination.Write(data)
	if err != nil {
		logrus.Errorf("streaming error: %v", err)
	}
}
