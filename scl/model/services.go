/*
 *  services.go
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

import "ModelGenerator/scl"

// Services holds IED service configuration.
type Services struct {
	node *scl.XMLNode
}

// NewServices creates a Services from a node.
func NewServices(node *scl.XMLNode) *Services {
	return &Services{node: node}
}

// GetReportSettings returns the ReportSettings, or nil if not present.
func (s *Services) GetReportSettings() *ReportSettings {
	child := scl.GetChildNodeWithTag(s.node, "ReportSettings")
	if child == nil {
		return nil
	}
	return NewReportSettings(child)
}
