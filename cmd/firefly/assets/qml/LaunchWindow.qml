import QtQuick 2.12
import QtQuick.Controls 2.3
import QtQuick.Layouts 1.3
import Dark 1.0

Rectangle {
    id: root
    color: Shared.backgroundColor

    Item {
        id: leftSection
        width: parent.width / 2
        anchors.bottom: parent.bottom
        anchors.bottomMargin: 0
        anchors.top: parent.top
        anchors.topMargin: 0
        anchors.left: parent.left
        anchors.leftMargin: 0

        Image {
            id: logo
            y: 20
            width: parent.width * 0.6
            height: width
            anchors.horizontalCenter: parent.horizontalCenter
            fillMode: Image.PreserveAspectFit
            source: "qrc:/assets/images/logo.png"
        }

        Label {
            id: version
            text: Model.version
            width: logo.width
            horizontalAlignment: Text.AlignHCenter
            anchors.top: logo.bottom
            anchors.horizontalCenter: parent.horizontalCenter
        }

        Button {
            id: newProject
            width: logo.width * 0.9
            height: 32
            text: qsTr("New Project ...")
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: version.bottom
            anchors.topMargin: 20

            onClicked: {
                Model.newProject()
            }
        }

        Button {
            id: openProject
            text: qsTr("Open Project ...")
            width: logo.width * 0.9
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.top: newProject.bottom
            anchors.topMargin: 20

            onClicked: {
                Model.openProject()
            }
        }

        Button {
            id: button
            y: 669
            text: qsTr("Settings")
            anchors.bottom: parent.bottom
            anchors.bottomMargin: 20
            anchors.left: parent.left
            anchors.leftMargin: 20

            onClicked: {
                Model.openSettings()
            }
        }
    }

    Rectangle {
        id: separator
        width: 1
        height: parent.height
        color: Shared.interactableColor
        anchors.left: leftSection.right
        anchors.leftMargin: 0
    }

    Item {
        id: rightSection
        width: parent.width / 2
        anchors.top: parent.top
        anchors.topMargin: 0
        anchors.bottom: parent.bottom
        anchors.bottomMargin: 0
        anchors.right: parent.right
        anchors.rightMargin: 0

        GroupBox {
            id: recentProjects
            anchors.fill: parent
            title: qsTr("Recent Projects")

            ListView {
                id: listView
                boundsBehavior: Flickable.StopAtBounds
                anchors.fill: parent
                delegate: Rectangle {
                    id: listItem
                    width: parent.width
                    height: 40
                    color: "transparent"
                    radius: 4

                    property color fontColor: Shared.textColor

                    MouseArea {
                        anchors.fill: parent
                        hoverEnabled: true
                        onEntered: {
                            listItem.color = Shared.interactableColor
                            listItem.fontColor = Shared.lighterTextColor
                        }
                        onExited: {
                            listItem.color = "transparent"
                            listItem.fontColor = Shared.textColor
                        }
                        onClicked: {
                            Model.openProjectPath(model.display.path)
                        }
                    }

                    Column {
                        id: row1
                        width: parent.width
                        padding: 4

                        Text {
                            id: major
                            width: parent.width - parent.padding*3
                            elide: Text.ElideRight
                            text: model.display.songTitle + " - " + model.display.songAuthor
                            font.bold: true
                            color: listItem.fontColor
                        }

                        Text {
                            width: parent.width - parent.padding*3
                            elide: Text.ElideLeft
                            text: model.display.path
                            color: listItem.fontColor
                        }
                    }
                }
                model: Model.recentFiles
            }
        }
    }
}





























/*##^## Designer {
    D{i:4;anchors_x:116}D{i:5;anchors_x:20}D{i:1;anchors_height:720}
}
 ##^##*/
