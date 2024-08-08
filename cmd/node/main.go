package main

import (
	ggen "github.com/PlayerR9/go-generator/generator"
	pkg "github.com/PlayerR9/grammar/cmd/node/pkg"
)

func main() {
	type_name, node_name, err := pkg.ParseFlags()
	if err != nil {
		ggen.PrintFlags()

		pkg.Logger.Fatalf("Failed to parse flags: %s", err.Error())
	}

	data := &pkg.GenData{
		NodeName: node_name,
		TypeName: type_name,
	}

	res, err := pkg.Generator.Generate(pkg.OutputLocFlag, type_name+"_node.go", data)
	if err != nil {
		pkg.Logger.Fatalf("Failed to generate: %s", err.Error())
	}

	dest, err := res.WriteFile("")
	if err != nil {
		pkg.Logger.Fatal(err.Error())
	}

	pkg.Logger.Printf("Successfully generated: %q", dest)
}
