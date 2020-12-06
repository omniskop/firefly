import QtQuick 2.12
import QtQuick.Window 2.2
import QtQuick.Templates 2.12 as T

T.Label {
    id: control

    verticalAlignment: lineCount > 1 ? Text.AlignTop : Text.AlignVCenter

    activeFocusOnTab: false
    // from: https://github.com/mauikit/qqc2-desktop-style-maui/blob/master/Label.qml
    //Text.NativeRendering is broken on non integer pixel ratios
    renderType: Window.devicePixelRatio % 1 !== 0 ? Text.QtRendering : Text.NativeRendering
    
    color: Shared.textColor

    Accessible.role: Accessible.StaticText
    Accessible.name: text    
}