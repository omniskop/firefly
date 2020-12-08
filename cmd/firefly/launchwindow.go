package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/omniskop/firefly/cmd/firefly/settings"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/quick"

	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"

	"github.com/therecipe/qt/widgets"
)

var launchWindow *quick.QQuickView

// launchWindowModel is the model for the LaunchWindow
type launchWindowModel struct {
	core.QObject

	recentFilesData []*recentFile

	_ *core.QAbstractListModel `property:"recentFiles"`
	_ string                   `property:"version"`

	_ func()       `constructor:"init"`
	_ func()       `signal:"newProject,auto"`
	_ func()       `signal:"openProject,auto"`
	_ func()       `signal:"openSettings,auto"`
	_ func(string) `slot:"openProjectPath,auto"`
}

func (m *launchWindowModel) init() {
	m.recentFilesData = getRecentFiles()
	recentFiles := core.NewQAbstractListModel(nil)
	recentFiles.ConnectData(m.getRecentFile)
	recentFiles.ConnectRowCount(m.getRecentFilesCount)
	m.SetRecentFiles(recentFiles)
	m.SetVersion(settings.GetString("version"))
}

// newProject opens the project creation dialog
func (m *launchWindowModel) newProject() {
	dialog, err := NewProjectSetupWindow(nil)
	if err != nil {
		logrus.WithField("err", err).Error("unable to create project setup window")
		return
	}
	dialog.ConnectFinished(func(result int) {
		if result == int(widgets.QDialog__Rejected) {
			launchWindow.Show()
		} else {
			launchWindow.Hide()
		}
	})
}

// openProject asks the user to open a project file
func (m *launchWindowModel) openProject() {
	err := openProject()
	if err == errUserAbort {
		return
	} else if err != nil {
		msgBox := widgets.NewQMessageBox2(widgets.QMessageBox__NoIcon, "Open Project", "The project could not be opened", widgets.QMessageBox__Ok, nil, core.Qt__Dialog)
		msgBox.SetInformativeText(err.Error())
		msgBox.Exec()
		return
	}
	launchWindow.Hide()
}

// openSettings opens the settings window
func (m *launchWindowModel) openSettings() {
	err := NewAppSettingsWindow()
	if err != nil {
		logrus.WithField("err", err).Error("unable to create app settings window")
		return
	}
}

// openProjectPath opens the project at the given path
func (m *launchWindowModel) openProjectPath(path string) {
	err := openProjectPath(path)
	if err != nil {
		removeRecentFile(path)
		launchWindow.RootContext().SetContextProperty("Model", NewLaunchWindowModel(nil)) // update the model
		msgBox := widgets.NewQMessageBox2(widgets.QMessageBox__NoIcon, "Open Project", "The project could not be opened", widgets.QMessageBox__Ok, nil, core.Qt__Dialog)
		msgBox.SetInformativeText(err.Error())
		msgBox.Exec()
	} else {
		launchWindow.Hide()
	}
}

// getRecentFilesCount implements part of the QAbtractListModel
func (m *launchWindowModel) getRecentFilesCount(*core.QModelIndex) int {
	return len(m.recentFilesData)
}

// getRecentFile implements part of the QAbtractListModel
func (m *launchWindowModel) getRecentFile(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() || index.Row() >= len(m.recentFilesData) {
		return core.NewQVariant()
	}
	i := len(m.recentFilesData) - index.Row() - 1 // reverse the index to show the newest file first
	return m.recentFilesData[i].ToVariant()
}

// OpenLaunchWindow opens the launch window. It will reuse a previously created window if available.
func OpenLaunchWindow() {
	if launchWindow == nil {
		var err error
		launchWindow, err = newLaunchWindow()
		if err != nil {
			logrus.Errorf("unable to create launch window: %v", err)
		}
	}
	launchWindow.Show()
}

// newLaunchWindow creates a new launch window without showing it. For opening a new window OpenLaunchWindow should be used.
func newLaunchWindow() (*quick.QQuickView, error) {
	view := quick.NewQQuickView(nil)

	// setup model
	model := NewLaunchWindowModel(nil)

	// setup view
	view.RootContext().SetContextProperty("Model", model)
	view.Engine().AddImportPath(":/assets/qml/Styles")
	view.SetSource(core.NewQUrl3("qrc:/assets/qml/LaunchWindow.qml", core.QUrl__TolerantMode))
	if view.Status() == quick.QQuickView__Error {
		errs := view.Errors()
		out := make([]string, len(errs))
		for i, err := range errs {
			out[i] = err.ToString()
		}
		return nil, fmt.Errorf("launch window could not be created: \r\n%s", strings.Join(out, "\r\n"))
	}
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	view.SetMaximumSize(core.NewQSize2(700, 450))
	view.SetMinimumSize(core.NewQSize2(700, 450))
	view.ConnectKeyPressEvent(func(ev *gui.QKeyEvent) {
		if core.Qt__Key(ev.Key()) == core.Qt__Key_R && ev.Modifiers() == core.Qt__ControlModifier {
			fmt.Println("reload")
			view.Engine().ClearComponentCache()
			view.SetSource(core.NewQUrl3("./cmd/firefly/assets/qml/LaunchWindow.qml", core.QUrl__TolerantMode))
			return
		}
		view.KeyPressEventDefault(ev)
	})

	return view, nil
}

// === recent files ===

// recentFile contains information about a project file
type recentFile struct {
	core.QObject

	_ string `property:"songTitle"`
	_ string `property:"songAuthor"`
	_ string `property:"path"`
	_ uint64 `property:"lastOpened"` // unix timestamp
}

// newRecentFileDetailed creates a new recentFile object with the given information
func newRecentFileDetailed(title, author, path string) *recentFile {
	rf := NewRecentFile(nil)
	rf.SetSongTitle(title)
	rf.SetSongAuthor(author)
	rf.SetPath(path)
	rf.SetLastOpened(uint64(time.Now().Unix()))
	return rf
}

// serialzeRecentFiles takes a list of files and serializes them in a string
func serializeRecentFiles(files []*recentFile) string {
	var data = make([]map[string]interface{}, len(files))
	for i, rf := range files {
		data[i] = map[string]interface{}{
			"songTitle":  rf.SongTitle(),
			"songAuthor": rf.SongAuthor(),
			"path":       rf.Path(),
			"lastOpened": rf.LastOpened(),
		}
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(encoded)
}

// deserializeRecentFiles parses a stored list of files
func deserializeRecentFiles(input string) (rfiles []*recentFile) {
	var data []struct {
		SongTitle  string
		SongAuthor string
		Path       string
		LastOpened uint64
	}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	rfiles = make([]*recentFile, len(data))
	for i, d := range data {
		rf := NewRecentFile(nil)
		rf.SetSongTitle(d.SongTitle)
		rf.SetSongAuthor(d.SongAuthor)
		rf.SetPath(d.Path)
		rf.SetLastOpened(d.LastOpened)
		rfiles[i] = rf
	}
	return
}

// getRecentFiles loads the recently used files from the settings
func getRecentFiles() []*recentFile {
	raw, ok := settings.Get("recentFiles").(string)
	if !ok {
		return nil
	}
	return deserializeRecentFiles(raw)
}

// saveRecentFiles stores the file list in the settings
func saveRecentFiles(files []*recentFile) {
	encoded := serializeRecentFiles(files)
	settings.Set("recentFiles", core.NewQVariant12(encoded))
}

// addRecentFile adds the file to the list of recently used files
func addRecentFile(newRF *recentFile) {
	files := getRecentFiles()
	for i, f := range files {
		if f.Path() == newRF.Path() {
			copy(files[i:], files[i+1:])
			files[len(files)-1] = nil
			files = files[:len(files)-1]
		}
	}
	files = append(files, newRF)
	saveRecentFiles(files)
}

// removeRecentFile removes the file with the given path from the list of recently used files
func removeRecentFile(filePath string) {
	files := getRecentFiles()
	for i, f := range files {
		if f.Path() == filePath {
			copy(files[i:], files[i+1:])
			files[len(files)-1] = nil
			files = files[:len(files)-1]
		}
	}
	saveRecentFiles(files)
}
