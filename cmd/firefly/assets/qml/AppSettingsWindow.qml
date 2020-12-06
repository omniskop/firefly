import QtQuick 2.12
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.3
import QtQuick.Window 2.2
import Dark 1.0

Rectangle {
    id: root

//    width: 350
//    height: childrenRect.height
    color: Shared.backgroundColor

    property bool saveOnClose: false

    ColumnLayout {
        id: columnLayout
        width: parent.width
        spacing: 20

        GroupBox {
            id: groupBox
            title: "Live Preview"
            Layout.fillWidth: true

            ColumnLayout {
                width: parent.width

                GridLayout {
                    columns: 2
                    columnSpacing: 10

                    Label {
                        id: spacer2
                        text: "Enabled"
                        Layout.alignment: Qt.AlignRight | Qt.AlignVCenter
                    }

                    CheckBox {
                        id: checkBox
                        onClicked: Model.liveLedStripEnabled = checked
                        Component.onCompleted: checked = Model.liveLedStripEnabled
                        Layout.fillWidth: true
                    }

                    Label {
                        id: label3
                        text: qsTr("Address")
                        Layout.alignment: Qt.AlignRight | Qt.AlignVCenter
                    }

                    TextField {
                        id: textField1
                        text: Model.liveLedStripAddress
                        placeholderText: "192.168.178.0"
                        validator: RegExpValidator { regExp: /\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/ }
                        onTextChanged: Model.liveLedStripAddress = text
                        Layout.fillWidth: true
                    }

                    Label {
                        id: label4
                        text: qsTr("Port")
                        Layout.alignment: Qt.AlignRight | Qt.AlignVCenter
                    }

                    TextField {
                        id: textField2
                        text: Model.liveLedStripPort
                        placeholderText: "20202"
                        validator: IntValidator {bottom: 0; top: 65535}
                        onTextChanged: Model.liveLedStripPort = parseInt(text)
                        Layout.fillWidth: true
                    }

                    Label {
                        text: qsTr("Mapping")
                        Layout.alignment: Qt.AlignRight | Qt.AlignVCenter
                    }

                    ComboBox {
                        model: [qsTr("Simple")/*, qsTr("Custom")*/]
                        Layout.fillWidth: true

                        Component.onCompleted: {
                            activated(currentIndex)
                        }

                        onActivated: {
                            if(index == 0) {
                                pixelsLabel.visible = true
                                pixelsInput.visible = true
                                mappingEditorLabel.visible = false
                                mappingEditor.visible = false
                            } else if(index == 1) {
                                pixelsLabel.visible = false
                                pixelsInput.visible = false
                                mappingEditorLabel.visible = true
                                mappingEditor.visible = true
                            }
                        }    
                    }

                    Label {
                        id: pixelsLabel
                        text: qsTr("Pixels")
                        Layout.alignment: Qt.AlignRight | Qt.AlignVCenter
                        horizontalAlignment: Text.AlignLeft
                    }

                    TextField {
                        id: pixelsInput
                        text: Model.ledCount
                        placeholderText: "60"
                        validator: IntValidator {bottom: 1; top: 1000000}
                        onTextChanged: Model.ledCount = text == "" ? 0 : parseInt(text)
                        Layout.fillWidth: true
                    }

                } // PropertyGridLayout

                Label {
                    id: mappingEditorLabel
                    text: "Custom Mapping"
                    font.weight: Font.DemiBold
                    topPadding: 10
                }

                MappingEditor {
                    id: mappingEditor
                    Layout.fillWidth: true
                }
            } // ColumnLayout
        } // GroupBox

        Row {
            Layout.alignment: Qt.AlignRight | Qt.AlignVCenter
            spacing: 10
            rightPadding: 10

            Button {
                id: "cancelButton"
                text: qsTr("Cancel")
                onClicked: Model.cancel()
            }

            Button {
                id: "okButton"
                text: qsTr("OK")

                onClicked: Model.ok()
            }
        }
    }
}






/*##^## Designer {
    D{i:0;width:400}D{i:19;invisible:true}
}
 ##^##*/
