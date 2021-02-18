import QtQuick 2.12
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.3
import Dark 1.0

ColumnLayout {
    id: columnLayout

    GroupBox {
        title: "Editor"
        Layout.fillWidth: true

        ColumnLayout {
            width: parent.width

            GridLayout {
                columns: 2
                columnSpacing: 10

                Label {
                    text: "Paste Elements:"
                }

                ComboBox {
                    model: [qsTr("Auto"), qsTr("at Mouse"), qsTr("at Needle")]
                    Layout.fillWidth: true
                    currentIndex: ["auto", "mouse", "needle"].indexOf(Model.editorPasteMode)

                    onActivated: Model.editorPasteMode = ["auto", "mouse", "needle"][index]
                }
            }
        }
    }

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
                    id: mappingSelection
                    model: [qsTr("Simple"), qsTr("Custom")]
                    Layout.fillWidth: true
                    currentIndex: Model.liveLedStripMappingMode

                    Component.onCompleted: {
                        activated(currentIndex)
                    }

                    onActivated: {
                        Model.liveLedStripMappingMode = index
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
                visible: false
                text: "Custom Mapping"
                font.weight: Font.DemiBold
                topPadding: 10
            }

            MappingEditor {
                id: mappingEditor
                visible: false
                Layout.fillWidth: true
            }
        } // ColumnLayout
    } // GroupBox
}