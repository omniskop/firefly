package audio

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/multimedia"

	"github.com/omniskop/firefly/pkg/project"

	"github.com/sirupsen/logrus"
)

func init() {
	Register(new(fileProvider))
}

func getPathPrefix() string {
	return core.QDir_CurrentPath()
}

type fileProvider struct{}

func (p *fileProvider) CanProvide(audio project.Audio) bool {
	currentFolder := getPathPrefix()
	files, errs := GetAllFiles([]string{path.Join(currentFolder, "AudioFiles")})
	for _, err := range errs {
		logrus.Error("[Audio][FileProvider] ", err)
	}

	audio.Title = strings.ToLower(audio.Title)
	audio.Author = strings.ToLower(audio.Author)

	for _, file := range files {
		name := strings.ToLower(file.Info.Name())
		if strings.Contains(name, audio.Title) && strings.Contains(name, audio.Author) {
			return true
		}
	}
	return false
}

func (p *fileProvider) Provide(audio project.Audio) (Player, bool) {
	currentFolder := getPathPrefix()
	files, errs := GetAllFiles([]string{path.Join(currentFolder, "AudioFiles")})
	for _, err := range errs {
		logrus.Error("[Audio][FileProvider] ", err)
	}

	audio.Title = strings.ToLower(audio.Title)
	audio.Author = strings.ToLower(audio.Author)

	for _, file := range files {
		name := strings.ToLower(file.Info.Name())
		if strings.Contains(name, audio.Title) && strings.Contains(name, audio.Author) {
			return NewFilePlayer(file.Path), true
		}
	}
	return nil, false
}

type potentialFile struct {
	Path string
	Info os.FileInfo
}

func GetFiles(basepath string) ([]potentialFile, []error) {
	var files []potentialFile
	var errs []error
	filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			errs = append(errs, err)
			if info != nil && info.IsDir() {
				return filepath.SkipDir
			}
		} else if !info.IsDir() {
			files = append(files, potentialFile{
				Path: path,
				Info: info,
			})
		}
		return nil
	})

	return files, errs
}

func GetAllFiles(paths []string) ([]potentialFile, []error) {
	var files []potentialFile
	var errs []error
	for _, basepath := range paths {
		tmp, err := GetFiles(basepath)
		if err != nil {
			errs = append(errs, err...)
		}
		files = append(files, tmp...)
	}
	return files, errs
}

type FilePlayer struct {
	*multimedia.QMediaPlayer
	onReady func()
	onError func(error)
}

func NewFilePlayer(mediapath string) *FilePlayer {
	player := FilePlayer{
		QMediaPlayer: multimedia.NewQMediaPlayer(nil, multimedia.QMediaPlayer__LowLatency),
	}

	player.SetNotifyInterval(16)
	player.ConnectMediaStatusChanged(player.MediaStatusChangedEvent)
	player.ConnectError2(player.ErrorEvent)

	player.SetFile(mediapath)

	return &player
}

func (p *FilePlayer) SetFile(file string) {
	p.SetMedia(multimedia.NewQMediaContent2(core.QUrl_FromLocalFile(file)), nil)
}

func (p *FilePlayer) Play() {
	p.QMediaPlayer.Play()
}

func (p *FilePlayer) Pause() {
	p.QMediaPlayer.Pause()
	p.QMediaPlayer.SetPosition(p.QMediaPlayer.Position())
}

func (p *FilePlayer) Time() float64 {
	return float64(p.Position()) / 1000 // milliseconds to seconds
}

func (p *FilePlayer) SetTime(value float64) {
	p.SetPosition(int64(value * 1000))
}

func (p *FilePlayer) Duration() float64 {
	return float64(p.QMediaPlayer.Duration()) / 1000 // milliseconds to seconds
}

func (p *FilePlayer) Volume() float64 {
	return float64(p.QMediaPlayer.Volume()) / 100
}

func (p *FilePlayer) SetVolume(value float64) {
	p.QMediaPlayer.SetVolume(int(value * 100))
}

func (p *FilePlayer) OnReady(f func()) {
	p.onReady = f
}

func (p *FilePlayer) OnError(f func(error)) {
	p.onError = f
}

func (p *FilePlayer) ErrorEvent(qerr multimedia.QMediaPlayer__Error) {
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
	if p.onError != nil {
		p.onError(err)
	}
}

func (p *FilePlayer) MediaStatusChangedEvent(status multimedia.QMediaPlayer__MediaStatus) {
	switch status {
	case multimedia.QMediaPlayer__UnknownMediaStatus:
		logrus.Info("[Audio] Media Status: Unknown Media Status")
	case multimedia.QMediaPlayer__NoMedia:
		logrus.Info("[Audio] Media Status: No Media")
	case multimedia.QMediaPlayer__LoadingMedia:
		logrus.Info("[Audio] Media Status: Loading Media")
	case multimedia.QMediaPlayer__LoadedMedia:
		logrus.Info("[Audio] Media Status: Loaded Media")
		if p.onReady != nil {
			p.onReady()
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
