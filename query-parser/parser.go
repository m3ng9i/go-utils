package queryparser

import "fmt"
import "strings"
import "unicode"


/*
scan state.
out means out of key or value.
eg. "abc":123
    when the scanned char is a、b、c or 1、2、3，the state is in
    when the scanned char is quote, colon or comma，the state is out
*/
type scanState bool
const s_out scanState = false
const s_in scanState = true

type valueType bool
const v_key valueType = true
const v_value valueType = false

type quoteType string
const q_single quoteType = `'`
const q_double quoteType = `"`
const q_none quoteType = ``


type Node struct {
    Key string
    Values []string
    Negative bool
}


type Nodes []Node


func (nodes *Nodes) append(node Node) {
    if len(node.Values) > 0 {
        var n Node
        n.Key = node.Key
        n.Negative = node.Negative
        n.Values = reduceDupString(node.Values)
        if nodes != nil {
            *nodes = append(*nodes, n)
        }
    }
}


type InvalidCharError struct {
    Char string
    Pos int
    Msg string
}


func (e *InvalidCharError) Error() string {
    return fmt.Sprintf("%s: \"%s\" at position %d", e.Msg, e.Char, e.Pos)
}


func reduceDupString(s []string) (n []string) {
    var in = func(all []string, element string) bool {
        for _, i := range all {
            if i == element {
                return true
            }
        }
        return false
    }

    for _, i := range s {
        if !in(n, i) {
            n = append(n, i)
        }
    }

    return
}


func isSpecialChar(s string) (t bool, err error) {

    r := []rune(s)

    l := len(r)
    if l == 0 {
        return
    }
    if l != 1 {
        err = fmt.Errorf("string must contains only one character, got %s", s)
        return
    }

    if unicode.IsPunct(r[0]) || unicode.IsSymbol(r[0]) || unicode.IsSpace(r[0]) || unicode.IsControl(r[0]) {
        t = true
        return
    }
    return
}


func Parse(s string) (nodes *Nodes, err error) {

    defer func() {
        if e := recover(); e != nil {
            err = fmt.Errorf("PANIC: %s", e)
            return
        }
    }()

    var phrase []string
    var values []string
    var node Node
    nodes = &Nodes{}
    state := s_out
    vType := v_key
    quote := q_none

    for pos, item := range []rune(s) {

        c := string(item)

        if state == s_out {

            switch c {

                case `"`:
                    quote = q_double
                    state = s_in

                case `'`:
                    quote = q_single
                    state = s_in

                case `,`: continue
                case `:`: continue

                case `-`:
                    node.Negative = true

                case ` `:
                    if len(values) > 0 {
                        node.Values = values
                        values = []string{}
                    } else if len(node.Key) > 0 {
                        node.Values = []string{node.Key}
                        node.Key = ""
                    }
                    if len(node.Values) > 0 {
                        nodes.append(node)
                    }
                    node = Node{}
                    vType = v_key
                    quote = q_none

                default:
                    special, e := isSpecialChar(c)
                    if e != nil {
                        err = e
                        return
                    }
                    if special {
                        err = &InvalidCharError{c, pos, "Invalid character"}
                        return
                    }
                    phrase = append(phrase, c)
                    state = s_in

            } // end of switch

        } else {

            switch c {

                case `"`: fallthrough
                case `'`:
                    if quote == q_none {
                        quote = quoteType(c)
                        if len(phrase) == 0 {
                            continue
                        }
                        switch phrase[len(phrase) - 1] {
                            case `,`:
                                vType = v_value
                                values = append(values, strings.Join(phrase, ""))
                                phrase = []string{}
                            case `:`:
                                vType = v_value
                            default:
                                err = &InvalidCharError{c, pos, "Invalid character"}
                                return
                        }
                    } else if quote == quoteType(c) {
                        if vType == v_value {
                            values = append(values, strings.Join(phrase, ""))
                            phrase = []string{}
                            state = s_out
                        }
                        quote = q_none
                    } else {
                        phrase = append(phrase, c)
                    }

                case `,`:
                    if vType == v_key {
                        if quote == q_none {
                            vType = v_value
                            values = append(values, strings.Join(phrase, ""))
                            phrase = []string{}
                        } else {
                            phrase = append(phrase, c)
                        }
                    } else {
                        if quote == q_none {
                            values = append(values, strings.Join(phrase, ""))
                            phrase = []string{}
                            state = s_out
                        } else {
                            phrase = append(phrase, c)
                        }
                    }

                case `:`:
                    if quote != q_none {
                        phrase = append(phrase, c)
                    } else {
                        if vType == v_key {
                            vType = v_value
                            node.Key = strings.Join(phrase, "")
                            phrase = []string{}
                            state = s_out
                        } else {
                            err = &InvalidCharError{c, pos, "Cannot appear more than once"}
                            return
                        }
                    }

                case `-`:
                    phrase = append(phrase, c)

                case ` `:
                    if quote != q_none {
                        phrase = append(phrase, c)
                    } else {

                        if vType == v_key {
                            node.Values = []string{strings.Join(phrase, "")}
                        } else {
                            if len(phrase) > 0 {
                                node.Values = append(values, strings.Join(phrase, ""))
                            } else {
                                node.Values = values
                            }
                        }

                        if (vType == v_key && node.Key == "") || vType == v_value {
                            nodes.append(node)
                            node = Node{}
                        }

                        state = s_out
                        vType = v_key
                        values = []string{}
                        phrase = []string{}
                    }

                default:
                    if quote == q_none {
                        special, e := isSpecialChar(c)
                        if e != nil {
                            err = e
                            return
                        }
                        if special {
                            err = &InvalidCharError{c, pos, "Invalid character"}
                            return
                        }
                    }
                    phrase = append(phrase, c)

            } // end of switch
        } // end of else

        // DEBUG
        //fmt.Println("c:", c)
        //fmt.Println("phrase:", phrase)
        //fmt.Println("values:", values)
        //fmt.Println("node", node)
        //fmt.Println("nodes", nodes)
        //if state {
        //    fmt.Println("state:", "in")
        //} else {
        //    fmt.Println("state:", "out")
        //}
        //if vType {
        //    fmt.Println("vType:", "key")
        //} else {
        //    fmt.Println("vType:", "value")
        //}
        //fmt.Println("quote:", quote)
        //fmt.Println("=====")

    } // end of for

    if len(phrase) > 0 {
        values = append(values, strings.Join(phrase, ""))
    }

    if len(values) > 0 {
        node.Values = values
        nodes.append(node)
    }

    return
}

