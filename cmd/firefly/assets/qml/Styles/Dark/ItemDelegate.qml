import QtQuick 2.5
import QtQuick.Templates 2.3 as T

T.ItemDelegate {
    id: control

    implicitWidth: contentItem.implicitWidth + leftPadding + rightPadding
    implicitHeight: Math.max(contentItem.implicitHeight,
                                      indicator ? indicator.implicitHeight : 0) + topPadding + bottomPadding
    hoverEnabled: true

    padding: 10
    spacing: 4
    rightPadding: 20

    contentItem: Label {
        leftPadding: control.mirrored ? (control.indicator ? control.indicator.width : 0) + control.spacing : 0
        rightPadding: !control.mirrored ? (control.indicator ? control.indicator.width : 0) + control.spacing : 0

        text: control.text
        font: control.font
        color: control.highlighted || control.checked || (control.pressed && !control.checked && !control.sectionDelegate) ? Shared.lighterTextColor : (control.enabled ? Shared.textColor : Shared.disabledTextColor)
        elide: Text.ElideRight
        visible: control.text
        horizontalAlignment: Text.AlignLeft
        verticalAlignment: Text.AlignVCenter
    }

    background: Rectangle {
        color: control.highlighted ? Shared.interactableColor : "transparent"
        radius: 4
    }
}