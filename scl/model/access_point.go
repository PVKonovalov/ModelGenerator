/*
 *  access_point.go
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

// AccessPoint represents an IED access point.
type AccessPoint struct {
	name      string
	server    *Server
	parentIED *IED
}

// NewAccessPoint parses an AccessPoint XML node.
func NewAccessPoint(node *scl.XMLNode, typeDeclarations *types.TypeDeclarations, parent *IED) (*AccessPoint, error) {
	ap := &AccessPoint{parentIED: parent}
	ap.name = scl.ParseAttribute(node, "name")
	if ap.name == "" {
		return nil, fmt.Errorf("required attribute \"name\" is missing in AccessPoint")
	}

	serverNode := scl.GetChildNodeWithTag(node, "Server")
	if serverNode != nil {
		s, err := NewServer(serverNode, typeDeclarations, ap)
		if err != nil {
			return nil, fmt.Errorf("Server: %w", err)
		}
		ap.server = s
	}

	return ap, nil
}

func (ap *AccessPoint) GetName() string    { return ap.name }
func (ap *AccessPoint) GetServer() *Server { return ap.server }
func (ap *AccessPoint) GetParentIED() *IED { return ap.parentIED }
