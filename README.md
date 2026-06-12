# ModelGenerator

IEC 61850 SCL model code generator - a Go port of the Java-based `model_generator` tool from [libIEC61850](https://github.com/mz-automation/libiec61850).

## Overview

ModelGenerator reads IEC 61850 Substation Configuration Language (SCL) files (`.icd`, `.cid`, `.scd`) and generates C source code or configuration files for use with the libIEC61850 C library.

It provides four commands, matching the original Java toolset:

| Command | Original Java class | Output |
|---|---|---|
| `genmodel` | `StaticModelGenerator` | `.c` / `.h` static model files |
| `genconfig` | `DynamicModelGenerator` | Text model config file |
| `gencode` | `DynamicCodeGenerator` | C stubs for dynamic model creation |
| `viewmodel` | `ModelViewer` | Human-readable model inspection |

## Building

Requires Go 1.21 or later.

```bash
go build -o modelgenerator .
```

## Usage

### genmodel - Static C model

Generates a `.c` and `.h` file pair with pre-initialized C structs for direct linking into a libIEC61850 application.

```bash
modelgenerator genmodel <ICD file> [options]

Options:
  -ied <name>           Select IED by name (default: first IED with a Server)
  -ap <name>            Select AccessPoint by name (default: first)
  -out <filename>       Output filename without extension (default: static_model)
  -modelprefix <prefix> C symbol prefix (default: iedModel)
  -initializeonce       Emit IedModel_init() instead of static initializer
```

Example:

```bash
modelgenerator genmodel device.icd -out static_model
# produces static_model.c and static_model.h
```

### genconfig - Dynamic text model

Generates a text-format model configuration file consumed by libIEC61850's dynamic model loader.

```bash
modelgenerator genconfig <ICD file> [options] [output file]

Options:
  -ied <name>   Select IED by name
  -ap <name>    Select AccessPoint by name
```

Example:

```bash
modelgenerator genconfig device.icd model.cfg
```

### gencode - Dynamic C stubs

Generates C helper functions (`LN_*_createInstance`, `DO_*_createInstance`, `DA_*_createInstance`) for building a model at runtime.

```bash
modelgenerator gencode <ICD file> [output file]
```

Example:

```bash
modelgenerator gencode device.icd model_stubs.c
```

### viewmodel - Model inspection

Prints a human-readable view of the IEC 61850 data model parsed from an SCL file.

```bash
modelgenerator viewmodel <ICD file> [options] [output file]

Options:
  -ied <name>   Select IED by name
  -ap <name>    Select AccessPoint by name
  -t            Print type list (all SCL types)
  -u            Print unused types only
  -s            Print model structure (LD / LN / DO / DA tree)
  -a            Print flat data attribute list with functional constraints
```

Example:

```bash
modelgenerator viewmodel device.icd -a
```

## Project Structure

```
.
├── main.go                         # CLI entry point, command dispatch
├── parser/
│   └── parser.go                   # SclParser - top-level SCL file parser
├── scl/
│   ├── xml_node.go                 # XML DOM representation
│   ├── parser_utils.go             # Attribute parsing helpers
│   ├── data_attribute_definition.go
│   ├── data_object_definition.go
│   ├── communication/              # SCL Communication section
│   │   ├── connected_ap.go
│   │   ├── gse.go
│   │   ├── smv.go
│   │   ├── phy_com_address.go
│   │   └── ...
│   ├── model/                      # SCL IED / data model
│   │   ├── ied.go
│   │   ├── logical_device.go
│   │   ├── logical_node.go
│   │   ├── data_object.go
│   │   ├── data_attribute.go
│   │   ├── data_set.go
│   │   ├── report_control_block.go
│   │   ├── log_control.go
│   │   └── ...
│   └── types/                      # SCL type declarations (LNType, DOType, DAType, EnumType)
│       ├── type_declarations.go
│       ├── logical_node_type.go
│       ├── data_object_type.go
│       ├── data_attribute_type.go
│       └── enumeration_type.go
└── tools/
    ├── static_model_generator.go   # genmodel implementation
    ├── dynamic_model_generator.go  # genconfig implementation
    ├── dynamic_code_generator.go   # gencode implementation
    └── model_viewer.go             # viewmodel implementation
```

## Go Port Notes

This tool is a port of the Java `model_generator` utility distributed with [libIEC61850](https://github.com/mz-automation/libiec61850), originally written by Michael Zillgith.

### Key differences from the Java version

**Attribute parsing semantics**

Java's `ParserUtils.parseAttribute()` returns `null` when an XML attribute is absent. The Go equivalent `scl.ParseAttribute()` returns `""` (empty string) for absent attributes. Throughout the codebase, null checks (`!= null`) are replaced with empty-string checks (`!= ""`).

The one place where this distinction matters is `LN0` (Logical Node Zero), which legitimately carries `inst=""` in IEC 61850. The Java guard `if (inst == null)` allows the empty string through; the Go port explicitly exempts `lnClass == "LLN0"` from the non-empty `inst` requirement.

Similarly, `LogControl.logName` may be present in SCL with an empty value - Java nullifies it and continues; Go uses `ParseAttributePtr` to distinguish truly absent (error) from present-with-empty (treated as no log name).

**No import cycles**

The Java codebase has no module boundary between the parser and the model. In Go, `scl/model` imports `scl`, so a parser living inside `scl` cannot also import `scl/model`. The top-level SCL parser is therefore placed in its own package (`parser/`).

**Nullable integers**

Java uses boxed `Integer` (nullable) for optional integer fields such as `ReportControlBlock.intgPd`. Go uses `*int` (pointer to int) for the same semantics - `nil` means the attribute was absent.

**GSE / SMV min/max time**

Java uses `-1` as the sentinel for "not set" on `GSE.minTime` / `GSE.maxTime`. The Go port preserves this convention with `int` fields defaulting to `-1`.

## License

Copyright 2014-2024 Michael Zillgith  
Copyright 2026 Pavel Konovalov - Golang port

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

See [COPYING](https://github.com/mz-automation/libiec61850/blob/v1.5/COPYING) for the full license text.
