package path

import "os"
import "errors"


var (
    ErrNotExist       = errors.New("file does not exist")
    ErrNotFileButDir  = errors.New("not a file, but a directory")
    ErrEmptyFile      = errors.New("file is empty")
)


// IsExistFile check if a path is exist and is a file (not a directory), if true, return nil
func IsExistFile(p string) error {
    info, err := os.Stat(p)

    if os.IsNotExist(err) {
        return ErrNotExist
    }
    if err != nil {
        return err
    }

    if info.IsDir() {
        return ErrNotFileButDir
    }

    return nil
}


// IsNonEmptyFile check if a file is exist and not empty, if true, return nil
func IsNonEmptyFile(p string) error {
    info, err := os.Stat(p)

    if os.IsNotExist(err) {
        return ErrNotExist
    }
    if err != nil {
        return err
    }

    if info.IsDir() {
        return ErrNotFileButDir
    }

    if info.Size() == 0 {
        return ErrEmptyFile
    }

    return nil
}
