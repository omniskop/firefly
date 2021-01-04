package streamer

import (
	"io"

	"github.com/sirupsen/logrus"

	"github.com/omniskop/firefly/pkg/scanner"
)

type WLEDStreamer struct {
	destination io.Writer
}

// NewWLED creates a new Streamer that can control a WLED Device.
// The Streamer does not perform gamma correction as that is handled by WLED iself.
func NewWLED(dst io.Writer) *WLEDStreamer {
	return &WLEDStreamer{
		destination: dst,
	}
}

func (s *WLEDStreamer) Stream(frame scanner.Frame) {
	// Currently only the DRGB protocol is used.
	// In the future it could be expanded to add support for more wled protocols.

	var packet = make([]byte, 2+3*len(frame.Pixels))
	packet[0] = 2   // 2 = DRGB protocol
	packet[1] = 255 // 255 = keep this frame until told otherwise
	for i, pixel := range frame.Pixels {
		r, g, b, _ := pixel.RGBA()
		// map from 0xffff to 0xff
		packet[2+i*3+0] = byte(float64(r) / 0xffff * 0xff)
		packet[2+i*3+1] = byte(float64(g) / 0xffff * 0xff)
		packet[2+i*3+2] = byte(float64(b) / 0xffff * 0xff)
	}
	_, err := s.destination.Write(packet)
	if err != nil {
		logrus.Errorf("streaming error: %v", err)
	}
}
