import QtQuick 2.12
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.3
import QtQuick.Window 2.2
import Dark 1.0

Rectangle {
    id: root

    color: Shared.backgroundColor

    property bool saveOnClose: false

    ColumnLayout {
        anchors.fill: parent

        TabBar {
            id: bar
            Layout.fillWidth: true
            TabButton {
                text: qsTr("General")
                width: implicitWidth
            }
            TabButton {
                text: qsTr("Editor")
                width: implicitWidth
            }
        }

        StackLayout {
            id: settingsStack
            Layout.fillWidth: true
            currentIndex: bar.currentIndex

            GeneralSettings {}
            EditorSettings {}
        }

        Row {
            id: buttonRow
            Layout.alignment: Qt.AlignRight | Qt.AlignBottom
            spacing: 10
            rightPadding: 10
            bottomPadding: 10

            Button {
                id: "cancelButton"
                text: qsTr("Cancel")
                onClicked: Model.cancel()
            }

            Button {
                id: "okButton"
                text: qsTr("OK")

                onClicked: Model.ok()
            }
        } // Row
    } // ColumnLayout
}






/*##^## Designer {
    D{i:0;width:400}D{i:19;invisible:true}
}
 ##^##*/
