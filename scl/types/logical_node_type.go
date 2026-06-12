/*
 *  logical_node_type.go
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

// LogicalNodeType represents an LNodeType from DataTypeTemplates.
type LogicalNodeType struct {
	*SclType
	LnClass     string
	DataObjects []*scl.DataObjectDefinition
}

// NewLogicalNodeType parses an LNodeType XML node.
func NewLogicalNodeType(node *scl.XMLNode) (*LogicalNodeType, error) {
	base, err := NewSclType(node)
	if err != nil {
		return nil, err
	}
	lnClass := scl.ParseAttribute(node, "lnClass")
	if lnClass == "" {
		return nil, fmt.Errorf("no lnClass attribute in LNodeType %q", base.ID)
	}

	lnt := &LogicalNodeType{SclType: base, LnClass: lnClass}

	for _, child := range scl.GetChildNodesWithTag(node, "DO") {
		dod, err := scl.NewDataObjectDefinition(child)
		if err != nil {
			return nil, err
		}
		if lnt.getByName(dod.Name) != nil {
			return nil, fmt.Errorf("logical node contains multiple data objects with name %q", dod.Name)
		}
		lnt.DataObjects = append(lnt.DataObjects, dod)
	}

	return lnt, nil
}

func (lnt *LogicalNodeType) getByName(name string) *scl.DataObjectDefinition {
	for _, d := range lnt.DataObjects {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (lnt *LogicalNodeType) GetDataObjectDefinitions() []*scl.DataObjectDefinition {
	return lnt.DataObjects
}
