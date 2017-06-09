# dom

[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![GoDoc](https://godoc.org/github.com/santhosh-tekuri/dom?status.svg)](https://godoc.org/github.com/santhosh-tekuri/dom)
[![Go Report Card](https://goreportcard.com/badge/github.com/santhosh-tekuri/dom)](https://goreportcard.com/report/github.com/santhosh-tekuri/dom)
[![Build Status](https://travis-ci.org/santhosh-tekuri/dom.svg?branch=master)](https://travis-ci.org/santhosh-tekuri/dom)
[![codecov.io](https://codecov.io/github/santhosh-tekuri/dom/coverage.svg?branch=master)](https://codecov.io/github/santhosh-tekuri/dom?branch=master)

Package dom provides document object model for xml.

It does not strictly follow DOM interfaces, but has everything needed for xml processing library.

## Example

```go
str := `
<developer xmlns="www.jroller.com/santhosh/">
<name>Santhosh Kumar Tekuri</name>
<email>santhosh.tekuri@gmail.com</email>
</developer>
`

doc, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(str)))
if err != nil {
    fmt.Println(err)
    return
}

root := doc.RootElement()
fmt.Printf("rootElement: {%s}%s\n", root.URI, root.Local)
buf := new(bytes.Buffer)
if err = dom.Marshal(doc, buf); err != nil {
    fmt.Println(err)
    return
}
fmt.Printf("xml:\n%s", buf.String())
```

Output:
```
rootElement: {www.jroller.com/santhosh/}developer
xml:
<developer xmlns="www.jroller.com/santhosh/">
    <name>Santhosh Kumar Tekuri</name>
    <email>santhosh.tekuri@gmail.com</email>
</developer>
```

