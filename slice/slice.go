package slice

import "fmt"
import "reflect"


/* Determine if b is in slice a
Parameters:
    a   a slice (e.g. []int, []string, etc.)
    b   an element (e.g. int, string, etc.)
*/
func InWithError(a, b interface{}) (exist bool, err error) {

    va := reflect.ValueOf(a)

    if va.Kind() != reflect.Slice {
        err = fmt.Errorf("Parameter a must be a slice.")
        return
    }

    if reflect.TypeOf(a).String()[2:] != reflect.TypeOf(b).String() {
        err = fmt.Errorf("Type of parameter b not match with parameter a.")
        return
    }

    for i := 0; i < va.Len(); i++ {
        if va.Index(i).Interface() == b {
            exist = true
            return
        }
    }

    return
}


// Determine if b is in slice a
func In(a, b interface{}) bool {
    exist, _ := InWithError(a, b)
    return exist
}


// Remove duplicate element of a slice.
// If error occurs, return value of r will be nil.
func UniqueWithError(s interface{}) (r interface{}, err error) {

    v := reflect.ValueOf(s)

    if v.Kind() != reflect.Slice {
        err = fmt.Errorf("Parameter a must be a slice.")
        return
    }

    rest := reflect.MakeSlice(reflect.TypeOf(s), 0, 0)

    for i := v.Len() - 1; i >= 0; i-- {
        current := reflect.ValueOf(s).Slice(0, i)

        exist, e := InWithError(current.Interface(), v.Index(i).Interface())
        if e != nil {
            err = e
            return
        }
        if !exist {
            rest = reflect.Append(rest, v.Index(i))
        }
        s = current.Interface()
    }

    n := reflect.ValueOf(s)
    for i := rest.Len() - 1; i >= 0; i-- {
        n = reflect.Append(n, rest.Index(i))
    }

    r = n.Interface()

    return
}


// Remove duplicate element of a slice.
func Unique(s interface{}) interface{} {
    result, err := UniqueWithError(s)
    if err != nil {
        result = s
    }
    return result
}

