/*
 *  data_attribute.go
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

// DataAttribute is an instance of a data attribute in the model hierarchy.
type DataAttribute struct {
	name              string
	atType            AttributeType
	count             int
	parent            DataModelNode
	fc                FunctionalConstraint
	isBasicAttribute  bool
	value             *DataModelValue
	shortAddress      string
	subDataAttributes []*DataAttribute
	sclType           types.SclTypeIface
	triggerOptions    *TriggerOptions
	definition        *scl.DataAttributeDefinition
}

// NewDataAttribute creates a DataAttribute from a definition.
func NewDataAttribute(
	daDef *scl.DataAttributeDefinition,
	typeDeclarations *types.TypeDeclarations,
	fc FunctionalConstraint,
	parent DataModelNode,
) (*DataAttribute, error) {

	da := &DataAttribute{
		name:             daDef.Name,
		count:            daDef.Count,
		parent:           parent,
		definition:       daDef,
		isBasicAttribute: true,
	}

	// Parse FC from definition, override if passed
	if daDef.FCString != "" {
		parsed, err := FCFromString(daDef.FCString)
		if err != nil {
			return nil, err
		}
		da.fc = parsed
	} else {
		da.fc = FC_NONE
	}
	if fc != FC_NONE {
		da.fc = fc
	}

	// Parse attribute type
	atType, err := ATFromSclString(daDef.BType)
	if err != nil {
		return nil, fmt.Errorf("unsupported attribute type %q for DA %q", daDef.BType, daDef.Name)
	}
	da.atType = atType

	// Trigger options
	if parentDA, ok := parent.(*DataAttribute); ok {
		da.triggerOptions = parentDA.triggerOptions
	} else {
		da.triggerOptions = NewTriggerOptions(daDef.Dchg, daDef.Qchg, daDef.Dupd, false, false)
	}

	// Handle constructed or enumerated
	switch da.atType {
	case AT_CONSTRUCTED:
		da.isBasicAttribute = false
		if err := da.createConstructedAttribute(daDef, typeDeclarations); err != nil {
			return nil, err
		}
	case AT_ENUMERATED:
		if err := da.createEnumeratedAttribute(daDef, typeDeclarations); err != nil {
			return nil, err
		}
	}

	return da, nil
}

func (da *DataAttribute) createEnumeratedAttribute(daDef *scl.DataAttributeDefinition, td *types.TypeDeclarations) error {
	t := td.LookupType(daDef.Type)
	if t == nil {
		return fmt.Errorf("missing type definition for enumerated data attribute: %q", daDef.BType)
	}
	if _, ok := t.(*types.EnumerationType); !ok {
		return fmt.Errorf("wrong type definition for enumerated data attribute %q", daDef.Name)
	}
	t.SetUsed(true)
	da.sclType = t
	return nil
}

func (da *DataAttribute) createConstructedAttribute(daDef *scl.DataAttributeDefinition, td *types.TypeDeclarations) error {
	t := td.LookupDataAttributeType(daDef.Type)
	if t == nil {
		return fmt.Errorf("missing type definition for constructed data attribute: %q", daDef.BType)
	}
	t.SetUsed(true)
	da.sclType = t

	for _, subDef := range t.GetSubDataAttributes() {
		sub, err := NewDataAttribute(subDef, td, da.fc, da)
		if err != nil {
			return err
		}
		da.subDataAttributes = append(da.subDataAttributes, sub)
	}
	return nil
}

func (da *DataAttribute) GetName() string                             { return da.name }
func (da *DataAttribute) GetFC() FunctionalConstraint                 { return da.fc }
func (da *DataAttribute) GetType() AttributeType                      { return da.atType }
func (da *DataAttribute) GetSubDataAttributes() []*DataAttribute      { return da.subDataAttributes }
func (da *DataAttribute) IsBasicAttribute() bool                      { return da.isBasicAttribute }
func (da *DataAttribute) GetCount() int                               { return da.count }
func (da *DataAttribute) GetValue() *DataModelValue                   { return da.value }
func (da *DataAttribute) SetValue(v *DataModelValue)                  { da.value = v }
func (da *DataAttribute) GetShortAddress() string                     { return da.shortAddress }
func (da *DataAttribute) SetShortAddress(s string)                    { da.shortAddress = s }
func (da *DataAttribute) GetTriggerOptions() *TriggerOptions          { return da.triggerOptions }
func (da *DataAttribute) GetDefinition() *scl.DataAttributeDefinition { return da.definition }

// GetEffectiveValue returns the overridden value, falling back to the definition's value.
// If the definition value has no computed ordinal, callers must call UpdateEnumOrdValue.
func (da *DataAttribute) GetEffectiveValue(td *types.TypeDeclarations) *DataModelValue {
	if da.value != nil {
		return da.value
	}
	if da.definition == nil || !da.definition.HasValue {
		return nil
	}
	var v *DataModelValue
	var err error
	if da.atType == AT_ENUMERATED {
		typeID := ""
		if da.sclType != nil {
			typeID = da.sclType.GetID()
		}
		v = NewDataModelValueFromEnumType(typeID, da.definition.ValueText)
	} else {
		v, err = NewDataModelValue(da.atType, da.sclType, da.definition.ValueText)
		if err != nil {
			return nil
		}
	}
	if v != nil && v.GetValue() == nil {
		v.UpdateEnumOrdValue(td)
	}
	return v
}
func (da *DataAttribute) GetSclType() types.SclTypeIface { return da.sclType }
func (da *DataAttribute) GetParent() DataModelNode       { return da.parent }

func (da *DataAttribute) GetChildByName(name string) DataModelNode {
	for _, sub := range da.subDataAttributes {
		if sub.name == name {
			return sub
		}
	}
	return nil
}
