package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/omniskop/firefly/cmd/firefly/settings"

	"github.com/omniskop/firefly/cmd/firefly/audio"

	"github.com/sirupsen/logrus"

	"github.com/omniskop/firefly/cmd/firefly/editor"
	"github.com/omniskop/firefly/pkg/project"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func NewProjectSetupWindow(parent *widgets.QWidget) (*widgets.QDialog, error) {
	window, err := loadUI(":assets/ui/projectSetup.ui", parent)
	if err != nil {
		return nil, err
	}

	window.SetFixedSize(window.Geometry().Size())
	dialog := widgets.NewQDialogFromPointer(window.Pointer())

	dialog.ConnectCloseEvent(func(*gui.QCloseEvent) {
		dialog.SetResult(int(widgets.QDialog__Rejected))
	})

	audioFileButton := widgets.NewQPushButtonFromPointer(window.FindChild("audioFileButton", core.Qt__FindChildrenRecursively).Pointer())
	audioFileName := widgets.NewQLabelFromPointer(window.FindChild("audioFileName", core.Qt__FindChildrenRecursively).Pointer())
	interpretText := widgets.NewQLineEditFromPointer(window.FindChild("interpretText", core.Qt__FindChildrenRecursively).Pointer())
	titleText := widgets.NewQLineEditFromPointer(window.FindChild("titleText", core.Qt__FindChildrenRecursively).Pointer())
	okButton := widgets.NewQDialogButtonBoxFromPointer(window.FindChild("buttonBox", core.Qt__FindChildrenRecursively).Pointer()).Button(widgets.QDialogButtonBox__Ok)

	var selectedAudioFile string

	validate := func() {
		okButton.SetEnabled(interpretText.Text() != "" && titleText.Text() != "" && selectedAudioFile != "")
	}
	validate()

	audioFileButton.ConnectClicked(func(bool) {
		selectedAudioFile = widgets.QFileDialog_GetOpenFileName(window, "Choose Audio File", ".", "Audio Files (*.mp3 *.wav)", "", 0)
		if selectedAudioFile == "" {
			audioFileName.SetText("No File")
		} else {
			audioFileName.SetText(path.Base(selectedAudioFile))
			validate()
		}
	})

	interpretText.ConnectTextChanged(func(string) {
		validate()
	})

	titleText.ConnectTextChanged(func(string) {
		validate()
	})

	dialog.ConnectAccepted(func() {
		if interpretText.Text() == "" || titleText.Text() == "" || selectedAudioFile == "" {
			return
		}
		err := createProject(interpretText.Text(), titleText.Text(), selectedAudioFile)
		if err != nil {
			logrus.Error(err)
			if errors.Is(err, audio.NoProviderErr) {
				err = fmt.Errorf("The audio file could not be opened.")
			}
			widgets.NewQMessageBox2(widgets.QMessageBox__Warning, "Create Project", fmt.Sprintf("The project could not be created.\n%v", err), widgets.QMessageBox__Ok, nil, core.Qt__Dialog).Exec()
		}
	})

	window.Show()

	return dialog, nil
}

func createProject(interpretText string, titleText string, selectedAudioFile string) error {
	audioFileSources := settings.GetStrings("audio/fileSources")
	var audioFilePath string
	if len(audioFileSources) > 0 && settings.GetString("audio/newProjectAudioCopy") == "audioSources" {
		audioFolder := audioFileSources[0]
		err := os.MkdirAll(audioFolder, 0755|os.ModeDir)
		if err != nil {
			return fmt.Errorf("create project: create audio folder: %w", err)
		}

		audioFilePath = path.Join(audioFolder, fmt.Sprintf("%s - %s%s", titleText, interpretText, path.Ext(selectedAudioFile)))

		err = copyFile(audioFilePath, selectedAudioFile)
		if err != nil {
			return fmt.Errorf("create project: copy audio file: %w", err)
		}
	} else {
		audioFilePath = selectedAudioFile
	}

	projAudio := project.Audio{
		Title:  titleText,
		Author: interpretText,
		Genres: nil,
		File:   nil,
	}

	// load a player to make sure that the file is valid
	player := audio.NewFilePlayer(audioFilePath)

	player.OnReady(func() {
		editor.New(&project.Project{
			Title:          fmt.Sprintf("%s - %s", titleText, interpretText),
			Author:         "Firefly Default Author",
			Tags:           nil,
			AdditionalInfo: nil,
			Duration:       player.Duration(),
			Scene:          project.Scene{},
			Audio:          projAudio,
			AudioOffset:    0,
		}, editor.Options{
			AudioLocation:   audioFilePath,
			CopyAudioOnSave: settings.GetString("audio/newProjectAudioCopy") == "projectFile",
		}, ApplicationCallbacks)
	})

	player.OnError(func(err error) {
		widgets.NewQMessageBox2(widgets.QMessageBox__Warning, "Create Project", fmt.Sprintf("The file could not be loaded.\n%v", err), widgets.QMessageBox__Ok, nil, core.Qt__Dialog).Exec()
	})

	return nil
}

func copyFile(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil && !os.IsExist(err) {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
