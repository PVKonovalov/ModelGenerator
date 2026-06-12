/*
 *  sub_network.go
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

import "ModelGenerator/scl"

// SubNetwork represents a communication sub-network.
type SubNetwork struct {
	Name         string
	Desc         string
	Type         string
	ConnectedAPs []*ConnectedAP
}

// NewSubNetwork parses a SubNetwork XML node.
func NewSubNetwork(node *scl.XMLNode) (*SubNetwork, error) {
	sn := &SubNetwork{}
	sn.Name = scl.ParseAttribute(node, "name")
	sn.Desc = scl.ParseAttribute(node, "desc")
	sn.Type = scl.ParseAttribute(node, "type")

	for _, child := range scl.GetChildNodesWithTag(node, "ConnectedAP") {
		cap, err := NewConnectedAP(child)
		if err != nil {
			return nil, err
		}
		sn.ConnectedAPs = append(sn.ConnectedAPs, cap)
	}

	return sn, nil
}

// GetConnectedAP returns the ConnectedAP for the given IED name and AP name, or nil.
func (sn *SubNetwork) GetConnectedAP(iedName, apName string) *ConnectedAP {
	for _, cap := range sn.ConnectedAPs {
		if cap.IEDName == iedName && cap.APName == apName {
			return cap
		}
	}
	return nil
}
