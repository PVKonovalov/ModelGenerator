/*
 *  client_ln.go
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

// ClientLN holds a reference to a client logical node.
type ClientLN struct {
	node *scl.XMLNode
}

// NewClientLN creates a ClientLN from a node.
func NewClientLN(node *scl.XMLNode) *ClientLN {
	return &ClientLN{node: node}
}

func (c *ClientLN) GetIedName() string { return scl.ParseAttribute(c.node, "iedName") }
func (c *ClientLN) GetApRef() string   { return scl.ParseAttribute(c.node, "apRef") }
func (c *ClientLN) GetLdInst() string  { return scl.ParseAttribute(c.node, "ldInst") }
func (c *ClientLN) GetPrefix() string  { return scl.ParseAttribute(c.node, "prefix") }
func (c *ClientLN) GetLnClass() string { return scl.ParseAttribute(c.node, "lnClass") }
func (c *ClientLN) GetLnInst() string  { return scl.ParseAttribute(c.node, "lnInst") }
func (c *ClientLN) GetDesc() string    { return scl.ParseAttribute(c.node, "desc") }
