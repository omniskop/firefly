package streamer

import (
	"github.com/omniskop/firefly/pkg/scanner"
)

// Pipeline contains a scanner and a streamer that operate asynchronously.
// Pipe a time into the Update channel to start a scan and a subsequent stream.
type Pipeline struct {
	Scanner   scanner.Scanner
	Streamers []Streamer
	LastFrame scanner.Frame
	Update    chan float64
}

// NewPipeline creates a new Pipeline
func NewPipeline(sca scanner.Scanner, str Streamer) *Pipeline {
	sp := &Pipeline{
		Scanner:   sca,
		Streamers: []Streamer{str},
		LastFrame: scanner.Frame{},
		Update:    make(chan float64),
	}
	go sp.routine()
	return sp
}

func (sp *Pipeline) AddStreamer(streamer Streamer) {
	sp.Streamers = append(sp.Streamers, streamer)
}

// Stop the pipeline. After calling Stop the Update channel will be closed.
func (sp *Pipeline) Stop() {
	close(sp.Update)
}

// routine listens on the Update channel to start the scanner and streamer
func (sp *Pipeline) routine() {
	for time := range sp.Update {
		sp.LastFrame = sp.Scanner.Scan(time)
		for _, streamer := range sp.Streamers {
			streamer.Stream(sp.LastFrame)
		}
	}
}
