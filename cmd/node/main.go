package main

import (
	"os"
	"path/filepath"

	pkg "github.com/PlayerR9/grammar/cmd/node/pkg"
	ggen "github.com/PlayerR9/lib_units/generator"
)

func main() {
	node_name, err := pkg.ParseFlags()
	if err != nil {
		ggen.PrintFlags()

		pkg.Logger.Fatalf("Failed to parse flags: %s", err.Error())
	}

	type_name, err := ggen.TypeListFlag.Type(0)
	if err != nil {
		pkg.Logger.Fatalf("Failed to get type: %s", err.Error())
	}

	type_name, err = ggen.FixVariableName(type_name, nil, ggen.Exported)
	if err != nil {
		pkg.Logger.Fatalf("Failed to fix variable name: %s", err.Error())
	}

	data := &pkg.GenData{
		NodeName: node_name,
		TypeName: type_name,
	}

	dest, err := pkg.Generator.Generate(type_name, "_node.go", data)
	if err != nil {
		pkg.Logger.Fatalf("Failed to generate: %s", err.Error())
	}

	dir := filepath.Dir(dest.DestLoc)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		pkg.Logger.Fatalf("Failed to create directory: %s", err.Error())
	}

	err = os.WriteFile(dest.DestLoc, dest.Data, 0644)
	if err != nil {
		pkg.Logger.Fatalf("Failed to write file: %s", err.Error())
	}

	pkg.Logger.Printf("Successfully generated: %q", dest.DestLoc)
}
