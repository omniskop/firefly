package editor

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/omniskop/firefly/cmd/firefly/settings"

	"github.com/omniskop/firefly/cmd/firefly/audio"

	"github.com/omniskop/firefly/pkg/project"

	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/multimedia"
)

type audioPlayer struct {
	*multimedia.QMediaPlayer
	probe          *multimedia.QAudioProbe
	mediaPath      string
	onTimeChangedF func(t float64)
	onReadyF       func()
	onErrorF       func(error)
}

func NewAudioPlayer(mediapath string) *audioPlayer {
	player := audioPlayer{
		QMediaPlayer: multimedia.NewQMediaPlayer(nil, multimedia.QMediaPlayer__LowLatency),
		probe:        multimedia.NewQAudioProbe(nil),
		mediaPath:    mediapath,
	}

	player.SetNotifyInterval(16)
	player.ConnectPositionChanged(player.positionChangedEvent)
	player.ConnectMediaStatusChanged(player.mediaStatusChangedEvent)
	player.ConnectError2(player.errorEvent)

	player.setFile(mediapath)

	return &player
}

func (p *audioPlayer) setFile(file string) {
	p.SetMedia(multimedia.NewQMediaContent2(core.QUrl_FromLocalFile(file)), nil)
}

func (p *audioPlayer) play() {
	p.QMediaPlayer.Play()
}

func (p *audioPlayer) pause() {
	p.QMediaPlayer.Pause()
	p.QMediaPlayer.SetPosition(p.QMediaPlayer.Position())
}

func (p *audioPlayer) time() float64 {
	return float64(p.Position()) / 1000 // milliseconds to seconds
}

func (p *audioPlayer) setTime(value float64) {
	p.SetPosition(int64(value * 1000))
}

func (p *audioPlayer) duration() float64 {
	return float64(p.QMediaPlayer.Duration()) / 1000 // milliseconds to seconds
}

func (p *audioPlayer) onTimeChanged(f func(float64)) {
	p.onTimeChangedF = f
}

func (p *audioPlayer) onReady(f func()) {
	p.onReadyF = f
}

func (p *audioPlayer) onError(f func(error)) {
	p.onErrorF = f
}

func (p *audioPlayer) positionChangedEvent(position int64) {
	// position in ms
	if p.onTimeChangedF != nil {
		p.onTimeChangedF(float64(position) / 1000)
	}
}

func (p *audioPlayer) errorEvent(qerr multimedia.QMediaPlayer__Error) {
	var err error
	switch qerr {
	case multimedia.QMediaPlayer__NoError:
		return
	case multimedia.QMediaPlayer__ResourceError:
		err = fmt.Errorf("fileplayer: resource error: %s", p.ErrorString())
	case multimedia.QMediaPlayer__FormatError:
		err = fmt.Errorf("fileplayer: format error: %s", p.ErrorString())
	case multimedia.QMediaPlayer__NetworkError:
		err = fmt.Errorf("fileplayer: network error: %s", p.ErrorString())
	case multimedia.QMediaPlayer__AccessDeniedError:
		err = fmt.Errorf("fileplayer: access denied: %s", p.ErrorString())
	case multimedia.QMediaPlayer__ServiceMissingError:
		err = fmt.Errorf("fileplayer: service missing: %s", p.ErrorString())
	default:
		err = fmt.Errorf("fileplayer: unknown error: %s", p.ErrorString())
	}
	logrus.Errorf("[Audio] %v", err)
	if p.onErrorF != nil {
		p.onErrorF(err)
	}
}

func (p *audioPlayer) mediaStatusChangedEvent(status multimedia.QMediaPlayer__MediaStatus) {
	switch status {
	case multimedia.QMediaPlayer__UnknownMediaStatus:
		logrus.Info("[Audio] Media Status: Unknown Media Status")
	case multimedia.QMediaPlayer__NoMedia:
		logrus.Info("[Audio] Media Status: No Media")
	case multimedia.QMediaPlayer__LoadingMedia:
		logrus.Info("[Audio] Media Status: Loading Media")
	case multimedia.QMediaPlayer__LoadedMedia:
		logrus.Info("[Audio] Media Status: Loaded Media")
		if p.onReadyF != nil {
			p.onReadyF()
		}
	case multimedia.QMediaPlayer__StalledMedia:
		logrus.Info("[Audio] Media Status: Stalled Media")
	case multimedia.QMediaPlayer__BufferingMedia:
		logrus.Info("[Audio] Media Status: Buffering Media")
	case multimedia.QMediaPlayer__BufferedMedia:
		logrus.Info("[Audio] Media Status: Buffered Media")
	case multimedia.QMediaPlayer__EndOfMedia:
		logrus.Info("[Audio] Media Status: End Of Media")
	case multimedia.QMediaPlayer__InvalidMedia:
		logrus.Info("[Audio] Media Status: Invalid Media")
	}
}

// LocateAudioFile searched for an audio file that matches the given Audio struct and returns it's path.
// The function will check all locations set by the user but can optionally also take additional ones.
func LocateAudioFile(audioInfo project.Audio, locations ...string) (string, error) {
	locations = append(locations, settings.GetStrings("audio/fileSources")...)
	files, errs := audio.GetAllFiles(locations)
	for _, err := range errs {
		var pe *os.PathError
		if errors.As(err, &pe) && os.IsNotExist(err) {
			logrus.Warnf("audio path %q does not exist", pe.Path)
		} else {
			logrus.Errorf("locate audio file: %v", err)
		}
	}

	audioInfo.Title = strings.ToLower(audioInfo.Title)
	audioInfo.Author = strings.ToLower(audioInfo.Author)

	for _, file := range files {
		name := strings.ToLower(file.Info.Name())
		if strings.Contains(name, audioInfo.Title) && strings.Contains(name, audioInfo.Author) {
			return file.Path, nil
		}
	}
	return "", errors.New("no matching files found")
}
