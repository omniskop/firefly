import QtQuick 2.12
import QtQuick.Window 2.2
import Dark 1.0

Item {
    id: element
    width: 30
    height: 30
    property alias label: label.text
    property Item sliderElement
    signal edited(double value)
    signal grabbed()
    signal released()
    signal deleteGrabber()

    TextInput {
        id: label
        text: "100%"
        font.pixelSize: 10
        color: Shared.textColor
        anchors.bottom: image.top
        anchors.bottomMargin: -3
        anchors.horizontalCenter: image.horizontalCenter
        horizontalAlignment: Text.AlignHCenter
        onTextEdited: {
            let v = parseInt(text) / 100;
            if(isNaN(v)) v = 0;
            v = Math.max(0, Math.min(v, 1));
            text = Math.floor(v*100) + "%";
            if(v != 0) parent.edited(v);
        }
    }

    Image {
        id: image
        width: 20
        height: 20
        anchors.bottom: parent.bottom
        anchors.horizontalCenter: parent.horizontalCenter
        smooth: false
        fillMode: Image.PreserveAspectFit
        source: "qrc:/assets/images/grabber.svg"

        // fixes svg resolution
        sourceSize.height: height * Screen.devicePixelRatio
        sourceSize.width: width * Screen.devicePixelRatio

        MouseArea {
            cursorShape: Qt.OpenHandCursor
            anchors.fill: parent
            acceptedButtons: Qt.LeftButton | Qt.RightButton
            onPressed: {
                if(mouse.button == Qt.RightButton) {
                    element.deleteGrabber()
                } else {
                    cursorShape = Qt.ClosedHandCursor
                }
            }
            onReleased: {
                cursorShape = Qt.OpenHandCursor
            }
            onPositionChanged: {
                let position = mapToItem(sliderElement, mouse.x, mouse.y).x / sliderElement.width;
                position = Math.max(0, Math.min(position, 1));
                element.edited(position)
            }
        }
    }
}
