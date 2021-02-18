import QtQuick 2.12
import QtQuick.Templates 2.12 as T
import QtGraphicalEffects 1.13

T.Button {
    id: control

    property bool fancy: false
    property bool small: false
    property bool square: false

    font: Shared.font

    implicitWidth: Math.max(background ? background.implicitWidth : 0,
                                         contentItem.implicitWidth + leftPadding + rightPadding)
    implicitHeight: Math.max(background ? background.implicitHeight : 0,
                                          contentItem.implicitHeight + topPadding + bottomPadding)

    states: [
        State { name: "small"; when: small
            PropertyChanges { target: control
                padding: 3
                width: square ? control.height : undefined
            }
        },
        State { name: "fancy"; when: fancy
            PropertyChanges { target: control
                padding: 10
            }
        },
        State { name: "normal"; when: !small && !fancy
            PropertyChanges { target: control
                leftPadding: 10
                rightPadding: 10
                topPadding: 5
                bottomPadding: 5
            }
        }
    ]

//    leftPadding: 10
//    rightPadding: 10
//    topPadding: 5
//    bottomPadding: 5
//    padding: 10
//    hoverEnabled: true

    background: Rectangle {
        id: buttonBackground
//        implicitWidth: 80
//        implicitHeight: 30
        color: Shared.interactableColor
        radius: 4

        states: [
            State {
                name: "normal"
                when: !control.down
                PropertyChanges {
                    target: buttonBackground
                    color: Shared.interactableColor
                }
            },
            State {
                name: "down"
                when: control.down
                PropertyChanges {
                    target: buttonBackground
                    color: Qt.darker(Shared.interactableColor, 1.3)
                }
            }
        ]
    }

    layer.enabled: fancy

    layer.effect: DropShadow {
        transparentBorder: true
        radius: 10
        samples: 21
        horizontalOffset: 0
        verticalOffset: 1
        color: Qt.rgba(0, 0, 0, 0.5)
    }

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
}

