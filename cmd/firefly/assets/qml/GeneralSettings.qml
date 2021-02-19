import QtQuick 2.12
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.3
import Dark 1.0
import QtQuick.Dialogs 1.3
//import Qt.labs.platform 1.1

ColumnLayout {
    id: columnLayout

    GroupBox {
        title: qsTr("Audio Sources")
        Layout.fillWidth: true

        ColumnLayout {
            width: parent.width

            Label {
                text: "Firefly will search these locations (including subfolders) for audio files when opening a project:"
            }

            ListView {
                id: audioSourceList
                Layout.fillWidth: true
                Layout.preferredHeight: 150
                clip: true
                highlightMoveDuration: 0  // disables animation...
                highlightMoveVelocity: -1 // ...of highlight
                ScrollBar.vertical: ScrollBar { }
                focus: true

                model: Model.audioSources

                delegate: Text {
                    text: {
                        if (Model.audioSources.newProjectAudioCopy == "audioSources" && index == 0) {
                            // if we use audio sources as a copy target and we are at index zero
                            // add a star to indicate that this is the primary source
                            return "★ " + model.display
                        } else {
                            return model.display
                        }
                    }
                    color: parent.ListView.isCurrentItem ? Shared.lighterTextColor : Shared.textColor
                    width: parent.width
                    padding: 3
                    elide: Text.ElideMiddle
                    MouseArea {
                        anchors.fill: parent
                        onClicked: parent.ListView.view.currentIndex = index
                    }
                }

                highlight:  Rectangle {
                    width: parent == null ? 0 : parent.width
                    color: Qt.darker(Shared.interactableColor, 1.2)
                    radius: 4
                }
            } // ListView

            Row {
                spacing: 5

                Button {
                    id: removeAudioSourceButton
                    small: true; square: true
                    text: "-"

                    onClicked: Model.audioSources.remove(audioSourceList.currentIndex)
                }

                Button {
                    id: addAudioSourceButton
                    small: true; square: true
                    text: "+"
                    onClicked: audioSourceFileDialog.visible = true
                }

                Button {
                    id: primaryAudioSourceButton
                    small: true
                    text: qsTr("Primary")
                    onClicked: Model.audioSources.setPrimary(audioSourceList.currentIndex)
                    visible: Model.audioSources.newProjectAudioCopy == "audioSources" // only show if needed
                }
            } // Row

            FileDialog {
                id: audioSourceFileDialog
                title: qsTr("Choose a directory to be searched for audio files.")
                selectFolder: true
                folder: shortcuts.home
                onAccepted: {
                    // get rid of "file://"
                    var path = String(audioSourceFileDialog.fileUrl).slice(7);
                    Model.audioSources.add(path)
                }
            }

            Label {
                text: "Copy the audio file of new projects:"
            }

            ComboBox {
                model: ["Into the primary audio location", "Next to the project file"]
                readonly property var settingsEnum: ["audioSources", "projectFile"]
                currentIndex: {
                    var index = settingsEnum.indexOf(Model.audioSources.newProjectAudioCopy)
                    if(index == -1) {
                        console.log("unknown value for setting regarding the copy behaviour of audio files for new projects")
                        enabled = false
                        model = ["Option disabled due to an error"]
                        return 0
                    }
                    return index
                }
                onActivated: {
                    Model.audioSources.newProjectAudioCopy = settingsEnum[index]
                }
            }
        } // ColumnLayout
    } // GroupBox
} // ColumnLayout