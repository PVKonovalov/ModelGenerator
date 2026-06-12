/*
 *  trigger_options.go
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

// TriggerOptions holds the trigger flags for a report or log control block.
type TriggerOptions struct {
	Dchg   bool // 1
	Qchg   bool // 2
	Dupd   bool // 4
	Period bool // 8
	GI     bool // 16 (general interrogation)
}

// NewTriggerOptionsFromNode parses a TrgOps XML node.
func NewTriggerOptionsFromNode(node *scl.XMLNode) (*TriggerOptions, error) {
	t := &TriggerOptions{GI: true}

	if b, err := scl.ParseBooleanAttribute(node, "dchg"); err != nil {
		return nil, err
	} else if b != nil {
		t.Dchg = *b
	}

	if b, err := scl.ParseBooleanAttribute(node, "qchg"); err != nil {
		return nil, err
	} else if b != nil {
		t.Qchg = *b
	}

	if b, err := scl.ParseBooleanAttribute(node, "dupd"); err != nil {
		return nil, err
	} else if b != nil {
		t.Dupd = *b
	}

	if b, err := scl.ParseBooleanAttribute(node, "period"); err != nil {
		return nil, err
	} else if b != nil {
		t.Period = *b
	}

	if b, err := scl.ParseBooleanAttribute(node, "gi"); err != nil {
		return nil, err
	} else if b != nil {
		t.GI = *b
	}

	return t, nil
}

// NewTriggerOptionsDefault creates TriggerOptions with default values (gi=true).
func NewTriggerOptionsDefault() *TriggerOptions {
	return &TriggerOptions{GI: true}
}

// NewTriggerOptions creates TriggerOptions with explicit values.
func NewTriggerOptions(dchg, qchg, dupd, period, gi bool) *TriggerOptions {
	return &TriggerOptions{Dchg: dchg, Qchg: qchg, Dupd: dupd, Period: period, GI: gi}
}

// GetIntValue returns the integer bitmask.
func (t *TriggerOptions) GetIntValue() int {
	v := 0
	if t.Dchg {
		v += 1
	}
	if t.Qchg {
		v += 2
	}
	if t.Dupd {
		v += 4
	}
	if t.Period {
		v += 8
	}
	if t.GI {
		v += 16
	}
	return v
}

func (t *TriggerOptions) IsDchg() bool   { return t.Dchg }
func (t *TriggerOptions) IsQchg() bool   { return t.Qchg }
func (t *TriggerOptions) IsDupd() bool   { return t.Dupd }
func (t *TriggerOptions) IsPeriod() bool { return t.Period }
func (t *TriggerOptions) IsGI() bool     { return t.GI }
