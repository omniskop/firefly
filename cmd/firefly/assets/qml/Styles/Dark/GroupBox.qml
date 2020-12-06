import QtQuick 2.12
import QtQuick.Templates 2.12 as T

T.GroupBox {
    id: control
    title: qsTr("GroupBox")
    implicitWidth: contentWidth + leftPadding + rightPadding
    implicitHeight: contentHeight + topPadding + bottomPadding

    contentWidth: contentItem.implicitWidth || (contentChildren.length === 1 ? contentChildren[0].implicitWidth : 0)
    contentHeight: contentItem.implicitHeight || (contentChildren.length === 1 ? contentChildren[0].implicitHeight : 0)

    padding: 20
//    topPadding: padding + (label && label.implicitWidth > 0 ? label.implicitHeight + spacing : 0)
    topPadding: 10 + (label && label.implicitWidth > 0 ? label.implicitHeight + spacing : 0)

    background: Rectangle {
        anchors.fill: control
        anchors.topMargin: (label && label.implicitWidth > 0 ? label.implicitHeight + spacing : 0)
        anchors.leftMargin: 10
        anchors.rightMargin: 10
        anchors.bottomMargin: 10
        color: Shared.areaColor
        radius: 10
        border.width: 1
        border.color: "#50000000"
    }

    label: Text {
        id: title
        text: control.title
        font.pixelSize: 15
        color: Shared.textColor
        font.weight: Font.DemiBold
        topPadding: 5
        leftPadding: 15
        bottomPadding: 2
    }
}
