/*
 *  data_attribute_type.go
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

// DataAttributeType represents a DAType from DataTypeTemplates.
type DataAttributeType struct {
	*SclType
	SubDataAttributes []*scl.DataAttributeDefinition
}

// NewDataAttributeType parses a DAType XML node.
func NewDataAttributeType(node *scl.XMLNode) (*DataAttributeType, error) {
	base, err := NewSclType(node)
	if err != nil {
		return nil, err
	}
	dat := &DataAttributeType{SclType: base}

	for _, child := range node.Children {
		if child.Name == "BDA" {
			dad, err := scl.NewDataAttributeDefinition(child)
			if err != nil {
				return nil, err
			}
			if dat.getByName(dad.Name) != nil {
				return nil, fmt.Errorf("DA type definition contains multiple elements of name %q", dad.Name)
			}
			dat.SubDataAttributes = append(dat.SubDataAttributes, dad)
		}
	}

	return dat, nil
}

func (dat *DataAttributeType) getByName(name string) *scl.DataAttributeDefinition {
	for _, d := range dat.SubDataAttributes {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (dat *DataAttributeType) GetSubDataAttributes() []*scl.DataAttributeDefinition {
	return dat.SubDataAttributes
}
