package logger

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// ANSI 颜色代码
const (
	colorReset  = "\x1b[0m"
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorBlue   = "\x1b[34m"
	colorPurple = "\x1b[35m"
	colorCyan   = "\x1b[36m"
	colorWhite  = "\x1b[37m"
	colorGray   = "\x1b[90m"
)

type CustomFormatter struct {
	// DisableColor 禁用所有 ANSI 颜色代码
	DisableColor bool
	// ColorFullMessage 在 Error/Warn 等级别时，对整个消息进行着色
	ColorFullMessage bool
	// TimestampFormat 时间戳格式
	TimestampFormat string
}

func NewCustomFormatter(disableColor, colorFullMessage bool) *CustomFormatter {
	return &CustomFormatter{
		DisableColor:     disableColor,
		ColorFullMessage: colorFullMessage,
		TimestampFormat:  "2006-01-02 15:04:05.000",
	}
}

func (f *CustomFormatter) getColor(colorCode string) string {
	if f.DisableColor {
		return ""
	}
	return colorCode
}

func (f *CustomFormatter) getReset() string {
	if f.DisableColor {
		return ""
	}
	return colorReset
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	isFrpPkg := entry.Data["pkg"] == "frp"
	resetCode := f.getReset()

	var levelColorCode string
	var msgColor string

	switch entry.Level {
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColorCode = colorRed
	case logrus.WarnLevel:
		levelColorCode = colorYellow
	case logrus.InfoLevel:
		levelColorCode = colorBlue
	case logrus.DebugLevel:
		levelColorCode = colorCyan
	case logrus.TraceLevel:
		levelColorCode = colorPurple
	default:
		levelColorCode = colorWhite
	}

	if !f.DisableColor && f.ColorFullMessage {
		switch entry.Level {
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			msgColor = f.getColor(colorRed)
		case logrus.WarnLevel:
			msgColor = f.getColor(colorYellow)
		}
	}

	if !isFrpPkg {
		timestampFormat := f.TimestampFormat
		if timestampFormat == "" {
			timestampFormat = "2006-01-02 15:04:05.000"
		}
		b.WriteString(entry.Time.Format(timestampFormat))

		levelColor := f.getColor(levelColorCode)
		fmt.Fprintf(b, " [%s%s%s]", levelColor, entry.Level.String(), resetCode)

		if entry.HasCaller() {
			fileName := filepath.Base(entry.Caller.File)
			fatherDir := filepath.Base(filepath.Dir(entry.Caller.File))
			callerColor := f.getColor(colorGray)
			fmt.Fprintf(b, " [%s%s:%d%s]", callerColor, filepath.Join(fatherDir, fileName), entry.Caller.Line, resetCode)
		}
		b.WriteString(" ")
	}

	for key, val := range entry.Data {
		fmt.Fprintf(b, "[%v: %v] ", key, val)
	}

	if entry.Message != "" {
		if !isFrpPkg {
			b.WriteString(" ")
		}
		b.WriteString(msgColor)
		b.WriteString(entry.Message)
		if msgColor != "" {
			b.WriteString(resetCode)
		}
	}

	if !isFrpPkg {
		b.WriteByte('\n')
	}

	return b.Bytes(), nil
}
