// Copyright 2017 Santhosh Kumar Tekuri. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dom_test

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/dom"
)

func TestIdentity(t *testing.T) {
	tests := []string{
		`<test a1="v1" a2="v2"/>`,
		`<x:test xmlns:x="ns1" x:a1="v1" a2="v2"/>`,
		`<test xmlns="ns1" a1="v1" a2="v2"/>`,
		`<x:one xmlns:x="ns1"><x:two/></x:one>`,
		`<x><!--ignore me--></x>`,
		`<!--ignore me--><x/>`,
		`<x><?abcd hello world?></x>`,
		`<?abcd hello world?><x/>`,
		`<e a="v&amp;"/>`,
		`<e a="a&lt;b"/>`,
		`<e a="a&gt;b"/>`,
		"<e>\n</e>",
		`<e xml:lang="en-us">\n</e>`,
	}
	for i, test := range tests {
		d, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(test)))
		if err != nil {
			t.Errorf("#%d: %s", i, err)
			continue
		}
		buf := new(bytes.Buffer)
		if err := dom.Marshal(d, buf); err != nil {
			t.Errorf("#%d: %s", i, err)
		}
		if s := buf.String(); s != test {
			t.Errorf("expected:\n%s\nbut got:\n%s\n", test, s)
		}
	}
}

func TestNormalized(t *testing.T) {
	tests := []struct {
		raw, normalized string
	}{
		{`<?xml version="1.0" encoding="UTF-8" standalone="no" ?><test/>`, `<test/>`},
		{` <a/>`, `<a/>`},
		{`<e a='v'/>`, `<e a="v"/>`},
		{`<e a='v"'/>`, `<e a="v&quot;"/>`},
		{`<e a="v'"/>`, `<e a="v&apos;"/>`},
		{"<e>\t</e>", `<e>&#x9;</e>`},
		{"<e>&#xD;</e>", `<e>&#xD;</e>`},
		{"<e a='&#xA;'/>", `<e a="&#xA;"/>`},
		{`<a>one<![CDATA[two]]>three<![CDATA[four]]>five</a>`, `<a>onetwothreefourfive</a>`},
	}
	for i, test := range tests {
		d, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(test.raw)))
		if err != nil {
			t.Errorf("#%d: %s", i, err)
			continue
		}
		buf := new(bytes.Buffer)
		if err := dom.Marshal(d, buf); err != nil {
			t.Errorf("#%d: %s", i, err)
		}
		if s := buf.String(); s != test.normalized {
			t.Errorf("#%d: expected:\n%s\nbut got:\n%s\n", i, test.normalized, s)
		}
	}
}

func TestInvalidXML(t *testing.T) {
	tests := []string{
		``,                                             // no root element
		`<e1`,                                          // incomplete start element
		`<e1>`,                                         // missing end element
		`<e1/><e2/>`,                                   // more than one root element
		`<ns1:e1/>`,                                    // unresolved element prefix
		`<e1 ns1:p1="v1"/>`,                            // unresolved attribute prefix
		`<ns1:x xmlns:ns1=""/>`,                        // empty namespace bound to prefix
		`<e1>hai</e2>`,                                 // wrong end element
		`<x:e xmlns:x="x" xmlns:y="x"></y:e>`,          // wrong prefix in end element
		`hai<e1/>`,                                     // text outside root element
		`<e a="v" a="v"/>`,                             // duplicate attribute
		`<e xmlns:x="x" xmlns:y="x" x:a="v" y:a="v"/>`, // duplicate attribute
		`<!--comment--xyz--><e1/>`,                     // "--" not allowed in comments
	}

	for i, test := range tests {
		if _, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(test))); err == nil {
			t.Errorf("#%d: FAIL: error expected", i)
			continue
		} else {
			t.Logf("#%d: %v", i, err)
		}
	}
}

func TestDefaultNamespace(t *testing.T) {
	doc, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(`<e xmlns="x" a="v"/>`)))
	if err != nil {
		t.Error(err)
	} else {
		if doc.RootElement().URI != "x" {
			t.Errorf("root namesame: got %q, want %q", doc.RootElement().URI, "x")
		}
		if doc.RootElement().Attrs[0].URI != "" {
			t.Errorf("attribute namesame: got %q, want %q", doc.RootElement().Attrs[0].URI, "")
		}
	}
}
