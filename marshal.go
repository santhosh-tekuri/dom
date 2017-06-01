// Copyright 2017 Santhosh Kumar Tekuri. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dom

import (
	"bufio"
	"encoding/xml"
	"io"
)

// Marshal writes the XML encoding of d into w.
func Marshal(w io.Writer, d *Document) error {
	p := &printer{bufio.NewWriter(w)}
	p.printNode(d)
	return p.Flush()
}

type printer struct {
	*bufio.Writer
}

func (p *printer) printNode(n Node) {
	switch n := n.(type) {
	case *Document:
		for _, c := range n.Children() {
			p.printNode(c)
		}
	case *Element:
		p.WriteByte('<')
		p.printName(n.Name)
		for prefix, uri := range n.NSDecl {
			p.WriteByte(' ')
			p.WriteString("xmlns")
			if prefix != "" {
				p.WriteByte(':')
				p.WriteString(prefix)
			}
			p.WriteByte('=')
			p.printValue(uri)
		}
		for _, attr := range n.Attrs {
			p.WriteByte(' ')
			p.printName(attr.Name)
			p.WriteByte('=')
			p.printValue(attr.Value)
		}
		if len(n.Children()) == 0 {
			p.WriteByte('/')
			p.WriteByte('>')
		} else {
			p.WriteByte('>')
			for _, c := range n.Children() {
				p.printNode(c)
			}
			p.WriteByte('<')
			p.WriteByte('/')
			p.printName(n.Name)
			p.WriteByte('>')
		}

	case *Text:
		xml.EscapeText(p.Writer, []byte(n.Data))
	case *Comment:
		p.WriteString("<!--")
		p.WriteString(n.Data)
		p.WriteString("-->")
	case *ProcInst:
		p.WriteString("<?")
		p.WriteString(n.Target)
		if len(n.Data) > 0 {
			p.WriteByte(' ')
			p.WriteString(n.Data)
		}
		p.WriteString("?>")
	}
}

func (m *printer) printName(n *Name) {
	if n.Prefix != "" {
		m.WriteString(n.Prefix)
		m.WriteByte(':')
	}
	m.WriteString(n.Local)
}

func (m *printer) printValue(s string) {
	m.WriteByte('"')
	xml.EscapeText(m.Writer, []byte(s))
	m.WriteByte('"')
}
