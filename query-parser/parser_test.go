package queryparser

import "fmt"
import "testing"

func TestParser(t *testing.T) {

    var p = func(data string) (s string, err error) {
        var nodes *Nodes
        nodes, err = Parse(data)
        if err != nil {
            return
        }
        s = fmt.Sprintf("%v", *nodes)
        return
    }

    // input and expect output values
    data := []string {
        `k:v`,                                  `[{k [v] false}]`,
        `k:'v'`,                                `[{k [v] false}]`,
        `k:"v"`,                                `[{k [v] false}]`,
        `'k':v`,                                `[{k [v] false}]`,
        `"k":'v'`,                              `[{k [v] false}]`,
        `'k':'v'`,                              `[{k [v] false}]`,
        `"k":"v"`,                              `[{k [v] false}]`,
        `k:'v'`,                                `[{k [v] false}]`,

        `k:v `,                                 `[{k [v] false}]`,
        `k:'v' `,                               `[{k [v] false}]`,
        `k:"v" `,                               `[{k [v] false}]`,
        `'k':'v' `,                             `[{k [v] false}]`,
        `'k':"v" `,                             `[{k [v] false}]`,
        `"k":"v" `,                             `[{k [v] false}]`,
        `"k":'v' `,                             `[{k [v] false}]`,

        `' k':'v'`,                             `[{ k [v] false}]`,
        `' k':"v"`,                             `[{ k [v] false}]`,
        `" k":"v"`,                             `[{ k [v] false}]`,
        `" k":'v'`,                             `[{ k [v] false}]`,

        `'k ':'v'`,                             `[{k  [v] false}]`,
        `'k ':"v"`,                             `[{k  [v] false}]`,
        `"k ":"v"`,                             `[{k  [v] false}]`,
        `"k ":'v'`,                             `[{k  [v] false}]`,
        `'k':' v'`,                             `[{k [ v] false}]`,
        `'k':" v"`,                             `[{k [ v] false}]`,
        `"k":" v"`,                             `[{k [ v] false}]`,
        `"k":' v'`,                             `[{k [ v] false}]`,

        `'k':'v '`,                             `[{k [v ] false}]`,
        `'k':"v "`,                             `[{k [v ] false}]`,
        `"k":"v "`,                             `[{k [v ] false}]`,
        `"k":'v '`,                             `[{k [v ] false}]`,
        `k1:"v1",v3`,                           `[{k1 [v1 v3] false}]`,

        `键:值`,                                `[{键 [值] false}]`,
        `-键:值`,                               `[{键 [值] true}]`,
        `"键":"值"`,                            `[{键 [值] false}]`,
        `-"键":"值"`,                           `[{键 [值] true}]`,
        `"-键":"值"`,                           `[{-键 [值] false}]`,

        `-k1:v1,v2`,                            `[{k1 [v1 v2] true}]`,
        `-k:"a b c"`,                           `[{k [a b c] true}]`,
        `"k":v1,v2`,                            `[{k [v1 v2] false}]`,
        `"k":"v1,v2"`,                          `[{k [v1,v2] false}]`,
        `-"k":"v1,v2"`,                         `[{k [v1,v2] true}]`,
        `-'k,!@#%':"v,!:@"`,                    `[{k,!@#% [v,!:@] true}]`,

        `k1:v1 k2:v2`,                          `[{k1 [v1] false} {k2 [v2] false}]`,
        `k1:v1 k2:v2,v3`,                       `[{k1 [v1] false} {k2 [v2 v3] false}]`,
        `k1:"v1" k2:"v2"`,                      `[{k1 [v1] false} {k2 [v2] false}]`,
        `k1:v1,v4 k2:v2,v3`,                    `[{k1 [v1 v4] false} {k2 [v2 v3] false}]`,
        `k1:"v1","v3" k2:'v2',v4`,              `[{k1 [v1 v3] false} {k2 [v2 v4] false}]`,
        `k1:"v1",'v3' k2:'v2',v4`,              `[{k1 [v1 v3] false} {k2 [v2 v4] false}]`,

        `k1:"v1",'v3' k2:'v2',"v4"`,            `[{k1 [v1 v3] false} {k2 [v2 v4] false}]`,
        `k1:v1,v3 k2:v2`,                       `[{k1 [v1 v3] false} {k2 [v2] false}]`,
        `k999:"v1",v3 k2:v2`,                   `[{k999 [v1 v3] false} {k2 [v2] false}]`,
        `uuuu:v1,"v3" k2:v2`,                   `[{uuuu [v1 v3] false} {k2 [v2] false}]`,
        `k1:"v1,v3" k2:"v2",v3`,                `[{k1 [v1,v3] false} {k2 [v2 v3] false}]`,

        `'k1':v1,v3 -"k2":v2`,                  `[{k1 [v1 v3] false} {k2 [v2] true}]`,
        `键:v1,值 -"k2":值`,                    `[{键 [v1 值] false} {k2 [值] true}]`,
        `-k1:v1,v3 k2:v2`,                      `[{k1 [v1 v3] true} {k2 [v2] false}]`,

        `k1:v1,v3 -k2:v2 "k3":v4 "k5":'v6',v7`, `[{k1 [v1 v3] false} {k2 [v2] true} {k3 [v4] false} {k5 [v6 v7] false}]`,

        `k:"v1!v2" -"k2":v3`,                   `[{k [v1!v2] false} {k2 [v3] true}]`,
        `"k":v 'k3':'haha"hehe'`,               `[{k [v] false} {k3 [haha"hehe] false}]`,
        `-"k2":v3 "k":v 'k3':'ha"he'`,          `[{k2 [v3] true} {k [v] false} {k3 [ha"he] false}]`,
        `"ab":'cd' -ef:gh i:"jk"`,              `[{ab [cd] false} {ef [gh] true} {i [jk] false}]`,
        `"xxxxx":'cd' -fffff:'gh' i:"jk"`,      `[{xxxxx [cd] false} {fffff [gh] true} {i [jk] false}]`,
        `"a:b":'c!d' k:v -"e,f":'gh' i:"j:k"`,  `[{a:b [c!d] false} {k [v] false} {e,f [gh] true} {i [j:k] false}]`,
        `'a:b':"c'd" -'e!f':"g:h" "i j":"k h"`, `[{a:b [c'd] false} {e!f [g:h] true} {i j [k h] false}]`,

        `a:b a:c`,                              `[{a [b] false} {a [c] false}]`,
        `a:b a:c a:c`,                          `[{a [b] false} {a [c] false} {a [c] false}]`,
        `a:b a:c x:y a:c`,                      `[{a [b] false} {a [c] false} {x [y] false} {a [c] false}]`,
        `a b:c cd`,                             `[{ [a] false} {b [c] false} { [cd] false}]`,
        `a "b" "c:d"`,                          `[{ [a] false} { [b] false} { [c:d] false}]`,

        `a "b" "c:d" h:i`,                      `[{ [a] false} { [b] false} { [c:d] false} {h [i] false}]`,
        `-a -"b"`,                              `[{ [a] true} { [b] true}]`,
        `-a:b -a:c`,                            `[{a [b] true} {a [c] true}]`,
        `a9999: b ,c`,                          `[{ [a9999] false} { [b] false} { [c] false}]`,
        `: b c`,                                `[{ [b] false} { [c] false}]`,
        `a: b :c`,                              `[{ [a] false} { [b] false} { [c] false}]`,

        `a111,x a,b 'c`,                        `[{ [a111 x] false} { [a b] false} { [c] false}]`,
        `a111, ,b 'c`,                          `[{ [a111] false} { [b] false} { [c] false}]`,
        `a111, ,b a:'c`,                        `[{ [a111] false} { [b] false} {a [c] false}]`,
        `a222, a,b 'c`,                         `[{ [a222] false} { [a b] false} { [c] false}]`,
        `,a111, ,b ac`,                         `[{ [a111] false} { [b] false} { [ac] false}]`,
        `:aaa: ,bbb, 'ccc' "ddd"`,              `[{ [aaa] false} { [bbb] false} { [ccc] false} { [ddd] false}]`,
        `-"not this" -"not that"`,              `[{ [not this] true} { [not that] true}]`,

        `k1,v1,'v1' k2,v2`,                     `[{ [k1 v1] false} { [k2 v2] false}]`,
        `k1,v1,v1 k1,"v2"`,                     `[{ [k1 v1] false} { [k1 v2] false}]`,
        `k1:v1,v1`,                             `[{k1 [v1] false}]`,
        `k1:v1,v1 k2:v2,"v2"`,                  `[{k1 [v1] false} {k2 [v2] false}]`,
        `k1:v1,"v1" k1:v1`,                     `[{k1 [v1] false} {k1 [v1] false}]`,

        `fid:36,37 tag:blog orderby:date author:"aa bb" -title:tt`,
        `[{fid [36 37] false} {tag [blog] false} {orderby [date] false} {author [aa bb] false} {title [tt] true}]`,

        `k1,v1,'v1',"v1",v2 "k1","v2",v2,'v2',v1,"v1"`,
        `[{ [k1 v1 v2] false} { [k1 v2 v1] false}]`,

        `from:"one@example.com" to:"another@example.com" order:asc`,
        `[{from [one@example.com] false} {to [another@example.com] false} {order [asc] false}]`,
    }

    for i := 0 ; i < len(data); i += 2 {
        input := data[i]
        expect := data[i+1]
        output, err := p(input)
        if err != nil {
            t.Error(err)
        }

        success := "SUCCESS"
        if expect != output {
            success = "FAILED"
        }
        if success == "SUCCESS" {
            t.Logf("%s:\tinput: %s\toutput: %s", success, input, output)
        } else {
            t.Errorf("%s:\tinput: %s\texpect: %s\toutput: %s", success, input, expect, output)
        }

    }
}

