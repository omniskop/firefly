package streamer

import (
	"encoding/gob"
	"fmt"
	"image/color"
	"io"

	"github.com/omniskop/firefly/pkg/scanner"
)

func init() {
	gob.Register(color.RGBA{})
	gob.Register(color.RGBA64{})
	gob.Register(color.NRGBA{})
	gob.Register(color.NRGBA64{})
	gob.Register(color.Gray{})
	gob.Register(color.Gray16{})
}

type GobStreamer struct {
	gamma   float64
	encoder *gob.Encoder
}

func NewGob(dst io.Writer) *GobStreamer {
	return &GobStreamer{
		gamma:   2.2,
		encoder: gob.NewEncoder(dst),
	}
}

func (gs *GobStreamer) Stream(frame scanner.Frame) {
	err := gs.encoder.Encode(frame)
	if err != nil {
		fmt.Println("gob streamer:", err)
	}
}
