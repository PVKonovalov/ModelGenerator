/*
 *  parser_utils.go
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
	"fmt"
	"strings"
)

// ParseAttribute returns the attribute value from a node, or "" if not found.
func ParseAttribute(node *XMLNode, name string) string {
	if node == nil {
		return ""
	}
	v, ok := node.Attrs[name]
	if !ok {
		return ""
	}
	return v
}

// ParseAttributePtr returns a pointer to the attribute value, or nil if not present.
func ParseAttributePtr(node *XMLNode, name string) *string {
	if node == nil {
		return nil
	}
	v, ok := node.Attrs[name]
	if !ok {
		return nil
	}
	return &v
}

// GetChildNodeWithTag returns the first child element with the given tag name, or nil.
func GetChildNodeWithTag(node *XMLNode, tag string) *XMLNode {
	if node == nil {
		return nil
	}
	for _, child := range node.Children {
		if child.Name == tag {
			return child
		}
	}
	return nil
}

// GetChildNodesWithTag returns all child elements with the given tag name.
func GetChildNodesWithTag(node *XMLNode, tag string) []*XMLNode {
	var result []*XMLNode
	if node == nil {
		return result
	}
	for _, child := range node.Children {
		if child.Name == tag {
			result = append(result, child)
		}
	}
	return result
}

// ParseBooleanAttribute parses a boolean attribute ("true"/"false"). Returns nil if not present.
func ParseBooleanAttribute(node *XMLNode, name string) (*bool, error) {
	ptr := ParseAttributePtr(node, name)
	if ptr == nil {
		return nil, nil
	}
	s := strings.ToUpper(*ptr)
	var v bool
	if s == "TRUE" {
		v = true
	} else if s == "FALSE" {
		v = false
	} else {
		return nil, fmt.Errorf("illegal value for boolean attribute %q: %q", name, *ptr)
	}
	return &v, nil
}
