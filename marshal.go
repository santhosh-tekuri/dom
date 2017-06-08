// Copyright 2017 Santhosh Kumar Tekuri. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dom

import (
	"bufio"
	"io"
	"unicode/utf8"
)

var (
	esc_quot = []byte("&quot;")
	esc_apos = []byte("&apos;")
	esc_amp  = []byte("&amp;")
	esc_lt   = []byte("&lt;")
	esc_gt   = []byte("&gt;")
	esc_tab  = []byte("&#x9;")
	esc_nl   = []byte("&#xA;")
	esc_cr   = []byte("&#xD;")
	esc_fffd = []byte("\uFFFD") // Unicode replacement character
)

// Marshal writes the XML encoding of d into w.
func Marshal(d *Document, w io.Writer) error {
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
		p.escapeString(n.Data, false)
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

func (p *printer) printName(n *Name) {
	if n.Prefix != "" {
		p.WriteString(n.Prefix)
		p.WriteByte(':')
	}
	p.WriteString(n.Local)
}

func (p *printer) escapeString(s string, escapeNewline bool) {
	var esc []byte
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width
		switch r {
		case '"':
			esc = esc_quot
		case '\'':
			esc = esc_apos
		case '&':
			esc = esc_amp
		case '<':
			esc = esc_lt
		case '>':
			esc = esc_gt
		case '\t':
			esc = esc_tab
		case '\n':
			if !escapeNewline {
				continue
			}
			esc = esc_nl
		case '\r':
			esc = esc_cr
		default:
			if !isInCharacterRange(r) || (r == 0xFFFD && width == 1) {
				esc = esc_fffd
				break
			}
			continue
		}
		p.WriteString(s[last : i-width])
		p.Write(esc)
		last = i
	}
	p.WriteString(s[last:])
}

func (p *printer) printValue(s string) {
	p.WriteByte('"')
	p.escapeString(s, true)
	p.WriteByte('"')
}

// Decide whether the given rune is in the XML Character Range, per
// the Char production of http://www.xml.com/axml/testaxml.htm,
// Section 2.2 Characters.
func isInCharacterRange(r rune) bool {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xDF77 ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}
