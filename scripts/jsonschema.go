package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/floffah/culprit/internal/recipe"
	"github.com/invopop/jsonschema"
)

//go:generate go run jsonschema.go

func main() {
	r := &jsonschema.Reflector{
		FieldNameTag: "toml",
	}
	s := r.ReflectFromType(reflect.TypeOf(recipe.Recipe{}))

	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		panic(err.Error())
	}

	outFile, err := os.OpenFile("../recipe.schema.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer outFile.Close()

	_, err = outFile.Write(data)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Generated recipe.schema.json")
}
