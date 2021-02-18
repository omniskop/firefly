import QtQuick 2.12
import QtQuick.Window 2.2
import QtQuick.Controls 2.12 as Controls
import QtQuick.Templates 2.12 as T
import QtGraphicalEffects 1.13

T.ComboBox {
    id: control

    implicitHeight: 24
    implicitWidth: Math.max(mainText.contentWidth + 50, 200)

    hoverEnabled: true
    property int popupPadding: 5

    delegate: ItemDelegate {
        width: control.popup.width - popupPadding*2
        text: control.textRole ? (Array.isArray(control.model) ? modelData[control.textRole] : model[control.textRole]) : modelData
        highlighted: control.highlightedIndex == index
        font.weight: control.currentIndex == index ? Font.DemiBold : Font.Regular;
        property bool separatorVisible: false
        padding: 5
    }

    indicator: Image {
        x: control.width - width - control.rightPadding
        y: control.topPadding + (control.availableHeight - height) / 2
        width: 22
        height: 22
        sourceSize.height: height * Screen.devicePixelRatio
        sourceSize.width: width * Screen.devicePixelRatio
        source: "qrc:/assets/images/ComboBox_indicator.svg"
        smooth: false
    }

    contentItem: MouseArea {
        onPressed: mouse.accepted = false;
//        onWheel: {
//            if (wheel.pixelDelta.y < 0 || wheel.angleDelta.y < 0) {
//                control.currentIndex = (control.currentIndex + 1) % delegateModel.count
//            } else {
//                control.currentIndex = (control.currentIndex - 1 + delegateModel.count) % delegateModel.count
//            }
//        }
        T.TextField {
            id: mainText
            anchors {
                fill: parent
                leftMargin: control.mirrored ? 12 : 1
                rightMargin: !control.mirrored ? 12 : 1
            }

            text: control.editText

            visible: typeof(control.editable) != "undefined" && control.editable
            readOnly: control.popup.visible
            inputMethodHints: control.inputMethodHints
            validator: control.validator
            renderType: Window.devicePixelRatio % 1 !== 0 ? Text.QtRendering : Text.NativeRendering
            color: Shared.lighterTextColor
            selectByMouse: true

            font: control.font
            horizontalAlignment: Text.AlignLeft
            verticalAlignment: Text.AlignVCenter
            opacity: control.enabled ? 1 : 0.3
        }
    }

    background: Rectangle {
        anchors.fill: parent
        radius: 5
        color: Shared.interactableColor

        LinearGradient {
            anchors.fill: parent
            cached: false
            source: parent
            gradient: Gradient {
                GradientStop { position: 0.0; color: Qt.lighter(Shared.interactableColor, 1.2) }
                GradientStop { position: 1.0; color: Qt.darker(Shared.interactableColor, 1.2) }
            }
        }

        Text {
            anchors.fill: parent
            anchors.leftMargin: 5
            text: control.displayText
            horizontalAlignment: Text.AlignLeft
            verticalAlignment: Text.AlignVCenter
            color: Shared.lighterTextColor
            elide: Text.ElideRight
        }
    }

    popup: Popup {
        y: control.height
        width: Math.max(control.width, 150)
        implicitHeight: contentItem.implicitHeight
        topMargin: 6
        bottomMargin: 6
        padding: popupPadding

        contentItem: ListView {
            id: listview
            clip: true
            implicitHeight: Math.max(control.count*30, 14) + 2
            model: control.popup.visible ? control.delegateModel : null
            currentIndex: control.highlightedIndex
            highlightRangeMode: ListView.ApplyRange
            highlightMoveDuration: 0
            T.ScrollBar.vertical: Controls.ScrollBar { }
        }
    }
}