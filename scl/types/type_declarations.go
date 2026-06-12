/*
 *  type_declarations.go
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

import "fmt"

// TypeDeclarations is a registry of all SCL type definitions.
type TypeDeclarations struct {
	types []SclTypeIface
}

// SclTypeIface is the interface for all SCL type variants.
type SclTypeIface interface {
	GetID() string
	GetDesc() string
	GetIsUsed() bool
	SetUsed(bool)
}

// NewTypeDeclarations creates a new registry pre-populated with default enum types.
func NewTypeDeclarations() *TypeDeclarations {
	td := &TypeDeclarations{}
	for _, et := range GetDefaultEnumTypes() {
		td.types = append(td.types, et)
	}
	return td
}

// AddType registers a type.
func (td *TypeDeclarations) AddType(t SclTypeIface) {
	td.types = append(td.types, t)
}

// LookupType finds a type by ID, optionally filtered by a type name ("LogicalNodeType", etc.).
func (td *TypeDeclarations) LookupType(typeID string) SclTypeIface {
	for _, t := range td.types {
		if t.GetID() == typeID {
			return t
		}
	}
	fmt.Printf("Cannot find type %s\n", typeID)
	return nil
}

// LookupLogicalNodeType finds a LogicalNodeType by ID.
func (td *TypeDeclarations) LookupLogicalNodeType(typeID string) *LogicalNodeType {
	for _, t := range td.types {
		if t.GetID() == typeID {
			if lnt, ok := t.(*LogicalNodeType); ok {
				return lnt
			}
		}
	}
	fmt.Printf("Cannot find type %s\n", typeID)
	return nil
}

// LookupDataObjectType finds a DataObjectType by ID.
func (td *TypeDeclarations) LookupDataObjectType(typeID string) *DataObjectType {
	for _, t := range td.types {
		if t.GetID() == typeID {
			if dot, ok := t.(*DataObjectType); ok {
				return dot
			}
		}
	}
	fmt.Printf("Cannot find type %s\n", typeID)
	return nil
}

// LookupDataAttributeType finds a DataAttributeType by ID.
func (td *TypeDeclarations) LookupDataAttributeType(typeID string) *DataAttributeType {
	for _, t := range td.types {
		if t.GetID() == typeID {
			if dat, ok := t.(*DataAttributeType); ok {
				return dat
			}
		}
	}
	fmt.Printf("Cannot find type %s\n", typeID)
	return nil
}

// LookupEnumerationType finds an EnumerationType by ID.
func (td *TypeDeclarations) LookupEnumerationType(typeID string) *EnumerationType {
	for _, t := range td.types {
		if t.GetID() == typeID {
			if et, ok := t.(*EnumerationType); ok {
				return et
			}
		}
	}
	fmt.Printf("Cannot find type %s\n", typeID)
	return nil
}

// GetTypeDeclarations returns all registered types.
func (td *TypeDeclarations) GetTypeDeclarations() []SclTypeIface {
	return td.types
}
