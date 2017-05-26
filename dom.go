// Copyright 2017 Santhosh Kumar Tekuri. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dom

import (
	"errors"
	"fmt"
)

// A Node is an interface holding one of the types:
// *Document, *Element, *Text, *Comment, *ProcInst, *Attr or *NS.
type Node interface {
	Parent() Parent
	SetParent(p Parent)
}

type NS struct {
	Prefix string
	URI    string
}

func (*NS) Parent() Parent {
	return nil
}

func (*NS) SetParent(p Parent) {}

type Name struct {
	*NS
	Local string
}

func (n *Name) String() string {
	if n.Prefix == "" {
		return n.Local
	}
	return fmt.Sprintf("%s:%s", n.Prefix, n.Local)
}

// A ParentNode is an interface holding one of the types:
// *Document or *Element.
type Parent interface {
	Node
	Append(child Node) error
	Children() []Node
}

type Text struct {
	ParentNode Parent
	Data       string
}

func (t *Text) Parent() Parent {
	return t.ParentNode
}

func (t *Text) SetParent(p Parent) {
	t.ParentNode = p
}

type Comment struct {
	ParentNode Parent
	Data       string
}

func (c *Comment) Parent() Parent {
	return c.ParentNode
}

func (c *Comment) SetParent(p Parent) {
	c.ParentNode = p
}

type ProcInst struct {
	ParentNode Parent
	Target     string
	Data       string
}

func (pi *ProcInst) Parent() Parent {
	return pi.ParentNode
}

func (pi *ProcInst) SetParent(p Parent) {
	pi.ParentNode = p
}

type Attr struct {
	Owner *Element
	*Name
	Value string
}

func (*Attr) Parent() Parent {
	return nil
}

func (*Attr) SetParent(p Parent) {}

type Element struct {
	ParentNode Parent
	*Name
	NSDecl     []*NS
	Attrs      []*Attr
	ChildNodes []Node
}

func (e *Element) Parent() Parent {
	return e.ParentNode
}

func (e *Element) SetParent(p Parent) {
	e.ParentNode = p
}

func (e *Element) Append(child Node) error {
	switch child.(type) {
	case *Element, *Text, *Comment, *ProcInst:
		e.ChildNodes = append(e.ChildNodes, child)
		child.SetParent(e)
		return nil
	default:
		return fmt.Errorf("child of type %T is not allowed in *dom.Element", child)
	}
}

func (e *Element) Children() []Node {
	return e.ChildNodes
}

func (e *Element) resolvePrefix(prefix string) (string, bool) {
	if prefix == "xml" {
		return "http://www.w3.org/XML/1998/namespace", true
	}
	for {
		for _, ns := range e.NSDecl {
			if prefix == ns.Prefix {
				return ns.URI, true
			}
		}
		if _, ok := e.Parent().(*Element); ok {
			e = e.Parent().(*Element)
		} else {
			break
		}
	}
	return "", prefix == ""
}

type Document struct {
	ChildNodes []Node
}

func (*Document) Parent() Parent {
	return nil
}

func (*Document) SetParent(p Parent) {}

func (d *Document) Append(child Node) error {
	switch child.(type) {
	case *Element:
		if d.RootElement() != nil {
			return errors.New("document cannot have more than one element")
		}
	case *ProcInst, *Comment:
		// allowed
	default:
		return fmt.Errorf("child of type %T is not allowed in *dom.Document", child)
	}
	d.ChildNodes = append(d.ChildNodes, child)
	child.SetParent(d)
	return nil
}

func (d *Document) Children() []Node {
	return d.ChildNodes
}

func (d *Document) RootElement() *Element {
	for _, c := range d.ChildNodes {
		if e, ok := c.(*Element); ok {
			return e
		}
	}
	return nil
}
