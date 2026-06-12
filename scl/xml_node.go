/*
 *  xml_node.go
 *
 *  Copyright 2014-2024 Michael Zillgith
 *  Copyright 2026 Pavel Konovalov Golang port
 *
 *  This file is part of libIEC61850.
 *
 *  libIEC61850 is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  libIEC61850 is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with libIEC61850.  If not, see <http://www.gnu.org/licenses/>.
 *
 *  See COPYING file for the complete license text.
 */

package scl

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// XMLNode is a DOM-like XML node mirroring org.w3c.dom.Node.
type XMLNode struct {
	Name       string
	Attrs      map[string]string
	Children   []*XMLNode
	Text       string
	LineNumber int
	ColNumber  int
}

// ParseXMLDocument parses an XML document into a tree of XMLNodes.
func ParseXMLDocument(r io.Reader) (*XMLNode, error) {
	decoder := xml.NewDecoder(r)
	decoder.Strict = false

	var root *XMLNode
	var stack []*XMLNode

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("XML parse error: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			node := &XMLNode{
				Name:  t.Name.Local,
				Attrs: make(map[string]string),
			}
			for _, attr := range t.Attr {
				node.Attrs[attr.Name.Local] = attr.Value
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, node)
			}
			stack = append(stack, node)

		case xml.EndElement:
			if len(stack) == 0 {
				return nil, fmt.Errorf("unexpected end element %s", t.Name.Local)
			}
			node := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				root = node
			}

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && len(stack) > 0 {
				stack[len(stack)-1].Text += string(t)
			}
		}
	}

	return root, nil
}
