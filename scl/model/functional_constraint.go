/*
 *  functional_constraint.go
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

import "fmt"

// FunctionalConstraint represents an IEC 61850 functional constraint.
type FunctionalConstraint int

const (
	FC_ST   FunctionalConstraint = 0
	FC_MX   FunctionalConstraint = 1
	FC_SP   FunctionalConstraint = 2
	FC_SV   FunctionalConstraint = 3
	FC_CF   FunctionalConstraint = 4
	FC_DC   FunctionalConstraint = 5
	FC_SG   FunctionalConstraint = 6
	FC_SE   FunctionalConstraint = 7
	FC_SR   FunctionalConstraint = 8
	FC_OR   FunctionalConstraint = 9
	FC_BL   FunctionalConstraint = 10
	FC_EX   FunctionalConstraint = 11
	FC_CO   FunctionalConstraint = 12
	FC_ALL  FunctionalConstraint = 99
	FC_NONE FunctionalConstraint = -1
)

var fcNames = map[FunctionalConstraint]string{
	FC_ST:   "ST",
	FC_MX:   "MX",
	FC_SP:   "SP",
	FC_SV:   "SV",
	FC_CF:   "CF",
	FC_DC:   "DC",
	FC_SG:   "SG",
	FC_SE:   "SE",
	FC_SR:   "SR",
	FC_OR:   "OR",
	FC_BL:   "BL",
	FC_EX:   "EX",
	FC_CO:   "CO",
	FC_ALL:  "ALL",
	FC_NONE: "NONE",
}

func (fc FunctionalConstraint) String() string {
	if s, ok := fcNames[fc]; ok {
		return s
	}
	return fmt.Sprintf("FC(%d)", int(fc))
}

func (fc FunctionalConstraint) GetIntValue() int {
	return int(fc)
}

// FCFromString parses a functional constraint string.
func FCFromString(s string) (FunctionalConstraint, error) {
	for fc, name := range fcNames {
		if name == s {
			return fc, nil
		}
	}
	return FC_NONE, fmt.Errorf("unknown functional constraint %q", s)
}
