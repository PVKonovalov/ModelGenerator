/*
 *  functional_constraint_data.go
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

// FunctionalConstraintData is an FCDA entry in a data set.
type FunctionalConstraintData struct {
	LdInstance string
	Prefix     string
	LnClass    string
	LnInstance string
	DoName     string
	DaName     string
	FC         FunctionalConstraint
	Ix         *int
}

// NewFunctionalConstraintData parses an FCDA XML node.
func NewFunctionalConstraintData(node *scl.XMLNode) (*FunctionalConstraintData, error) {
	f := &FunctionalConstraintData{}
	f.LdInstance = scl.ParseAttribute(node, "ldInst")
	f.Prefix = scl.ParseAttribute(node, "prefix")
	f.LnClass = scl.ParseAttribute(node, "lnClass")
	f.LnInstance = scl.ParseAttribute(node, "lnInst")
	f.DoName = scl.ParseAttribute(node, "doName")
	f.DaName = scl.ParseAttribute(node, "daName")

	fcStr := scl.ParseAttribute(node, "fc")
	if fcStr != "" {
		fc, err := FCFromString(fcStr)
		if err != nil {
			return nil, fmt.Errorf("invalid FC %q: %w", fcStr, err)
		}
		f.FC = fc
	}

	if indexStr := scl.ParseAttribute(node, "ix"); indexStr != "" {
		n, err := strconv.Atoi(indexStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ix %q: %w", indexStr, err)
		}
		f.Ix = &n
	}

	return f, nil
}

func (f *FunctionalConstraintData) GetLdInstance() string       { return f.LdInstance }
func (f *FunctionalConstraintData) GetLnClass() string          { return f.LnClass }
func (f *FunctionalConstraintData) GetLnInstance() string       { return f.LnInstance }
func (f *FunctionalConstraintData) GetDoName() string           { return f.DoName }
func (f *FunctionalConstraintData) GetDaName() string           { return f.DaName }
func (f *FunctionalConstraintData) GetFC() FunctionalConstraint { return f.FC }
func (f *FunctionalConstraintData) GetIx() *int                 { return f.Ix }
func (f *FunctionalConstraintData) GetPrefix() string           { return f.Prefix }

func (f *FunctionalConstraintData) String() string {
	s := ""
	if f.LdInstance != "" {
		s = f.LdInstance + "/"
	}
	if f.LnClass != "" {
		if f.Prefix != "" {
			s += f.Prefix
		}
		s += f.LnClass
		if f.LnInstance == "" {
			s += "."
		}
	}
	if f.LnInstance != "" {
		s += f.LnInstance + "."
	}
	if f.DoName != "" {
		s += f.DoName
	}
	if f.DaName != "" {
		s += "." + f.DaName
	}
	return s
}
