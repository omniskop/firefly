// Package streamer contains everything to stream an animation to a client
// It will write the protocol data to an io.Writer and maybe also read for bidirectional communication.
package streamer

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"sync"

	"github.com/omniskop/firefly/pkg/scanner"
	"github.com/sirupsen/logrus"
)

type Streamer struct {
	destination io.Writer
	gamma       float64
	Version     int
	mutex       sync.Mutex
}

func New(dst io.Writer) Streamer {
	return Streamer{
		destination: dst,
		gamma:       2.2,
		Version:     1,
	}
}

func (s *Streamer) SetDestination(writer io.Writer) {
	s.mutex.Lock()
	s.destination = writer
	s.mutex.Unlock()
}

func (s Streamer) Stream(frame scanner.Frame) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.destination == nil {
		return
	}
	switch s.Version {
	case 0:
		s.streamVersion0(frame)
	case 1:
		s.streamVersion1(frame)
	}
}

func (s Streamer) streamVersion1(frame scanner.Frame) {
	const maxPixelsPerPacket = 300
	// split the data in multiple packets
	packetCount := int(math.Ceil(float64(len(frame.Pixels)) / maxPixelsPerPacket))
	for i := 0; i < packetCount; i++ {
		// get number of pixels in this packet
		pixelCount := len(frame.Pixels) - i*maxPixelsPerPacket
		if pixelCount > maxPixelsPerPacket {
			pixelCount = maxPixelsPerPacket
		}
		//fmt.Printf("packet %d (%d pixels)\n", i, pixelCount)

		packet := bytes.NewBuffer(make([]byte, 0, 6+3*pixelCount))
		packet.WriteByte(1) // packet type
		if i == 0 {
			// write header in first packet
			writeHeader(packet, frame)
		} else {
			packet.WriteByte(0) // header length
		}
		binary.Write(packet, binary.LittleEndian, uint16(i*maxPixelsPerPacket)) // pixel offset
		binary.Write(packet, binary.LittleEndian, uint16(pixelCount*3))         // data length
		for c := 0; c < pixelCount; c++ {
			r, g, b, _ := frame.Pixels[i*maxPixelsPerPacket+c].RGBA()
			packet.Write([]byte{
				byte(math.Pow(float64(r)/0xffff, s.gamma) * 0xff),
				byte(math.Pow(float64(g)/0xffff, s.gamma) * 0xff),
				byte(math.Pow(float64(b)/0xffff, s.gamma) * 0xff),
			})
		}

		_, err := io.Copy(s.destination, packet)
		if err != nil {
			logrus.Errorf("streaming error: %v", err)
		}
	}
}

func writeHeader(packet *bytes.Buffer, frame scanner.Frame) {
	buffer := new(bytes.Buffer)
	// total led count
	buffer.WriteByte(0)
	binary.Write(buffer, binary.LittleEndian, uint16(len(frame.Pixels)))

	if buffer.Len() > 256 {
		logrus.Errorf("streamer: packet 1 header too big (%d)", buffer.Len())
		packet.WriteByte(0)
		return
	}
	packet.WriteByte(byte(buffer.Len()))
	io.Copy(packet, buffer)
}

func (s Streamer) streamVersion0(frame scanner.Frame) {
	var data = make([]byte, 1+3*len(frame.Pixels))
	data[0] = 0
	for i, pixel := range frame.Pixels {
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
