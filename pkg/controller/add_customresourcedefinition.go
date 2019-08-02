package controller

import (
	"github.com/jiuchen1986/crd-controller/pkg/controller/customresourcedefinition"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, customresourcedefinition.Add)
}
