package logger

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	basicLogger = log.New()
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// type Logger log.Entry

func init() {
	//formatter := &log.TextFormatter{DisableColors:true, FullTimestamp:true, TimestampFormat:"01-02 15:04:05.000"}
	//formatter = &log.TextFormatter{}
	//formatter := &log.JSONFormatter{TimestampFormat:"15:04:05.000"}
	formatter := &MyFomatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"
	basicLogger.Level = log.InfoLevel
	log.SetFormatter(formatter)
	basicLogger.Formatter = formatter
}

type Debugger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

type GameDebugger interface {
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
}

func NewLogger() *BasicLogger {
	return &BasicLogger{
		Logger: basicLogger,
	}
}

//BasicLogger 不需打印src
type BasicLogger struct {
	*log.Logger
}

func PrintRaw(content string) {
	basicLogger.Println(content)
}

func GetLogger() *log.Entry {
	fileName, line, funcName := "?file?", 0, "?func?"
	pc, fileName, line, ok := runtime.Caller(2)
	if ok {
		dir, file := filepath.Split(fileName)
		dir = filepath.Base(dir)
		fileName = filepath.Join(dir, file)

		funcName = runtime.FuncForPC(pc).Name()
		idx := strings.LastIndex(funcName, ".")
		funcName = funcName[idx+1:]
	}

	return basicLogger.WithField("__src", fmt.Sprintf("[%s:%d %s]", fileName, line, funcName))
}

func Trace(v ...interface{}) {
	GetLogger().Trace(v...)
}

func Tracef(format string, v ...interface{}) {
	GetLogger().Tracef(format, v...)
}

func Traceln(v ...interface{}) {
	GetLogger().Traceln(v...)
}

func Debug(v ...interface{}) {
	GetLogger().Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	GetLogger().Debugf(format, v...)
}

func Debugln(v ...interface{}) {
	GetLogger().Debugln(v...)
}

func Info(v ...interface{}) {
	GetLogger().Info(v...)
}

func Infof(format string, v ...interface{}) {
	GetLogger().Infof(format, v...)
}

func Infoln(v ...interface{}) {
	GetLogger().Infoln(v...)
}

func Warn(v ...interface{}) {
	GetLogger().Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	GetLogger().Warnf(format, v...)
}

func Warnln(v ...interface{}) {
	GetLogger().Warnln(v...)
}

func Error(v ...interface{}) {
	GetLogger().Error(v...)
}

func Errorf(format string, v ...interface{}) {
	GetLogger().Errorf(format, v...)
}

func Errorln(v ...interface{}) {
	GetLogger().Errorln(v...)
}

func Fatal(v ...interface{}) {
	GetLogger().Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	GetLogger().Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	GetLogger().Fatalln(v...)
}

func Panic(v ...interface{}) {
	GetLogger().Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	GetLogger().Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	GetLogger().Panicln(v...)
}

func WithError(err error) *log.Entry {
	return GetLogger().WithError(err)
}

func WithField(key string, value interface{}) *log.Entry {
	return GetLogger().WithField(key, value)
}

type Fields map[string]interface{}

func WithFields(fields Fields) *log.Entry {
	return GetLogger().WithFields(log.Fields(fields))
}

func WithSrc(entry *log.Entry) *log.Entry {
	fileName, line, funcName := "?file?", 0, "?func?"
	pc, fileName, line, ok := runtime.Caller(2)
	if ok {
		dir, file := filepath.Split(fileName)
		dir = filepath.Base(dir)
		fileName = filepath.Join(dir, file)

		funcName = runtime.FuncForPC(pc).Name()
		idx := strings.LastIndex(funcName, ".")
		funcName = funcName[idx+1:]
	}

	return entry.WithField("src", fmt.Sprintf("[%s:%d %s]", fileName, line, funcName))
}

type MyFomatter struct {
	log.TextFormatter
}

func (f *MyFomatter) Format(entry *log.Entry) ([]byte, error) {
	var src string
	if value, ok := entry.Data["__src"]; ok {
		src = value.(string)
		delete(entry.Data, "__src")
	}
	var ctx string
	if value, ok := entry.Data["__ctx"]; ok {
		ctx = value.(string)
		delete(entry.Data, "__ctx")
	}

	var b *bytes.Buffer
	var keys []string = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixFieldClashes(entry.Data)
	timestampFormat := f.TimestampFormat

	if !f.DisableTimestamp {
		f.appendValue(b, entry.Time.Format(timestampFormat))
	}

	f.appendValue(b, printLogLevel(entry.Level))

	if ctx != "" {
		f.appendValue(b, ctx)
	}
	if src != "" {
		f.appendValue(b, src)
	}

	if entry.Message != "" {
		f.appendValue(b, entry.Message)
	}

	for _, key := range keys {
		f.appendKeyValue(b, key, entry.Data[key])
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func printLogLevel(level log.Level) string {
	switch level {
	case log.DebugLevel:
		return "DEBUG"
	case log.InfoLevel:
		// return Reset + Green + "INFO" + Reset
		return "INFO"
	case log.WarnLevel:
		return Reset + Yellow + "WARN" + Reset
	case log.ErrorLevel:
		return Reset + Red + "ERROR" + Reset
	case log.FatalLevel:
		return Reset + Red + "FATAL" + Reset
	case log.PanicLevel:
		return Reset + Red + "PANIC" + Reset
	}

	return "unknown"
}

func prefixFieldClashes(data log.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}
}

func needsQuoting(text string) bool {
	return false
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return true
		}
	}
	return false
}

func (f *MyFomatter) appendValue(b *bytes.Buffer, value interface{}) {
	switch value := value.(type) {
	case string:
		b.WriteString(value)
	case error:
		errmsg := value.Error()
		if !needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	default:
		fmt.Fprint(b, value)
	}

	b.WriteByte(' ')
}

func (f *MyFomatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	b.WriteString(key)
	b.WriteByte('=')

	switch value := value.(type) {
	case string:
		if !needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	case error:
		errmsg := value.Error()
		if !needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	default:
		fmt.Fprint(b, value)
	}

	b.WriteByte(' ')
}

// SetLogLevel ..
func SetLogLevel(level log.Level) {
	basicLogger.Level = level
}
