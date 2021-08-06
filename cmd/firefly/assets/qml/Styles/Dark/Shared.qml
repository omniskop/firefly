pragma Singleton

import QtQuick 2.12

QtObject {
    readonly property color voidBackgroundColor: "#202020" // behind the background
    readonly property color backgroundColor: "#38393A"
    readonly property color areaColor: "#444444"
    readonly property color interactableColor: "#7D7D7D"
    readonly property color textColor: "#E0E0E0"
    readonly property color lighterTextColor: "#F0F0F0"
    readonly property color disabledTextColor: Qt.darker(textColor, 1.3)

    property font font
    font.bold: true
    font.underline: false
    font.pixelSize: 20
    font.family: "arial"
}
