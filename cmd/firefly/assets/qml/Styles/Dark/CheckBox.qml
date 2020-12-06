import QtQuick 2.12
import QtQuick.Window 2.2
import QtQuick.Templates 2.12 as T

T.CheckBox {
    id: control

    implicitWidth: Math.max(background ? background.implicitWidth : 0,
                                         contentItem.implicitWidth + leftPadding + rightPadding)
    implicitHeight: Math.max(background ? background.implicitHeight : 0,
                                          Math.max(contentItem.implicitHeight,
                                                   indicator ? indicator.implicitHeight : 0) + topPadding + bottomPadding)


    indicator: Rectangle {
        id: checkboxHandle
        implicitWidth: 18
        implicitHeight: 18
        anchors.verticalCenter: parent.verticalCenter
        radius: 3
        color: Shared.interactableColor

        Rectangle {
            id: rectangle
            anchors.fill: parent
            anchors.margins: 3
            radius: 2
            visible: false
            color: Shared.textColor
        }

        states: [
            State {
                name: "unchecked"
                when: !control.checked
            },
            State {
                name: "checked"
                when: control.checked

                PropertyChanges {
                    target: rectangle
                    visible: true
                }
            }
        ]
    }

    contentItem: Text {
        leftPadding: control.indicator.width + 7

        text: control.text
        color: Shared.textColor
        elide: Text.ElideRight
        visible: control.text
        horizontalAlignment: Text.AlignLeft
        verticalAlignment: Text.AlignVCenter
    }
}