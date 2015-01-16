package set

import "testing"
import "fmt"

func TestIsLegal(t *testing.T) {

    varBool := true
    if !IsLegal(varBool) {
        t.Error("bool should be legal.")
        t.Fail()
    }

    varInt := 123
    if !IsLegal(varInt) {
        t.Error("int should be legal.")
        t.Fail()
    }

    varSlice := []int{1,2,3}
    if IsLegal(varSlice) {
        t.Error("slice should not be legal.")
        t.Fail()
    }
}


func ExampleAdd() {
    s := New()
    fmt.Println(s.Add(123))
    fmt.Println(s.Add([]int{1,2,3}))
    // Output: true
    // false
}


func ExampleEquals() {

    s1 := New()
    s1.Add("abc", 123, 456, "xyz")
    s2 := MustNew(123, "xyz", 456, "abc")

    fmt.Println(Equals(s1, s2))
    // Output: true
}


func ExampleIsSuperset() {
    s1 := MustNew("123", "abc", 456)
    s2 := MustNew("abc", 456)

    fmt.Println(IsSuperset(s1, s2))
    // Output: true
}


func ExampleLen() {
    fmt.Println(MustNew("123", "abc", 456).Len())
    // Output: 3
}


func ExampleRemove() {
    s1 := MustNew("123", "abc", 456)
    s1.Remove("xyz")
    s1.Remove(456)
    s2 := MustNew("abc", "123")

    fmt.Println(Equals(s1, s2))
    // Output: true
}


func ExampleHas() {
    s := MustNew("123", "abc", 456)
    fmt.Println(s.Has(123))
    fmt.Println(s.Has("123"))
    fmt.Println(s.Has(456))
    // Output: false
    // true
    // true
}


func ExampleClear() {

    s := MustNew("123", "abc", 456)
    s.Clear()
    fmt.Println(s.Len())
    // Output: 0
}


func ExampleIntersect() {
    s1 := MustNew(1,2,3,4,5)
    s2 := MustNew("xyz", 123, 2, 3)
    s3 := MustNew(30, 40, 50, 3)
    s4 := MustNew(3)
    fmt.Println(Equals(s4, Intersect(s1, s2, s3)))
    // Output: true
}


func ExampleUnion() {
    s1 := MustNew(1,2)
    s2 := MustNew(3,4)
    s3 := s2.Clone()
    s3.Remove(4)
    if !s3.IsEmpty() {
        s3.MustAdd(5,6)
    }
    s4 := MustNew(1,2,3,4,5,6)
    if Equals(Union(s1, s2, s3), s4) {
        fmt.Println("yes")
    }
    // Output: yes
}


func ExampleString() {
    fmt.Println(MustNew(1).String())
    fmt.Println(MustNew("1").String())
    // Output: Set{1}
    // Set{1}
}


