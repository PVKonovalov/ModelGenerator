/*
 *  enumeration_type.go
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

// EnumerationType represents an IEC 61850 enumeration type definition.
type EnumerationType struct {
	*SclType
	EnumValues []*EnumerationValue
}

// NewEnumerationType parses an EnumType XML node.
func NewEnumerationType(node *scl.XMLNode) (*EnumerationType, error) {
	base, err := NewSclType(node)
	if err != nil {
		return nil, err
	}
	et := &EnumerationType{SclType: base}
	for _, child := range node.Children {
		if child.Name == "EnumVal" {
			ev, err := NewEnumerationValue(child)
			if err != nil {
				return nil, err
			}
			et.EnumValues = append(et.EnumValues, ev)
		}
	}
	return et, nil
}

// NewEnumerationTypeFromName creates an empty enum type with a known name.
func NewEnumerationTypeFromName(name string) *EnumerationType {
	return &EnumerationType{SclType: NewSclTypeFromName(name, "")}
}

// GetOrdByEnumString returns the ordinal for the given symbolic name.
func (et *EnumerationType) GetOrdByEnumString(enumString string) (int, error) {
	for _, v := range et.EnumValues {
		if v.SymbolicName == enumString {
			return v.Ord, nil
		}
	}
	return 0, fmt.Errorf("enum %q has no value %q", et.ID, enumString)
}

// IsValidOrdValue reports whether the ordinal value exists in the enum.
func (et *EnumerationType) IsValidOrdValue(ordValue int) bool {
	for _, v := range et.EnumValues {
		if v.Ord == ordValue {
			return true
		}
	}
	return false
}

// GetDefaultEnumTypes returns the default enum types (empty list per Java source).
func GetDefaultEnumTypes() []*EnumerationType {
	return nil
}
