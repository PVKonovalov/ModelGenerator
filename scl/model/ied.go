/*
 *  ied.go
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

// IED represents an Intelligent Electronic Device.
type IED struct {
	name         string
	desc         string
	accessPoints []*AccessPoint
	services     *Services
}

// NewIED parses an IED XML node.
func NewIED(node *scl.XMLNode, typeDeclarations *types.TypeDeclarations) (*IED, error) {
	ied := &IED{}
	ied.name = scl.ParseAttribute(node, "name")
	ied.desc = scl.ParseAttribute(node, "desc")

	if ied.name == "" {
		return nil, fmt.Errorf("required attribute \"name\" is missing in IED")
	}

	if servicesNode := scl.GetChildNodeWithTag(node, "Services"); servicesNode != nil {
		ied.services = NewServices(servicesNode)
	}

	for _, child := range scl.GetChildNodesWithTag(node, "AccessPoint") {
		ap, err := NewAccessPoint(child, typeDeclarations, ied)
		if err != nil {
			return nil, fmt.Errorf("AccessPoint: %w", err)
		}
		ied.accessPoints = append(ied.accessPoints, ap)
	}

	return ied, nil
}

func (ied *IED) GetName() string                 { return ied.name }
func (ied *IED) GetDesc() string                 { return ied.desc }
func (ied *IED) GetAccessPoints() []*AccessPoint { return ied.accessPoints }
func (ied *IED) GetServices() *Services          { return ied.services }

func (ied *IED) GetFirstAccessPoint() *AccessPoint {
	if len(ied.accessPoints) == 0 {
		return nil
	}
	return ied.accessPoints[0]
}

func (ied *IED) GetAccessPointByName(name string) *AccessPoint {
	for _, ap := range ied.accessPoints {
		if ap.GetName() == name {
			return ap
		}
	}
	return nil
}
