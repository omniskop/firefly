import QtQuick 2.12
import QtQuick.Window 2.2
import QtQuick.Controls 2.3 as Controls
import QtQuick.Templates 2.12 as T
import QtGraphicalEffects 1.12

T.TextField {
    id: control

    implicitWidth: Math.max(200,
                            placeholderText ? placeholder.implicitWidth + leftPadding + rightPadding : 0)
                            || contentWidth + leftPadding + rightPadding
    implicitHeight: Math.max(contentHeight + topPadding + bottomPadding,
                             background ? 20 : 0,
                             placeholder.implicitHeight + topPadding + bottomPadding)

    padding: 4
    color: Shared.lighterTextColor

    // color: control.enabled ? Kirigami.Theme.textColor : Kirigami.Theme.disabledTextColor
    // selectionColor: Kirigami.Theme.highlightColor
    // selectedTextColor: Kirigami.Theme.highlightedTextColor
    verticalAlignment: TextInput.AlignVCenter
    //Text.NativeRendering is broken on non integer pixel ratios
    renderType: Window.devicePixelRatio % 1 !== 0 ? Text.QtRendering : Text.NativeRendering
    selectByMouse: true

    Controls.Label {
        id: placeholder
        x: control.leftPadding
        y: control.topPadding
        width: control.width - (control.leftPadding + control.rightPadding)
        height: control.height - (control.topPadding + control.bottomPadding)

        text: control.placeholderText
        font: control.font
        horizontalAlignment: control.horizontalAlignment
        verticalAlignment: control.verticalAlignment
        visible: !control.length && !control.preeditText && (!control.activeFocus || control.horizontalAlignment !== Qt.AlignHCenter)
        elide: Text.ElideRight
        color: Shared.textColor
    }

    background: Rectangle {
        id: textFieldRectangle
        radius: 5
        color: Shared.interactableColor
        border.width: 1
        border.color: "#30000000"
        layer.enabled: true

        layer.effect: InnerShadow {
            anchors.fill: parent
            radius: 5
            samples: 11
            horizontalOffset: 0
            verticalOffset: 1
            color: "#40000000"
            cached: false
        }

        states: [
            State {
                name: "unfocused"
                when: !control.activeFocus
            },
            State {
                name: "focused"
                when: control.activeFocus

                PropertyChanges {
                    target: textFieldRectangle
                    border.color: "#60000000"
                    color: Qt.darker(Shared.interactableColor, 1.2)
                }
            }
        ]
    }
}