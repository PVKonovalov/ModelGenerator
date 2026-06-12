/*
 *  gse_control.go
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
	"fmt"
	"strconv"
)

// GSEControl holds GSEControl (GOOSE) configuration.
type GSEControl struct {
	Name      string
	Desc      string
	DataSet   string
	AppID     string
	ConfRev   int
	FixedOffs bool
}

// NewGSEControl parses a GSEControl XML node.
func NewGSEControl(node *scl.XMLNode) (*GSEControl, error) {
	g := &GSEControl{ConfRev: 1}
	g.Name = scl.ParseAttribute(node, "name")
	g.Desc = scl.ParseAttribute(node, "desc")
	g.DataSet = scl.ParseAttribute(node, "datSet")
	g.AppID = scl.ParseAttribute(node, "appID")

	if confRevStr := scl.ParseAttribute(node, "confRev"); confRevStr != "" {
		n, err := strconv.Atoi(confRevStr)
		if err == nil {
			g.ConfRev = n
		}
	}

	if b, err := scl.ParseBooleanAttribute(node, "fixedOffs"); err != nil {
		return nil, err
	} else if b != nil {
		g.FixedOffs = *b
	}

	if typeStr := scl.ParseAttribute(node, "type"); typeStr != "" {
		if typeStr != "GOOSE" {
			return nil, fmt.Errorf("GSEControl of type %q not supported", typeStr)
		}
	}

	return g, nil
}

func (g *GSEControl) GetName() string    { return g.Name }
func (g *GSEControl) GetDesc() string    { return g.Desc }
func (g *GSEControl) GetDataSet() string { return g.DataSet }
func (g *GSEControl) GetAppID() string   { return g.AppID }
func (g *GSEControl) GetConfRev() int    { return g.ConfRev }
func (g *GSEControl) IsFixedOffs() bool  { return g.FixedOffs }
