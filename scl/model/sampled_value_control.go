/*
 *  sampled_value_control.go
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

// SampledValueControl holds SMV control block configuration.
type SampledValueControl struct {
	Name      string
	Desc      string
	DatSet    string
	SmvID     string
	SmvOpts   *SmvOpts
	ConfRev   int
	SmpRate   int
	NofASDU   int
	Multicast bool
	SmpMod    SmpMod
}

// NewSampledValueControl parses a SampledValueControl XML node.
func NewSampledValueControl(node *scl.XMLNode) (*SampledValueControl, error) {
	s := &SampledValueControl{
		ConfRev: 1,
		SmpMod:  SmpMod_SMP_PER_PERIOD,
	}
	s.Name = scl.ParseAttribute(node, "name")
	s.Desc = scl.ParseAttribute(node, "desc")
	s.DatSet = scl.ParseAttribute(node, "datSet")
	s.SmvID = scl.ParseAttribute(node, "smvID")

	if confRevStr := scl.ParseAttribute(node, "confRev"); confRevStr != "" {
		if n, err := strconv.Atoi(confRevStr); err == nil {
			s.ConfRev = n
		}
	}
	if smpRateStr := scl.ParseAttribute(node, "smpRate"); smpRateStr != "" {
		if n, err := strconv.Atoi(smpRateStr); err == nil {
			s.SmpRate = n
		}
	}
	if nofASDUStr := scl.ParseAttribute(node, "nofASDU"); nofASDUStr != "" {
		if n, err := strconv.Atoi(nofASDUStr); err == nil {
			s.NofASDU = n
		}
	}

	if b, err := scl.ParseBooleanAttribute(node, "multicast"); err != nil {
		return nil, err
	} else if b != nil {
		s.Multicast = *b
	}

	smvOptsNode := scl.GetChildNodeWithTag(node, "SmvOpts")
	var err error
	s.SmvOpts, err = NewSmvOpts(smvOptsNode)
	if err != nil {
		return nil, err
	}

	if smpModStr := scl.ParseAttribute(node, "smpMod"); smpModStr != "" {
		switch smpModStr {
		case "SmpPerPeriod":
			s.SmpMod = SmpMod_SMP_PER_PERIOD
		case "SmpPerSec":
			s.SmpMod = SmpMod_SMP_PER_SECOND
		case "SecPerSmp":
			s.SmpMod = SmpMod_SEC_PER_SMP
		default:
			return nil, fmt.Errorf("invalid smpMod value %q", smpModStr)
		}
	}

	return s, nil
}

func (s *SampledValueControl) GetName() string      { return s.Name }
func (s *SampledValueControl) GetDesc() string      { return s.Desc }
func (s *SampledValueControl) GetDatSet() string    { return s.DatSet }
func (s *SampledValueControl) GetSmvID() string     { return s.SmvID }
func (s *SampledValueControl) GetConfRev() int      { return s.ConfRev }
func (s *SampledValueControl) GetSmpRate() int      { return s.SmpRate }
func (s *SampledValueControl) GetNofASDI() int      { return s.NofASDU }
func (s *SampledValueControl) IsMulticast() bool    { return s.Multicast }
func (s *SampledValueControl) GetSmvOpts() *SmvOpts { return s.SmvOpts }
func (s *SampledValueControl) GetSmpMod() SmpMod    { return s.SmpMod }
