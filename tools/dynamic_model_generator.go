/*
 *  dynamic_model_generator.go
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
	"fmt"
	"io"
)

// RunDynamicModelGenerator generates the dynamic text model to output.
func RunDynamicModelGenerator(r io.Reader, icdFile string, output io.Writer, iedName, accessPointName string) error {
	sclParser, err := parser.NewSclParser(r, true)
	if err != nil {
		return err
	}

	var ied *model.IED
	if iedName == "" {
		ied = sclParser.GetFirstIed()
	} else {
		ied = sclParser.GetIedByName(iedName)
	}
	if ied == nil {
		return fmt.Errorf("no data model present in SCL file")
	}

	hasOwner := false
	services := ied.GetServices()
	if services != nil {
		rptSettings := services.GetReportSettings()
		if rptSettings != nil {
			hasOwner = rptSettings.HasOwner()
		}
	}

	var accessPoint *model.AccessPoint
	if accessPointName != "" {
		accessPoint = ied.GetAccessPointByName(accessPointName)
	} else {
		accessPoint = ied.GetFirstAccessPoint()
	}
	if accessPoint == nil {
		return fmt.Errorf("no valid access point found")
	}

	connectedAP := sclParser.GetConnectedAP(ied.GetName(), accessPoint.GetName())
	connectedAPs := sclParser.GetConnectedAPsAll()

	ldevices := accessPoint.GetServer().GetLogicalDevices()

	fmt.Fprintf(output, "MODEL(%s){\n", ied.GetName())
	for _, ld := range ldevices {
		fmt.Fprintf(output, "LD(%s", ld.GetInst())
		if ld.GetLdName() != "" {
			fmt.Fprintf(output, " %s", ld.GetLdName())
		}
		fmt.Fprintf(output, "){\n")

		for _, ln := range ld.GetLogicalNodes() {
			fmt.Fprintf(output, "LN(%s){\n", ln.GetName())
			exportLogicalNode(output, ln, ld, ied, connectedAP, connectedAPs, hasOwner, sclParser.GetTypeDeclarations())
			fmt.Fprintln(output, "}")
		}

		fmt.Fprintln(output, "}")
	}
	fmt.Fprintln(output, "}")
	return nil
}

func exportLogicalNode(output io.Writer, ln *model.LogicalNode, ld *model.LogicalDevice, ied *model.IED,
	connectedAP *communication.ConnectedAP, connectedAPs []*communication.ConnectedAP,
	hasOwner bool, td interface{}) {

	for _, sgcb := range ln.GetSettingGroupControlBlocks() {
		fmt.Fprintf(output, "SG(%d %d)\n", sgcb.GetActSG(), sgcb.GetNumOfSGs())
	}

	for _, do := range ln.GetDataObjects() {
		fmt.Fprintf(output, "DO(%s %d){\n", do.GetName(), do.GetCount())
		exportDataObject(output, do, false)
		fmt.Fprintln(output, "}")
	}

	for _, ds := range ln.GetDataSets() {
		exportDataSet(output, ds, ln)
	}

	for _, rcb := range ln.GetReportControlBlocks() {
		if rcb.IsIndexed() {
			maxInstances := 1
			if rcb.GetRptEna() != nil {
				maxInstances = rcb.GetRptEna().GetMaxInstances()
			}
			for i := 0; i < maxInstances; i++ {
				index := fmt.Sprintf("%02d", i+1)
				printRCBInstance(output, rcb, index, hasOwner)
			}
		} else {
			printRCBInstance(output, rcb, "", hasOwner)
		}
	}

	for _, lcb := range ln.GetLogControlBlocks() {
		printLCB(output, lcb, ln, ld)
	}

	for _, log := range ln.GetLogs() {
		fmt.Fprintf(output, "LOG(%s);\n", log.GetName())
	}

	for _, svcb := range ln.GetSampledValueControlBlocks() {
		var smv *communication.SMV
		var smvAddress *communication.PhyComAddress

		if connectedAP != nil {
			smv = connectedAP.LookupSMV(ld.GetInst(), svcb.GetName())
			if smv == nil {
				for _, ap := range connectedAPs {
					smv = ap.LookupSMV(ld.GetInst(), svcb.GetName())
					if smv != nil {
						break
					}
				}
			}
			if smv == nil {
				fmt.Println("ConnectedAP not found for SMV")
			}
			if smv != nil {
				smvAddress = smv.Address
			}
		} else {
			fmt.Printf("WARNING: IED %q has no connected access point!\n", ied.GetName())
		}

		fmt.Fprintf(output, "SMVC(%s %s %s %d %d %d %d %s)",
			svcb.GetName(), svcb.GetSmvID(), svcb.GetDatSet(),
			svcb.GetConfRev(), svcb.GetSmpMod().GetValue(),
			svcb.GetSmpRate(), svcb.GetSmvOpts().GetIntValue(),
			boolToStr(!svcb.IsMulticast(), "1", "0"),
		)

		if smvAddress != nil {
			fmt.Fprintln(output, "{")
			fmt.Fprintf(output, "PA(%d %d %d ",
				smvAddress.VlanPriority, smvAddress.VlanID, smvAddress.AppID)
			for i := 0; i < 6; i++ {
				fmt.Fprintf(output, "%02x", smvAddress.MacAddress[i])
			}
			fmt.Fprintln(output, ");")
			fmt.Fprintln(output, "}")
		} else {
			fmt.Fprintln(output, ";")
		}
	}

	for _, gcb := range ln.GetGSEControlBlocks() {
		var gse *communication.GSE
		var gseAddress *communication.PhyComAddress

		if connectedAP != nil {
			gse = connectedAP.LookupGSE(ld.GetInst(), gcb.GetName())
			if gse == nil {
				for _, ap := range connectedAPs {
					gse = ap.LookupGSE(ld.GetInst(), gcb.GetName())
					if gse != nil {
						break
					}
				}
			}
			if gse == nil {
				fmt.Println("ConnectedAP not found for GSE")
			}
			if gse != nil {
				gseAddress = gse.GetAddress()
			}
		} else {
			fmt.Printf("WARNING: IED %q has no connected access point!\n", ied.GetName())
		}

		minTime, maxTime := -1, -1
		if gse != nil {
			minTime = gse.GetMinTime()
			maxTime = gse.GetMaxTime()
		}

		fmt.Fprintf(output, "GC(%s %s %s %d %s %d %d ",
			gcb.GetName(), gcb.GetAppID(), gcb.GetDataSet(),
			gcb.GetConfRev(),
			boolToStr(gcb.IsFixedOffs(), "1", "0"),
			minTime, maxTime,
		)

		if gseAddress == nil {
			fmt.Fprintln(output, ");")
		} else {
			fmt.Fprintln(output, "){")
			fmt.Fprintf(output, "PA(%d %d %d ",
				gseAddress.VlanPriority, gseAddress.VlanID, gseAddress.AppID)
			for i := 0; i < 6; i++ {
				fmt.Fprintf(output, "%02x", gseAddress.MacAddress[i])
			}
			fmt.Fprintln(output, ");")
			fmt.Fprintln(output, "}")
		}
	}
}

func boolToStr(b bool, ifTrue, ifFalse string) string {
	if b {
		return ifTrue
	}
	return ifFalse
}

func printLCB(output io.Writer, lcb *model.LogControl, ln *model.LogicalNode, ld *model.LogicalDevice) {
	fmt.Fprintf(output, "LC(%s ", lcb.GetName())

	if lcb.GetDataSet() != "" {
		fmt.Fprintf(output, "%s ", lcb.GetDataSet())
	} else {
		fmt.Fprint(output, "- ")
	}

	if lcb.GetLogName() != "" {
		logRef := ld.GetInst() + "/" + ln.GetName() + "$" + lcb.GetLogName()
		fmt.Fprintf(output, "%s ", logRef)
	} else {
		fmt.Fprint(output, "- ")
	}

	fmt.Fprintf(output, "%d ", lcb.GetTriggerOptions().GetIntValue())
	fmt.Fprintf(output, "%d ", lcb.GetIntgPd())
	fmt.Fprintf(output, "%s ", boolToStr(lcb.IsLogEna(), "1", "0"))
	fmt.Fprintf(output, "%s);\n", boolToStr(lcb.IsReasonCode(), "1", "0"))
}

func printRCBInstance(output io.Writer, rcb *model.ReportControlBlock, index string, hasOwner bool) {
	fmt.Fprintf(output, "RC(%s%s ", rcb.GetName(), index)

	if rcb.GetRptID() != "" {
		fmt.Fprintf(output, "%s ", rcb.GetRptID())
	} else {
		fmt.Fprint(output, "- ")
	}

	fmt.Fprintf(output, "%s ", boolToStr(rcb.IsBuffered(), "1", "0"))

	if rcb.GetDataSet() != "" {
		fmt.Fprintf(output, "%s ", rcb.GetDataSet())
	} else {
		fmt.Fprint(output, "- ")
	}

	fmt.Fprintf(output, "%v ", rcb.GetConfRef())

	trigOps := rcb.GetTriggerOptions().GetIntValue()
	if hasOwner {
		trigOps += 64
	}
	fmt.Fprintf(output, "%d ", trigOps)
	fmt.Fprintf(output, "%d ", rcb.GetOptionFields().GetIntValue())
	fmt.Fprintf(output, "%d ", rcb.GetBufferTime())

	if rcb.GetIntegrityPeriod() != nil {
		fmt.Fprintf(output, "%d", *rcb.GetIntegrityPeriod())
	} else {
		fmt.Fprint(output, "0")
	}
	fmt.Fprintln(output, ");")
}

func exportDataSet(output io.Writer, ds *model.DataSet, ln *model.LogicalNode) {
	fmt.Fprintf(output, "DS(%s){\n", ds.GetName())
	for _, fcda := range ds.GetFCDA() {
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

		// cross-LD ref
		if fcda.GetLdInstance() != "" && fcda.GetLdInstance() != ln.GetParentLogicalDevice().GetInst() {
			variableName = fcda.GetLdInstance() + "/" + variableName
		}

		if variableName != "" && arrayIndex != -1 && componentName != "" {
			fmt.Fprintf(output, "DE(%s %d %s);\n", variableName, arrayIndex, componentName)
		} else if variableName != "" && arrayIndex != -1 {
			fmt.Fprintf(output, "DE(%s %d);\n", variableName, arrayIndex)
		} else if variableName != "" {
			fmt.Fprintf(output, "DE(%s);\n", variableName)
		}
	}
	fmt.Fprintln(output, "}")
}

func exportDataObject(output io.Writer, do *model.DataObject, isTransient bool) {
	if do.IsTransient() {
		isTransient = true
	}

	if do.GetCount() > 0 {
		for i := 0; i < do.GetCount(); i++ {
			fmt.Fprintf(output, "[%d]{\n", i)
			exportDataObjectChild(output, do, isTransient)
			fmt.Fprint(output, "}\n")
		}
	} else {
		exportDataObjectChild(output, do, isTransient)
	}
}

func exportDataObjectChild(output io.Writer, do *model.DataObject, isTransient bool) {
	for _, sdo := range do.GetSubDataObjects() {
		fmt.Fprintf(output, "DO(%s %d){\n", sdo.GetName(), sdo.GetCount())
		exportDataObject(output, sdo, isTransient)
		fmt.Fprintln(output, "}")
	}
	for _, da := range do.GetDataAttributes() {
		exportDataAttribute(output, da, isTransient)
	}
}

func exportDataAttribute(output io.Writer, da *model.DataAttribute, isTransient bool) {
	fmt.Fprintf(output, "DA(%s %d %d %d ",
		da.GetName(), da.GetCount(),
		da.GetType().GetIntValue(),
		da.GetFC().GetIntValue())

	trgOpsVal := da.GetTriggerOptions().GetIntValue()
	if isTransient {
		trgOpsVal += 128
	}
	fmt.Fprintf(output, "%d ", trgOpsVal)

	sAddr := int64(0)
	if da.GetShortAddress() != "" {
		if v, err := fmt.Sscanf(da.GetShortAddress(), "%d", &sAddr); v != 1 || err != nil {
			fmt.Printf("WARNING: short address %q is not valid for libIEC61850!\n", da.GetShortAddress())
			sAddr = 0
		}
	}
	fmt.Fprintf(output, "%d)", sAddr)

	if da.GetCount() > 0 {
		fmt.Fprint(output, "{\n")
		for i := 0; i < da.GetCount(); i++ {
			fmt.Fprintf(output, "[%d]", i)
			printDataAttributeValue(output, da, isTransient)
		}
		fmt.Fprint(output, "}\n")
	} else {
		printDataAttributeValue(output, da, isTransient)
	}
}

func printDataAttributeValue(output io.Writer, da *model.DataAttribute, isTransient bool) {
	if da.IsBasicAttribute() {
		value := da.GetEffectiveValue(nil)

		if value != nil {
			switch da.GetType() {
			case model.AT_ENUMERATED, model.AT_INT8, model.AT_INT16, model.AT_INT32, model.AT_INT64:
				fmt.Fprintf(output, "=%d", value.GetIntValue())
			case model.AT_INT8U, model.AT_INT16U, model.AT_INT24U, model.AT_INT32U:
				fmt.Fprintf(output, "=%d", value.GetLongValue())
			case model.AT_BOOLEAN:
				if b, ok := value.GetValue().(bool); ok && b {
					fmt.Fprint(output, "=1")
				}
			case model.AT_UNICODE_STRING_255, model.AT_CURRENCY,
				model.AT_VISIBLE_STRING_32, model.AT_VISIBLE_STRING_64,
				model.AT_VISIBLE_STRING_129, model.AT_VISIBLE_STRING_255, model.AT_VISIBLE_STRING_65:
				fmt.Fprintf(output, "=\"%v\"", value.GetValue())
			case model.AT_FLOAT32, model.AT_FLOAT64:
				fmt.Fprintf(output, "=%v", value.GetValue())
			case model.AT_TIMESTAMP, model.AT_ENTRY_TIME:
				fmt.Fprintf(output, "=%d", value.GetLongValue())
			default:
				fmt.Printf("Unknown default value for %s type: %v\n", da.GetName(), da.GetType())
			}
		}
		fmt.Fprintln(output, ";")
	} else {
		fmt.Fprintln(output, "{")
		for _, sub := range da.GetSubDataAttributes() {
			exportDataAttribute(output, sub, isTransient)
		}
		fmt.Fprintln(output, "}")
	}
}
