/*
 *  report_control_block.go
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

// ReportControlBlock holds RCB configuration.
type ReportControlBlock struct {
	Name            string
	Desc            string
	DataSet         string
	IntegrityPeriod *int64
	RptID           string
	ConfRef         *int64
	Buffered        bool
	BufferTime      int64
	TriggerOptions  *TriggerOptions
	OptionFields    *OptionFields
	Indexed         bool
	RptEna          *RptEnabled
}

// NewReportControlBlock parses a ReportControl XML node.
func NewReportControlBlock(node *scl.XMLNode) (*ReportControlBlock, error) {
	r := &ReportControlBlock{Indexed: true}
	r.Name = scl.ParseAttribute(node, "name")
	r.Desc = scl.ParseAttribute(node, "desc")
	r.DataSet = scl.ParseAttribute(node, "datSet")

	if intgPdStr := scl.ParseAttribute(node, "intgPd"); intgPdStr != "" {
		n, err := strconv.ParseInt(intgPdStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid intgPd %q: %w", intgPdStr, err)
		}
		r.IntegrityPeriod = &n
	}

	r.RptID = scl.ParseAttribute(node, "rptID")
	if r.RptID == "" {
		r.RptID = ""
	}

	confRefStr := scl.ParseAttribute(node, "confRev")
	if confRefStr == "" {
		return nil, fmt.Errorf("missing required attribute \"confRef\"")
	}
	n, err := strconv.ParseInt(confRefStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid confRev %q: %w", confRefStr, err)
	}
	r.ConfRef = &n

	if b, err := scl.ParseBooleanAttribute(node, "buffered"); err != nil {
		return nil, err
	} else if b != nil {
		r.Buffered = *b
	}

	if bufTimeStr := scl.ParseAttribute(node, "bufTime"); bufTimeStr != "" {
		n, err := strconv.ParseInt(bufTimeStr, 10, 64)
		if err == nil {
			r.BufferTime = n
		}
	}

	if trgOpsNode := scl.GetChildNodeWithTag(node, "TrgOps"); trgOpsNode != nil {
		r.TriggerOptions, err = NewTriggerOptionsFromNode(trgOpsNode)
		if err != nil {
			return nil, err
		}
	} else {
		r.TriggerOptions = NewTriggerOptionsDefault()
	}

	optFieldsNode := scl.GetChildNodeWithTag(node, "OptFields")
	r.OptionFields, err = NewOptionFields(optFieldsNode)
	if err != nil {
		return nil, err
	}

	if b, err := scl.ParseBooleanAttribute(node, "indexed"); err != nil {
		return nil, err
	} else if b != nil {
		r.Indexed = *b
	}

	if rptEnaNode := scl.GetChildNodeWithTag(node, "RptEnabled"); rptEnaNode != nil {
		rptEna := NewRptEnabled(rptEnaNode)
		if !r.Indexed && rptEna.GetMaxInstances() != 1 {
			return nil, fmt.Errorf("RptEnabled.max != 1 is not allowed when indexed=\"false\"")
		}
		r.RptEna = rptEna
	}

	return r, nil
}

func (r *ReportControlBlock) GetName() string                    { return r.Name }
func (r *ReportControlBlock) GetDesc() string                    { return r.Desc }
func (r *ReportControlBlock) GetDataSet() string                 { return r.DataSet }
func (r *ReportControlBlock) GetIntegrityPeriod() *int64         { return r.IntegrityPeriod }
func (r *ReportControlBlock) GetRptID() string                   { return r.RptID }
func (r *ReportControlBlock) GetConfRef() *int64                 { return r.ConfRef }
func (r *ReportControlBlock) IsBuffered() bool                   { return r.Buffered }
func (r *ReportControlBlock) GetBufferTime() int64               { return r.BufferTime }
func (r *ReportControlBlock) GetTriggerOptions() *TriggerOptions { return r.TriggerOptions }
func (r *ReportControlBlock) GetOptionFields() *OptionFields     { return r.OptionFields }
func (r *ReportControlBlock) IsIndexed() bool                    { return r.Indexed }
func (r *ReportControlBlock) GetRptEna() *RptEnabled             { return r.RptEna }
