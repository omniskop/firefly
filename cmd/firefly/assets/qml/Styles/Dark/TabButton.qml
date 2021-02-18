import QtQuick 2.12
import QtQuick.Templates 2.12 as T
import QtQuick.Controls 2.12

T.TabButton {
    id: control

    font: Shared.font

    implicitWidth: Math.max(implicitBackgroundWidth + leftInset + rightInset,
                            implicitContentWidth + leftPadding + rightPadding)
    implicitHeight: Math.max(implicitBackgroundHeight + topInset + bottomInset,
                             implicitContentHeight + topPadding + bottomPadding)

    padding: 6
    spacing: 6

    contentItem: Text {
         id: textItem
         text: control.text

         font: control.font
         color: Shared.textColor
         horizontalAlignment: Text.AlignHCenter
         verticalAlignment: Text.AlignVCenter
         elide: Text.ElideRight

         states: [
             State {
                 name: "normal"
                 when: !control.down
             },
             State {
                 name: "down"
                 when: control.down
                 PropertyChanges {
                     target: textItem
                 }
             }
         ]
     }

    background: Rectangle {
        id: buttonBackground
        implicitWidth: 80
        implicitHeight: 30
        color: Shared.backgroundColor
        radius: 4

        states: [
            State {
                name: "down"
                when: control.down
                PropertyChanges {
                    target: buttonBackground
                    color: Qt.darker(Shared.backgroundColor, 1.5)
                }
            },
            State {
                name: "checked"
                when: control.checked
                PropertyChanges {
                    target: buttonBackground
                    color: Shared.backgroundColor
                }
            },
            State {
                name: "normal"
                when: !control.down
                PropertyChanges {
                    target: buttonBackground
                    color: Qt.darker(Shared.backgroundColor, 1.3)
                }
            }
        ]

        Rectangle {
            // blocks the bottom radius
            id: bottomBlock
            color: parent.color
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.bottom: parent.bottom
            height: parent.radius
        }
    }
}
