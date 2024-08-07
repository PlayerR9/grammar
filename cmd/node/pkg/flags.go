package pkg

import (
	"flag"
	"fmt"

	ggen "github.com/PlayerR9/go-generator/generator"
)

var (
	TypeNameFlag *string
)

func init() {
	TypeNameFlag = flag.String("name", "", "The name of the node. This flag is required.")

	ggen.SetTypeListFlag("type", true, 1, "The type of the node to generate.")
	ggen.SetOutputFlag("<type>_node.go", true)
	ggen.SetGenericsSignFlag("g", false, 1)
}

func ParseFlags() (string, string, error) {
	err := ggen.ParseFlags()
	if err != nil {
		return "", "", err
	}

	if *TypeNameFlag == "" {
		return "", "", fmt.Errorf("type flag is required")
	}

	node_name, err := ggen.FixVariableName(*TypeNameFlag, nil, ggen.Exported)
	if err != nil {
		return "", "", fmt.Errorf("invalid type name: %w", err)
	}

	type_name, err := ggen.TypeListFlag.Type(0)
	if err != nil {
		return "", "", fmt.Errorf("invalid type name: %w", err)
	}

	type_name, err = ggen.FixVariableName(type_name, nil, ggen.Exported)
	if err != nil {
		return "", "", fmt.Errorf("invalid type name: %w", err)
	}

	return type_name, node_name, nil
}
