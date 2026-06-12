/*
 *  model_viewer.go
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

package tools

import (
	"ModelGenerator/parser"
	"ModelGenerator/scl/model"
	"ModelGenerator/scl/types"
	"fmt"
	"io"
)

// RunModelViewer prints model information based on the selected mode flags.
func RunModelViewer(r io.Reader, output io.Writer, iedName, accessPointName string,
	printTypeList, printUnusedTypes, printModelStructure, printDataAttributes bool) error {

	withOutput := false
	sclParser, err := parser.NewSclParser(r, withOutput)
	if err != nil {
		return err
	}

	if printTypeList {
		showTypes(sclParser.GetTypeDeclarations(), output, false)
	}

	if printUnusedTypes {
		showTypes(sclParser.GetTypeDeclarations(), output, true)
	}

	if printModelStructure {
		ied := resolveIED(sclParser, iedName)
		if ied == nil {
			return fmt.Errorf("IED not found")
		}
		printStructure(ied, output)
	}

	if printDataAttributes {
		ied := resolveMainIED(sclParser, iedName)
		if ied == nil {
			return fmt.Errorf("IED not found")
		}
		printAttributeList(ied, output)
	}

	return nil
}

func resolveIED(p *parser.SclParser, iedName string) *model.IED {
	if iedName == "" {
		return p.GetFirstIed()
	}
	return p.GetIedByName(iedName)
}

func resolveMainIED(p *parser.SclParser, iedName string) *model.IED {
	if iedName == "" {
		return p.GetMainIed()
	}
	return p.GetIedByName(iedName)
}

func showTypes(typeDecl *types.TypeDeclarations, output io.Writer, unusedOnly bool) {
	for _, t := range typeDecl.GetTypeDeclarations() {
		if unusedOnly && t.GetIsUsed() {
			continue
		}

		fmt.Fprint(output, t.GetID())

		switch t.(type) {
		case *types.LogicalNodeType:
			fmt.Fprint(output, " : LogicalNode")
		case *types.DataObjectType:
			fmt.Fprint(output, " : DataObject")
		case *types.DataAttributeType:
			fmt.Fprint(output, " : DataAttribute")
		case *types.EnumerationType:
			fmt.Fprint(output, " : Enumeration")
		}

		if !unusedOnly && !t.GetIsUsed() {
			fmt.Fprint(output, "  UNUSED TYPE!")
		}

		fmt.Fprintln(output)
	}
}

func printStructure(ied *model.IED, output io.Writer) {
	ap := ied.GetFirstAccessPoint()
	if ap == nil || ap.GetServer() == nil {
		return
	}
	for _, ld := range ap.GetServer().GetLogicalDevices() {
		fmt.Fprintln(output, ld.GetInst())
		for _, ln := range ld.GetLogicalNodes() {
			fmt.Fprintln(output, "  "+ln.GetName())
			printDataObjects(ln.GetDataObjects(), output, "    ", "  ")
		}
	}
}

func printDataObjects(dos []*model.DataObject, output io.Writer, indent, add string) {
	for _, do := range dos {
		fmt.Fprintln(output, indent+do.GetName())
		printDataObjects(do.GetSubDataObjects(), output, indent+add, add)
		for _, da := range do.GetDataAttributes() {
			fmt.Fprintln(output, indent+add+da.GetName()+" ["+da.GetFC().String()+"]")
			printSubAttributes(da, output, indent+add, add)
		}
	}
}

func printSubAttributes(da *model.DataAttribute, output io.Writer, indent, add string) {
	for _, sub := range da.GetSubDataAttributes() {
		fmt.Fprintln(output, indent+add+sub.GetName())
		printSubAttributes(sub, output, indent+add, add)
	}
}

func printAttributeList(ied *model.IED, output io.Writer) {
	ap := ied.GetFirstAccessPoint()
	if ap == nil || ap.GetServer() == nil {
		return
	}
	for _, ld := range ap.GetServer().GetLogicalDevices() {
		devPrefix := ied.GetName() + ld.GetInst() + "/"
		for _, ln := range ld.GetLogicalNodes() {
			lnPrefix := devPrefix + ln.GetName()
			printObjectList(ln.GetDataObjects(), output, lnPrefix)
		}
	}
}

func printObjectList(dos []*model.DataObject, output io.Writer, prefix string) {
	for _, do := range dos {
		doPrefix := prefix + "." + do.GetName()
		printObjectList(do.GetSubDataObjects(), output, doPrefix)
		for _, da := range do.GetDataAttributes() {
			daPrefix := doPrefix + "." + da.GetName()
			fmt.Fprintln(output, daPrefix+" ["+da.GetFC().String()+"]")
			printSubAttributeList(da, output, daPrefix)
		}
	}
}

func printSubAttributeList(da *model.DataAttribute, output io.Writer, prefix string) {
	for _, sub := range da.GetSubDataAttributes() {
		nextPrefix := prefix + "." + sub.GetName()
		fmt.Fprintln(output, nextPrefix+" ["+sub.GetFC().String()+"]")
		printSubAttributeList(sub, output, nextPrefix)
	}
}
