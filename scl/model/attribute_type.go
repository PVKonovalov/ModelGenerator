/*
 *  attribute_type.go
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

// AttributeType represents the basic data type of a data attribute.
type AttributeType int

const (
	AT_BOOLEAN            AttributeType = 0
	AT_INT8               AttributeType = 1
	AT_INT16              AttributeType = 2
	AT_INT32              AttributeType = 3
	AT_INT64              AttributeType = 4
	AT_INT128             AttributeType = 5
	AT_INT8U              AttributeType = 6
	AT_INT16U             AttributeType = 7
	AT_INT24U             AttributeType = 8
	AT_INT32U             AttributeType = 9
	AT_FLOAT32            AttributeType = 10
	AT_FLOAT64            AttributeType = 11
	AT_ENUMERATED         AttributeType = 12
	AT_OCTET_STRING_64    AttributeType = 13
	AT_OCTET_STRING_6     AttributeType = 14
	AT_OCTET_STRING_8     AttributeType = 15
	AT_VISIBLE_STRING_32  AttributeType = 16
	AT_VISIBLE_STRING_64  AttributeType = 17
	AT_VISIBLE_STRING_65  AttributeType = 18
	AT_VISIBLE_STRING_129 AttributeType = 19
	AT_VISIBLE_STRING_255 AttributeType = 20
	AT_UNICODE_STRING_255 AttributeType = 21
	AT_TIMESTAMP          AttributeType = 22
	AT_QUALITY            AttributeType = 23
	AT_CHECK              AttributeType = 24
	AT_CODEDENUM          AttributeType = 25
	AT_GENERIC_BITSTRING  AttributeType = 26
	AT_CONSTRUCTED        AttributeType = 27
	AT_ENTRY_TIME         AttributeType = 28
	AT_PHYCOMADDR         AttributeType = 29
	AT_CURRENCY           AttributeType = 30
	AT_OPTFLDS            AttributeType = 31
	AT_TRGOPS             AttributeType = 32
)

var atNames = map[AttributeType]string{
	AT_BOOLEAN:            "BOOLEAN",
	AT_INT8:               "INT8",
	AT_INT16:              "INT16",
	AT_INT32:              "INT32",
	AT_INT64:              "INT64",
	AT_INT128:             "INT128",
	AT_INT8U:              "INT8U",
	AT_INT16U:             "INT16U",
	AT_INT24U:             "INT24U",
	AT_INT32U:             "INT32U",
	AT_FLOAT32:            "FLOAT32",
	AT_FLOAT64:            "FLOAT64",
	AT_ENUMERATED:         "ENUMERATED",
	AT_OCTET_STRING_64:    "OCTET_STRING_64",
	AT_OCTET_STRING_6:     "OCTET_STRING_6",
	AT_OCTET_STRING_8:     "OCTET_STRING_8",
	AT_VISIBLE_STRING_32:  "VISIBLE_STRING_32",
	AT_VISIBLE_STRING_64:  "VISIBLE_STRING_64",
	AT_VISIBLE_STRING_65:  "VISIBLE_STRING_65",
	AT_VISIBLE_STRING_129: "VISIBLE_STRING_129",
	AT_VISIBLE_STRING_255: "VISIBLE_STRING_255",
	AT_UNICODE_STRING_255: "UNICODE_STRING_255",
	AT_TIMESTAMP:          "TIMESTAMP",
	AT_QUALITY:            "QUALITY",
	AT_CHECK:              "CHECK",
	AT_CODEDENUM:          "CODEDENUM",
	AT_GENERIC_BITSTRING:  "GENERIC_BITSTRING",
	AT_CONSTRUCTED:        "CONSTRUCTED",
	AT_ENTRY_TIME:         "ENTRY_TIME",
	AT_PHYCOMADDR:         "PHYCOMADDR",
	AT_CURRENCY:           "CURRENCY",
	AT_OPTFLDS:            "OPTFLDS",
	AT_TRGOPS:             "TRGOPS",
}

func (at AttributeType) String() string {
	if s, ok := atNames[at]; ok {
		return s
	}
	return fmt.Sprintf("AT(%d)", int(at))
}

func (at AttributeType) GetIntValue() int { return int(at) }

// ATFromSclString parses a bType string (e.g. "BOOLEAN", "VisString32") to an AttributeType.
func ATFromSclString(s string) (AttributeType, error) {
	switch s {
	case "BOOLEAN":
		return AT_BOOLEAN, nil
	case "INT8":
		return AT_INT8, nil
	case "INT16":
		return AT_INT16, nil
	case "INT32":
		return AT_INT32, nil
	case "INT64":
		return AT_INT64, nil
	case "INT128":
		return AT_INT128, nil
	case "INT8U":
		return AT_INT8U, nil
	case "INT16U":
		return AT_INT16U, nil
	case "INT24U":
		return AT_INT24U, nil
	case "INT32U":
		return AT_INT32U, nil
	case "FLOAT32":
		return AT_FLOAT32, nil
	case "FLOAT64":
		return AT_FLOAT64, nil
	case "Enum":
		return AT_ENUMERATED, nil
	case "Dbpos", "Tcmd":
		return AT_CODEDENUM, nil
	case "Check":
		return AT_CHECK, nil
	case "Octet64":
		return AT_OCTET_STRING_64, nil
	case "Quality":
		return AT_QUALITY, nil
	case "Timestamp":
		return AT_TIMESTAMP, nil
	case "Currency":
		return AT_CURRENCY, nil
	case "VisString32":
		return AT_VISIBLE_STRING_32, nil
	case "VisString64":
		return AT_VISIBLE_STRING_64, nil
	case "VisString65":
		return AT_VISIBLE_STRING_65, nil
	case "VisString129", "ObjRef":
		return AT_VISIBLE_STRING_129, nil
	case "VisString255":
		return AT_VISIBLE_STRING_255, nil
	case "Unicode255":
		return AT_UNICODE_STRING_255, nil
	case "OptFlds":
		return AT_OPTFLDS, nil
	case "TrgOps":
		return AT_TRGOPS, nil
	case "EntryID":
		return AT_OCTET_STRING_8, nil
	case "EntryTime":
		return AT_ENTRY_TIME, nil
	case "PhyComAddr":
		return AT_PHYCOMADDR, nil
	case "Struct":
		return AT_CONSTRUCTED, nil
	default:
		return AT_BOOLEAN, fmt.Errorf("unsupported attribute type %q", s)
	}
}
