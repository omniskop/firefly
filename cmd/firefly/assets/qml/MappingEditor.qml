import QtQuick 2.12
import QtQuick.Controls 2.3
import QtQuick.Layouts 1.3
import Dark 1.0

Item {
    id: stripSetupContainer
    width: 295
    height: 100

    property color lineColor: Qt.darker(Shared.textColor, 1.3)
    signal updatePositions

    Label {
        id: label_ledStrip
        height: 15
        text: qsTr("LED-Strip")
        anchors.top: stripGraphic.bottom
        anchors.topMargin: 5
        anchors.right: parent.right
        anchors.rightMargin: 0
        anchors.left: parent.left
        anchors.leftMargin: 0
        font.pixelSize: 12
    }

    Label {
        text: qsTr("Left Click: Add Step\nRight Click: Remove Step")
        anchors.top: label_ledStrip.top
        anchors.right: parent.right
        font.pixelSize: 11
        horizontalAlignment: Text.AlignRight
    }

    Label {
        id: label_firefly
        height: 15
        text: "Firefly Animation"
        anchors.top: parent.top
        anchors.topMargin: 0
        anchors.right: parent.right
        anchors.rightMargin: 0
        anchors.left: parent.left
        anchors.leftMargin: 0
        font.pixelSize: 12
    }

    Item {
        id: stripGraphic
        height: 40
        anchors.right: parent.right
        anchors.rightMargin: 0
        anchors.left: parent.left
        anchors.leftMargin: 0
        anchors.top: fireflyGraphic.bottom
        anchors.topMargin: 0

        Rectangle {
            height: 1
            anchors.top: parent.bottom
            anchors.topMargin: -10
            anchors.left: parent.left
            anchors.right: parent.right
            color: stripSetupContainer.lineColor
        }

        Item { // left strip offset
            id: leftStripOffset
            width: 35
            anchors.left: parent.left
            anchors.leftMargin: 0
            anchors.bottom: parent.bottom
            anchors.bottomMargin: 0
            anchors.top: parent.top
            anchors.topMargin: 0

            Rectangle {
                id: stripVertialLineLeft
                width: 1
                anchors.top: parent.bottom
                anchors.topMargin: -10
                anchors.right: parent.right
                anchors.bottom: parent.bottom
                color: stripSetupContainer.lineColor
            }

            // number of led's
            TextInput {
                width: parent.width
                anchors.baseline: parent.bottom
                text: Model.mapping.startOffset
                horizontalAlignment: Text.AlignHCenter
                font.pixelSize: 10
                color: Shared.textColor
                validator: IntValidator {bottom: 0; top: 1000000}

                onTextEdited: {
                    if(text == "") text = "0";
                    Model.mapping.startOffset = parseInt(text)
                }
            }
        } // left strip offset

        MouseArea {
            cursorShape: Qt.PointingHandCursor
            anchors.fill: stripContent

            onClicked: {
                Model.mapping.addPoint(mouse.x / stripContent.width);
                stripSetupContainer.updatePositions();
            }
        }

        Row {
            id: stripContent
            spacing: 0
            anchors.left: leftStripOffset.right
            anchors.right: rightStripOffset.left
            anchors.top: parent.top
            anchors.bottom: parent.bottom

            Repeater {
                id: stripRepeater
                model: Model.mapping.pointCount()

                Component.onCompleted: {
                    stripSetupContainer.updatePositions.connect(updateModel)
                }

                Component.onDestruction: {
                    stripSetupContainer.updatePositions.disconnect(updateModel);
                }

                function updateModel() {
                    if(stripRepeater.model != Model.mapping.pointCount()) {
                        stripRepeater.model = Model.mapping.pointCount();
                        stripSetupContainer.updatePositions();
                    }
                }

                Item {
                    height: parent.height

                    Component.onCompleted: {
                        stripSetupContainer.updatePositions.connect(updateWidth);
                        stripContent.onWidthChanged.connect(updateWidth);
                    }

                    Component.onDestruction: {
                        stripSetupContainer.updatePositions.disconnect(updateWidth);
                        stripContent.onWidthChanged.disconnect(updateWidth);
                    }

                    function updateWidth() {
                        let prevPosition = index == 0 ? 0 : Model.mapping.getPosition(index-1);
                        width = stripContent.width * (Model.mapping.getPosition(index) - prevPosition);
                    }

                    // vertical line
                    Rectangle {
                        id: stripVerticalLine
                        width: 1
                        anchors.top: parent.bottom
                        anchors.topMargin: -10
                        anchors.left: parent.right
                        anchors.bottom: parent.bottom
                        color: stripSetupContainer.lineColor

                        Grabber {
                            anchors.bottom: parent.top
                            anchors.bottomMargin: 0
                            anchors.horizontalCenter: parent.horizontalCenter
                            label: Math.floor(Model.mapping.getPosition(index)*100) + "%"
                            sliderElement: stripContent

                            onEdited: {
                                Model.mapping.setPosition(index, value);
                                stripSetupContainer.updatePositions();
                            }

                            onDeleteGrabber: {
                                Model.mapping.deletePoint(index);
                                stripSetupContainer.updatePositions();
                            }

                            Component.onCompleted: {
                                stripSetupContainer.updatePositions.connect(updateWidth)
                            }

                            Component.onDestruction: {
                                stripSetupContainer.updatePositions.disconnect(updateWidth);
                            }

                            function updateWidth() {
                                label = Math.floor(Model.mapping.getPosition(index)*100) + "%"
                            }
                        }
                    }

                    // number of led's
                    TextInput {
                        id: textInput
                        width: parent.width
                        anchors.baseline: parent.bottom
                        text: Model.mapping.getLeds(index)
                        horizontalAlignment: Text.AlignHCenter
                        font.pixelSize: 10
                        color: Shared.textColor
                        validator: IntValidator {bottom: 0; top: 1000000}

                        onTextEdited: {
                            if(text == "") text = "0";
                            Model.mapping.setLeds(index, parseInt(text))
                        }
                    }
                }
            }

            // last led segment
            Item {
                height: parent.height

                Component.onCompleted: {
                    stripSetupContainer.updatePositions.connect(updateWidth)
                    stripContent.onWidthChanged.connect(updateWidth);
                    updateWidth();
                }

                Component.onDestruction: {
                    stripSetupContainer.updatePositions.disconnect(updateWidth);
                    stripContent.onWidthChanged.disconnect(updateWidth);
                }

                function updateWidth() {
                    if(Model.mapping.pointCount() == 0) {
                        var lastPosition = 0;
                    } else {
                        var lastPosition = Model.mapping.getPosition(Model.mapping.pointCount()-1);
                    }
                    width = stripContent.width * (1-lastPosition);
                }

                // number of led's
                TextInput {
                    id: textInput
                    width: parent.width
                    anchors.baseline: parent.bottom
                    text: Model.mapping.getLeds(Model.mapping.pointCount())
                    horizontalAlignment: Text.AlignHCenter
                    font.pixelSize: 10
                    color: Shared.textColor
                    validator: IntValidator {bottom: 0; top: 1000000}

                    Component.onCompleted: {
                        stripSetupContainer.updatePositions.connect(updateText);
                    }

                    Component.onDestruction: {
                        stripSetupContainer.updatePositions.disconnect(updateText);
                    }

                    function updateText() {
                        text = Model.mapping.getLeds(Model.mapping.pointCount())
                    }

                    onTextEdited: {
                        if(text == "") text = "0";
                        Model.mapping.setLeds(Model.mapping.pointCount(), parseInt(text))
                        stripSetupContainer.updatePositions();
                    }
                }
            }
        } // strip center row

        Item {
            id: rightStripOffset
            width: 35
            anchors.right: parent.right
            anchors.rightMargin: 0
            anchors.bottom: parent.bottom
            anchors.bottomMargin: 0
            anchors.top: stripContent.top
            anchors.topMargin: 0

            Rectangle {
                width: 1
                anchors.top: parent.bottom
                anchors.topMargin: -10
                anchors.left: parent.left
                anchors.bottom: parent.bottom
                color: stripSetupContainer.lineColor
            }

            // number of led's
            TextInput {
                width: parent.width
                anchors.baseline: parent.bottom
                text: Model.mapping.stopOffset
                horizontalAlignment: Text.AlignHCenter
                font.pixelSize: 10
                color: Shared.textColor
                validator: IntValidator {bottom: 0; top: 1000000}

                onTextEdited: {
                    if(text == "") text = "0";
                    Model.mapping.stopOffset = parseInt(text)
                }
            }
        } // right strip offset
    }

    // firefly side
    Rectangle {
        id: fireflyGraphic
        height: 10
        color: Shared.interactableColor
        radius: 5
        anchors.right: parent.right
        anchors.rightMargin: 35
        anchors.left: parent.left
        anchors.leftMargin: 35
        anchors.top: label_firefly.bottom
        anchors.topMargin: 5

        Row {
            id: fireflyRow
            x: 0
            y: 0
            width: parent.width
            height: parent.height

            Repeater {
                model: Model.mapping.pointCount()

                Component.onCompleted: {
                    stripSetupContainer.updatePositions.connect(updateModel)
                }

                Component.onDestruction: {
                    stripSetupContainer.updatePositions.disconnect(updateModel);
                }

                function updateModel() {
                    model = Model.mapping.pointCount();
                }

                Item {
                    id: rectangle
                    height: parent.height

                    Component.onCompleted: {
                        stripSetupContainer.updatePositions.connect(updateWidth);
                        fireflyRow.onWidthChanged.connect(updateWidth);
                        updateWidth();
                    }

                    Component.onDestruction: {
                        stripSetupContainer.updatePositions.disconnect(updateWidth);
                        fireflyRow.onWidthChanged.disconnect(updateWidth);
                    }

                    function updateWidth() {
                        let prevPosition = index == 0 ? 0 : Model.mapping.getPosition(index-1);
                        width = fireflyRow.width * (Model.mapping.getPosition(index) - prevPosition);
                    }

                    Rectangle {
                        id: rectangle1
                        anchors.right: parent.right
                        anchors.rightMargin: 0
                        width: 1
                        height: parent.height
                        color: "#000000"
                        transformOrigin: Item.Center
                    }
                }
            }
        }
    }
}
