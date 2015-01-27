package log

import "testing"
import "os"
import "io/ioutil"
import "fmt"
import "time"

func TestLoggerStdout(t *testing.T) {
    var config Config
    config.Layout        = LY_DEFAULT
    config.LayoutStyle   = LS_DEFAULT
    config.Level         = INFO
    config.TimeFormat    = TF_LONG
    config.Utc           = false

    logger, err := New(os.Stdout, config)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    defer logger.Wait()

    logger.Debug("Test Debug(). This line will not output.")
    logger.Debugf("Test Debugf(). This line will not output.")
    logger.Info("Test Info().")
    logger.Infof("Test Infof(): %s, %d", "string", 123)
    logger.Warn("Test Warn(). ", "warning")
    logger.Warnf("Test Warnf(): %s", "warning")
    logger.Error("Test Error().")
    logger.Errorf("Test Errorf().")

    logger.Print(WARN, "Test Print().")
    logger.Printf(WARN, "Test Printf().")
    logger.Write([]byte("Test Write()."))
}


// Test to write log to a temporary file.
func TestLoggerFile(t *testing.T) {
    file, err := ioutil.TempFile("", "test_log_")
    if err != nil {
        t.Error(err)
        t.Fail()
    }

    filename := file.Name()
    t.Logf("Create temp file: %s", filename)

    var config Config
    config.Layout        = LY_DEFAULT
    config.LayoutStyle   = LS_DEFAULT
    config.Level         = INFO
    config.TimeFormat    = TF_TIMELONG
    config.Utc           = true

    logger, err := New(file, config)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    logger.Debug("Test Debug().")
    logger.Info("Test Info().")
    logger.Write([]byte("Test Write()."))

    t.Logf("Write log message to temp file: %s", filename)

    err = file.Close()
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    t.Logf("Close temp file: %s", filename)

    err = os.Remove(filename)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    t.Logf("Remove temp file: %s", filename)

    logger.Wait()
}


func TestSimpleLogger(t *testing.T) {
    Output("std logger")
    Outputf("std %s", "logger")
}


func ExampleLogger() {

    var config Config
    config.Layout        = LY_LEVEL
    config.LayoutStyle   = "{level}: {msg}"
    config.Level         = DEBUG
    config.TimeFormat    = TF_DEFAULT
    config.Utc           = false

    logger, err := New(os.Stdout, config)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    logger.Debug("Test Debug().")
    logger.Info("Test Info(). ", "abc", 123)
    logger.Warnf("Test Warn(). %s", "string")
    logger.Wait()
    // Output: DEBUG: Test Debug().
    // INFO: Test Info(). abc123
    // WARN: Test Warn(). string
}


func ExampleLogger_ifRotate() {

    file, err := ioutil.TempFile("", "test_log_")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    logger, err := New(file, Config{Rotate:R_DAYLY})
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    t1, _ := time.Parse("2006-01-02 15:04:05", "2015-01-26 13:58:44")
    t2, _ := time.Parse("2006-01-02 15:04:05", "2015-02-27 13:59:44")

    fmt.Println(logger.ifRotate(t1, t2))
    // Output: 2015-01-26
}


// Benchmark for writing log to stdout.
func BenchmarkLoggerStdout(b *testing.B) {

    b.StopTimer()

        var config Config
        config.Layout        = LY_TIME
        config.LayoutStyle   = LS_SIMPLE
        config.Level         = DEBUG
        config.TimeFormat    = TF_DEFAULT
        config.Utc           = true

        logger, err := New(os.Stdout, config)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

    b.StartTimer()

        for i := 0; i < b.N; i++ {
            logger.Debug("Test writing log to stdout. The more you write, the slower the speed is.")
        }

    b.StopTimer()

        logger.Wait()
}


// Benchmark for writing log to file.
func BenchmarkLoggerFile(b *testing.B) {

    b.StopTimer()

        file, err := ioutil.TempFile("", "test_log_")
        if err != nil {
            b.Error(err)
            b.Fail()
        }

        filename := file.Name()

        var config Config
        config.Layout        = LY_TIME
        config.LayoutStyle   = LS_SIMPLE
        config.Level         = DEBUG
        config.TimeFormat    = TF_DEFAULT
        config.Utc           = true

        logger, err := New(file, config)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

    b.StartTimer()

        for i := 0; i < b.N; i++ {
            logger.Debug("Test writing log to file. The more you write, the slower the speed is.")
        }

    b.StopTimer()

        err = os.Remove(filename)
        if err != nil {
            b.Error(err)
            b.Fail()
        }

        logger.Wait()
}

