/*
 *  logical_node.go
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
	"strconv"
)

// LogicalNode is an instance of a logical node.
type LogicalNode struct {
	lnClass                   string
	lnType                    string
	inst                      string
	desc                      string
	prefix                    string
	sclType                   types.SclTypeIface
	dataObjects               []*DataObject
	dataSets                  []*DataSet
	reportControlBlocks       []*ReportControlBlock
	gseControlBlocks          []*GSEControl
	smvControlBlocks          []*SampledValueControl
	logControlBlocks          []*LogControl
	logs                      []*Log
	settingGroupControlBlocks []*SettingControl
	parentLogicalDevice       *LogicalDevice
}

// NewLogicalNode parses an LN or LN0 XML node.
func NewLogicalNode(node *scl.XMLNode, typeDeclarations *types.TypeDeclarations, parent *LogicalDevice) (*LogicalNode, error) {
	ln := &LogicalNode{parentLogicalDevice: parent}
	ln.lnClass = scl.ParseAttribute(node, "lnClass")
	ln.lnType = scl.ParseAttribute(node, "lnType")
	ln.inst = scl.ParseAttribute(node, "inst")
	ln.desc = scl.ParseAttribute(node, "desc")
	ln.prefix = scl.ParseAttribute(node, "prefix")

	if ln.lnClass == "" {
		return nil, fmt.Errorf("required attribute \"lnClass\" is missing in logical node")
	}
	if ln.lnType == "" {
		return nil, fmt.Errorf("required attribute \"lnType\" is missing in logical node")
	}
	// LN0 always has inst="" (empty string) which is valid in IEC 61850.
	// Only non-LN0 nodes require a non-empty inst.
	if ln.inst == "" && ln.lnClass != "LLN0" {
		return nil, fmt.Errorf("required attribute \"inst\" is missing in logical node")
	}

	t := typeDeclarations.LookupLogicalNodeType(ln.lnType)
	if t == nil {
		return nil, fmt.Errorf("missing type declaration %q", ln.lnType)
	}
	t.SetUsed(true)
	ln.sclType = t

	for _, doDef := range t.GetDataObjectDefinitions() {
		do, err := NewDataObject(doDef, typeDeclarations, ln)
		if err != nil {
			return nil, fmt.Errorf("data object %q: %w", doDef.Name, err)
		}
		ln.dataObjects = append(ln.dataObjects, do)
	}

	// DataSets
	for _, child := range scl.GetChildNodesWithTag(node, "DataSet") {
		ds, err := NewDataSet(child)
		if err != nil {
			return nil, err
		}
		ln.dataSets = append(ln.dataSets, ds)
	}

	// ReportControlBlocks
	for _, child := range scl.GetChildNodesWithTag(node, "ReportControl") {
		rcb, err := NewReportControlBlock(child)
		if err != nil {
			return nil, err
		}
		ln.reportControlBlocks = append(ln.reportControlBlocks, rcb)
	}

	// GSEControl
	for _, child := range scl.GetChildNodesWithTag(node, "GSEControl") {
		gse, err := NewGSEControl(child)
		if err != nil {
			return nil, err
		}
		ln.gseControlBlocks = append(ln.gseControlBlocks, gse)
	}

	// SampledValueControl
	for _, child := range scl.GetChildNodesWithTag(node, "SampledValueControl") {
		smv, err := NewSampledValueControl(child)
		if err != nil {
			return nil, err
		}
		ln.smvControlBlocks = append(ln.smvControlBlocks, smv)
	}

	// LogControl
	for _, child := range scl.GetChildNodesWithTag(node, "LogControl") {
		lc, err := NewLogControl(child)
		if err != nil {
			return nil, err
		}
		ln.logControlBlocks = append(ln.logControlBlocks, lc)
	}

	// Logs
	for _, child := range scl.GetChildNodesWithTag(node, "Log") {
		ln.logs = append(ln.logs, NewLog(child))
	}

	// SettingControl
	sgNodes := scl.GetChildNodesWithTag(node, "SettingControl")
	if ln.lnClass != "LLN0" && len(sgNodes) > 0 {
		return nil, fmt.Errorf("LN other than LN0 is not allowed to contain SettingControl")
	}
	if len(sgNodes) > 1 {
		return nil, fmt.Errorf("LN contains more than one SettingControl")
	}
	for _, child := range sgNodes {
		sg, err := NewSettingControl(child)
		if err != nil {
			return nil, err
		}
		ln.settingGroupControlBlocks = append(ln.settingGroupControlBlocks, sg)
	}

	// Parse DOI instances (override default values)
	for _, doiNode := range scl.GetChildNodesWithTag(node, "DOI") {
		doiName := scl.ParseAttribute(doiNode, "name")

		var idx int = -1
		if ixAttr := scl.ParseAttribute(doiNode, "ix"); ixAttr != "" {
			n, err := strconv.Atoi(ixAttr)
			if err != nil {
				return nil, fmt.Errorf("invalid ix attribute in %q: %w", doiName, err)
			}
			idx = n
		}

		doNode := ln.GetChildByName(doiName)
		if doNode == nil {
			return nil, fmt.Errorf("missing data object with name %q", doiName)
		}

		if err := parseDataAttributeNodes(doiNode, doNode, idx); err != nil {
			return nil, err
		}
		if err := parseSubDataInstances(doiNode, doNode); err != nil {
			return nil, err
		}
	}

	return ln, nil
}

func parseDataAttributeNodes(doiNode *scl.XMLNode, parent DataModelNode, parentIdx int) error {
	for _, daiNode := range scl.GetChildNodesWithTag(doiNode, "DAI") {
		daiName := scl.ParseAttribute(daiNode, "name")

		idx := parentIdx
		if ixAttr := scl.ParseAttribute(daiNode, "ix"); ixAttr != "" {
			n, err := strconv.Atoi(ixAttr)
			if err != nil {
				return fmt.Errorf("invalid ix in %q: %w", daiName, err)
			}
			idx = n
		}
		_ = idx

		daNode := parent.GetChildByName(daiName)
		if daNode == nil {
			return fmt.Errorf("missing data attribute with name %q", daiName)
		}
		da, ok := daNode.(*DataAttribute)
		if !ok {
			continue
		}

		if valNode := scl.GetChildNodeWithTag(daiNode, "Val"); valNode != nil {
			valueText := valNode.Text
			var dmv *DataModelValue
			var err error
			if da.GetType() == AT_ENUMERATED {
				dmv = NewDataModelValueFromEnumType(
					func() string {
						if da.GetSclType() != nil {
							return da.GetSclType().GetID()
						}
						return ""
					}(),
					valueText,
				)
			} else {
				dmv, err = NewDataModelValue(da.GetType(), da.GetSclType(), valueText)
				if err != nil {
					return fmt.Errorf("DAI %q value: %w", daiName, err)
				}
			}
			da.SetValue(dmv)
		}

		if shortAddr := scl.ParseAttribute(daiNode, "sAddr"); shortAddr != "" {
			da.SetShortAddress(shortAddr)
		}
	}
	return nil
}

func parseSubDataInstances(doiNode *scl.XMLNode, parent DataModelNode) error {
	for _, sdiNode := range scl.GetChildNodesWithTag(doiNode, "SDI") {
		sdiName := scl.ParseAttribute(sdiNode, "name")
		child := parent.GetChildByName(sdiName)
		if child == nil {
			fmt.Printf("subelement with name %q not found!\n", sdiName)
			continue
		}
		if err := parseDataAttributeNodes(sdiNode, child, -1); err != nil {
			return err
		}
		if err := parseSubDataInstances(sdiNode, child); err != nil {
			return err
		}
	}
	return nil
}

func (ln *LogicalNode) GetLnClass() string                            { return ln.lnClass }
func (ln *LogicalNode) GetLnType() string                             { return ln.lnType }
func (ln *LogicalNode) GetInst() string                               { return ln.inst }
func (ln *LogicalNode) GetDesc() string                               { return ln.desc }
func (ln *LogicalNode) GetPrefix() string                             { return ln.prefix }
func (ln *LogicalNode) GetDataObjects() []*DataObject                 { return ln.dataObjects }
func (ln *LogicalNode) GetDataSets() []*DataSet                       { return ln.dataSets }
func (ln *LogicalNode) GetReportControlBlocks() []*ReportControlBlock { return ln.reportControlBlocks }
func (ln *LogicalNode) GetGSEControlBlocks() []*GSEControl            { return ln.gseControlBlocks }
func (ln *LogicalNode) GetSampledValueControlBlocks() []*SampledValueControl {
	return ln.smvControlBlocks
}
func (ln *LogicalNode) GetLogControlBlocks() []*LogControl { return ln.logControlBlocks }
func (ln *LogicalNode) GetLogs() []*Log                    { return ln.logs }
func (ln *LogicalNode) GetSettingGroupControlBlocks() []*SettingControl {
	return ln.settingGroupControlBlocks
}
func (ln *LogicalNode) GetParentLogicalDevice() *LogicalDevice { return ln.parentLogicalDevice }
func (ln *LogicalNode) GetSclType() types.SclTypeIface         { return ln.sclType }
func (ln *LogicalNode) GetParent() DataModelNode               { return nil }

// GetName returns the full LN name (prefix + lnClass + inst).
func (ln *LogicalNode) GetName() string {
	name := ""
	if ln.prefix != "" {
		name += ln.prefix
	}
	name += ln.lnClass
	name += ln.inst
	return name
}

func (ln *LogicalNode) GetChildByName(name string) DataModelNode {
	for _, do := range ln.dataObjects {
		if do.name == name {
			return do
		}
	}
	return nil
}
