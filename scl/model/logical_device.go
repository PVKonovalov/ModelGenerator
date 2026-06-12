/*
 *  logical_device.go
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

// LogicalDevice represents an IED logical device instance.
type LogicalDevice struct {
	inst         string
	desc         string
	ldName       string
	logicalNodes []*LogicalNode
	parentServer *Server
}

// NewLogicalDevice parses an LDevice XML node.
func NewLogicalDevice(node *scl.XMLNode, typeDeclarations *types.TypeDeclarations, parent *Server) (*LogicalDevice, error) {
	ld := &LogicalDevice{parentServer: parent}
	ld.inst = scl.ParseAttribute(node, "inst")
	ld.desc = scl.ParseAttribute(node, "desc")
	ld.ldName = scl.ParseAttribute(node, "ldName")

	if ld.inst == "" {
		return nil, fmt.Errorf("required attribute \"inst\" is missing in LDevice")
	}

	// LN0 must come first
	for _, child := range scl.GetChildNodesWithTag(node, "LN0") {
		ln, err := NewLogicalNode(child, typeDeclarations, ld)
		if err != nil {
			return nil, fmt.Errorf("LN0: %w", err)
		}
		ld.logicalNodes = append(ld.logicalNodes, ln)
	}

	for _, child := range scl.GetChildNodesWithTag(node, "LN") {
		ln, err := NewLogicalNode(child, typeDeclarations, ld)
		if err != nil {
			return nil, fmt.Errorf("LN: %w", err)
		}
		ld.logicalNodes = append(ld.logicalNodes, ln)
	}

	return ld, nil
}

func (ld *LogicalDevice) GetInst() string                 { return ld.inst }
func (ld *LogicalDevice) GetDesc() string                 { return ld.desc }
func (ld *LogicalDevice) GetLdName() string               { return ld.ldName }
func (ld *LogicalDevice) GetLogicalNodes() []*LogicalNode { return ld.logicalNodes }
func (ld *LogicalDevice) GetParentServer() *Server        { return ld.parentServer }
