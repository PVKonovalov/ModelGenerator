/*
 *  data_model_value.go
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
	"ModelGenerator/scl/types"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DataModelValue holds the default value of a data attribute.
type DataModelValue struct {
	value            interface{}
	unknownEnumValue string
	enumType         string
}

// NewDataModelValueFromEnumType creates a value for an ENUMERATED attribute whose type is not yet resolved.
func NewDataModelValueFromEnumType(enumType, valueText string) *DataModelValue {
	return &DataModelValue{unknownEnumValue: valueText, enumType: enumType}
}

// NewDataModelValue parses a value string for the given attribute type.
func NewDataModelValue(atType AttributeType, sclType types.SclTypeIface, valueText string) (*DataModelValue, error) {
	dmv := &DataModelValue{}

	switch atType {
	case AT_ENUMERATED:
		et, ok := sclType.(*types.EnumerationType)
		if !ok {
			return nil, fmt.Errorf("wrong type definition for enumerated attribute")
		}
		ord, err := et.GetOrdByEnumString(valueText)
		if err != nil {
			// try parsing as integer
			n, err2 := strconv.Atoi(valueText)
			if err2 != nil {
				return nil, fmt.Errorf("%s is not a valid value: %w", valueText, err)
			}
			if !et.IsValidOrdValue(n) {
				return nil, fmt.Errorf("%s is not a valid ordinal value", valueText)
			}
			fmt.Println("WARNING: Initialization of ENUM with ordinal value!")
			dmv.value = n
		} else {
			dmv.value = ord
		}

	case AT_INT8, AT_INT16, AT_INT32, AT_INT8U, AT_INT16U, AT_INT32U, AT_INT24U, AT_INT64:
		trimmed := strings.TrimSpace(valueText)
		if trimmed != valueText {
			fmt.Println("WARNING: value initializer contains leading or trailing whitespace")
		}
		if trimmed == "" {
			dmv.value = int64(0)
		} else {
			n, err := strconv.ParseInt(trimmed, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer value %q: %w", trimmed, err)
			}
			dmv.value = n
		}

	case AT_BOOLEAN:
		trimmed := strings.TrimSpace(valueText)
		if strings.EqualFold(trimmed, "true") {
			dmv.value = true
		} else {
			dmv.value = false
		}

	case AT_FLOAT32:
		trimmed := strings.TrimSpace(strings.ReplaceAll(valueText, ",", "."))
		if trimmed == "" {
			dmv.value = float32(0)
		} else {
			f, err := strconv.ParseFloat(trimmed, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid float32 value %q: %w", trimmed, err)
			}
			dmv.value = float32(f)
		}

	case AT_FLOAT64:
		trimmed := strings.TrimSpace(strings.ReplaceAll(valueText, ",", "."))
		if trimmed == "" {
			dmv.value = float64(0)
		} else {
			f, err := strconv.ParseFloat(trimmed, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid float64 value %q: %w", trimmed, err)
			}
			dmv.value = f
		}

	case AT_UNICODE_STRING_255:
		dmv.value = valueText

	case AT_OCTET_STRING_64:
		b, err := base64.StdEncoding.DecodeString(valueText)
		if err != nil {
			return nil, fmt.Errorf("Val element for Octet64 type does not contain a valid base64 encoded string")
		}
		dmv.value = b

	case AT_VISIBLE_STRING_32, AT_VISIBLE_STRING_64, AT_VISIBLE_STRING_65,
		AT_VISIBLE_STRING_129, AT_VISIBLE_STRING_255, AT_CURRENCY:
		dmv.value = valueText

	case AT_CHECK:
		fmt.Println("Warning: Initialization of CHECK is unsupported!")
		dmv.value = nil

	case AT_CODEDENUM:
		dmv.value = nil
		switch valueText {
		case "intermediate-state":
			dmv.value = 0
		case "off":
			dmv.value = 1
		case "on":
			dmv.value = 2
		case "bad-state":
			dmv.value = 4
		case "stop":
			dmv.value = 0
		case "lower":
			dmv.value = 1
		case "higher":
			dmv.value = 2
		case "reserved":
			dmv.value = 4
		default:
			fmt.Printf("Warning: CODEDENUM is initialized with unsupported value %q\n", valueText)
		}

	case AT_QUALITY:
		dmv.value = nil
		fmt.Println("Warning: Initialization of QUALITY is unsupported!")

	case AT_TIMESTAMP, AT_ENTRY_TIME:
		modValue := strings.ReplaceAll(valueText, ",", ".")
		layouts := []string{
			"2006-01-2T15:04:05.000",
			"2006-01-2T15:04:05",
		}
		var parsed bool
		for _, layout := range layouts {
			if t, err := time.Parse(layout, modValue); err == nil {
				dmv.value = t.UnixMilli()
				parsed = true
				break
			}
		}
		if !parsed {
			dmv.value = nil
			fmt.Printf("Warning: Val element does not contain a valid time stamp: %q\n", valueText)
		}

	default:
		return nil, fmt.Errorf("unsupported type %v for value %q", atType, valueText)
	}

	return dmv, nil
}

// GetValue returns the underlying value as interface{}.
func (d *DataModelValue) GetValue() interface{} { return d.value }

// GetUnknownEnumValue returns the unresolved enum string.
func (d *DataModelValue) GetUnknownEnumValue() string { return d.unknownEnumValue }

// UpdateEnumOrdValue resolves the enum string to an ordinal using the type registry.
func (d *DataModelValue) UpdateEnumOrdValue(typeDecls *types.TypeDeclarations) {
	if d.enumType == "" {
		return
	}
	et := typeDecls.LookupEnumerationType(d.enumType)
	if et == nil {
		fmt.Printf("  failed to find enum type %q!\n", d.enumType)
		return
	}
	ord, err := et.GetOrdByEnumString(d.unknownEnumValue)
	if err != nil {
		fmt.Printf("  failed to resolve enum value %q: %v\n", d.unknownEnumValue, err)
		return
	}
	d.value = ord
}

// GetIntValue returns the value as int.
func (d *DataModelValue) GetIntValue() int {
	switch v := d.value.(type) {
	case int64:
		return int(v)
	case int:
		return v
	case bool:
		if v {
			return 1
		}
		return 0
	}
	return 0
}

// GetLongValue returns the value as int64.
func (d *DataModelValue) GetLongValue() int64 {
	switch v := d.value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	}
	return 0
}
