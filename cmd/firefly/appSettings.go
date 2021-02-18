package main

import (
	"encoding/json"
	"fmt"
	"math"
	"path"
	"strings"

	"github.com/omniskop/firefly/pkg/scanner"

	"github.com/therecipe/qt/quick"

	"github.com/omniskop/firefly/cmd/firefly/settings"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

type mappingModel struct {
	core.QObject

	_         int `property:"startOffset"`
	_         int `property:"stopOffset"`
	positions []float32
	leds      []int

	_ func()             `constructor:"init"`
	_ func(float32)      `slot:"addPoint,auto"`
	_ func(int)          `slot:"deletePoint,auto"`
	_ func() int         `slot:"pointCount,auto"`
	_ func(int) float32  `slot:"getPosition,auto"`
	_ func(int, float32) `slot:"setPosition,auto"`
	_ func(int) int      `slot:"getLeds,auto"`
	_ func(int, int)     `slot:"setLeds,auto"`
}

func (m *mappingModel) init() {
	m.SetStartOffset(0)
	m.SetStopOffset(0)
	m.positions = []float32{}
	m.leds = []int{60}
}

func (m *mappingModel) addPoint(newPos float32) {
	var newPositions []float32
	var newLeds []int
	didAdd := false
	for i, p := range m.positions {
		if p > newPos {
			didAdd = true
			newPositions = append(newPositions, newPos, p)
			a := int(math.Ceil(float64(m.leds[i]) / 2))
			b := int(math.Floor(float64(m.leds[i]) / 2))
			newLeds = append(newLeds, a, b)
		} else {
			newPositions = append(newPositions, p)
			newLeds = append(newLeds, m.leds[i])
		}
	}
	if !didAdd {
		newPositions = append(newPositions, newPos)
		leds := float64(m.leds[len(m.leds)-1]) / 2
		a := int(math.Ceil(leds))
		b := int(math.Floor(leds))
		newLeds = append(newLeds, a, b)
	} else {
		newLeds = append(newLeds, m.leds[len(m.leds)-1])
	}
	m.positions = newPositions
	m.leds = newLeds
}

func (m *mappingModel) deletePoint(index int) {
	var newPositions []float32
	var newLeds []int
	var additionalLeds = 0
	for i, p := range m.positions {
		if i == index {
			additionalLeds = m.leds[i]
		} else {
			newPositions = append(newPositions, p)
			newLeds = append(newLeds, m.leds[i]+additionalLeds)
			additionalLeds = 0
		}
	}
	newLeds = append(newLeds, m.leds[len(m.leds)-1]+additionalLeds)
	m.positions = newPositions
	m.leds = newLeds
}

func (m *mappingModel) pointCount() int {
	return len(m.positions)
}

func (m *mappingModel) getPosition(i int) float32 {
	if i >= len(m.positions) {
		return 0.1
	}
	return m.positions[i]
}

func (m *mappingModel) setPosition(i int, v float32) {
	if i < len(m.positions) {
		m.positions[i] = v
	}
}

func (m *mappingModel) getLeds(i int) int {
	if i >= len(m.leds) {
		return -1
	}
	return m.leds[i]
}

func (m *mappingModel) setLeds(i int, v int) {
	if i < len(m.leds) {
		m.leds[i] = v
	}
}

func (m *mappingModel) scannerMapping() *scanner.Mapping {
	mapping := &scanner.Mapping{
		Reversed:    false,
		StartOffset: m.StartOffset(),
		EndOffset:   m.StopOffset(),
		Segments:    nil,
	}
	for i, p := range m.positions {
		mapping.AddStop(float64(p), m.leds[i])
	}
	mapping.AddStop(1, m.leds[len(m.leds)-1])
	return mapping
}

func (m *mappingModel) setScannerMapping(mapping *scanner.Mapping) {
	m.SetStartOffset(mapping.StartOffset)
	m.SetStopOffset(mapping.EndOffset)
	m.positions = make([]float32, len(mapping.Segments)-1)
	m.leds = make([]int, len(mapping.Segments))
	for i, seg := range mapping.Segments {
		if i < len(mapping.Segments)-1 {
			m.positions[i] = float32(seg.To)
		}
		m.leds[i] = seg.PixelSize
	}
}

// ====

type audioSource struct {
	core.QObject

	_ string `property:"path"`
}

type audioSourceModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ []string `property:"paths"`
	_ string   `property:"newProjectAudioCopy"`

	_ func(string)  `slot:"add,auto"`
	_ func(row int) `slot:"remove,auto"`
	_ func(int)     `slot:"setPrimary,auto"`
}

func (m *audioSourceModel) init() {
	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.SetPaths(settings.GetStrings("audio/fileSources"))
	m.SetNewProjectAudioCopy(settings.GetString("audio/newProjectAudioCopy"))
}

func (m *audioSourceModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(m.Paths()) {
		return core.NewQVariant()
	}

	var p = m.Paths()[index.Row()]

	return core.NewQVariant1(p)
}

func (m *audioSourceModel) rowCount(parent *core.QModelIndex) int {
	return len(m.Paths())
}

func (m *audioSourceModel) add(new string) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.Paths()), len(m.Paths()))
	m.SetPaths(append(m.Paths(), new))
	fmt.Println(m.Paths())
	m.EndInsertRows()
}

func (m *audioSourceModel) remove(row int) {
	if row < 0 || row >= len(m.Paths()) {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), row, row)
	m.SetPaths(append(m.Paths()[:row], m.Paths()[row+1:]...))
	m.EndRemoveRows()
}

func (m *audioSourceModel) setPrimary(row int) {
	if row <= 0 || row >= len(m.Paths()) {
		// don't do anything if the index is out of range or 0
		return
	}
	m.BeginMoveRows(core.NewQModelIndex(), row, row, core.NewQModelIndex(), 0)
	paths := m.Paths()
	newPaths := append([]string{paths[row]}, paths[:row]...)
	newPaths = append(newPaths, paths[row+1:]...)
	m.SetPaths(newPaths)
	m.EndMoveRows()
}

// ====

type appSettingsModel struct {
	core.QObject

	_ audioSourceModel `property:"audioSources"`

	_ string `property:"editorPasteMode"`

	_ bool   `property:"liveLedStripEnabled"`
	_ string `property:"liveLedStripAddress"`
	_ int    `property:"liveLedStripPort"`

	_ int          `property:"liveLedStripMappingMode"` // 0 = simple/linear; 1 = custom
	_ int          `property:"ledCount"`
	_ mappingModel `property:"mapping"`

	_ func() `constructor:"init"`
	_ func() `slot:"ok"`
	_ func() `slot:"cancel"`
}

func (m *appSettingsModel) init() {
	m.SetAudioSources(NewAudioSourceModel(m))

	m.SetEditorPasteMode(settings.GetString("editor/pasteMode"))
	m.SetLiveLedStripEnabled(settings.GetBool("liveLedStrip/enabled"))
	m.SetLiveLedStripAddress(settings.GetString("liveLedStrip/address"))
	m.SetLiveLedStripPort(settings.GetInt("liveLedStrip/port"))

	var mapping *scanner.Mapping
	err := json.Unmarshal([]byte(settings.GetString("liveLedStrip/mapping")), &mapping)
	if err != nil {
		mapping = scanner.NewLinearMapping(30)
	}
	m.SetMapping(NewMappingModel(m))
	m.Mapping().setScannerMapping(mapping)
	if mappingIsLinear(mapping) {
		m.SetLiveLedStripMappingMode(0)
		m.SetLedCount(mapping.Segments[0].PixelSize)
	} else {
		m.SetLiveLedStripMappingMode(1)
		m.SetLedCount(0)
	}
}

func mappingIsLinear(m *scanner.Mapping) bool {
	return m.StartOffset == 0 && m.EndOffset == 0 && m.Reversed == false && len(m.Segments) == 1
}

func (m *appSettingsModel) save() {
	settings.Set("audio/fileSources", m.AudioSources().Paths())
	fmt.Println(m.AudioSources().NewProjectAudioCopy())
	settings.Set("audio/newProjectAudioCopy", m.AudioSources().NewProjectAudioCopy())

	settings.Set("editor/pasteMode", m.EditorPasteMode())
	settings.Set("liveLedStrip/enabled", m.IsLiveLedStripEnabled())
	settings.Set("liveLedStrip/address", m.LiveLedStripAddress())
	settings.Set("liveLedStrip/port", m.LiveLedStripPort())
	if m.LiveLedStripMappingMode() == 0 {
		data, _ := json.Marshal(scanner.NewLinearMapping(m.LedCount()))
		settings.Set("liveLedStrip/mapping", string(data))
	} else {
		data, _ := json.Marshal(m.Mapping().scannerMapping())
		settings.Set("liveLedStrip/mapping", string(data))
	}
}

func NewAppSettingsWindow() error {
	view := quick.NewQQuickView(nil)

	// setup model
	model := NewAppSettingsModel(nil)
	model.ConnectOk(func() {
		model.save()
		view.Close()
	})
	model.ConnectCancel(func() {
		view.Close()
	})
	view.RootContext().SetContextProperty("Model", model)
	view.Engine().AddImportPath(":/assets/qml/Styles")
	view.Engine().AddImportPath(":/assets/qml/Styles/Dark")

	// setup view
	view.SetSource(core.NewQUrl3("qrc:/assets/qml/AppSettingsWindow.qml", core.QUrl__TolerantMode))
	if view.Status() == quick.QQuickView__Error {
		errs := view.Errors()
		out := make([]string, len(errs))
		for i, err := range errs {
			out[i] = err.ToString()
		}
		return fmt.Errorf("settings view could not be created: \r\n%s", strings.Join(out, "\r\n"))
	}
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	view.SetMinimumSize(core.NewQSize2(650, 470))
	view.ConnectKeyPressEvent(func(ev *gui.QKeyEvent) {
		if core.Qt__Key(ev.Key()) == core.Qt__Key_R && ev.Modifiers() == core.Qt__ControlModifier {
			fmt.Println("reload")
			view.Engine().ClearComponentCache()
			view.SetSource(core.NewQUrl3("./cmd/firefly/assets/qml/AppSettingsWindow.qml", core.QUrl__TolerantMode))
			return
		}
		view.KeyPressEventDefault(ev)
	})
	view.Show()

	return nil
}

func restoreDefaultSettings() {
	musicLocation := core.QStandardPaths_WritableLocation(core.QStandardPaths__MusicLocation)
	if musicLocation == "" {
		settings.Set("audio/fileSources", []string{})
	} else {
		settings.Set("audio/fileSources", []string{path.Join(musicLocation, "Firefly Audio Files")})
	}
	settings.Set("editor/pasteMode", "auto")
	settings.Set("liveLedStrip/enabled", false)
	settings.Set("liveLedStrip/address", "127.0.0.1")
	settings.Set("liveLedStrip/port", "20202")
	mapping, _ := json.Marshal(scanner.NewLinearMapping(30))
	settings.Set("liveLedStrip/mapping", string(mapping))
}
