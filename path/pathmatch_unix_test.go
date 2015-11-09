// +build !windows

package path

import "fmt"

var root = "/tmp/dqdoc/published"

func ExamplePathMatch() {
    matchList := []string{"*.ignore", ".*", "_*", "/tmp"}

    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/.git/kk"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/.git/"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/.git"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/kkk/k2.ignore"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/kkk/lkasjdflk/_kk"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/kkk/lkasjdflk/kk"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/tmp"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/tmp/k"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/kkk/tmp/k"))
    // Output: true <nil>
    // true <nil>
    // true <nil>
    // true <nil>
    // true <nil>
    // false <nil>
    // true <nil>
    // true <nil>
    // false <nil>
}


func ExamplePathMatch1() {
    matchList := []string{"*.ignore", ".*", "_*", "/tmp", "/"}

    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/.git/kk"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/.git/"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/.git"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/kkk/k2.ignore"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/kkk/lkasjdflk/_kk"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/kkk/lkasjdflk/kk"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/tmp"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/tmp/k"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/kkk/tmp/k"))
    // Output: true <nil>
    // true <nil>
    // true <nil>
    // true <nil>
    // true <nil>
    // true <nil>
    // true <nil>
    // true <nil>
    // true <nil>
}


func ExamplePathMatch2() {
    matchList := []string{"*.ignore"}

    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/xyz/other.ignore"))
    fmt.Println(PathMatch(matchList, root, "/tmp/dqdoc/published/xyz/other.ignore/photo.jpg"))
    // Output: true <nil>
    // true <nil>
}
