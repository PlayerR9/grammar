package main

import (
	pkg "github.com/PlayerR9/grammar/cmd/token/pkg"
	luc "github.com/PlayerR9/lib_units/common"
	ggen "github.com/PlayerR9/lib_units/generator"
)

func main() {
	err := ggen.ParseFlags()
	if err != nil {
		pkg.Logger.Fatalf("Failed to parse flags: %s", err.Error())
	}

	name := luc.AssertDerefNil(pkg.NameFlag, "pkg.NameFlag")

	name, err = ggen.FixVariableName(name, nil, ggen.Exported)
	if err != nil {
		pkg.Logger.Fatalf("Failed to fix variable name: %s", err.Error())
	}

	g := &pkg.Gen{
		TypeName: name,
	}

	dest, err := pkg.Generator.Generate(name, "_token.go", g)
	if err != nil {
		pkg.Logger.Fatalf("Failed to generate: %s", err.Error())
	}

	pkg.Logger.Printf("Successfully generated: %q", dest)
}
