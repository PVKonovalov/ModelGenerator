/*
 *  address.go
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

// Address holds a collection of P parameters.
type Address struct {
	Parameters []*P
}

// NewAddress parses an Address XML node.
func NewAddress(node *scl.XMLNode) *Address {
	a := &Address{}
	for _, child := range scl.GetChildNodesWithTag(node, "P") {
		a.Parameters = append(a.Parameters, NewP(child))
	}
	return a
}

// GetP returns the P parameter value for the given type, or empty string.
func (a *Address) GetP(pType string) string {
	for _, p := range a.Parameters {
		if p.Type == pType {
			return p.Value
		}
	}
	return ""
}
