package main

import (
	"fmt"
	"math"
	"strings"

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
			a := int(math.Ceil(float64(m.leds[i] / 2)))
			b := int(math.Floor(float64(m.leds[i] / 2)))
			newLeds = append(newLeds, a, b)
		} else {
			newPositions = append(newPositions, p)
			newLeds = append(newLeds, m.leds[i])
		}
	}
	if !didAdd {
		newPositions = append(newPositions, newPos)
		leds := float64(m.leds[len(m.leds)-1] / 2)
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
	newLeds = append(newLeds, m.leds[len(m.leds)-1])
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

type appSettingsModel struct {
	core.QObject

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
	m.load()

	m.SetLiveLedStripMappingMode(0)
	m.SetMapping(NewMappingModel(nil))
}

func (m *appSettingsModel) load() {
	m.SetLiveLedStripEnabled(settings.GetBool("liveLedStrip/enabled"))
	m.SetLedCount(settings.GetInt("ledCount"))
	m.SetLiveLedStripAddress(settings.GetString("liveLedStrip/address"))
	m.SetLiveLedStripPort(settings.GetInt("liveLedStrip/port"))
}

func (m *appSettingsModel) save() {
	settings.Set("liveLedStrip/enabled", m.IsLiveLedStripEnabled())
	settings.Set("ledCount", m.LedCount())
	settings.Set("liveLedStrip/address", m.LiveLedStripAddress())
	settings.Set("liveLedStrip/port", m.LiveLedStripPort())
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
	view.SetMinimumSize(core.NewQSize2(300, 360))
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
	settings.Set("ledCount", 30)
	settings.Set("liveLedStrip/enabled", false)
	settings.Set("liveLedStrip/address", "127.0.0.1")
	settings.Set("liveLedStrip/port", "20202")
}
