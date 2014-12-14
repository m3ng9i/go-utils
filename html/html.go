package html

import "html"
import "strings"

/* Convert text to html

Example:

text := `
First line

Second line
Third line
    
<b id="html">line contains html</b>
`
html := Text2Html(text)

now html will be:
<p>First line</p><p>Second line<br>Third line</p><p>&lt;b id=&#34;html&#34;&gt;line contains html&lt;/b&gt;</p>

*/
func Text2Html(text string) string {
    text = html.EscapeString(text)
    var h []string

    newPara := true
    for _, line := range(strings.Split(text, "\n")) {
        l := strings.Trim(line, "\r\n\t ")
        if newPara {
            if l == "" {
                continue
            }
            newPara = false
            h = append(h, "<p>" + l)
        } else {
            if l == "" {
                h = append(h, "</p>")
                newPara = true
            } else {
                h = append(h, "<br>" + l)
            }
        }
    }

    if newPara == false {
        h = append(h, "</p>")
    }

    return strings.Join(h, "")
}

