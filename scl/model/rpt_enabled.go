/*
 *  rpt_enabled.go
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
	"strconv"
)

// RptEnabled holds report enablement settings.
type RptEnabled struct {
	MaxInstances int
	Desc         string
	ClientLNs    []*ClientLN
}

// NewRptEnabled parses an RptEnabled XML node.
func NewRptEnabled(node *scl.XMLNode) *RptEnabled {
	r := &RptEnabled{MaxInstances: 1}
	r.Desc = scl.ParseAttribute(node, "desc")
	if maxStr := scl.ParseAttribute(node, "max"); maxStr != "" {
		if n, err := strconv.Atoi(maxStr); err == nil {
			r.MaxInstances = n
		}
	}
	for _, child := range scl.GetChildNodesWithTag(node, "ClientLN") {
		r.ClientLNs = append(r.ClientLNs, NewClientLN(child))
	}
	return r
}

func (r *RptEnabled) GetMaxInstances() int      { return r.MaxInstances }
func (r *RptEnabled) GetDesc() string           { return r.Desc }
func (r *RptEnabled) GetClientLNs() []*ClientLN { return r.ClientLNs }
