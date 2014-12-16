// Package name declaration
package logger

// Import packages
import (
 	"bytes"
 	"fmt"
 	"log"
 	"time"
 	"os"
 	"sync/atomic"
)

var (
	// Map for te various codes of colors
	colors map[string]string
	
	// Contains color strings for stdout
	logNo uint64
)

// Color numbers for stdout
const (
	Black = (iota + 30)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Info class, Contains all the info on what has to logged, time is the current time, 
// level is the state, importance and type of message logged,
// Message contains the string to be logged, format is the format of string to be passed to sprintf
type Info struct {
	Id uint64
	Time string
	Level string
	Message string
	format string
}

// Returns a proper string to be outputted for a particular info
func (r *Info) Output() string {
	msg := fmt.Sprintf(r.format, r.Id, r.Time, r.Level, r.Message )
	return msg
}

// Logger class that is an interface to user to log messages,
// worker is variable of Worker class that is used in bottom layers to log the message
type Logger struct {
	Prefix string
    FullPath string
    Out *os.File
    Color int
	Minion *log.Logger
}

// Returns a proper string to output for colored logging
func colorString(color int) string {
	return fmt.Sprintf("\033[%dm", int(color))
}

// Initializes the map of colors
func initColors() {
	colors = map[string]string{
		"CRITICAL": colorString(Magenta),
		"ERROR":    colorString(Red),
		"WARNING":  colorString(Yellow),
		"NOTICE":   colorString(Green),
		"DEBUG":    colorString(Cyan),
		"INFO" :    colorString(White),
	}
}

// Returns a new instance of logger class, module is the specific module for which we are logging
// , color defines whether the output is to be colored or not, out is instance of type io.Writer defaults
// to os.Stderr
func New(prefix string, color int) (*Logger, error) {
    initColors()

    fullPath := fmt.Sprintf("%v.%v", prefix, time.Now().Format("2006-01-02"))
    out, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }

    logger := &Logger {
        Prefix: prefix,
        FullPath: fullPath,
        Out: out,
        Color: color,
        Minion: log.New(out, "", 0),
    }

    return logger, nil
}

// The log commnand is the function available to user to log message, lvl specifies
// the degree of the messagethe user wants to log, message is the info user wants to log
func (l *Logger) Log(lvl string, message string) error {

    //Checking this is a new day
    fullPath := fmt.Sprintf("%v.%v", l.Prefix, time.Now().Format("2006-01-02"))
    if fullPath != l.FullPath {
        if l.Out != nil {
            l.Out.Close()
            l.Out = nil
        }

        out, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            return err
        }

        l.FullPath = fullPath
        l.Out = out
        l.Minion = log.New(out, "", 0)
    }

	var formatString string = "#%d %s ▶ %.3s %s"
	info := &Info{
		Id:      atomic.AddUint64(&logNo, 1),
		Time:    time.Now().Format("2006-01-02 15:04:05") ,
		Level:   lvl,
		Message: message,
		format:  formatString,
	}
	return l.Logging(lvl, 2, info)
}

// Function of Worker class to log a string based on level
func (l *Logger) Logging(level string, calldepth int, info *Info) error {
	if l.Color != 0 {
		buf := &bytes.Buffer{}
		buf.Write([]byte(colors[level]))
		buf.Write([]byte(info.Output()))
		buf.Write([]byte("\033[0m"))
		return l.Minion.Output(calldepth+1, buf.String())
	} else {
		return l.Minion.Output(calldepth+1, info.Output())
	}
}

// Fatal is just like func l,Cr.tical logger except that it is followed by exit to program
func (l *Logger) Fatal(message string) {
	l.Log("CRITICAL", message)
	os.Exit(1)
}

// Panic is just like func l.Critical except that it is followed by a call to panic
func (l *Logger) Panic(message string) {
	l.Log("CRITICAL", message)
	panic(message)
}

// Critical logs a message at a Critical Level
func (l *Logger) Critical(message string) {
	l.Log("CRITICAL", message)
}

// Error logs a message at Error level
func (l *Logger) Error(message string) {
	l.Log("ERROR", message)
}

// Warning logs a message at Warning level
func (l *Logger) Warning(message string) {
	l.Log("WARNING", message)
}

// Notice logs a message at Notice level
func (l *Logger) Notice(message string) {
	l.Log("NOTICE", message)
}

// Info logs a message at Info level
func (l *Logger) Info(message string) {
	l.Log("INFO", message)
}

// Debug logs a message at Debug level
func (l *Logger) Debug(message string) {
	l.Log("DEBUG", message)
}

// Fatal is just like func l,Cr.tical logger except that it is followed by exit to program
func (l *Logger) Fatalf(format string, a ...interface{}) {
    message := fmt.Sprintf(format, a)
    l.Fatal(message)
}

// Panic is just like func l.Critical except that it is followed by a call to panic
func (l *Logger) Panicf(format string, a ...interface{}) {
    message := fmt.Sprintf(format, a)
    l.Panic(message)
}

// Critical logs a message at a Critical Level
func (l *Logger) Criticalf(format string, a ...interface{}) {
    message := fmt.Sprintf(format, a)
    l.Critical(message)
}

// Error logs a message at Error level
func (l *Logger) Errorf(format string, a ...interface{}) {
    message := fmt.Sprintf(format, a)
    l.Error(message)
}

// Warning logs a message at Warning level
func (l *Logger) Warningf(format string, a ...interface{}) {
    message := fmt.Sprintf(format, a)
    l.Warning(message)
}

// Notice logs a message at Notice level
func (l *Logger) Noticef(format string, a ...interface{}) {
    message := fmt.Sprintf(format, a)
    l.Notice(message)
}

// Info logs a message at Info level
func (l *Logger) Infof(format string, a ...interface{}) {
    message := fmt.Sprintf(format, a)
    l.Info(message)
}

// Debug logs a message at Debug level
func (l *Logger) Debugf(format string, a ...interface{}) {
    message := fmt.Sprintf(format, a)
    l.Debug(message)
}
