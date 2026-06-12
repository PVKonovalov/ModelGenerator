/*
 *  report_settings.go
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

// ReportSettings holds report settings.
type ReportSettings struct {
	node *scl.XMLNode
}

// NewReportSettings creates a ReportSettings from a node.
func NewReportSettings(node *scl.XMLNode) *ReportSettings {
	return &ReportSettings{node: node}
}

// HasOwner returns true if the owner attribute is set to true.
func (r *ReportSettings) HasOwner() bool {
	b, err := scl.ParseBooleanAttribute(r.node, "owner")
	if err != nil || b == nil {
		return false
	}
	return *b
}
