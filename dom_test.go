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
	}
	for i, test := range tests {
		d, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(test)))
		if err != nil {
			t.Errorf("#%d: %s", i, err)
			continue
		}
		buf := new(bytes.Buffer)
		if err := dom.Marshal(buf, d); err != nil {
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
		{`<e a='v"'/>`, `<e a="v&#34;"/>`},
		{`<a>one<![CDATA[two]]>three<![CDATA[four]]>five</a>`, `<a>onetwothreefourfive</a>`},
	}
	for i, test := range tests {
		d, err := dom.Unmarshal(xml.NewDecoder(strings.NewReader(test.raw)))
		if err != nil {
			t.Errorf("#%d: %s", i, err)
			return
		}
		buf := new(bytes.Buffer)
		if err := dom.Marshal(buf, d); err != nil {
			t.Errorf("#%d: %s", i, err)
		}
		if s := buf.String(); s != test.normalized {
			t.Errorf("expected:\n%s\nbut got:\n%s\n", test.normalized, s)
		}
	}
}
