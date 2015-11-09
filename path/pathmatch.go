package path

import "fmt"
import "path/filepath"


// check if a path is match the match list
func PathMatch(matchList []string, root, p string) (match bool, err error) {

    root = filepath.Clean(root)
    p = filepath.Clean(p)

    base := filepath.Base(p)

    if !filepath.IsAbs(root) {
        err = fmt.Errorf("Root must be a absolute")
        return
    }

    if !filepath.IsAbs(p) {
        err = fmt.Errorf("Path must be a absolute")
        return
    }

    for _, item := range matchList {
        i := filepath.Join(root, item)

        // example: /root/path/page.html (p) is match to /root/path/*.html (i)
        m, e := filepath.Match(i, p)
        if e != nil {
            err = e
            return
        }
        if m {
            return true, nil
        }

        // example: /root/path/page.html (newp) is match to /root/path (i)
        newp := p
        for {
            newp = filepath.Dir(newp)
            if len(newp) < len(root) {
                break
            }

            m, e := filepath.Match(i, newp)
            if e != nil {
                err = e
                return
            }
            if m {
                return true, nil
            }

            // example: /root/path/xyz/other.ignore/path1/path2 (newp) is match to *.ignore
            if !filepath.IsAbs(item) {
                m, e := filepath.Match(item, filepath.Base(newp))
                if e != nil {
                    err = e
                    return
                }
                if m {
                    return true, nil
                }
            }
        }

        // example: page.html (base) is match to *.html (item)
        if !filepath.IsAbs(item) {
            m, e := filepath.Match(item, base)
            if e != nil {
                err = e
                return
            }
            if m {
                return true, nil
            }
        }
    }

    return false, nil
}

