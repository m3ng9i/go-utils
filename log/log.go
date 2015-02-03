package log

import "io"
import "errors"
import "strings"
import "time"
import "fmt"
import "sync"
import "os"
import "path/filepath"


// ------------------------------------------------
// const


// Log level.
const (
    DEBUG = iota  // Record everything
    INFO
    WARN
    ERROR
    FATAL
)


// Time format.
const (
    TF_DEFAULT  = "2006-01-02 15:04:05.000000"
    TF_NORMAL   = "2006-01-02 15:04:05"
    TF_LONG     = "2006-01-02 15:04:05.000000 -0700"
    TF_TIME     = "15:04:05"
    TF_TIMELONG = "15:04:05.000000"
)


// Log message layout style. Available marks are: {time}, {level} and {msg}.
const (
    LS_DEFAULT = "{time} {level}: {msg}"
    LS_SIMPLE = "{time}: {msg}"
)


// Log message layout elements.
const (
    LY_TIME = 1 << iota
    LY_LEVEL
    LY_DEFAULT = LY_TIME | LY_LEVEL
    LY_MSGONLY = -1 // message only
)


// Rotate config for log file.
const (
    R_NONE = iota   // Don't rename log file.
    R_HOURLY        // Rename log file every hour.
    R_DAILY         // Rename log file every day.
    R_MONTHLY       // Rename log file every month.
)


// Default filename rotate pattern of a file log. Available marks are: {time}, {basename} and {ext}.
const RP_DEFAULT = "{time}_{basename}{ext}"


// Max jobs in Logger.
const maxJobs = 1024


// ------------------------------------------------
// Utils functions


func isLevelLegal(level int) bool {
    return (level >= DEBUG && level <= FATAL)
}


func isLayoutLegal(layout int, style string) error {

    if layout == 0 {
        return errors.New("Layout cannot be zero.")
    }

    if !strings.Contains(style, "{msg}") {
        return errors.New("Layout style should include {msg}.")
    }

    if layout == LY_MSGONLY {
        return nil
    }

    if layout & LY_TIME > 0 && !strings.Contains(style, "{time}") {
        return errors.New("Layout style should include {time}.")
    }

    if layout & LY_LEVEL > 0 && !strings.Contains(style, "{level}") {
        return errors.New("Layout style should include {level}.")
    }

    return nil
}


func isRotateLegal(rotate int) bool {
    return (rotate >= R_NONE && rotate <= R_MONTHLY)
}


func isRotatePatternLegal(pattern string) error {
    if  !strings.Contains(pattern, "{time}") ||
        !strings.Contains(pattern, "{basename}") ||
        !strings.Contains(pattern, "{ext}") {

        return errors.New("Rotate pattern should include {time}, {basename} and {ext}")
    }

    return nil
}


// Check if a writer is legal. Argument "rotate" is a value of Config.Rotate.
func ifWriterLegal(w io.Writer, rotate int) error {
    if rotate > R_NONE {
        file, ok := w.(*os.File)
        if !ok {
            return errors.New("Writer is not file, could not be rotated.")
        }
        if file == os.Stdout || file == os.Stderr {
            return errors.New("Stdout or stderr could not be rotated.")
        }
    }
    return nil
}


func level2string(level int) string {
    switch level {
        case DEBUG: return "DEBUG"
        case INFO:  return "INFO"
        case WARN:  return "WARN"
        case ERROR: return "ERROR"
        case FATAL: return "FATAL"
    }
    return ""
}


// If a config element is zero value, set it to default.
func setConfigDefault(config *Config) {

    if config.TimeFormat == "" {
        config.TimeFormat = TF_DEFAULT
    }

    if config.LayoutStyle == "" {
        config.LayoutStyle = LS_DEFAULT
    }

    if config.Layout == 0 {
        config.Layout = LY_DEFAULT
    }

    if config.Rotate > R_NONE {
        config.RotatePattern = RP_DEFAULT
    }
}


// Open a file to write log
func OpenFile(filename string, filemode ...os.FileMode) (*os.File, error) {
    var mode os.FileMode = 0640
    if len(filemode) > 0 {
        mode = filemode[0]
    }
    return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, mode)
}


// ------------------------------------------------
// Config


type Config struct {
    TimeFormat      string          // Time format in output log message.
    LayoutStyle     string          // Log messae layout style.
    Layout          int             // Log message layout element.
    Level           int             // Log level.
    Utc             bool            // If use utc time in output.
    Rotate          int             // How to rotate file log. See R_NONE, R_HOURLY, R_DAILY and R_MONTHLY.
    RotatePattern   string          // Filename rotate pattern of a file log. See RP_DEFAULT for example.
}


// ------------------------------------------------
// Message


type Message struct {
    Msg string
    Time time.Time
    Level int
}


// Make a new Message structure.
func newMsg(s string, level int) Message {

    var m Message
    m.Time = time.Now()

    m.Level = level
    m.Msg = s

    return m
}


// ------------------------------------------------
// Handle


// User defined function.
type Handle struct {
    Func func(Message)  // User defined function for log message process.
    Level int           // Messages of "Level" and above could be processed by "Func" function.
}


// ------------------------------------------------
// Logger

type Logger struct {
    Config
    w io.Writer
    jobs chan Message
    wg  sync.WaitGroup
    handle *Handle
}


// Create a Logger.
func New(w io.Writer, config Config, handle ...Handle) (logger *Logger, err error) {
    logger = new(Logger)

    setConfigDefault(&config)

    if !isLevelLegal(config.Level) {
        err = fmt.Errorf("Log level is not legal: %d", config.Level)
        return
    }

    err = isLayoutLegal(config.Layout, config.LayoutStyle)
    if err != nil {
        return
    }

    if !isRotateLegal(config.Rotate) {
        err = fmt.Errorf("Rotate is not legal: %d", config.Rotate)
        return
    }

    err = ifWriterLegal(w, config.Rotate) 
    if err != nil {
        return
    }

    if config.Rotate > R_NONE {
        err = isRotatePatternLegal(config.RotatePattern)
        if err != nil {
            return
        }
    } else {
        if len(config.RotatePattern) > 0 {
            err = errors.New("Rotate pattern does not match rotate value.")
            return
        }
    }

    // set user defined function
    if len(handle) > 0 {
        if handle[0].Func == nil {
            err = errors.New("Handle function could not be nil.")
            return
        }
        if !isLevelLegal(handle[0].Level) {
            err = fmt.Errorf("Log level of handle function is not legal: %d", handle[0].Level)
            return
        }
        logger.handle = &handle[0]
    }

    logger.Config   = config
    logger.w        = w
    logger.jobs     = make(chan Message, maxJobs)

    logger.start()

    return
}


// Start to receive logging jobs, then write log message to io.Writer. This function will rotate file log as needed.
func (this *Logger) start() {
    go func() {
        var lastMsgTime time.Time
        for msg := range this.jobs {

            // Call user defined function as needed.
            if this.handle != nil && msg.Level >= this.handle.Level {
                this.wg.Add(1)
                go func(m Message) {
                    defer this.wg.Done()
                    (*this.handle).Func(m)
                }(msg)
            }

            // Check if need to rotate file log
            if this.Rotate > R_NONE && !lastMsgTime.IsZero() {
                timestr := this.ifRotate(lastMsgTime, msg.Time)

                // now rotate
                if timestr != "" {
                    file, _ := this.w.(*os.File)
                    filename := file.Name()

                    err := file.Close()
                    if err != nil {
                        fmt.Println(os.Stderr, err)
                    }

                    newFilename := this.rotateName(filename, timestr)

                    err = os.Rename(filename, newFilename)
                    if err != nil {
                        fmt.Println(os.Stderr, err)
                    }

                    this.w, err = OpenFile(filename)
                    if err != nil {
                        fmt.Println(os.Stderr, err)
                    }
                }
            }

            lastMsgTime = msg.Time
            this.w.Write(this.msg2bytes(msg))
            this.wg.Done()
        }
    }()
}


/* Check if need to rotate file log. If yes, return value is a time string used in filename. If not, return a empty string.

Parameters:
    last        last log message's time.
    current     current log message's time.
*/
func (this *Logger) ifRotate(last, current time.Time) string {

    if this.Rotate == R_NONE {
        return ""
    }

    var format string

    switch this.Rotate {
        case R_HOURLY:
            format = "2006-01-02_15"
        case R_DAILY:
            format = "2006-01-02"
        case R_MONTHLY:
            format = "2006-01"
    }

    lastTime := last.Format(format)
    currentTime := current.Format(format)

    if lastTime != currentTime {
        return lastTime
    }

    return ""
}


/* Generate a filename (include path) for file log rotate.

Parameters:
    filename    filename of current log
    timestr     time string that will be added to filename
*/
func (this *Logger) rotateName(filename, timestr string) string {

    basename := filepath.Base(filename)

    idx := strings.LastIndex(basename, ".")

    var base, ext string

    if idx >= 0 {
        base = basename[:idx]
        ext = basename[idx:]
    } else {
        base = basename
    }

    replacer := strings.NewReplacer("{time}", timestr, "{basename}", base, "{ext}", ext)
    return filepath.Join(filepath.Dir(filename), replacer.Replace(this.RotatePattern))
}


func (this *Logger) msg2bytes(m Message) []byte {

    if this.Utc {
        m.Time = m.Time.UTC()
    }

    replacer := strings.NewReplacer("{time}", m.Time.Format(this.TimeFormat), 
                    "{level}", level2string(m.Level),
                    "{msg}", m.Msg)
    s := replacer.Replace(this.LayoutStyle)

    var b []byte
    if len(s) > 0 && s[len(s)-1] != '\n' {
        b = []byte(s + "\n")
    } else {
        b = []byte(s)
    }

    return b
}


// implement for io.Writer
func (this *Logger) Write(b []byte) (int, error) {
    m := newMsg(string(b), INFO)
    this.wg.Add(1)
    this.jobs <- m
    return len(m.Msg), nil
}


func (this *Logger) Print(level int, v ...interface{}) (int, error) {
    if level < this.Level {
        return 0, nil
    }
    m := newMsg(fmt.Sprint(v...), level)
    this.wg.Add(1)
    this.jobs <- m
    return len(m.Msg), nil
}


func (this *Logger) Printf(level int, format string, v ...interface{}) (int, error) {
    if level < this.Level {
        return 0, nil
    }
    m := newMsg(fmt.Sprintf(format, v...), level)
    this.wg.Add(1)
    this.jobs <- m
    return len(m.Msg), nil
}


// Wait termination of the log writing goroutine.
func (this *Logger) Wait() {
    this.wg.Wait()
}


// --------------------------------------------
// 10 convenient method to output log message.


func (this *Logger) Debug(v ...interface{}) {
    this.Print(DEBUG, v...)
}


func (this *Logger) Debugf(format string, v ...interface{}) {
    this.Printf(DEBUG, format, v...)
}


func (this *Logger) Info(v ...interface{}) {
    this.Print(INFO, v...)
}


func (this *Logger) Infof(format string, v ...interface{}) {
    this.Printf(INFO, format, v...)
}


func (this *Logger) Warn(v ...interface{}) {
    this.Print(WARN, v...)
}


func (this *Logger) Warnf(format string, v ...interface{}) {
    this.Printf(WARN, format, v...)
}


func (this *Logger) Error(v ...interface{}) {
    this.Print(ERROR, v...)
}


func (this *Logger) Errorf(format string, v ...interface{}) {
    this.Printf(ERROR, format, v...)
}


func (this *Logger) Fatal(v ...interface{}) {
    this.Print(FATAL, v...)
    this.Wait()
    os.Exit(1)
}


func (this *Logger) Fatalf(format string, v ...interface{}) {
    this.Printf(FATAL, format, v...)
    this.Wait()
    os.Exit(1)
}


// --------------------------------------------
// stdLogger


var stdLogger *Logger


func init() {

    var config Config
    config.LayoutStyle  = LS_SIMPLE
    config.Layout       = LY_TIME
    config.Level        = INFO

    var err error

    stdLogger, err = New(os.Stdout, config)
    if err != nil {
        panic(err.Error())
    }
}


// Output log message directly into stdout.
func Output(v ...interface{}) {
    stdLogger.Info(v...)
    stdLogger.Wait()
}


// Output log message directly into stdout, like fmt.Printf.
func Outputf(format string, v ...interface{}) {
    stdLogger.Infof(format, v...)
    stdLogger.Wait()
}
