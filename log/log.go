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


// Rotate config when outputing to a regular file.
const (
    R_NONE = iota   // Don't rotate.
    R_HOURLY        // Rotate every hour.
    R_DAYLY         // Rotate every day.
    R_MONTHLY       // Rotate every month.
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


// open a file to write log
func OpenFile(filename string, filemode ...os.FileMode) (*os.File, error) {
    var mode os.FileMode = 0640
    if len(filemode) > 0 {
        mode = filemode[0]
    }
    return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, mode)
}


// ------------------------------------------------
// Logger


type Config struct {
    TimeFormat      string          // Time format in output log message.
    LayoutStyle     string          // Log messae layout style.
    Layout          int             // Log message layout element.
    Level           int             // Log level.
    Utc             bool            // If use utc time in output.
    Rotate          int             // How to rotate file log.
    RotatePattern   string          // Filename rotate pattern of a file log. See RP_DEFAULT for example.
}


type Logger struct {
    Config
    w io.Writer
    jobs chan message
    wg  sync.WaitGroup
}


type message struct {
    msg []byte
    mtime time.Time
    level int
}


// Create a Logger.
func New(w io.Writer, config Config) (logger *Logger, err error) {
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

    if config.Rotate > R_NONE {
        _, ok := w.(*os.File)
        if !ok {
            config.Rotate = R_NONE
            config.RotatePattern = ""
        } else {
            err = isRotatePatternLegal(config.RotatePattern)
            if err != nil {
                return
            }
        }
    }


    logger.Config   = config
    logger.w        = w
    logger.jobs     = make(chan message, maxJobs)

    logger.start()

    return
}


// Start to receive logging jobs, then write log message to io.Writer. This function will rotate file log as needed.
func (this *Logger) start() {
    go func() {
        var lastMsgTime time.Time
        for {
            msg := <- this.jobs

            // check if need to rotate file log
            if this.Rotate > R_NONE && !lastMsgTime.IsZero() {
                timestr := this.ifRotate(lastMsgTime, msg.mtime)

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

            lastMsgTime = msg.mtime
            this.w.Write(msg.msg)
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
        case R_DAYLY:
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


func (this *Logger) msg(s string, level ...int) message {

    var m message
    m.mtime = time.Now()

    if this.Utc {
        m.mtime = m.mtime.UTC()
    }

    var nowstring, lev string

    if this.Layout & LY_TIME > 0 {
        nowstring = m.mtime.Format(this.TimeFormat)
    }

    if this.Layout & LY_LEVEL > 0 {
        if len(level) > 0 {
            lev = level2string(level[0])
            m.level = level[0]
        } else {
            lev = level2string(DEBUG)
            m.level = DEBUG
        }
    }

    replacer := strings.NewReplacer("{time}", nowstring, "{level}", lev, "{msg}", s)
    s = replacer.Replace(this.LayoutStyle)
    if len(s) > 0 && s[len(s)-1] != '\n' {
        m.msg = []byte(s + "\n")
    } else {
        m.msg = []byte(s)
    }

    return m
}


// implement for io.Writer
func (this *Logger) Write(b []byte) (int, error) {
    msg := this.msg(string(b), INFO)
    this.wg.Add(1)
    this.jobs <- msg
    return len(msg.msg), nil
}


func (this *Logger) Print(level int, v ...interface{}) (int, error) {
    if level < this.Level {
        return 0, nil
    }
    msg := this.msg(fmt.Sprint(v...), level)
    this.wg.Add(1)
    this.jobs <- msg
    return len(msg.msg), nil
}


func (this *Logger) Printf(level int, format string, v ...interface{}) (int, error) {
    if level < this.Level {
        return 0, nil
    }
    msg := this.msg(fmt.Sprintf(format, v...), level)
    this.wg.Add(1)
    this.jobs <- msg
    return len(msg.msg), nil
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
