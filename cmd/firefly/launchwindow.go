package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/core"

	"github.com/therecipe/qt/widgets"
)

var launchWindow *widgets.QWidget

func OpenLaunchwindow() {
	if launchWindow == nil {
		var err error
		launchWindow, err = NewLaunchWindow()
		if err != nil {
			logrus.Errorf("unable to create launch window: %v", err)
		}
	}
	launchWindow.Show()
}

func NewLaunchWindow() (*widgets.QWidget, error) {

	formWidget, err := loadUI(":assets/ui/launchwindow.ui", nil)
	if err != nil {
		return nil, err
	}

	formWidget.SetWindowTitle("Firefly")
	formWidget.SetFixedSize(formWidget.Geometry().Size())

	newProjectButton := widgets.NewQPushButtonFromPointer(formWidget.FindChild("newProjectButton", core.Qt__FindChildrenRecursively).Pointer())
	openProjectButton := widgets.NewQPushButtonFromPointer(formWidget.FindChild("openProjectButton", core.Qt__FindChildrenRecursively).Pointer())
	projectList := widgets.NewQListWidgetFromPointer(formWidget.FindChild("projectList", core.Qt__FindChildrenRecursively).Pointer())

	newProjectButton.ConnectClicked(func(bool) {
		dialog, err := NewProjectSetupWindow(formWidget)
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
	})

	openProjectButton.ConnectClicked(func(bool) {
		err := openProject()
		if err == errUserAbort {
			return
		} else if err != nil {
			msgBox := widgets.NewQMessageBox2(widgets.QMessageBox__NoIcon, "Open Project", fmt.Sprintf("The project could not be opened.\n%v", err), widgets.QMessageBox__Ok, formWidget, core.Qt__Dialog)
			msgBox.Exec()
			return
		}
		launchWindow.Hide()
	})

	projectList.AddItem("No Recent Projects")

	return formWidget, nil
}
