/*
 *  scl_type.go
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

// SclType is the base type for all SCL type definitions.
type SclType struct {
	ID          string
	Description string
	IsUsed      bool
}

// NewSclType parses common SCL type attributes (id, desc).
func NewSclType(node *scl.XMLNode) (*SclType, error) {
	id := scl.ParseAttribute(node, "id")
	if id == "" {
		return nil, fmt.Errorf("id is missing in type definition")
	}
	desc := scl.ParseAttribute(node, "desc")
	return &SclType{ID: id, Description: desc}, nil
}

// NewSclTypeFromName creates an SclType with a known name.
func NewSclTypeFromName(id, description string) *SclType {
	return &SclType{ID: id, Description: description}
}

func (t *SclType) GetID() string     { return t.ID }
func (t *SclType) GetDesc() string   { return t.Description }
func (t *SclType) GetIsUsed() bool   { return t.IsUsed }
func (t *SclType) SetUsed(used bool) { t.IsUsed = used }
