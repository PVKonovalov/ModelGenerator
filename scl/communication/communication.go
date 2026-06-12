/*
 *  communication.go
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

// Communication holds all SubNetworks from the SCL Communication section.
type Communication struct {
	SubNetworks []*SubNetwork
}

// NewCommunication parses a Communication XML node.
func NewCommunication(node *scl.XMLNode) (*Communication, error) {
	c := &Communication{}
	for _, child := range scl.GetChildNodesWithTag(node, "SubNetwork") {
		sn, err := NewSubNetwork(child)
		if err != nil {
			return nil, err
		}
		c.SubNetworks = append(c.SubNetworks, sn)
	}
	return c, nil
}

// GetConnectedAP searches all sub-networks for a ConnectedAP with the given IED and AP name.
func (c *Communication) GetConnectedAP(iedName, apName string) *ConnectedAP {
	for _, sn := range c.SubNetworks {
		cap := sn.GetConnectedAP(iedName, apName)
		if cap != nil {
			return cap
		}
	}
	return nil
}
