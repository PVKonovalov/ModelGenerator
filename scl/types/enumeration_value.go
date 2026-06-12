/*
 *  enumeration_value.go
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

package types

import (
	"ModelGenerator/scl"
	"fmt"
	"strconv"
	"strings"
)

// EnumerationValue is a single ordinal-name pair in an enumeration type.
type EnumerationValue struct {
	SymbolicName string
	Ord          int
}

// NewEnumerationValue parses an EnumVal node.
func NewEnumerationValue(node *scl.XMLNode) (*EnumerationValue, error) {
	ordStr := scl.ParseAttribute(node, "ord")
	if ordStr == "" {
		return nil, fmt.Errorf("ord attribute missing in EnumVal")
	}
	ord, err := strconv.Atoi(ordStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ord value %q: %w", ordStr, err)
	}
	symbolicName := strings.TrimSpace(node.Text)
	return &EnumerationValue{SymbolicName: symbolicName, Ord: ord}, nil
}
