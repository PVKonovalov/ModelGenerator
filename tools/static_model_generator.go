/*
 *  static_model_generator.go
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
	"ModelGenerator/scl/communication"
	"ModelGenerator/scl/model"
	"ModelGenerator/scl/types"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

// StaticModelGenerator generates a static C/H model from an SCL file.
type StaticModelGenerator struct {
	cOut io.Writer
	hOut io.Writer

	initializerBuffer         strings.Builder
	reportControlBlocks       strings.Builder
	rcbVariableNames          []string
	logControlBlocks          strings.Builder
	lcbVariableNames          []string
	logs                      strings.Builder
	logVariableNames          []string
	gseControlBlocks          strings.Builder
	gseVariableNames          []string
	smvControlBlocks          strings.Builder
	smvVariableNames          []string
	settingGroupControlBlocks strings.Builder
	sgcbVariableNames         []string

	connectedAP    *communication.ConnectedAP
	outputFileName string
	modelPrefix    string
	initializeOnce bool
	sclParser      *parser.SclParser
	hDefineName    string
	hasOwner       bool

	currentRcbVariableNumber  int
	currentLcbVariableNumber  int
	currentLogVariableNumber  int
	currentGseVariableNumber  int
	currentSvCBVariableNumber int
	currentSGCBVariableNumber int

	dataSetNames  []string
	ied           *model.IED
	accessPoint   *model.AccessPoint
	variablesList []string

	typeDeclarations *types.TypeDeclarations
}

// NewStaticModelGenerator runs the static model generation.
func NewStaticModelGenerator(
	r io.Reader, icdFile string,
	cOut, hOut io.Writer,
	outputFileName, iedName, accessPointName, modelPrefix string,
	initializeOnce bool,
) error {
	g := &StaticModelGenerator{
		cOut:           cOut,
		hOut:           hOut,
		outputFileName: outputFileName,
		modelPrefix:    modelPrefix,
		initializeOnce: initializeOnce,
	}

	g.hDefineName = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(outputFileName, ".", "_"), "-", "_")) + "_H_"
	if idx := strings.LastIndex(g.hDefineName, "/"); idx >= 0 {
		g.hDefineName = g.hDefineName[idx+1:]
	}

	var err error
	g.sclParser, err = parser.NewSclParser(r, true)
	if err != nil {
		return err
	}

	g.typeDeclarations = g.sclParser.GetTypeDeclarations()

	if iedName == "" {
		g.ied = g.sclParser.GetFirstIed()
	} else {
		g.ied = g.sclParser.GetIedByName(iedName)
	}

	if g.ied == nil {
		fmt.Fprintln(cOut, "IED model not found in SCL file! Exit.")
		return fmt.Errorf("IED not found")
	}

	services := g.ied.GetServices()
	if services != nil {
		rptSettings := services.GetReportSettings()
		if rptSettings != nil {
			g.hasOwner = rptSettings.HasOwner()
		}
	}

	if accessPointName == "" {
		g.accessPoint = g.ied.GetFirstAccessPoint()
	} else {
		g.accessPoint = g.ied.GetAccessPointByName(accessPointName)
	}

	if g.accessPoint == nil {
		return fmt.Errorf("access point not found")
	}

	g.connectedAP = g.sclParser.GetConnectedAP(g.ied.GetName(), g.accessPoint.GetName())

	g.printCFileHeader(icdFile)
	g.printHeaderFileHeader(icdFile)

	server := g.accessPoint.GetServer()
	g.printForwardDeclarations(server)
	g.printDeviceModelDefinitions()
	g.printInitializerFunction()
	g.printVariablePointerDefines()
	g.printHeaderFileFooter()

	return nil
}

func fprintln(w io.Writer, s string) { fmt.Fprintln(w, s) }
func fprint(w io.Writer, s string)   { fmt.Fprint(w, s) }

func (g *StaticModelGenerator) printCFileHeader(filename string) {
	include := filepath.Base(g.outputFileName + ".h\"")
	fprintln(g.cOut, "/*")
	fprintln(g.cOut, " * "+g.outputFileName+".c")
	fprintln(g.cOut, " *")
	fprintln(g.cOut, " * automatically generated from "+filename)
	fprintln(g.cOut, " */")
	fprintln(g.cOut, "#include \""+include)
	fprintln(g.cOut, "")
}

func (g *StaticModelGenerator) printHeaderFileHeader(filename string) {
	fprintln(g.hOut, "/*")
	fprintln(g.hOut, " * "+g.outputFileName+".h")
	fprintln(g.hOut, " *")
	fprintln(g.hOut, " * automatically generated from "+filename)
	fprintln(g.hOut, " */\n")
	fprintln(g.hOut, "#ifndef "+g.hDefineName)
	fprintln(g.hOut, "#define "+g.hDefineName+"\n")
	fprintln(g.hOut, "#include <stdlib.h>")
	fprintln(g.hOut, "#include \"iec61850_model.h\"")
	fprintln(g.hOut, "")
}

func (g *StaticModelGenerator) printHeaderFileFooter() {
	fprintln(g.hOut, "")
	fprintln(g.hOut, "#endif /* "+g.hDefineName+" */\n")
}

func (g *StaticModelGenerator) printInitializerFunction() {
	fprintln(g.cOut, "\nstatic void\ninitializeValues()")
	fprintln(g.cOut, "{")
	fprint(g.cOut, g.initializerBuffer.String())
	fprintln(g.cOut, "}")
}

func (g *StaticModelGenerator) printVariablePointerDefines() {
	fprintln(g.hOut, "\n\n")
	for _, variableName := range g.variablesList {
		name := strings.ToUpper(g.modelPrefix) + variableName[len(g.modelPrefix):]
		fprintln(g.hOut, "#define "+name+" (&"+variableName+")")
	}
}

func (g *StaticModelGenerator) printForwardDeclarations(server *model.Server) {
	fprintln(g.cOut, "static void initializeValues();")
	fprintln(g.hOut, "extern IedModel "+g.modelPrefix+";")

	for _, ld := range server.GetLogicalDevices() {
		ldName := g.modelPrefix + "_" + ld.GetInst()
		fprintln(g.hOut, "extern LogicalDevice "+ldName+";")

		for _, ln := range ld.GetLogicalNodes() {
			lnName := ldName + "_" + ln.GetName()
			fprintln(g.hOut, "extern LogicalNode   "+lnName+";")
			g.printDataObjectForwardDeclarations(lnName, ln.GetDataObjects())
		}
	}
}

func (g *StaticModelGenerator) printDataObjectForwardDeclarations(prefix string, dataObjects []*model.DataObject) {
	for _, do := range dataObjects {
		doName := prefix + "_" + do.GetName()
		fprintln(g.hOut, "extern DataObject    "+doName+";")

		if do.GetCount() > 0 {
			for idx := 0; idx < do.GetCount(); idx++ {
				arrayElemName := doName + "_" + strconv.Itoa(idx)
				fprintln(g.hOut, "extern DataObject    "+arrayElemName+";")
				if do.GetSubDataObjects() != nil {
					g.printDataObjectForwardDeclarations(arrayElemName, do.GetSubDataObjects())
				}
				g.printDataAttributeForwardDeclarations(arrayElemName, do.GetDataAttributes())
			}
		} else {
			if do.GetSubDataObjects() != nil {
				g.printDataObjectForwardDeclarations(doName, do.GetSubDataObjects())
			}
			g.printDataAttributeForwardDeclarations(doName, do.GetDataAttributes())
		}
	}
}

func (g *StaticModelGenerator) printDataAttributeForwardDeclarations(doName string, dataAttributes []*model.DataAttribute) {
	for _, da := range dataAttributes {
		daName := doName + "_" + da.GetName()
		if da.GetFC() == model.FC_SE {
			if !strings.HasPrefix(daName, g.modelPrefix+"_SE_") {
				daName = daName[:9] + "SE_" + daName[9:]
			}
		}
		fprintln(g.hOut, "extern DataAttribute "+daName+";")

		if da.GetCount() > 0 {
			for idx := 0; idx < da.GetCount(); idx++ {
				arrayElemDaName := daName + "_" + strconv.Itoa(idx)
				fprintln(g.hOut, "extern DataAttribute "+arrayElemDaName+";")
				if da.GetSubDataAttributes() != nil {
					g.printDataAttributeForwardDeclarations(arrayElemDaName, da.GetSubDataAttributes())
				}
			}
		} else {
			if da.GetSubDataAttributes() != nil {
				g.printDataAttributeForwardDeclarations(daName, da.GetSubDataAttributes())
			}
		}
	}
}

func (g *StaticModelGenerator) printDeviceModelDefinitions() {
	g.printDataSets()

	logicalDevices := g.accessPoint.GetServer().GetLogicalDevices()

	g.createReportVariableList(logicalDevices)
	g.createLogControlVariableList(logicalDevices)
	g.createLogVariableList(logicalDevices)
	g.createGooseVariableList(logicalDevices)
	g.createSmvVariableList(logicalDevices)
	g.createSettingControlsVariableList(logicalDevices)

	for i, ld := range logicalDevices {
		ldName := g.modelPrefix + "_" + ld.GetInst()
		g.variablesList = append(g.variablesList, ldName)

		fprintln(g.cOut, "\nLogicalDevice "+ldName+" = {")
		fprintln(g.cOut, "    LogicalDeviceModelType,")
		fprintln(g.cOut, "    \""+ld.GetInst()+"\",")
		fprintln(g.cOut, "    (ModelNode*) &"+g.modelPrefix+",")

		if i < len(logicalDevices)-1 {
			fprintln(g.cOut, "    (ModelNode*) &"+g.modelPrefix+"_"+logicalDevices[i+1].GetInst()+",")
		} else {
			fprintln(g.cOut, "    NULL,")
		}

		firstChildName := ldName + "_" + ld.GetLogicalNodes()[0].GetName()
		fprintln(g.cOut, "    (ModelNode*) &"+firstChildName+",")

		if ld.GetLdName() != "" {
			fprintln(g.cOut, "    \""+ld.GetLdName()+"\"")
		} else {
			fprintln(g.cOut, "    NULL")
		}
		fprintln(g.cOut, "};\n")

		g.printLogicalNodeDefinitions(ldName, ld, ld.GetLogicalNodes())
	}

	for _, rcb := range g.rcbVariableNames {
		fprintln(g.cOut, "extern ReportControlBlock "+rcb+";")
	}
	fprintln(g.cOut, "")
	fprint(g.cOut, g.reportControlBlocks.String())

	for _, smv := range g.smvVariableNames {
		fprintln(g.cOut, "extern SVControlBlock "+smv+";")
	}
	fprint(g.cOut, g.smvControlBlocks.String())

	for _, gcb := range g.gseVariableNames {
		fprintln(g.cOut, "extern GSEControlBlock "+gcb+";")
	}
	fprint(g.cOut, g.gseControlBlocks.String())

	for _, sgcb := range g.sgcbVariableNames {
		fprintln(g.cOut, "extern SettingGroupControlBlock "+sgcb+";")
	}
	fprint(g.cOut, g.settingGroupControlBlocks.String())

	for _, lcb := range g.lcbVariableNames {
		fprintln(g.cOut, "extern LogControlBlock "+lcb+";")
	}
	fprint(g.cOut, g.logControlBlocks.String())

	for _, log := range g.logVariableNames {
		fprintln(g.cOut, "extern Log "+log+";")
	}
	fprint(g.cOut, g.logs.String())

	firstLd := logicalDevices[0]
	fprintln(g.cOut, "\nIedModel "+g.modelPrefix+" = {")
	fprintln(g.cOut, "    \""+g.ied.GetName()+"\",")
	fprintln(g.cOut, "    &"+g.modelPrefix+"_"+firstLd.GetInst()+",")

	if len(g.dataSetNames) > 0 {
		fprintln(g.cOut, "    &"+g.dataSetNames[0]+",")
	} else {
		fprintln(g.cOut, "    NULL,")
	}
	if len(g.rcbVariableNames) > 0 {
		fprintln(g.cOut, "    &"+g.rcbVariableNames[0]+",")
	} else {
		fprintln(g.cOut, "    NULL,")
	}
	if len(g.gseVariableNames) > 0 {
		fprintln(g.cOut, "    &"+g.gseVariableNames[0]+",")
	} else {
		fprintln(g.cOut, "    NULL,")
	}
	if len(g.smvVariableNames) > 0 {
		fprintln(g.cOut, "    &"+g.smvVariableNames[0]+",")
	} else {
		fprintln(g.cOut, "    NULL,")
	}
	if len(g.sgcbVariableNames) > 0 {
		fprintln(g.cOut, "    &"+g.sgcbVariableNames[0]+",")
	} else {
		fprintln(g.cOut, "    NULL,")
	}
	if len(g.lcbVariableNames) > 0 {
		fprintln(g.cOut, "    &"+g.lcbVariableNames[0]+",")
	} else {
		fprintln(g.cOut, "    NULL,")
	}
	if len(g.logVariableNames) > 0 {
		fprintln(g.cOut, "    &"+g.logVariableNames[0]+",")
	} else {
		fprintln(g.cOut, "    NULL,")
	}
	fprintln(g.cOut, "    initializeValues\n};")
}

func (g *StaticModelGenerator) createReportVariableList(lds []*model.LogicalDevice) {
	for _, ld := range lds {
		ldName := ld.GetInst()
		rcbCount := 0
		for _, ln := range ld.GetLogicalNodes() {
			for _, rcb := range ln.GetReportControlBlocks() {
				maxInstances := 1
				if rcb.GetRptEna() != nil {
					maxInstances = rcb.GetRptEna().GetMaxInstances()
				}
				for i := 0; i < maxInstances; i++ {
					name := g.modelPrefix + "_" + ldName + "_" + ln.GetName() + "_report" + strconv.Itoa(rcbCount)
					g.rcbVariableNames = append(g.rcbVariableNames, name)
					rcbCount++
				}
			}
		}
	}
}

func (g *StaticModelGenerator) createLogControlVariableList(lds []*model.LogicalDevice) {
	for _, ld := range lds {
		ldName := ld.GetInst()
		lcbCount := 0
		for _, ln := range ld.GetLogicalNodes() {
			for range ln.GetLogControlBlocks() {
				name := g.modelPrefix + "_" + ldName + "_" + ln.GetName() + "_lcb" + strconv.Itoa(lcbCount)
				g.lcbVariableNames = append(g.lcbVariableNames, name)
				lcbCount++
			}
		}
	}
}

func (g *StaticModelGenerator) createLogVariableList(lds []*model.LogicalDevice) {
	for _, ld := range lds {
		ldName := ld.GetInst()
		logCount := 0
		for _, ln := range ld.GetLogicalNodes() {
			for range ln.GetLogs() {
				name := g.modelPrefix + "_" + ldName + "_" + ln.GetName() + "_log" + strconv.Itoa(logCount)
				g.logVariableNames = append(g.logVariableNames, name)
				logCount++
			}
		}
	}
}

func (g *StaticModelGenerator) createGooseVariableList(lds []*model.LogicalDevice) {
	for _, ld := range lds {
		ldName := ld.GetInst()
		gseCount := 0
		for _, ln := range ld.GetLogicalNodes() {
			for range ln.GetGSEControlBlocks() {
				name := g.modelPrefix + "_" + ldName + "_" + ln.GetName() + "_gse" + strconv.Itoa(gseCount)
				g.gseVariableNames = append(g.gseVariableNames, name)
				gseCount++
			}
		}
	}
}

func (g *StaticModelGenerator) createSmvVariableList(lds []*model.LogicalDevice) {
	for _, ld := range lds {
		ldName := ld.GetInst()
		smvCount := 0
		for _, ln := range ld.GetLogicalNodes() {
			for range ln.GetSampledValueControlBlocks() {
				name := g.modelPrefix + "_" + ldName + "_" + ln.GetName() + "_smv" + strconv.Itoa(smvCount)
				g.smvVariableNames = append(g.smvVariableNames, name)
				smvCount++
			}
		}
	}
}

func (g *StaticModelGenerator) createSettingControlsVariableList(lds []*model.LogicalDevice) {
	for _, ld := range lds {
		ldName := ld.GetInst()
		for _, ln := range ld.GetLogicalNodes() {
			for range ln.GetSettingGroupControlBlocks() {
				name := g.modelPrefix + "_" + ldName + "_" + ln.GetName() + "_sgcb"
				g.sgcbVariableNames = append(g.sgcbVariableNames, name)
			}
		}
	}
}

func (g *StaticModelGenerator) printLogicalNodeDefinitions(ldName string, ld *model.LogicalDevice, lns []*model.LogicalNode) {
	for i, ln := range lns {
		lnName := ldName + "_" + ln.GetName()
		g.variablesList = append(g.variablesList, lnName)

		fprintln(g.cOut, "LogicalNode "+lnName+" = {")
		fprintln(g.cOut, "    LogicalNodeModelType,")
		fprintln(g.cOut, "    \""+ln.GetName()+"\",")
		fprintln(g.cOut, "    (ModelNode*) &"+ldName+",")
		if i < len(lns)-1 {
			fprintln(g.cOut, "    (ModelNode*) &"+ldName+"_"+lns[i+1].GetName()+",")
		} else {
			fprintln(g.cOut, "    NULL,")
		}
		firstChildName := lnName + "_" + ln.GetDataObjects()[0].GetName()
		fprintln(g.cOut, "    (ModelNode*) &"+firstChildName+",")
		fprintln(g.cOut, "};\n")

		g.printDataObjectDefinitions(lnName, ln.GetDataObjects(), "", false, -1)
		g.printReportControlBlocks(lnName, ln)
		g.printLogControlBlocks(ldName, lnName, ln, ld)
		g.printLogs(lnName, ln)
		g.printGSEControlBlocks(ldName, lnName, ln)
		g.printSVControlBlocks(ldName, lnName, ln)
		g.printSettingControlBlock(lnName, ln)
	}
}

func (g *StaticModelGenerator) printDataObjectDefinitions(lnName string, dos []*model.DataObject, dataAttributeSibling string, isTransient bool, arrayIdx int) {
	for i, do := range dos {
		isArray := do.GetCount() > 0

		doName := lnName + "_" + do.GetName()
		g.variablesList = append(g.variablesList, doName)

		fprintln(g.cOut, "DataObject "+doName+" = {")
		fprintln(g.cOut, "    DataObjectModelType,")
		fprintln(g.cOut, "    \""+do.GetName()+"\",")
		fprintln(g.cOut, "    (ModelNode*) &"+lnName+",")

		if i < len(dos)-1 {
			fprintln(g.cOut, "    (ModelNode*) &"+lnName+"_"+dos[i+1].GetName()+",")
		} else if dataAttributeSibling != "" {
			fprintln(g.cOut, "    (ModelNode*) &"+dataAttributeSibling+",")
		} else {
			fprintln(g.cOut, "    NULL,")
		}

		isDoTransient := isTransient || do.IsTransient()

		if isArray {
			childName := doName + "_0"
			fprintln(g.cOut, "    (ModelNode*) &"+childName+",")
			fprintln(g.cOut, "    "+strconv.Itoa(do.GetCount())+",")
			fprintln(g.cOut, "    0"+strconv.Itoa(arrayIdx))
			fprintln(g.cOut, "};\n")

			for idx := 0; idx < do.GetCount(); idx++ {
				elemName := doName + "_" + strconv.Itoa(idx)
				fprintln(g.cOut, "DataObject "+elemName+" = {")
				fprintln(g.cOut, "    DataObjectModelType,")
				fprintln(g.cOut, "    \""+do.GetName()+"\",")
				fprintln(g.cOut, "    (ModelNode*) &"+doName+",")

				if idx != do.GetCount()-1 {
					fprintln(g.cOut, "    (ModelNode*) &"+doName+"_"+strconv.Itoa(idx+1)+",")
				} else {
					fprintln(g.cOut, "    NULL,")
				}

				var firstSubDOName, firstDAName string
				if len(do.GetSubDataObjects()) > 0 {
					firstSubDOName = elemName + "_" + do.GetSubDataObjects()[0].GetName()
				}
				if len(do.GetDataAttributes()) > 0 {
					firstDAName = elemName + "_" + do.GetDataAttributes()[0].GetName()
				}

				if firstSubDOName != "" {
					fprintln(g.cOut, "    (ModelNode*) &"+firstSubDOName+",")
				} else if firstDAName != "" {
					fprintln(g.cOut, "    (ModelNode*) &"+firstDAName+",")
				} else {
					fprintln(g.cOut, "    NULL,")
				}

				fprintln(g.cOut, "    0,")
				fprintln(g.cOut, "    "+strconv.Itoa(idx))
				fprintln(g.cOut, "};\n")

				if do.GetSubDataObjects() != nil {
					g.printDataObjectDefinitions(elemName, do.GetSubDataObjects(), firstDAName, isDoTransient, -1)
				}
				if do.GetDataAttributes() != nil {
					g.printDataAttributeDefinitions(elemName, do.GetDataAttributes(), isDoTransient, -1)
				}
			}
		} else {
			var firstSubDOName, firstDAName string
			if len(do.GetSubDataObjects()) > 0 {
				firstSubDOName = doName + "_" + do.GetSubDataObjects()[0].GetName()
			}
			if len(do.GetDataAttributes()) > 0 {
				firstDAName = doName + "_" + do.GetDataAttributes()[0].GetName()
			}

			if firstSubDOName != "" {
				fprintln(g.cOut, "    (ModelNode*) &"+firstSubDOName+",")
			} else if firstDAName != "" {
				fprintln(g.cOut, "    (ModelNode*) &"+firstDAName+",")
			} else {
				fprintln(g.cOut, "    NULL,")
			}
			fprintln(g.cOut, "    0,")
			fprintln(g.cOut, "    "+strconv.Itoa(arrayIdx))
			fprintln(g.cOut, "};\n")

			if do.GetSubDataObjects() != nil {
				g.printDataObjectDefinitions(doName, do.GetSubDataObjects(), firstDAName, isDoTransient, -1)
			}
			if do.GetDataAttributes() != nil {
				g.printDataAttributeDefinitions(doName, do.GetDataAttributes(), isDoTransient, -1)
			}
		}
	}
}

func (g *StaticModelGenerator) seDaName(doName, daName string, da *model.DataAttribute) string {
	full := doName + "_" + daName
	if da.GetFC() == model.FC_SE {
		if !strings.HasPrefix(full, g.modelPrefix+"_SE_") {
			full = full[:9] + "SE_" + full[9:]
		}
	}
	return full
}

func (g *StaticModelGenerator) printDataAttributeDefinitions(doName string, das []*model.DataAttribute, isTransient bool, arrayIdx int) {
	for i, da := range das {
		isArray := da.GetCount() > 0

		daName := g.seDaName(doName, da.GetName(), da)
		g.variablesList = append(g.variablesList, daName)

		fprintln(g.cOut, "DataAttribute "+daName+" = {")
		fprintln(g.cOut, "    DataAttributeModelType,")
		fprintln(g.cOut, "    \""+da.GetName()+"\",")
		fprintln(g.cOut, "    (ModelNode*) &"+doName+",")

		if i < len(das)-1 {
			sibling := das[i+1]
			siblingName := g.seDaName(doName, sibling.GetName(), sibling)
			fprintln(g.cOut, "    (ModelNode*) &"+siblingName+",")
		} else {
			fprintln(g.cOut, "    NULL,")
		}

		if isArray {
			fprintln(g.cOut, "    (ModelNode*) &"+daName+"_0,")
			fprintln(g.cOut, "    "+strconv.Itoa(da.GetCount())+",")
			fprintln(g.cOut, "    -1,")
			fprintln(g.cOut, "    IEC61850_FC_"+da.GetFC().String()+",")
			fprintln(g.cOut, "    IEC61850_"+da.GetType().String()+",")
			g.printTriggerOptions(da.GetTriggerOptions(), isTransient)
			fprintln(g.cOut, "    NULL,")
			fprint(g.cOut, "    "+g.sAddrValue(da))
			fprintln(g.cOut, "};\n")

			for idx := 0; idx < da.GetCount(); idx++ {
				elemName := daName + "_" + strconv.Itoa(idx)
				g.variablesList = append(g.variablesList, elemName)

				fprintln(g.cOut, "DataAttribute "+elemName+" = {")
				fprintln(g.cOut, "    DataAttributeModelType,")
				fprintln(g.cOut, "    NULL,")
				fprintln(g.cOut, "    (ModelNode*) &"+daName+",")

				if idx != da.GetCount()-1 {
					fprintln(g.cOut, "    (ModelNode*) &"+daName+"_"+strconv.Itoa(idx+1)+",")
				} else {
					fprintln(g.cOut, "    NULL,")
				}

				if len(da.GetSubDataAttributes()) > 0 {
					fprintln(g.cOut, "    (ModelNode*) &"+elemName+"_"+da.GetSubDataAttributes()[0].GetName()+",")
				} else {
					fprintln(g.cOut, "    NULL,")
				}
				fprintln(g.cOut, "    0,")
				fprintln(g.cOut, "    "+strconv.Itoa(idx)+",")
				fprintln(g.cOut, "    IEC61850_FC_"+da.GetFC().String()+",")
				fprintln(g.cOut, "    IEC61850_"+da.GetType().String()+",")
				g.printTriggerOptions(da.GetTriggerOptions(), isTransient)
				fprintln(g.cOut, "    NULL,")
				fprint(g.cOut, "    0")
				fprintln(g.cOut, "};\n")

				if da.GetSubDataAttributes() != nil {
					g.printDataAttributeDefinitions(elemName, da.GetSubDataAttributes(), isTransient, -1)
				}
			}
		} else {
			if len(da.GetSubDataAttributes()) > 0 {
				fprintln(g.cOut, "    (ModelNode*) &"+daName+"_"+da.GetSubDataAttributes()[0].GetName()+",")
			} else {
				fprintln(g.cOut, "    NULL,")
			}
			fprintln(g.cOut, "    "+strconv.Itoa(da.GetCount())+",")
			fprintln(g.cOut, "    -1,")
			fprintln(g.cOut, "    IEC61850_FC_"+da.GetFC().String()+",")
			fprintln(g.cOut, "    IEC61850_"+da.GetType().String()+",")
			g.printTriggerOptions(da.GetTriggerOptions(), isTransient)
			fprintln(g.cOut, "    NULL,")
			fprint(g.cOut, "    "+g.sAddrValue(da))
			fprintln(g.cOut, "};\n")

			if da.GetSubDataAttributes() != nil {
				g.printDataAttributeDefinitions(daName, da.GetSubDataAttributes(), isTransient, -1)
			}

			value := da.GetEffectiveValue(g.typeDeclarations)
			if value != nil {
				g.printValue(daName, da, value)
			}
		}
	}
}

func (g *StaticModelGenerator) printTriggerOptions(trgOps *model.TriggerOptions, isTransient bool) {
	s := "    0"
	if trgOps.IsDchg() {
		s += " + TRG_OPT_DATA_CHANGED"
	}
	if trgOps.IsDupd() {
		s += " + TRG_OPT_DATA_UPDATE"
	}
	if trgOps.IsQchg() {
		s += " + TRG_OPT_QUALITY_CHANGED"
	}
	if isTransient {
		s += " + TRG_OPT_TRANSIENT"
	}
	fprintln(g.cOut, s+",")
}

func (g *StaticModelGenerator) sAddrValue(da *model.DataAttribute) string {
	if da.GetShortAddress() != "" {
		v, err := strconv.ParseInt(da.GetShortAddress(), 10, 64)
		if err != nil {
			fmt.Printf("WARNING: short address %q is not valid for libIEC61850!\n", da.GetShortAddress())
			return "0"
		}
		return strconv.FormatInt(v, 10)
	}
	return "0"
}

func (g *StaticModelGenerator) appendHexArray(buf *strings.Builder, data []byte) {
	buf.WriteString("{")
	for i, b := range data {
		if i == 0 {
			buf.WriteString(fmt.Sprintf("0x%02X", b))
		} else {
			buf.WriteString(fmt.Sprintf(", 0x%02X", b))
		}
	}
	buf.WriteString("}")
}

func (g *StaticModelGenerator) printValue(daName string, da *model.DataAttribute, value *model.DataModelValue) {
	buf := &g.initializerBuffer
	buf.WriteString("\n")
	if g.initializeOnce {
		buf.WriteString("if (!")
		buf.WriteString(daName)
		buf.WriteString(".mmsValue)\n")
	}
	buf.WriteString(daName)
	buf.WriteString(".mmsValue = ")

	switch da.GetType() {
	case model.AT_ENUMERATED, model.AT_INT8, model.AT_INT16, model.AT_INT32, model.AT_INT64:
		buf.WriteString("MmsValue_newIntegerFromInt32(" + strconv.Itoa(value.GetIntValue()) + ");")
	case model.AT_INT8U, model.AT_INT16U, model.AT_INT24U, model.AT_INT32U:
		buf.WriteString("MmsValue_newUnsignedFromUint32(" + strconv.FormatInt(value.GetLongValue(), 10) + ");")
	case model.AT_BOOLEAN:
		buf.WriteString("MmsValue_newBoolean(" + fmt.Sprintf("%v", value.GetValue()) + ");")
	case model.AT_OCTET_STRING_64:
		daValName := daName + "__val"
		buf.WriteString("MmsValue_newOctetString(0, 64);\n")
		buf.WriteString("uint8_t " + daValName + "[] = ")
		if b, ok := value.GetValue().([]byte); ok {
			g.appendHexArray(buf, b)
			buf.WriteString(";\n")
			buf.WriteString("MmsValue_setOctetString(")
			buf.WriteString(daName)
			buf.WriteString(".mmsValue, " + daValName + ", " + strconv.Itoa(len(b)) + ");\n")
		} else {
			buf.WriteString("{};\n")
		}
	case model.AT_CODEDENUM:
		buf.WriteString("MmsValue_newBitString(2);\n")
		buf.WriteString("MmsValue_setBitStringFromIntegerBigEndian(")
		buf.WriteString(daName)
		buf.WriteString(".mmsValue, ")
		buf.WriteString(fmt.Sprintf("%v", value.GetValue()))
		buf.WriteString(");\n")
	case model.AT_UNICODE_STRING_255:
		buf.WriteString("MmsValue_newMmsString(\"" + fmt.Sprintf("%v", value.GetValue()) + "\");")
	case model.AT_VISIBLE_STRING_32, model.AT_VISIBLE_STRING_64, model.AT_VISIBLE_STRING_129,
		model.AT_VISIBLE_STRING_255, model.AT_VISIBLE_STRING_65, model.AT_CURRENCY:
		buf.WriteString("MmsValue_newVisibleString(\"" + fmt.Sprintf("%v", value.GetValue()) + "\");")
	case model.AT_FLOAT32:
		buf.WriteString("MmsValue_newFloat(" + fmt.Sprintf("%v", value.GetValue()) + ");")
	case model.AT_FLOAT64:
		buf.WriteString("MmsValue_newDouble(" + fmt.Sprintf("%v", value.GetValue()) + ");")
	case model.AT_TIMESTAMP:
		buf.WriteString("MmsValue_newUtcTimeByMsTime(" + fmt.Sprintf("%v", value.GetValue()) + ");")
	default:
		fmt.Printf("Unknown default value for %s type: %v\n", daName, da.GetType())
		buf.WriteString("NULL;")
	}
	buf.WriteString("\n")
}

func (g *StaticModelGenerator) printReportControlBlocks(lnPrefix string, ln *model.LogicalNode) {
	reportsCount := len(ln.GetReportControlBlocks())
	reportNumber := 0
	_ = reportsCount

	for _, rcb := range ln.GetReportControlBlocks() {
		if rcb.IsIndexed() {
			maxInstances := 1
			var clientLNs []*model.ClientLN
			if rcb.GetRptEna() != nil {
				maxInstances = rcb.GetRptEna().GetMaxInstances()
				clientLNs = rcb.GetRptEna().GetClientLNs()
			}

			for i := 0; i < maxInstances; i++ {
				index := fmt.Sprintf("%02d", i+1)
				fmt.Printf("print report instance %s\n", index)
				clientAddress := make([]byte, 17)

				if clientLNs != nil && i < len(clientLNs) {
					cl := clientLNs[i]
					if cl != nil && cl.GetIedName() != "" {
						ip := g.getIpAddressByIedName(cl.GetIedName(), cl.GetApRef())
						if ip != "" {
							addr, err := net.ResolveIPAddr("ip", ip)
							if err == nil {
								if ipv4 := addr.IP.To4(); ipv4 != nil {
									clientAddress[0] = 4
									copy(clientAddress[1:], ipv4)
								} else if ipv6 := addr.IP.To16(); ipv6 != nil {
									clientAddress[0] = 6
									copy(clientAddress[1:], ipv6)
								}
							}
						}
					}
				}

				g.printReportControlBlockInstance(lnPrefix, rcb, index, reportNumber, clientAddress)
				reportNumber++
			}
		} else {
			clientAddress := make([]byte, 17)
			g.printReportControlBlockInstance(lnPrefix, rcb, "", reportNumber, clientAddress)
			reportNumber++
		}
	}
}

func (g *StaticModelGenerator) getIpAddressByIedName(iedName, apRef string) string {
	comm := g.sclParser.GetCommunication()
	if comm == nil {
		return ""
	}
	for _, sn := range comm.SubNetworks {
		for _, cap := range sn.ConnectedAPs {
			if apRef != "" && cap.APName != apRef {
				continue
			}
			if cap.IEDName == iedName && cap.Address != nil {
				return cap.Address.GetP("IP")
			}
		}
	}
	return ""
}

func (g *StaticModelGenerator) printReportControlBlockInstance(lnPrefix string, rcb *model.ReportControlBlock, index string, reportNumber int, clientIpAddr []byte) {
	rcbVarName := lnPrefix + "_report" + strconv.Itoa(reportNumber)

	var sb strings.Builder
	sb.WriteString("ReportControlBlock " + rcbVarName + " = {")
	sb.WriteString("&" + lnPrefix + ", ")
	sb.WriteString("\"" + rcb.GetName() + index + "\", ")

	if rcb.GetRptID() == "" {
		sb.WriteString("NULL, ")
	} else {
		sb.WriteString("\"" + rcb.GetRptID() + "\", ")
	}

	if rcb.IsBuffered() {
		sb.WriteString("true, ")
	} else {
		sb.WriteString("false, ")
	}

	if rcb.GetDataSet() != "" {
		sb.WriteString("\"" + rcb.GetDataSet() + "\", ")
	} else {
		sb.WriteString("NULL, ")
	}

	if rcb.GetConfRef() != nil {
		sb.WriteString(strconv.FormatInt(*rcb.GetConfRef(), 10) + ", ")
	} else {
		sb.WriteString("0, ")
	}

	triggerOps := 16
	if rcb.GetTriggerOptions() != nil {
		triggerOps = rcb.GetTriggerOptions().GetIntValue()
	}
	if g.hasOwner {
		triggerOps += 64
	}
	sb.WriteString(strconv.Itoa(triggerOps) + ", ")

	options := 0
	if rcb.GetOptionFields() != nil {
		of := rcb.GetOptionFields()
		if of.IsSeqNum() {
			options += 1
		}
		if of.IsTimeStamp() {
			options += 2
		}
		if of.IsReasonCode() {
			options += 4
		}
		if of.IsDataSet() {
			options += 8
		}
		if of.IsDataRef() {
			options += 16
		}
		if of.IsBufOvfl() {
			options += 32
		}
		if of.IsEntryID() {
			options += 64
		}
		if of.IsConfigRef() {
			options += 128
		}
	} else {
		options = 32
	}
	sb.WriteString(strconv.Itoa(options) + ", ")
	sb.WriteString(strconv.FormatInt(rcb.GetBufferTime(), 10) + ", ")

	if rcb.GetIntegrityPeriod() != nil {
		sb.WriteString(strconv.FormatInt(*rcb.GetIntegrityPeriod(), 10) + ", ")
	} else {
		sb.WriteString("0, ")
	}

	sb.WriteString("{")
	for i := 0; i < 17; i++ {
		sb.WriteString("0x" + fmt.Sprintf("%x", clientIpAddr[i]&0xff))
		if i == 16 {
			sb.WriteString("}, ")
		} else {
			sb.WriteString(", ")
		}
	}

	g.currentRcbVariableNumber++
	if g.currentRcbVariableNumber < len(g.rcbVariableNames) {
		sb.WriteString("&" + g.rcbVariableNames[g.currentRcbVariableNumber])
	} else {
		sb.WriteString("NULL")
	}
	sb.WriteString("};\n")
	g.reportControlBlocks.WriteString(sb.String())
}

func (g *StaticModelGenerator) printLogControlBlocks(ldName, lnPrefix string, ln *model.LogicalNode, ld *model.LogicalDevice) {
	lcbNumber := 0
	for _, lcb := range ln.GetLogControlBlocks() {
		g.printLogControlBlock(ld, lnPrefix, lcb, lcbNumber)
		lcbNumber++
	}
}

func (g *StaticModelGenerator) printLogControlBlock(ld *model.LogicalDevice, lnPrefix string, lcb *model.LogControl, lcbNumber int) {
	lcbVarName := lnPrefix + "_lcb" + strconv.Itoa(lcbNumber)

	var sb strings.Builder
	sb.WriteString("LogControlBlock " + lcbVarName + " = {")
	sb.WriteString("&" + lnPrefix + ", ")
	sb.WriteString("\"" + lcb.GetName() + "\", ")

	if lcb.GetDataSet() != "" {
		sb.WriteString("\"" + lcb.GetDataSet() + "\", ")
	} else {
		sb.WriteString("NULL, ")
	}

	logRef := ""
	if lcb.GetLdInst() == "" {
		logRef = ld.GetInst() + "/"
	} else {
		logRef = lcb.GetLdInst() + "/"
	}
	if lcb.GetLnClass() == "LLN0" {
		logRef += "LLN0$"
	} else {
		logRef += lcb.GetLnClass() + "$"
	}

	if lcb.GetLogName() != "" {
		sb.WriteString("\"" + logRef + lcb.GetLogName() + "\", ")
	} else {
		sb.WriteString("NULL, ")
	}

	triggerOps := 0
	if lcb.GetTriggerOptions() != nil {
		triggerOps = lcb.GetTriggerOptions().GetIntValue()
	}
	if triggerOps >= 16 {
		triggerOps -= 16
	}
	sb.WriteString(strconv.Itoa(triggerOps) + ", ")

	if lcb.GetIntgPd() != 0 {
		sb.WriteString(strconv.Itoa(lcb.GetIntgPd()) + ", ")
	} else {
		sb.WriteString("0, ")
	}

	if lcb.IsLogEna() {
		sb.WriteString("true, ")
	} else {
		sb.WriteString("false, ")
	}
	if lcb.IsReasonCode() {
		sb.WriteString("true, ")
	} else {
		sb.WriteString("false, ")
	}

	g.currentLcbVariableNumber++
	if g.currentLcbVariableNumber < len(g.lcbVariableNames) {
		sb.WriteString("&" + g.lcbVariableNames[g.currentLcbVariableNumber])
	} else {
		sb.WriteString("NULL")
	}
	sb.WriteString("};\n")
	g.logControlBlocks.WriteString(sb.String())
}

func (g *StaticModelGenerator) printLogs(lnPrefix string, ln *model.LogicalNode) {
	logNumber := 0
	for _, log := range ln.GetLogs() {
		g.printLog(lnPrefix, log, logNumber)
		logNumber++
	}
}

func (g *StaticModelGenerator) printLog(lnPrefix string, log *model.Log, logNumber int) {
	logVarName := lnPrefix + "_log" + strconv.Itoa(logNumber)

	var sb strings.Builder
	sb.WriteString("Log " + logVarName + " = {")
	sb.WriteString("&" + lnPrefix + ", ")
	sb.WriteString("\"" + log.GetName() + "\", ")

	g.currentLogVariableNumber++
	if g.currentLogVariableNumber < len(g.logVariableNames) {
		sb.WriteString("&" + g.logVariableNames[g.currentLogVariableNumber])
	} else {
		sb.WriteString("NULL")
	}
	sb.WriteString("};\n")
	g.logs.WriteString(sb.String())
}

func (g *StaticModelGenerator) printGSEControlBlocks(ldName, lnPrefix string, ln *model.LogicalNode) {
	parts := strings.SplitN(ldName, "_", 2)
	logicalDeviceName := ""
	if len(parts) > 1 {
		logicalDeviceName = parts[1]
	}
	gseControlNumber := 0

	for _, gcb := range ln.GetGSEControlBlocks() {
		gse := g.connectedAP.LookupGSE(logicalDeviceName, gcb.GetName())

		if gse != nil {
			gseAddress := gse.GetAddress()
			var sb strings.Builder

			phyComAddrName := ""
			if gseAddress != nil {
				phyComAddrName = lnPrefix + "_gse" + strconv.Itoa(gseControlNumber) + "_address"
				sb.WriteString("\nstatic PhyComAddress " + phyComAddrName + " = {\n")
				sb.WriteString("  " + strconv.Itoa(gseAddress.VlanPriority) + ",\n")
				sb.WriteString("  " + strconv.Itoa(gseAddress.VlanID) + ",\n")
				sb.WriteString("  " + strconv.Itoa(gseAddress.AppID) + ",\n")
				sb.WriteString("  {")
				for i := 0; i < 6; i++ {
					sb.WriteString("0x" + fmt.Sprintf("%x", gseAddress.MacAddress[i]))
					if i == 5 {
						sb.WriteString("}\n")
					} else {
						sb.WriteString(", ")
					}
				}
				sb.WriteString("};\n\n")
			}

			gseVarName := lnPrefix + "_gse" + strconv.Itoa(gseControlNumber)
			sb.WriteString("GSEControlBlock " + gseVarName + " = {")
			sb.WriteString("&" + lnPrefix + ", ")
			sb.WriteString("\"" + gcb.GetName() + "\", ")

			if gcb.GetAppID() == "" {
				sb.WriteString("NULL, ")
			} else {
				sb.WriteString("\"" + gcb.GetAppID() + "\", ")
			}

			if gcb.GetDataSet() != "" {
				sb.WriteString("\"" + gcb.GetDataSet() + "\", ")
			} else {
				sb.WriteString("NULL, ")
			}

			sb.WriteString(strconv.Itoa(gcb.GetConfRev()) + ", ")

			if gcb.IsFixedOffs() {
				sb.WriteString("true, ")
			} else {
				sb.WriteString("false, ")
			}

			if gseAddress != nil {
				sb.WriteString("&" + phyComAddrName + ", ")
			} else {
				sb.WriteString("NULL, ")
			}

			sb.WriteString(strconv.Itoa(gse.GetMinTime()) + ", ")
			sb.WriteString(strconv.Itoa(gse.GetMaxTime()) + ", ")

			g.currentGseVariableNumber++
			if g.currentGseVariableNumber < len(g.gseVariableNames) {
				sb.WriteString("&" + g.gseVariableNames[g.currentGseVariableNumber])
			} else {
				sb.WriteString("NULL")
			}
			sb.WriteString("};\n")
			g.gseControlBlocks.WriteString(sb.String())
			gseControlNumber++
		} else {
			fmt.Printf("GSE not found for GoCB %s\n", gcb.GetName())
		}
	}
}

func (g *StaticModelGenerator) printSVControlBlocks(ldName, lnPrefix string, ln *model.LogicalNode) {
	parts := strings.SplitN(ldName, "_", 2)
	logicalDeviceName := ""
	if len(parts) > 1 {
		logicalDeviceName = parts[1]
	}
	smvControlNumber := 0

	for _, svCB := range ln.GetSampledValueControlBlocks() {
		fmt.Printf("SVCB: %s\n", svCB.GetName())

		var svAddress *communication.PhyComAddress
		if g.connectedAP != nil {
			svAddress = g.connectedAP.LookupSMVAddress(logicalDeviceName, svCB.GetName())
		}

		var sb strings.Builder
		phyComAddrName := ""

		if svAddress != nil {
			phyComAddrName = lnPrefix + "_smv" + strconv.Itoa(smvControlNumber) + "_address"
			sb.WriteString("\nstatic PhyComAddress " + phyComAddrName + " = {\n")
			sb.WriteString("  " + strconv.Itoa(svAddress.VlanPriority) + ",\n")
			sb.WriteString("  " + strconv.Itoa(svAddress.VlanID) + ",\n")
			sb.WriteString("  " + strconv.Itoa(svAddress.AppID) + ",\n")
			sb.WriteString("  {")
			for i := 0; i < 6; i++ {
				sb.WriteString("0x" + fmt.Sprintf("%x", svAddress.MacAddress[i]))
				if i == 5 {
					sb.WriteString("}\n")
				} else {
					sb.WriteString(", ")
				}
			}
			sb.WriteString("};\n\n")
		}

		smvVarName := lnPrefix + "_smv" + strconv.Itoa(smvControlNumber)
		sb.WriteString("SVControlBlock " + smvVarName + " = {")
		sb.WriteString("&" + lnPrefix + ", ")
		sb.WriteString("\"" + svCB.GetName() + "\", ")

		if svCB.GetSmvID() == "" {
			sb.WriteString("NULL, ")
		} else {
			sb.WriteString("\"" + svCB.GetSmvID() + "\", ")
		}

		if svCB.GetDatSet() != "" {
			sb.WriteString("\"" + svCB.GetDatSet() + "\", ")
		} else {
			sb.WriteString("NULL, ")
		}

		sb.WriteString(strconv.Itoa(svCB.GetSmvOpts().GetIntValue()) + ", ")
		sb.WriteString(strconv.Itoa(svCB.GetSmpMod().GetValue()) + ", ")
		sb.WriteString(strconv.Itoa(svCB.GetSmpRate()) + ", ")
		sb.WriteString(strconv.Itoa(svCB.GetConfRev()) + ", ")

		if svAddress != nil {
			sb.WriteString("&" + phyComAddrName + ", ")
		} else {
			sb.WriteString("NULL, ")
		}

		if svCB.IsMulticast() {
			sb.WriteString("false, ")
		} else {
			sb.WriteString("true, ")
		}

		sb.WriteString(strconv.Itoa(svCB.GetNofASDI()) + ", ")

		g.currentSvCBVariableNumber++
		if g.currentSvCBVariableNumber < len(g.smvVariableNames) {
			sb.WriteString("&" + g.smvVariableNames[g.currentSvCBVariableNumber])
		} else {
			sb.WriteString("NULL")
		}
		sb.WriteString("};\n")
		g.smvControlBlocks.WriteString(sb.String())
		smvControlNumber++
	}
}

func (g *StaticModelGenerator) printSettingControlBlock(lnPrefix string, ln *model.LogicalNode) {
	sgcbs := ln.GetSettingGroupControlBlocks()
	if len(sgcbs) == 0 {
		return
	}
	fmt.Printf("print SGCB for %s\n", lnPrefix)
	sgcb := sgcbs[0]

	var sb strings.Builder
	sgcbVarName := lnPrefix + "_sgcb"
	sb.WriteString("\nSettingGroupControlBlock " + sgcbVarName + " = {")
	sb.WriteString("&" + lnPrefix + ", ")
	sb.WriteString(strconv.Itoa(sgcb.GetActSG()) + ", " + strconv.Itoa(sgcb.GetNumOfSGs()) + ", 0, false, 0, 0, ")

	if g.currentSGCBVariableNumber < len(g.sgcbVariableNames)-1 {
		sb.WriteString("&" + g.sgcbVariableNames[g.currentSGCBVariableNumber+1])
	} else {
		sb.WriteString("NULL")
	}
	sb.WriteString("};\n")
	g.settingGroupControlBlocks.WriteString(sb.String())
	g.currentSGCBVariableNumber++
}

func (g *StaticModelGenerator) printDataSets() {
	ldevices := g.accessPoint.GetServer().GetLogicalDevices()

	g.dataSetNames = nil

	for _, ld := range ldevices {
		for _, ln := range ld.GetLogicalNodes() {
			for _, ds := range ln.GetDataSets() {
				name := g.modelPrefix + "ds_" + ld.GetInst() + "_" + ln.GetName() + "_" + ds.GetName()
				g.dataSetNames = append(g.dataSetNames, name)
			}
		}
	}

	fprintln(g.cOut, "")
	for _, n := range g.dataSetNames {
		fprintln(g.cOut, "extern DataSet "+n+";")
	}
	fprintln(g.cOut, "")

	dataSetNameIdx := 0

	for _, ld := range ldevices {
		for _, ln := range ld.GetLogicalNodes() {
			for _, ds := range ln.GetDataSets() {
				dsVarName := g.dataSetNames[dataSetNameIdx]
				dataSetNameIdx++

				fcdaCount := 0
				totalFCDA := len(ds.GetFCDA())

				fprintln(g.cOut, "")
				for range ds.GetFCDA() {
					entryName := dsVarName + "_fcda" + strconv.Itoa(fcdaCount)
					fprintln(g.cOut, "extern DataSetEntry "+entryName+";")
					fcdaCount++
				}
				fprintln(g.cOut, "")

				fcdaCount = 0
				for _, fcda := range ds.GetFCDA() {
					entryName := dsVarName + "_fcda" + strconv.Itoa(fcdaCount)

					mmsVarName := ""
					if fcda.GetPrefix() != "" {
						mmsVarName += fcda.GetPrefix()
					}
					mmsVarName += fcda.GetLnClass()
					if fcda.GetLnInstance() != "" {
						mmsVarName += fcda.GetLnInstance()
					}
					mmsVarName += "$" + fcda.GetFC().String()
					mmsVarName += "$" + toMmsString(fcda.GetDoName())
					if fcda.GetDaName() != "" {
						mmsVarName += "$" + toMmsString(fcda.GetDaName())
					}

					variableName, arrayIndex, componentName := parseArrayRef(mmsVarName)

					fprintln(g.cOut, "DataSetEntry "+entryName+" = {")
					fprintln(g.cOut, "  \""+fcda.GetLdInstance()+"\",")
					fprintln(g.cOut, "  false,")
					fprintln(g.cOut, "  \""+variableName+"\", ")
					fprintln(g.cOut, "  "+strconv.Itoa(arrayIndex)+",")
					if componentName == "" {
						fprintln(g.cOut, "  NULL,")
					} else {
						fprintln(g.cOut, "  \""+componentName+"\",")
					}
					fprintln(g.cOut, "  NULL,")

					if fcdaCount+1 < totalFCDA {
						fprintln(g.cOut, "  &"+dsVarName+"_fcda"+strconv.Itoa(fcdaCount+1))
					} else {
						fprintln(g.cOut, "  NULL")
					}
					fprintln(g.cOut, "};\n")
					fcdaCount++
				}

				fprintln(g.cOut, "DataSet "+dsVarName+" = {")
				fprintln(g.cOut, "  \""+ld.GetInst()+"\",")
				fprintln(g.cOut, "  \""+ln.GetName()+"$"+ds.GetName()+"\",")
				fprintln(g.cOut, "  "+strconv.Itoa(len(ds.GetFCDA()))+",")
				fprintln(g.cOut, "  &"+dsVarName+"_fcda0,")

				if dataSetNameIdx < len(g.dataSetNames) {
					fprintln(g.cOut, "  &"+g.dataSetNames[dataSetNameIdx])
				} else {
					fprintln(g.cOut, "  NULL")
				}
				fprintln(g.cOut, "};")
			}
		}
	}
}

func toMmsString(s string) string {
	return strings.ReplaceAll(s, ".", "$")
}

func parseArrayRef(mmsVarName string) (variable string, arrayIndex int, component string) {
	arrayIndex = -1
	arrayStart := strings.Index(mmsVarName, "(")
	if arrayStart == -1 {
		return mmsVarName, arrayIndex, ""
	}

	variable = mmsVarName[:arrayStart]
	arrayEnd := strings.Index(mmsVarName, ")")
	idxStr := mmsVarName[arrayStart+1 : arrayEnd]
	idx, _ := strconv.Atoi(idxStr)
	arrayIndex = idx

	rest := mmsVarName[arrayEnd+1:]
	if len(rest) > 0 && rest[0] == '$' {
		rest = rest[1:]
	}
	if len(rest) > 0 {
		component = rest
	}
	return
}
