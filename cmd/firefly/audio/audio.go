package audio

import (
	"errors"

	"github.com/omniskop/firefly/pkg/project"
)

type Player interface {
	Play()
	Pause()
	Time() float64     // the current playback time in seconds
	SetTime(float64)   // sets the playback time
	Duration() float64 // the duration of the song in seconds
	Volume() float64   // the volume of the song as a value between 0 and 1
	SetVolume(float64) // sets the volume of the song to a value between 0 and 1
}

type Provider interface {
	Provide(audio project.Audio) (Player, bool)
	CanProvide(audio project.Audio) bool
}

var NoProviderErr = errors.New("no provider for the audio found")

var providers []Provider

func Register(prov Provider) {
	providers = append(providers, prov)
}

func Open(audio project.Audio) (Player, error) {
	provider, ok := findProvider(audio)
	if !ok {
		return nil, NoProviderErr
	}
	player, _ := provider.Provide(audio)
	return player, nil
}

func findProvider(audio project.Audio) (Provider, bool) {
	for _, provider := range providers {
		ok := provider.CanProvide(audio)
		if ok {
			return provider, true
		}
	}
	return nil, false
}
