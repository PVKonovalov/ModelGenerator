/*
 *  data_set.go
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
)

// DataSet holds the definition of a data set.
type DataSet struct {
	Name string
	Desc string
	FCDA []*FunctionalConstraintData
}

// NewDataSet parses a DataSet XML node.
func NewDataSet(node *scl.XMLNode) (*DataSet, error) {
	ds := &DataSet{}
	ds.Name = scl.ParseAttribute(node, "name")
	if ds.Name == "" {
		return nil, fmt.Errorf("Dataset misses required attribute \"name\"")
	}
	ds.Desc = scl.ParseAttribute(node, "desc")

	for _, child := range scl.GetChildNodesWithTag(node, "FCDA") {
		fcda, err := NewFunctionalConstraintData(child)
		if err != nil {
			return nil, err
		}
		ds.FCDA = append(ds.FCDA, fcda)
	}

	return ds, nil
}

func (ds *DataSet) GetName() string                      { return ds.Name }
func (ds *DataSet) GetDesc() string                      { return ds.Desc }
func (ds *DataSet) GetFCDA() []*FunctionalConstraintData { return ds.FCDA }
