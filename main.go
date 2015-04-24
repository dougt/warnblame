package main

import (
    "flag"
    "regexp"
    "bufio"
    "fmt"
    "os"
    "sort"
)

// Shamelessly stolen from https://gist.github.com/kylelemons/1236125
// but made it sort from great to less
type ValSorter struct {
        Keys []string
        Vals []int
}
 
func NewValSorter(m map[string]int) *ValSorter {
        vs := &ValSorter{
                Keys: make([]string, 0, len(m)),
                Vals: make([]int, 0, len(m)),
        }
        for k, v := range m {
                vs.Keys = append(vs.Keys, k)
                vs.Vals = append(vs.Vals, v)
        }
        return vs
}
 
func (vs *ValSorter) Sort() {
        sort.Sort(vs)
}
 
func (vs *ValSorter) Len() int           { return len(vs.Vals) }
func (vs *ValSorter) Less(i, j int) bool { return vs.Vals[i] > vs.Vals[j] }
func (vs *ValSorter) Swap(i, j int) {
        vs.Vals[i], vs.Vals[j] = vs.Vals[j], vs.Vals[i]
        vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}


func main() {

    flag.Parse()
    fh, err := os.Open(flag.Arg(0))
    f := bufio.NewReader(fh)

    if err != nil {
        fmt.Println("outa")
        return
    }
    defer fh.Close()

    buf := make([]byte, 4096)
    m := make(map[string]int)
    filepath := "/Users/dougt/builds/mozilla-inbound/"
    re := regexp.MustCompile(".*WARNING: .*failed: file " + filepath + "(.*), line (.*)")

    for {
        buf, _ , err = f.ReadLine()
        if err != nil {
            break
        }

        segs := re.FindAllStringSubmatch(string(buf), -1)
        for i := 0; i < len(segs); i++ {
            key := segs[i][1] + "#l" + segs[i][2]
            j := m[key]
            m[key] = j+1
        }
    }

    vs := NewValSorter(m)
    vs.Sort()



    hgblame := "https://hg.mozilla.org/mozilla-central/annotate/tip/"

    fmt.Printf("<html><head><title>Warning Blame</title><style>");
    fmt.Printf("</style></head><body>");
    fmt.Printf("<h1> The Warning Count Is Too Damn High!</h1><h2>")

    for i := 0; i < len(vs.Keys); i++ {
        key := vs.Keys[i];
        fmt.Printf("<p>%d <a href=\"%s\">%s</a>\n", m[key], hgblame+key, key)
    }
    fmt.Printf("</h2></body></html>") 


}
