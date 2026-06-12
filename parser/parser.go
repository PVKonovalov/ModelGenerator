/*
 *  parser.go
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

package parser

import (
	"ModelGenerator/scl"
	"ModelGenerator/scl/communication"
	"ModelGenerator/scl/model"
	"ModelGenerator/scl/types"
	"fmt"
	"io"
)

// SclParser parses an IEC 61850 SCL/XML file.
type SclParser struct {
	WithOutput       bool
	typeDeclarations *types.TypeDeclarations
	ieds             []*model.IED
	comm             *communication.Communication
}

// NewSclParser parses the SCL document from the given reader.
func NewSclParser(r io.Reader, withOutput bool) (*SclParser, error) {
	p := &SclParser{WithOutput: withOutput}

	root, err := scl.ParseXMLDocument(r)
	if err != nil {
		return nil, fmt.Errorf("parse XML: %w", err)
	}

	sclRoot := findSCLRoot(root)
	if sclRoot == nil {
		return nil, fmt.Errorf("no SCL root element found")
	}

	if withOutput {
		fmt.Println("parse data type templates ...")
	}

	p.typeDeclarations, err = parseTypeDeclarations(sclRoot)
	if err != nil {
		return nil, err
	}

	if withOutput {
		fmt.Println("parse IED section ...")
	}

	if err := p.parseIedSections(sclRoot); err != nil {
		return nil, err
	}

	if withOutput {
		fmt.Println("parse communication section ...")
	}

	p.comm, err = parseCommunicationSection(sclRoot)
	if err != nil {
		return nil, err
	}

	if p.comm == nil && withOutput {
		fmt.Println("WARNING: No communication section found!")
	}

	return p, nil
}

func findSCLRoot(root *scl.XMLNode) *scl.XMLNode {
	if root.Name == "SCL" {
		return root
	}
	return scl.GetChildNodeWithTag(root, "SCL")
}

func parseTypeDeclarations(sclRoot *scl.XMLNode) (*types.TypeDeclarations, error) {
	td := types.NewTypeDeclarations()

	dtNode := scl.GetChildNodeWithTag(sclRoot, "DataTypeTemplates")
	if dtNode == nil {
		return nil, fmt.Errorf("no DataTypeTemplates section found in SCL file")
	}

	for _, child := range dtNode.Children {
		switch child.Name {
		case "LNodeType":
			t, err := types.NewLogicalNodeType(child)
			if err != nil {
				return nil, err
			}
			td.AddType(t)
		case "DOType":
			t, err := types.NewDataObjectType(child)
			if err != nil {
				return nil, err
			}
			td.AddType(t)
		case "DAType":
			t, err := types.NewDataAttributeType(child)
			if err != nil {
				return nil, err
			}
			td.AddType(t)
		case "EnumType":
			t, err := types.NewEnumerationType(child)
			if err != nil {
				return nil, err
			}
			td.AddType(t)
		}
	}

	return td, nil
}

func (p *SclParser) parseIedSections(sclRoot *scl.XMLNode) error {
	for _, child := range scl.GetChildNodesWithTag(sclRoot, "IED") {
		ied, err := model.NewIED(child, p.typeDeclarations)
		if err != nil {
			return fmt.Errorf("IED: %w", err)
		}
		p.ieds = append(p.ieds, ied)
	}
	return nil
}

func parseCommunicationSection(sclRoot *scl.XMLNode) (*communication.Communication, error) {
	comNode := scl.GetChildNodeWithTag(sclRoot, "Communication")
	if comNode == nil {
		return nil, nil
	}
	return communication.NewCommunication(comNode)
}

func (p *SclParser) GetTypeDeclarations() *types.TypeDeclarations   { return p.typeDeclarations }
func (p *SclParser) GetIeds() []*model.IED                          { return p.ieds }
func (p *SclParser) GetCommunication() *communication.Communication { return p.comm }

func (p *SclParser) GetIedByName(name string) *model.IED {
	for _, ied := range p.ieds {
		if ied.GetName() == name {
			return ied
		}
	}
	return nil
}

func (p *SclParser) GetFirstIed() *model.IED {
	if len(p.ieds) == 0 {
		return nil
	}
	return p.ieds[0]
}

// GetMainIed returns the first IED that has an AccessPoint with a Server.
func (p *SclParser) GetMainIed() *model.IED {
	for _, ied := range p.ieds {
		for _, ap := range ied.GetAccessPoints() {
			if ap.GetServer() != nil {
				return ied
			}
		}
	}
	return nil
}

// GetConnectedAPsAll returns all ConnectedAPs across all SubNetworks.
func (p *SclParser) GetConnectedAPsAll() []*communication.ConnectedAP {
	if p.comm == nil {
		return nil
	}
	var caps []*communication.ConnectedAP
	for _, sn := range p.comm.SubNetworks {
		caps = append(caps, sn.ConnectedAPs...)
	}
	return caps
}

// GetConnectedAP returns the ConnectedAP for an IED and access point name.
func (p *SclParser) GetConnectedAP(iedName, apName string) *communication.ConnectedAP {
	if p.comm == nil {
		return nil
	}
	cap := p.comm.GetConnectedAP(iedName, apName)
	if cap != nil && p.WithOutput {
		fmt.Printf("Found connectedAP %s for IED %s\n", apName, iedName)
	}
	return cap
}
