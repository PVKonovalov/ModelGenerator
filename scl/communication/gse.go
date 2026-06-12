/*
 *  gse.go
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
	"strconv"
	"strings"
)

// GSE holds GSE address and timing parameters.
type GSE struct {
	CBRef   string
	LDInst  string
	Address *PhyComAddress
	MinTime int
	MaxTime int
}

// NewGSE parses a GSE XML node.
func NewGSE(node *scl.XMLNode) (*GSE, error) {
	g := &GSE{MinTime: -1, MaxTime: -1}
	g.CBRef = scl.ParseAttribute(node, "cbName")
	g.LDInst = scl.ParseAttribute(node, "ldInst")

	if g.CBRef == "" || g.LDInst == "" {
		return nil, fmt.Errorf("GSE is missing required attribute")
	}

	if minNode := scl.GetChildNodeWithTag(node, "MinTime"); minNode != nil {
		v, err := strconv.Atoi(strings.TrimSpace(minNode.Text))
		if err == nil {
			g.MinTime = v
		}
	}

	if maxNode := scl.GetChildNodeWithTag(node, "MaxTime"); maxNode != nil {
		v, err := strconv.Atoi(strings.TrimSpace(maxNode.Text))
		if err == nil {
			g.MaxTime = v
		}
	}

	addrNode := scl.GetChildNodeWithTag(node, "Address")
	if addrNode == nil {
		return nil, fmt.Errorf("GSE is missing address definition")
	}
	phy, err := NewPhyComAddress(addrNode)
	if err != nil {
		return nil, fmt.Errorf("GSE Address: %w", err)
	}
	g.Address = phy

	return g, nil
}

func (g *GSE) GetMinTime() int            { return g.MinTime }
func (g *GSE) GetMaxTime() int            { return g.MaxTime }
func (g *GSE) GetAddress() *PhyComAddress { return g.Address }
