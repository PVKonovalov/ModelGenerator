/*
 *  data_attribute_definition.go
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
	"strconv"
)

// DataAttributeDefinition holds the definition of a DA or BDA element.
// The actual FunctionalConstraint and AttributeType types are defined in the model package,
// but to avoid import cycles we keep them as plain integers here and let the model package
// interpret them.
type DataAttributeDefinition struct {
	Name      string
	BType     string
	Type      string // enum type id or struct type id
	Count     int
	FCString  string // raw FC string, interpreted by model package
	Dchg      bool
	Qchg      bool
	Dupd      bool
	ValueText string // raw value text from <Val> child
	HasValue  bool
}

// NewDataAttributeDefinition parses a DA or BDA XML node.
func NewDataAttributeDefinition(node *XMLNode) (*DataAttributeDefinition, error) {
	d := &DataAttributeDefinition{}
	d.Name = ParseAttribute(node, "name")
	if d.Name == "" {
		return nil, fmt.Errorf("attribute name is missing")
	}

	d.BType = ParseAttribute(node, "bType")
	if d.BType == "" {
		return nil, fmt.Errorf("attribute bType is missing for %q", d.Name)
	}

	d.Type = ParseAttribute(node, "type")
	d.FCString = ParseAttribute(node, "fc")

	// Handle special bType values that override type
	switch d.BType {
	case "Tcmd":
		d.Type = "Tcmd"
	case "Dbpos":
		d.Type = "Dbpos"
	case "Check":
		d.Type = "Check"
	}

	var err error
	if countStr := ParseAttribute(node, "count"); countStr != "" {
		d.Count, err = strconv.Atoi(countStr)
		if err != nil {
			return nil, fmt.Errorf("invalid count value %q: %w", countStr, err)
		}
	}

	if dchgStr := ParseAttribute(node, "dchg"); dchgStr != "" {
		b, _ := strconv.ParseBool(dchgStr)
		d.Dchg = b
	}
	if dupdStr := ParseAttribute(node, "dupd"); dupdStr != "" {
		b, _ := strconv.ParseBool(dupdStr)
		d.Dupd = b
	}
	if qchgStr := ParseAttribute(node, "qchg"); qchgStr != "" {
		b, _ := strconv.ParseBool(qchgStr)
		d.Qchg = b
	}

	// Parse <Val> child
	for _, child := range node.Children {
		if child.Name == "Val" {
			d.ValueText = child.Text
			d.HasValue = true
			break
		}
	}

	return d, nil
}

func (d *DataAttributeDefinition) GetName() string     { return d.Name }
func (d *DataAttributeDefinition) GetBType() string    { return d.BType }
func (d *DataAttributeDefinition) GetType() string     { return d.Type }
func (d *DataAttributeDefinition) GetFCString() string { return d.FCString }
func (d *DataAttributeDefinition) GetCount() int       { return d.Count }

// IsConstructed returns true if this DA has bType "Struct" (constructed type).
func (d *DataAttributeDefinition) IsConstructed() bool { return d.BType == "Struct" }
