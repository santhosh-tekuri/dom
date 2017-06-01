// Copyright 2017 Santhosh Kumar Tekuri. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dom

import (
	"errors"
	"fmt"
)

// A Node is an interface holding one of the types:
// *Document, *Element, *Text, *Comment, *ProcInst, *Attr or *Namespace.
type Node interface {
	Parent() Parent
	SetParent(p Parent)
}

// A Name represents an XML name.
type Name struct {
	URI    string
	Prefix string
	Local  string
}

// String returns qualified name
func (n *Name) String() string {
	if n.Prefix == "" {
		return n.Local
	}
	return fmt.Sprintf("%s:%s", n.Prefix, n.Local)
}

// A Parent is an interface holding one of the types:
// *Document or *Element.
type Parent interface {
	Node
	Append(child Node) error
	Children() []Node
}

// A Text represents XML character data (raw text),
// in which XML escape sequences have been replaced by
// the characters they represent.
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

// A Comment represents an XML comment of the form <!--comment-->.
// The bytes do not include the <!-- and --> comment markers.
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

// A ProcInst represents an XML processing instruction of the form <?target data?>
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

// An Attr represents an attribute in an XML element (Name=Value).
type Attr struct {
	Owner *Element
	*Name
	Value string
	Type  string
}

func (*Attr) Parent() Parent {
	return nil
}

func (*Attr) SetParent(Parent) {}

// An Element represents an XML element.
type Element struct {
	ParentNode Parent
	*Name
	NSDecl     map[string]string
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

func (e *Element) declareNS(prefix, uri string) {
	if e.NSDecl == nil {
		e.NSDecl = make(map[string]string)
	}
	e.NSDecl[prefix] = uri
}

// ResolvePrefix returns the URI bound to given prefix.
// The second return value tells whether prefix is bound or not.
func (e *Element) ResolvePrefix(prefix string) (string, bool) {
	if prefix == "xml" {
		return "http://www.w3.org/XML/1998/namespace", true
	}
	for {
		if uri, ok := e.NSDecl[prefix]; ok {
			return uri, true
		}
		if _, ok := e.Parent().(*Element); ok {
			e = e.Parent().(*Element)
		} else {
			break
		}
	}
	return "", prefix == ""
}

// GetAttr returns the attribute with given uri and local.
// It returns null, if attribute is not found.
func (e *Element) GetAttr(uri, local string) *Attr {
	for _, attr := range e.Attrs {
		if attr.URI == uri && attr.Local == local {
			return attr
		}
	}
	return nil
}

// A Document represents XML Document.
type Document struct {
	ChildNodes []Node
}

func (*Document) Parent() Parent {
	return nil
}

func (*Document) SetParent(Parent) {}

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

// RootElement returns the root element of the document
func (d *Document) RootElement() *Element {
	for _, c := range d.ChildNodes {
		if e, ok := c.(*Element); ok {
			return e
		}
	}
	return nil
}

// A NameSpace represents namespace node.
// This is only used by xpath engines.
type NameSpace struct {
	Owner  *Element
	Prefix string
	URI    string
}

func (*NameSpace) Parent() Parent {
	return nil
}

func (n *NameSpace) SetParent(Parent) {}

// Owner returns the node who owns the node
//
// for *Attr and *Namespace it returns Owner field,
// for others it returns their parent Node
func Owner(n Node) Node {
	switch n := n.(type) {
	case *Attr:
		return n.Owner
	case *NameSpace:
		return n.Owner
	default:
		return n.Parent()
	}
}

// OwnerDocument returns The Document object associated with given node.
func OwnerDocument(n Node) *Document {
	for n != nil {
		if d, ok := n.(*Document); ok {
			return d
		}
		n = Owner(n)
	}
	return nil
}
