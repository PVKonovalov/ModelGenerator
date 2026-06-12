/*
 *  data_object.go
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

package model

import (
	"ModelGenerator/scl"
	"ModelGenerator/scl/types"
	"fmt"
)

// DataObject is an instance of a data object in the model hierarchy.
type DataObject struct {
	name           string
	count          int
	sclType        types.SclTypeIface
	parent         DataModelNode
	dataAttributes []*DataAttribute
	subDataObjects []*DataObject
	trans          bool
}

// NewDataObject creates a DataObject from a definition.
func NewDataObject(
	doDef *scl.DataObjectDefinition,
	typeDeclarations *types.TypeDeclarations,
	parent DataModelNode,
) (*DataObject, error) {

	do := &DataObject{
		name:   doDef.Name,
		count:  doDef.Count,
		parent: parent,
		trans:  doDef.IsTransient(),
	}

	t := typeDeclarations.LookupDataObjectType(doDef.Type)
	if t == nil {
		return nil, fmt.Errorf("type declaration missing for data object %q (type %q)", doDef.Name, doDef.Type)
	}
	t.SetUsed(true)
	do.sclType = t

	if err := do.createDataAttributes(typeDeclarations, t); err != nil {
		return nil, err
	}
	if err := do.createSubDataObjects(typeDeclarations, t); err != nil {
		return nil, err
	}

	return do, nil
}

func (do *DataObject) createDataAttributes(td *types.TypeDeclarations, doType *types.DataObjectType) error {
	for _, daDef := range doType.GetDataAttributes() {
		// SE attributes: add both SG and SE copies
		if daDef.FCString == "SE" {
			sgDA, err := NewDataAttribute(daDef, td, FC_SG, do)
			if err != nil {
				return err
			}
			do.dataAttributes = append(do.dataAttributes, sgDA)
		}
		da, err := NewDataAttribute(daDef, td, FC_NONE, do)
		if err != nil {
			return err
		}
		do.dataAttributes = append(do.dataAttributes, da)
	}
	return nil
}

func (do *DataObject) createSubDataObjects(td *types.TypeDeclarations, doType *types.DataObjectType) error {
	for _, sdoDef := range doType.GetSubDataObjects() {
		sdo, err := NewDataObject(sdoDef, td, do)
		if err != nil {
			return err
		}
		do.subDataObjects = append(do.subDataObjects, sdo)
	}
	return nil
}

func (do *DataObject) GetName() string                     { return do.name }
func (do *DataObject) GetCount() int                       { return do.count }
func (do *DataObject) IsTransient() bool                   { return do.trans }
func (do *DataObject) GetDataAttributes() []*DataAttribute { return do.dataAttributes }
func (do *DataObject) GetSubDataObjects() []*DataObject    { return do.subDataObjects }
func (do *DataObject) GetSclType() types.SclTypeIface      { return do.sclType }
func (do *DataObject) GetParent() DataModelNode            { return do.parent }

func (do *DataObject) GetChildByName(name string) DataModelNode {
	for _, da := range do.dataAttributes {
		if da.name == name {
			return da
		}
	}
	for _, sdo := range do.subDataObjects {
		if sdo.name == name {
			return sdo
		}
	}
	return nil
}
