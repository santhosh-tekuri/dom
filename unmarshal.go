// Copyright 2017 Santhosh Kumar Tekuri. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dom

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Unmarshal reads tokens from decoder and constructs
// Document object.
func Unmarshal(decoder *xml.Decoder) (*Document, error) {
	d := new(Document)
	var cur Parent = d
	var elem *Element
	for {
		t, err := decoder.RawToken()
		if err == io.EOF {
			switch {
			case elem != nil:
				return nil, fmt.Errorf("expected </%s>", elem.Name)
			case d.RootElement() == nil:
				return nil, errors.New("document is empty")
			}
			return d, nil
		} else if err != nil {
			return d, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			elem = new(Element)
			if err = cur.Append(elem); err != nil {
				return nil, err
			}
			cur = elem
			for _, a := range t.Attr {
				if a.Name.Space == "xmlns" {
					if a.Value == "" {
						return nil, errors.New("empty namespace is not allowed")
					}
					elem.declareNS(a.Name.Local, a.Value)
				} else if a.Name.Space == "" && a.Name.Local == "xmlns" {
					elem.declareNS("", a.Value)
				}
			}
			elem.Name = translate(elem, t.Name)
			if elem.Name == nil {
				return nil, errors.New("unresolved prefix: " + t.Name.Space)
			}
			for _, a := range t.Attr {
				if a.Name.Space == "xmlns" || (a.Name.Space == "" && a.Name.Local == "xmlns") {
					continue
				}
				var name *Name
				if a.Name.Space == "" {
					name = &Name{"", "", a.Name.Local}
				} else {
					name = translate(elem, a.Name)
					if name == nil {
						return nil, errors.New("unresolved prefix: " + a.Name.Space)
					}
				}
				elem.Attrs = append(elem.Attrs, &Attr{elem, name, a.Value, "CDATA"})
			}
		case xml.EndElement:
			if elem.Prefix != t.Name.Space || elem.Local != t.Name.Local {
				return nil, fmt.Errorf("expected </%s>", elem.Name)
			}
			cur = elem.Parent()
			if _, ok := cur.(*Element); ok {
				elem = cur.(*Element)
			} else {
				elem = nil
			}
		case xml.CharData:
			if cur == elem {
				if len(elem.Children()) > 0 {
					last := elem.Children()[len(elem.Children())-1]
					if text, ok := last.(*Text); ok {
						text.Data += string(t)
						break
					}
				}
				_ = cur.Append(&Text{Data: string(t)})
			} else if len(strings.TrimSpace(string(t))) > 0 {
				return nil, errors.New("child of type *dom.Text is not allowed in *dom.Document")
			}
		case xml.Comment:
			_ = cur.Append(&Comment{Data: string(t)})
		case xml.ProcInst:
			if cur == d && t.Target == "xml" {
				break // don't add xml declaration to document
			}
			_ = cur.Append(&ProcInst{Target: t.Target, Data: string(t.Inst)})
		}
	}
}

func translate(e *Element, name xml.Name) *Name {
	if uri, ok := e.ResolvePrefix(name.Space); ok {
		return &Name{uri, name.Space, name.Local}
	}
	return nil
}
