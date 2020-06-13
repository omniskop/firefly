package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var logConsoleReceiver = logHook{
	formatter: consoleFormatter{},
}

type logHook struct {
	bytes.Buffer
	formatter logrus.Formatter
	onChange  func(string)
}

func (h *logHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *logHook) Fire(entry *logrus.Entry) error {
	defer func() {
		if r := recover(); r != nil {
			if r == bytes.ErrTooLarge { // the buffer has become too large
				h.Buffer.Reset()
				h.Buffer.WriteString("[LOG BUFFER RESET]")
			} else {
				panic(r)
			}
		}
	}()
	out, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	h.Buffer.Write(out)
	if h.onChange != nil {
		h.onChange(string(out))
	}
	return nil
}

type consoleFormatter struct{}

func (cf consoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var fields = make([]string, 0, len(entry.Data))
	for k, v := range entry.Data {
		fields = append(fields, fmt.Sprintf("<%s>=%v", k, v))
	}

	s := fmt.Sprintf("%s [%s] %s %s \n",
		entry.Time.Format("15:04:05"),
		strings.ToUpper(entry.Level.String())[0:4],
		entry.Message,
		strings.Join(fields, " "),
	)
	return []byte(s), nil
}

func NewLogConsoleWindow() {
	dialog := widgets.NewQDialog(nil, core.Qt__Tool)
	dialog.SetWindowTitle("Console")
	dialog.SetMinimumSize2(300, 200)
	dialog.SetLayout(widgets.NewQVBoxLayout())

	textField := widgets.NewQPlainTextEdit(nil)
	textField.SetReadOnly(true)
	textField.SetMaximumBlockCount(100)
	textField.CenterOnScroll()
	dialog.Layout().AddWidget(textField)

	textField.SetPlainText(logConsoleReceiver.String())
	logConsoleReceiver.onChange = func(txt string) {
		textField.AppendPlainText(strings.TrimRight(txt, "\n"))
	}

	dialog.Show()
	dialog.Raise()
	dialog.ActivateWindow()
}
