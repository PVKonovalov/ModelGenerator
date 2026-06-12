/*
 *  connected_ap.go
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

package communication

import (
	"ModelGenerator/scl"
	"fmt"
)

// ConnectedAP represents a connected access point in the communication section.
type ConnectedAP struct {
	IEDName string
	APName  string
	Address *Address
	GSEs    []*GSE
	SMVs    []*SMV
}

// NewConnectedAP parses a ConnectedAP XML node.
func NewConnectedAP(node *scl.XMLNode) (*ConnectedAP, error) {
	cap := &ConnectedAP{}
	cap.IEDName = scl.ParseAttribute(node, "iedName")
	cap.APName = scl.ParseAttribute(node, "apName")

	if cap.IEDName == "" {
		return nil, fmt.Errorf("required attribute \"iedName\" is missing in ConnectedAP")
	}
	if cap.APName == "" {
		return nil, fmt.Errorf("required attribute \"apName\" is missing in ConnectedAP")
	}

	if addrNode := scl.GetChildNodeWithTag(node, "Address"); addrNode != nil {
		cap.Address = NewAddress(addrNode)
	}

	for _, child := range scl.GetChildNodesWithTag(node, "GSE") {
		gse, err := NewGSE(child)
		if err != nil {
			return nil, err
		}
		cap.GSEs = append(cap.GSEs, gse)
	}

	for _, child := range scl.GetChildNodesWithTag(node, "SMV") {
		smv, err := NewSMV(child)
		if err != nil {
			return nil, err
		}
		cap.SMVs = append(cap.SMVs, smv)
	}

	return cap, nil
}

// LookupGSE returns the GSE for the given ldInst and cbName, or nil.
func (cap *ConnectedAP) LookupGSE(ldInst, cbName string) *GSE {
	for _, gse := range cap.GSEs {
		if gse.LDInst == ldInst && gse.CBRef == cbName {
			return gse
		}
	}
	return nil
}

// LookupSMV returns the SMV for the given ldInst and cbName, or nil.
func (cap *ConnectedAP) LookupSMV(ldInst, cbName string) *SMV {
	for _, smv := range cap.SMVs {
		if smv.LDInst == ldInst && smv.CBRef == cbName {
			return smv
		}
	}
	return nil
}

// LookupSMVAddress returns the PhyComAddress for the given ldInst and cbName, or nil.
func (cap *ConnectedAP) LookupSMVAddress(ldInst, cbName string) *PhyComAddress {
	smv := cap.LookupSMV(ldInst, cbName)
	if smv == nil {
		return nil
	}
	return smv.Address
}
