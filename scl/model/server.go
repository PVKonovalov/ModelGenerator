/*
 *  server.go
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

// Server represents an IED server.
type Server struct {
	authentication    *Authentication
	logicalDevices    []*LogicalDevice
	parentAccessPoint *AccessPoint
}

// NewServer parses a Server XML node.
func NewServer(node *scl.XMLNode, typeDeclarations *types.TypeDeclarations, parent *AccessPoint) (*Server, error) {
	s := &Server{parentAccessPoint: parent}

	if authNode := scl.GetChildNodeWithTag(node, "Authentication"); authNode != nil {
		s.authentication = NewAuthentication(authNode)
	}

	for _, child := range scl.GetChildNodesWithTag(node, "LDevice") {
		ld, err := NewLogicalDevice(child, typeDeclarations, s)
		if err != nil {
			return nil, fmt.Errorf("LDevice: %w", err)
		}
		s.logicalDevices = append(s.logicalDevices, ld)
	}

	return s, nil
}

func (s *Server) GetAuthentication() *Authentication  { return s.authentication }
func (s *Server) GetLogicalDevices() []*LogicalDevice { return s.logicalDevices }
func (s *Server) GetParentAccessPoint() *AccessPoint  { return s.parentAccessPoint }
