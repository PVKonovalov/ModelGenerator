/*
 *  dynamic_code_generator.go
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
	"ModelGenerator/scl"
	"ModelGenerator/scl/types"
	"fmt"
	"io"
	"strings"
)

// RunDynamicCodeGenerator generates C code stubs for dynamic model creation.
func RunDynamicCodeGenerator(r io.Reader, output io.Writer) error {
	sclParser, err := parser.NewSclParser(r, false)
	if err != nil {
		return err
	}

	declarations := sclParser.GetTypeDeclarations()
	createDynamicCode(declarations, output)
	return nil
}

func replaceInvalidCharacters(s string) string {
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, "+", "_")
	return s
}

func createDynamicCode(declarations *types.TypeDeclarations, out io.Writer) {
	var doTypeDefs []*types.DataObjectType
	var lnTypeDefs []*types.LogicalNodeType
	var daTypeDefs []*types.DataAttributeType

	seen := map[string]bool{}

	for _, t := range declarations.GetTypeDeclarations() {
		id := t.GetID()
		if seen[id] {
			continue
		}
		seen[id] = true
		switch v := t.(type) {
		case *types.LogicalNodeType:
			lnTypeDefs = append(lnTypeDefs, v)
		case *types.DataObjectType:
			doTypeDefs = append(doTypeDefs, v)
		case *types.DataAttributeType:
			daTypeDefs = append(daTypeDefs, v)
		}
	}

	// Function prototypes
	var prototypes []string
	for _, lnt := range lnTypeDefs {
		proto := fmt.Sprintf("LogicalNode*\nLN_%s_createInstance(char* lnName, LogicalDevice* parent);",
			replaceInvalidCharacters(lnt.GetID()))
		prototypes = append(prototypes, proto)
	}
	for _, dot := range doTypeDefs {
		proto := fmt.Sprintf("DataObject*\nDO_%s_createInstance(char* doName, ModelNode* parent, int arrayCount);",
			replaceInvalidCharacters(dot.GetID()))
		prototypes = append(prototypes, proto)
	}
	for _, dat := range daTypeDefs {
		proto := fmt.Sprintf("DataAttribute*\nDA_%s_createInstance(char* daName, ModelNode* parent, FunctionalConstraint fc, uint8_t triggerOptions);",
			replaceInvalidCharacters(dat.GetID()))
		prototypes = append(prototypes, proto)
	}

	fmt.Fprintln(out, "#include \"iec61850_server.h\"")
	fmt.Fprintln(out, "")
	for _, p := range prototypes {
		fmt.Fprintln(out, p)
		fmt.Fprintln(out, "")
	}

	for _, lnt := range lnTypeDefs {
		fmt.Fprintln(out, "/**")
		fmt.Fprintf(out, " * LN: %s ", replaceInvalidCharacters(lnt.GetID()))
		if lnt.GetDesc() != "" {
			fmt.Fprintf(out, "(%s)", lnt.GetDesc())
		}
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, " */")
		fmt.Fprintln(out, "LogicalNode*")
		fmt.Fprintf(out, "LN_%s_createInstance(char* lnName, LogicalDevice* parent)\n", replaceInvalidCharacters(lnt.GetID()))
		fmt.Fprintln(out, "{")
		fmt.Fprint(out, "    LogicalNode* newLn = LogicalNode_create(lnName, parent);\n\n")

		for _, doDef := range lnt.GetDataObjectDefinitions() {
			fmt.Fprintf(out, "    DO_%s_createInstance(\"%s\", (ModelNode*) newLn, %d);\n",
				replaceInvalidCharacters(doDef.Type), doDef.Name, doDef.Count)
		}

		fmt.Fprint(out, "\n    return newLn;\n")
		fmt.Fprint(out, "}\n\n\n")
	}

	for _, dot := range doTypeDefs {
		fmt.Fprintln(out, "/**")
		fmt.Fprintf(out, " * DO: %s ", replaceInvalidCharacters(dot.GetID()))
		if dot.GetDesc() != "" {
			fmt.Fprintf(out, "(%s)", dot.GetDesc())
		}
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, " */")
		fmt.Fprintln(out, "DataObject*")
		fmt.Fprintf(out, "DO_%s_createInstance(char* doName, ModelNode* parent, int arrayCount)\n", replaceInvalidCharacters(dot.GetID()))
		fmt.Fprintln(out, "{")
		fmt.Fprint(out, "    DataObject* newDo = DataObject_create(doName, parent, arrayCount);\n\n")

		for _, dad := range dot.GetDataAttributes() {
			trgOpts := newTriggerOptionsForDA(dad)
			if dad.IsConstructed() {
				fmt.Fprintf(out, "    DA_%s_createInstance(\"%s\", (ModelNode*) newDo, IEC61850_FC_%s, %d);\n",
					replaceInvalidCharacters(dad.GetType()), dad.GetName(), dad.GetFCString(), trgOpts)
			} else {
				fmt.Fprintf(out, "    DataAttribute_create(\"%s\", (ModelNode*) newDo, IEC61850_%s, IEC61850_FC_%s, %d, %d, 0);\n",
					dad.GetName(), dad.GetBType(), dad.GetFCString(), trgOpts, dad.GetCount())
			}
		}

		for _, dod := range dot.GetSubDataObjects() {
			fmt.Fprintf(out, "    DO_%s_createInstance(\"%s\", (ModelNode*) newDo, %d);\n",
				replaceInvalidCharacters(dod.Type), dod.Name, dod.Count)
		}

		fmt.Fprint(out, "\n    return newDo;\n")
		fmt.Fprint(out, "}\n\n\n")
	}

	for _, dat := range daTypeDefs {
		fmt.Fprintln(out, "/**")
		fmt.Fprintf(out, " * DA: %s ", replaceInvalidCharacters(dat.GetID()))
		if dat.GetDesc() != "" {
			fmt.Fprintf(out, "(%s)", dat.GetDesc())
		}
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, " */")
		fmt.Fprintln(out, "DataAttribute*")
		fmt.Fprintf(out, "DA_%s_createInstance(char* daName, ModelNode* parent, FunctionalConstraint fc, uint8_t triggerOptions)\n",
			replaceInvalidCharacters(dat.GetID()))
		fmt.Fprintln(out, "{")
		fmt.Fprint(out, "    DataAttribute* newDa = DataAttribute_create(daName, parent, IEC61850_CONSTRUCTED, fc, triggerOptions, 0, 0);\n\n")

		for _, dad := range dat.GetSubDataAttributes() {
			if dad.IsConstructed() {
				fmt.Fprintf(out, "    DA_%s_createInstance(\"%s\", (ModelNode*) newDa, fc, triggerOptions);\n",
					replaceInvalidCharacters(dad.GetType()), dad.GetName())
			} else {
				fmt.Fprintf(out, "    DataAttribute_create(\"%s\", (ModelNode*) newDa, IEC61850_%s, fc, triggerOptions, %d, 0);\n",
					dad.GetName(), dad.GetBType(), dad.GetCount())
			}
		}

		fmt.Fprint(out, "\n    return newDa;\n")
		fmt.Fprint(out, "}\n\n\n")
	}
}

// newTriggerOptionsForDA creates a trigger options bitmask from a DA definition.
func newTriggerOptionsForDA(dad *scl.DataAttributeDefinition) int {
	val := 0
	if dad.Dchg {
		val += 1
	}
	if dad.Qchg {
		val += 2
	}
	if dad.Dupd {
		val += 4
	}
	return val
}
