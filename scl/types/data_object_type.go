/*
 *  data_object_type.go
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

package types

import (
	"ModelGenerator/scl"
	"fmt"
)

// DataObjectType represents a DOType from DataTypeTemplates.
type DataObjectType struct {
	*SclType
	CDC            string
	DataAttributes []*scl.DataAttributeDefinition
	SubDataObjects []*scl.DataObjectDefinition
}

// NewDataObjectType parses a DOType XML node.
func NewDataObjectType(node *scl.XMLNode) (*DataObjectType, error) {
	base, err := NewSclType(node)
	if err != nil {
		return nil, err
	}
	cdc := scl.ParseAttribute(node, "cdc")
	if cdc == "" {
		return nil, fmt.Errorf("cdc is missing in DOType %q", base.ID)
	}

	dot := &DataObjectType{SclType: base, CDC: cdc}

	for _, child := range node.Children {
		switch child.Name {
		case "DA":
			dad, err := scl.NewDataAttributeDefinition(child)
			if err != nil {
				return nil, err
			}
			// Check duplicate name (allow different FC)
			existing := dot.getDAByName(dad.Name)
			if existing != nil && existing.FCString == dad.FCString {
				return nil, fmt.Errorf("DO type definition contains multiple elements of name %q", dad.Name)
			}
			if dot.getDOByName(dad.Name) != nil {
				return nil, fmt.Errorf("DO type definition contains multiple elements of name %q", dad.Name)
			}
			dot.DataAttributes = append(dot.DataAttributes, dad)

		case "SDO":
			dod, err := scl.NewDataObjectDefinition(child)
			if err != nil {
				return nil, err
			}
			if dot.getDAByName(dod.Name) != nil || dot.getDOByName(dod.Name) != nil {
				return nil, fmt.Errorf("DO type definition contains multiple elements of name %q", dod.Name)
			}
			dot.SubDataObjects = append(dot.SubDataObjects, dod)
		}
	}

	return dot, nil
}

func (dot *DataObjectType) getDAByName(name string) *scl.DataAttributeDefinition {
	for _, d := range dot.DataAttributes {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (dot *DataObjectType) getDOByName(name string) *scl.DataObjectDefinition {
	for _, d := range dot.SubDataObjects {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (dot *DataObjectType) GetCDC() string { return dot.CDC }
func (dot *DataObjectType) GetDataAttributes() []*scl.DataAttributeDefinition {
	return dot.DataAttributes
}
func (dot *DataObjectType) GetSubDataObjects() []*scl.DataObjectDefinition { return dot.SubDataObjects }
