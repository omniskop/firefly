package main

import (
	"strconv"

	"github.com/omniskop/firefly/cmd/firefly/settings"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func NewAppSettingsWindow(parent *widgets.QWidget) (*widgets.QDialog, error) {
	window, err := loadUI(":assets/ui/appSettings.ui", parent)
	if err != nil {
		return nil, err
	}

	window.SetFixedSize(window.Geometry().Size())
	dialog := widgets.NewQDialogFromPointer(window.Pointer())

	dialog.ConnectCloseEvent(func(*gui.QCloseEvent) {
		dialog.SetResult(int(widgets.QDialog__Rejected))
	})

	ledStripEnabled := widgets.NewQCheckBoxFromPointer(window.FindChild("ledStripEnabled", core.Qt__FindChildrenRecursively).Pointer())
	ledCount := widgets.NewQLineEditFromPointer(window.FindChild("ledCount", core.Qt__FindChildrenRecursively).Pointer())
	ledStripAddress := widgets.NewQLineEditFromPointer(window.FindChild("ledStripAddress", core.Qt__FindChildrenRecursively).Pointer())
	ledStripPort := widgets.NewQLineEditFromPointer(window.FindChild("ledStripPort", core.Qt__FindChildrenRecursively).Pointer())
	//okButton := widgets.NewQDialogButtonBoxFromPointer(window.FindChild("dialogButtons", core.Qt__FindChildrenRecursively).Pointer()).Button(widgets.QDialogButtonBox__Ok)

	ledStripEnabled.SetChecked(settings.GetBool("liveLedStrip/enabled"))
	ledCount.SetText(strconv.Itoa(settings.GetInt("ledCount")))
	ledStripAddress.SetText(settings.GetString("liveLedStrip/address"))
	ledStripPort.SetText(strconv.Itoa(settings.GetInt("liveLedStrip/port")))

	dialog.ConnectAccepted(func() {
		settings.Set("liveLedStrip/enabled", ledStripEnabled.IsChecked())
		v, err := strconv.Atoi(ledCount.Text())
		if err == nil {
			settings.Set("ledCount", v)
		}
		settings.Set("liveLedStrip/address", ledStripAddress.Text())
		v, err = strconv.Atoi(ledStripPort.Text())
		if err == nil {
			settings.Set("liveLedStrip/port", v)
		}
	})

	window.Show()

	return dialog, nil
}

func restoreDefaultSettings() {
	settings.Set("ledCount", 30)
	settings.Set("liveLedStrip/enabled", false)
	settings.Set("liveLedStrip/address", "127.0.0.1")
	settings.Set("liveLedStrip/port", "20202")
}
