/*
 *  option_fields.go
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
)

// OptionFields holds report option flags.
type OptionFields struct {
	SeqNum     bool
	TimeStamp  bool
	DataSet    bool
	ReasonCode bool
	DataRef    bool
	EntryID    bool
	ConfigRef  bool
	BufOvfl    bool
}

// NewOptionFields parses an OptFields XML node.
func NewOptionFields(node *scl.XMLNode) (*OptionFields, error) {
	of := &OptionFields{BufOvfl: true}
	if node == nil {
		return of, nil
	}

	if b, err := scl.ParseBooleanAttribute(node, "seqNum"); err != nil {
		return nil, err
	} else if b != nil {
		of.SeqNum = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "timeStamp"); err != nil {
		return nil, err
	} else if b != nil {
		of.TimeStamp = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "dataSet"); err != nil {
		return nil, err
	} else if b != nil {
		of.DataSet = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "reasonCode"); err != nil {
		return nil, err
	} else if b != nil {
		of.ReasonCode = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "dataRef"); err != nil {
		return nil, err
	} else if b != nil {
		of.DataRef = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "entryID"); err != nil {
		return nil, err
	} else if b != nil {
		of.EntryID = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "configRef"); err != nil {
		return nil, err
	} else if b != nil {
		of.ConfigRef = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "bufOvfl"); err != nil {
		return nil, err
	} else if b != nil {
		of.BufOvfl = *b
	}

	return of, nil
}

// GetIntValue returns the option bitmask.
func (of *OptionFields) GetIntValue() int {
	v := 0
	if of.SeqNum {
		v += 1
	}
	if of.TimeStamp {
		v += 2
	}
	if of.ReasonCode {
		v += 4
	}
	if of.DataSet {
		v += 8
	}
	if of.DataRef {
		v += 16
	}
	if of.BufOvfl {
		v += 32
	}
	if of.EntryID {
		v += 64
	}
	if of.ConfigRef {
		v += 128
	}
	return v
}

func (of *OptionFields) IsSeqNum() bool     { return of.SeqNum }
func (of *OptionFields) IsTimeStamp() bool  { return of.TimeStamp }
func (of *OptionFields) IsDataSet() bool    { return of.DataSet }
func (of *OptionFields) IsReasonCode() bool { return of.ReasonCode }
func (of *OptionFields) IsDataRef() bool    { return of.DataRef }
func (of *OptionFields) IsEntryID() bool    { return of.EntryID }
func (of *OptionFields) IsConfigRef() bool  { return of.ConfigRef }
func (of *OptionFields) IsBufOvfl() bool    { return of.BufOvfl }
