package cmd

import "os"
import "os/exec"
import "bytes"
import "io"


// Run run a shell command, return content of stdout, stderr.
func Run(name string, arg ...string) (sout, serr []byte, err error) {
    c := exec.Command(name, arg ...)

    var so bytes.Buffer
    var se bytes.Buffer
    c.Stdout = &so
    c.Stderr = &se

    err = c.Run()
    sout = so.Bytes()
    serr = se.Bytes()

    return
}


// Call run a shell command and print the result to stdout and stderr.
func Call(name string, arg ...string) error {
    sout, serr, err := Run(name, arg...)
    io.Copy(os.Stdout, bytes.NewBuffer(sout))
    io.Copy(os.Stderr, bytes.NewBuffer(serr))
    return err
}
