/*
 *  data_object_definition.go
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

// DataObjectDefinition holds the definition of a data object within an LN or DO type.
type DataObjectDefinition struct {
	Name  string
	Type  string
	Trans bool
	Count int
}

// NewDataObjectDefinition parses a DO or SDO XML node.
func NewDataObjectDefinition(node *XMLNode) (*DataObjectDefinition, error) {
	d := &DataObjectDefinition{}
	d.Name = ParseAttribute(node, "name")
	d.Type = ParseAttribute(node, "type")

	if d.Name == "" || d.Type == "" {
		return nil, fmt.Errorf("DO misses required attribute (name=%q, type=%q)", d.Name, d.Type)
	}

	isTransient, err := ParseBooleanAttribute(node, "transient")
	if err != nil {
		return nil, err
	}
	if isTransient != nil {
		d.Trans = *isTransient
	}

	if countStr := ParseAttribute(node, "count"); countStr != "" {
		d.Count, err = strconv.Atoi(countStr)
		if err != nil {
			return nil, fmt.Errorf("invalid count value %q: %w", countStr, err)
		}
	}

	return d, nil
}

func (d *DataObjectDefinition) GetName() string   { return d.Name }
func (d *DataObjectDefinition) GetType() string   { return d.Type }
func (d *DataObjectDefinition) GetCount() int     { return d.Count }
func (d *DataObjectDefinition) IsTransient() bool { return d.Trans }
