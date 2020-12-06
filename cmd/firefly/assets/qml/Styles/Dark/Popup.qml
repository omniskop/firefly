import QtQuick 2.6
import QtGraphicalEffects 1.0
import QtQuick.Templates 2.3 as T

T.Popup {
    id: control

    implicitWidth: Math.max(background ? background.implicitWidth : 0,
                            contentWidth > 0 ? contentWidth + leftPadding + rightPadding : 0)
    implicitHeight: Math.max(background ? background.implicitHeight : 0,
                             contentWidth > 0 ? contentHeight + topPadding + bottomPadding : 0)

    contentWidth: contentItem.implicitWidth || (contentChildren.length === 1 ? contentChildren[0].implicitWidth : 0)
    contentHeight: contentItem.implicitHeight || (contentChildren.length === 1 ? contentChildren[0].implicitHeight : 0)

    enter: Transition {
        NumberAnimation {
            property: "opacity"
            from: 0
            to: 1
            easing.type: Easing.InOutQuad
            duration: 100
        }
    }

    exit: Transition {
        NumberAnimation {
            property: "opacity"
            from: 1
            to: 0
            easing.type: Easing.InOutQuad
            duration: 100
        }
    }

    contentItem: Item { }

    background: Rectangle {
        radius: 4
        color: Shared.areaColor
        property color borderColor: Shared.textColor
        // border.color: Qt.rgba(borderColor.r, borderColor.g, borderColor.b, 0.3)
        layer.enabled: true
        border.width: 1
        border.color: Shared.backgroundColor
        
        layer.effect: DropShadow {
            transparentBorder: true
            radius: 8
            samples: 16
            horizontalOffset: 0
            verticalOffset: 4
            color: Qt.rgba(0, 0, 0, 0.3)
        }
    }
}