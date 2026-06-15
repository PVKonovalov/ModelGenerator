/*
 *  main.go
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

package main

import (
	"ModelGenerator/tools"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "genmodel":
		runGenModel(os.Args[2:])
	case "genconfig":
		runGenConfig(os.Args[2:])
	case "gencode":
		runGenCode(os.Args[2:])
	case "viewmodel":
		runViewModel(os.Args[2:])
	case "gengolibmodel":
		runGenGoLibModel(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  genmodel      <ICD file> [-ied <name>] [-ap <name>] [-out <output>] [-modelprefix <prefix>] [-initializeonce]")
	fmt.Println("  genconfig     <ICD file> [-ied <name>] [-ap <name>] [<output file>]")
	fmt.Println("  gencode       <ICD file> [<output file>]")
	fmt.Println("  viewmodel     <ICD file> [-ied <name>] [-ap <name>] [-t] [-u] [-s] [-a] [<output file>]")
	fmt.Println("    -t  print type list")
	fmt.Println("    -u  print unused types")
	fmt.Println("    -s  print model structure")
	fmt.Println("    -a  print data attribute list")
	fmt.Println("  gengolibmodel <ICD file> [-ied <name>] [-ap <name>] [-pkg <package>] [<output file>]")
	fmt.Println("    Generates a Go source file with BuildModel() for github.com/PVKonovalov/libiec61850-Go")
}

func runGenModel(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: genmodel <ICD file> [-ied <name>] [-ap <name>] [-out <output>] [-modelprefix <prefix>] [-initializeonce]")
		os.Exit(1)
	}

	icdFile := args[0]
	outputFileName := "static_model"
	iedName := ""
	accessPointName := ""
	modelPrefix := "iedModel"
	initializeOnce := false

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-ied":
			i++
			iedName = args[i]
			fmt.Println("Select IED", iedName)
		case "-ap":
			i++
			accessPointName = args[i]
			fmt.Println("Select access point", accessPointName)
		case "-out":
			i++
			outputFileName = args[i]
			fmt.Println("Select Output File", outputFileName)
		case "-modelprefix":
			i++
			modelPrefix = args[i]
			fmt.Println("Select Model Prefix", modelPrefix)
		case "-initializeonce":
			initializeOnce = true
			fmt.Println("Select Initialize Once")
		default:
			fmt.Fprintf(os.Stderr, "Unknown option: %q\n", args[i])
		}
	}

	fmt.Println("Select ICD File", icdFile)

	stream, err := os.Open(icdFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot open %s: %v\n", icdFile, err)
		os.Exit(1)
	}
	defer stream.Close()

	cFile, err := os.Create(outputFileName + ".c")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer cFile.Close()

	hFile, err := os.Create(outputFileName + ".h")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer hFile.Close()

	err = tools.NewStaticModelGenerator(stream, icdFile, cFile, hFile, outputFileName, iedName, accessPointName, modelPrefix, initializeOnce)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func runGenConfig(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: genconfig <ICD file> [-ied <name>] [-ap <name>] [<output file>]")
		os.Exit(1)
	}

	icdFile := args[0]
	iedName := ""
	accessPointName := ""
	output := os.Stdout

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-ied":
			i++
			iedName = args[i]
			fmt.Println("Select IED", iedName)
		case "-ap":
			i++
			accessPointName = args[i]
			fmt.Println("Select access point", accessPointName)
		default:
			f, err := os.Create(args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()
			output = f
		}
	}

	stream, err := os.Open(icdFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer stream.Close()

	err = tools.RunDynamicModelGenerator(stream, icdFile, output, iedName, accessPointName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func runViewModel(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: viewmodel <ICD file> [-ied <name>] [-ap <name>] [-t] [-u] [-s] [-a] [<output file>]")
		os.Exit(1)
	}

	icdFile := args[0]
	iedName := ""
	accessPointName := ""
	printTypeList := false
	printUnusedTypes := false
	printModelStructure := false
	printDataAttributes := false
	output := os.Stdout

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-ied":
			i++
			iedName = args[i]
			fmt.Println("Select IED", iedName)
		case "-ap":
			i++
			accessPointName = args[i]
			fmt.Println("Select access point", accessPointName)
		case "-t":
			printTypeList = true
		case "-u":
			printUnusedTypes = true
		case "-s":
			printModelStructure = true
		case "-a":
			printDataAttributes = true
		default:
			if !startsWithDash(args[i]) {
				f, err := os.Create(args[i])
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
					os.Exit(1)
				}
				defer f.Close()
				output = f
			} else {
				fmt.Fprintf(os.Stderr, "Unknown option: %q\n", args[i])
			}
		}
	}

	stream, err := os.Open(icdFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer stream.Close()

	err = tools.RunModelViewer(stream, output, iedName, accessPointName,
		printTypeList, printUnusedTypes, printModelStructure, printDataAttributes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func runGenGoLibModel(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gengolibmodel <ICD file> [-ied <name>] [-ap <name>] [-pkg <package>] [<output file>]")
		os.Exit(1)
	}

	icdFile := args[0]
	iedName := ""
	accessPointName := ""
	packageName := "model"
	output := os.Stdout

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-ied":
			i++
			iedName = args[i]
		case "-ap":
			i++
			accessPointName = args[i]
		case "-pkg":
			i++
			packageName = args[i]
		default:
			if !startsWithDash(args[i]) {
				f, err := os.Create(args[i])
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
					os.Exit(1)
				}
				defer f.Close()
				output = f
			} else {
				fmt.Fprintf(os.Stderr, "Unknown option: %q\n", args[i])
			}
		}
	}

	stream, err := os.Open(icdFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot open %s: %v\n", icdFile, err)
		os.Exit(1)
	}
	defer stream.Close()

	err = tools.RunGoModelGenerator(stream, output, packageName, iedName, accessPointName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func startsWithDash(s string) bool {
	return len(s) > 0 && s[0] == '-'
}

func runGenCode(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gencode <ICD file> [<output file>]")
		os.Exit(1)
	}

	icdFile := args[0]
	output := os.Stdout

	if len(args) > 1 {
		f, err := os.Create(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		output = f
	}

	stream, err := os.Open(icdFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer stream.Close()

	err = tools.RunDynamicCodeGenerator(stream, output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
