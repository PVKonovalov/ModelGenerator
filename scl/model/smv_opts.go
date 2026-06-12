/*
 *  smv_opts.go
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

// SmvOpts holds SMV option flags.
type SmvOpts struct {
	RefreshTime        bool // 1
	SampleSynchronized bool // 2
	SampleRate         bool // 4
	DataSet            bool // 8
	Security           bool // 16
}

// NewSmvOpts parses a SmvOpts XML node.
func NewSmvOpts(node *scl.XMLNode) (*SmvOpts, error) {
	s := &SmvOpts{}
	if node == nil {
		return s, nil
	}

	if b, err := scl.ParseBooleanAttribute(node, "refreshTime"); err != nil {
		return nil, err
	} else if b != nil {
		s.RefreshTime = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "sampleRate"); err != nil {
		return nil, err
	} else if b != nil {
		s.SampleRate = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "dataSet"); err != nil {
		return nil, err
	} else if b != nil {
		s.DataSet = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "security"); err != nil {
		return nil, err
	} else if b != nil {
		s.Security = *b
	}
	if b, err := scl.ParseBooleanAttribute(node, "sampleSynchronized"); err != nil {
		return nil, err
	} else if b != nil {
		s.SampleSynchronized = *b
	}

	return s, nil
}

// GetIntValue returns the option bitmask.
func (s *SmvOpts) GetIntValue() int {
	v := 0
	if s.RefreshTime {
		v += 1
	}
	if s.SampleSynchronized {
		v += 2
	}
	if s.SampleRate {
		v += 4
	}
	if s.DataSet {
		v += 8
	}
	if s.Security {
		v += 16
	}
	return v
}
