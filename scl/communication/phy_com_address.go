/*
 *  phy_com_address.go
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

// PhyComAddress holds physical communication address parameters.
type PhyComAddress struct {
	VlanID       int
	VlanPriority int
	AppID        int
	MacAddress   [6]byte
}

// NewPhyComAddress extracts PhyComAddress from an Address node.
func NewPhyComAddress(node *scl.XMLNode) (*PhyComAddress, error) {
	addr := NewAddress(node)
	phy := &PhyComAddress{}

	vlanIDStr := addr.GetP("VLAN-ID")
	if vlanIDStr != "" {
		v, err := strconv.ParseInt(vlanIDStr, 16, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid VLAN-ID %q: %w", vlanIDStr, err)
		}
		phy.VlanID = int(v)
	}

	vlanPrioStr := addr.GetP("VLAN-PRIORITY")
	if vlanPrioStr != "" {
		v, err := strconv.Atoi(vlanPrioStr)
		if err != nil {
			return nil, fmt.Errorf("invalid VLAN-PRIORITY %q: %w", vlanPrioStr, err)
		}
		phy.VlanPriority = v
	}

	appIDStr := addr.GetP("APPID")
	if appIDStr != "" {
		v, err := strconv.ParseInt(appIDStr, 16, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid APPID %q: %w", appIDStr, err)
		}
		phy.AppID = int(v)
	}

	macStr := addr.GetP("MAC-Address")
	if macStr != "" {
		parts := strings.Split(macStr, "-")
		if len(parts) != 6 {
			return nil, fmt.Errorf("invalid MAC-Address %q", macStr)
		}
		for i, part := range parts {
			b, err := strconv.ParseUint(part, 16, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid MAC-Address byte %q: %w", part, err)
			}
			phy.MacAddress[i] = byte(b)
		}
	}

	return phy, nil
}
