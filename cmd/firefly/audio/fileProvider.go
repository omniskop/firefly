package audio

import (
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

type fileProvider struct{}

func (p *fileProvider) CanProvide(audio project.Audio) bool {
	currentFolder, err := os.Getwd()
	if err != nil {
		currentFolder = "/"
	}
	files, errs := GetAllFiles([]string{path.Join(currentFolder, "assets/audio")})
	for _, err := range errs {
		logrus.Error("[Audio][FileProvider] ", err)
	}

	audio.Title = strings.ToLower(audio.Title)
	audio.Author = strings.ToLower(audio.Author)

	for _, file := range files {
		name := strings.ToLower(file.info.Name())
		if strings.Contains(name, audio.Title) && strings.Contains(name, audio.Author) {
			return true
		}
	}
	return false
}

func (p *fileProvider) Provide(audio project.Audio) (Player, bool) {
	currentFolder, err := os.Getwd()
	if err != nil {
		currentFolder = "/"
	}
	files, errs := GetAllFiles([]string{path.Join(currentFolder, "assets/audio")})
	for _, err := range errs {
		logrus.Error("[Audio][FileProvider] ", err)
	}

	audio.Title = strings.ToLower(audio.Title)
	audio.Author = strings.ToLower(audio.Author)

	for _, file := range files {
		name := strings.ToLower(file.info.Name())
		if strings.Contains(name, audio.Title) && strings.Contains(name, audio.Author) {
			return NewFilePlayer(file.path), true
		}
	}
	return nil, false
}

type potentialFile struct {
	path string
	info os.FileInfo
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
				path: path,
				info: info,
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
}

func NewFilePlayer(mediapath string) *FilePlayer {
	player := FilePlayer{
		QMediaPlayer: multimedia.NewQMediaPlayer(nil, 0),
	}

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

func (p *FilePlayer) ErrorEvent(err multimedia.QMediaPlayer__Error) {
	switch err {
	case multimedia.QMediaPlayer__NoError:
		logrus.Error("[Audio] No Error")
	case multimedia.QMediaPlayer__ResourceError:
		logrus.Error("[Audio] Resource Error")
	case multimedia.QMediaPlayer__FormatError:
		logrus.Error("[Audio] Format Error")
	case multimedia.QMediaPlayer__NetworkError:
		logrus.Error("[Audio] Network Error")
	case multimedia.QMediaPlayer__AccessDeniedError:
		logrus.Error("[Audio] Access Denied Error")
	case multimedia.QMediaPlayer__ServiceMissingError:
		logrus.Error("[Audio] Service Missing Error")
	default:
		logrus.WithField("error", err).Error("[Audio] Unknown Error")
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
