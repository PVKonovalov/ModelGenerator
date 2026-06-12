/*
 *  smv.go
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

// SMV holds Sampled Value address parameters.
type SMV struct {
	CBRef   string
	LDInst  string
	Address *PhyComAddress
}

// NewSMV parses an SMV XML node.
func NewSMV(node *scl.XMLNode) (*SMV, error) {
	s := &SMV{}
	s.CBRef = scl.ParseAttribute(node, "cbName")
	s.LDInst = scl.ParseAttribute(node, "ldInst")

	if s.LDInst == "" || s.CBRef == "" {
		return nil, fmt.Errorf("SMV is missing required attribute")
	}

	addrNode := scl.GetChildNodeWithTag(node, "Address")
	if addrNode == nil {
		return nil, fmt.Errorf("SMV is missing address definition")
	}
	phy, err := NewPhyComAddress(addrNode)
	if err != nil {
		return nil, fmt.Errorf("SMV Address: %w", err)
	}
	s.Address = phy

	return s, nil
}
