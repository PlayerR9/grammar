package grammar

import _ "github.com/PlayerR9/tree"

//go:generate go run github.com/PlayerR9/tree/cmd -name=Node -fields=Type/N,Data/string -g=N/NodeTyper -o=ast/generic_node.go
