/*
 *  log_control.go
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

// LogControl holds log control block configuration.
type LogControl struct {
	Name           string
	Desc           string
	DataSet        string
	LdInst         string
	Prefix         string
	LnClass        string
	LnInst         string
	LogName        string
	LogEna         bool
	ReasonCode     bool
	IntgPd         int
	TriggerOptions *TriggerOptions
}

// NewLogControl parses a LogControl XML node.
func NewLogControl(node *scl.XMLNode) (*LogControl, error) {
	lc := &LogControl{
		LnClass:    "LLN0",
		LogEna:     true,
		ReasonCode: true,
	}

	lc.Name = scl.ParseAttribute(node, "name")
	if lc.Name == "" {
		return nil, fmt.Errorf("LogControl is missing required attribute \"name\"")
	}

	lc.Desc = scl.ParseAttribute(node, "desc")
	lc.DataSet = scl.ParseAttribute(node, "datSet")

	lc.LdInst = scl.ParseAttribute(node, "ldInst")
	lc.Prefix = scl.ParseAttribute(node, "prefix")

	if lnClassStr := scl.ParseAttribute(node, "lnClass"); lnClassStr != "" {
		lc.LnClass = lnClassStr
	}
	if lnInstStr := scl.ParseAttribute(node, "lnInst"); lnInstStr != "" {
		lc.LnInst = lnInstStr
	}

	logNamePtr := scl.ParseAttributePtr(node, "logName")
	if logNamePtr == nil {
		return nil, fmt.Errorf("LogControl is missing required attribute \"logName\"")
	}
	lc.LogName = *logNamePtr // may be "" if attribute present but empty — treated as no log name

	if intgPdStr := scl.ParseAttribute(node, "intgPd"); intgPdStr != "" {
		n, err := strconv.Atoi(intgPdStr)
		if err == nil {
			lc.IntgPd = n
		}
	}

	if b, err := scl.ParseBooleanAttribute(node, "logEna"); err != nil {
		return nil, err
	} else if b != nil {
		lc.LogEna = *b
	}

	if b, err := scl.ParseBooleanAttribute(node, "reasonCode"); err != nil {
		return nil, err
	} else if b != nil {
		lc.ReasonCode = *b
	}

	if trgOpsNode := scl.GetChildNodeWithTag(node, "TrgOps"); trgOpsNode != nil {
		var err error
		lc.TriggerOptions, err = NewTriggerOptionsFromNode(trgOpsNode)
		if err != nil {
			return nil, err
		}
	} else {
		lc.TriggerOptions = NewTriggerOptionsDefault()
	}

	return lc, nil
}

func (lc *LogControl) GetName() string                    { return lc.Name }
func (lc *LogControl) GetDesc() string                    { return lc.Desc }
func (lc *LogControl) GetDataSet() string                 { return lc.DataSet }
func (lc *LogControl) GetLdInst() string                  { return lc.LdInst }
func (lc *LogControl) GetPrefix() string                  { return lc.Prefix }
func (lc *LogControl) GetLnClass() string                 { return lc.LnClass }
func (lc *LogControl) GetLogName() string                 { return lc.LogName }
func (lc *LogControl) IsLogEna() bool                     { return lc.LogEna }
func (lc *LogControl) IsReasonCode() bool                 { return lc.ReasonCode }
func (lc *LogControl) GetIntgPd() int                     { return lc.IntgPd }
func (lc *LogControl) GetTriggerOptions() *TriggerOptions { return lc.TriggerOptions }
