/*
 *  setting_control.go
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
	"strconv"
)

// SettingControl holds setting group control block configuration.
type SettingControl struct {
	Desc     string
	NumOfSGs int
	ActSG    int
}

// NewSettingControl parses a SettingControl XML node.
func NewSettingControl(node *scl.XMLNode) (*SettingControl, error) {
	s := &SettingControl{NumOfSGs: 1, ActSG: 1}
	s.Desc = scl.ParseAttribute(node, "desc")

	if numStr := scl.ParseAttribute(node, "numOfSGs"); numStr != "" {
		if n, err := strconv.Atoi(numStr); err == nil {
			s.NumOfSGs = n
		}
	}
	if actStr := scl.ParseAttribute(node, "actSG"); actStr != "" {
		if n, err := strconv.Atoi(actStr); err == nil {
			s.ActSG = n
		}
	}

	return s, nil
}

func (s *SettingControl) GetDesc() string  { return s.Desc }
func (s *SettingControl) GetNumOfSGs() int { return s.NumOfSGs }
func (s *SettingControl) GetActSG() int    { return s.ActSG }
