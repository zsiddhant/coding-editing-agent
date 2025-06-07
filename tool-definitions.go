package main

import (
	"encoding/json"
	"os"

	"github.com/invopop/jsonschema"
)

var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

type ReadFileInput struct {
	Path     string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
	Dollar   string `json:"$" jsonschema_description:"The relative path of a file in the working directory (alternative key)."`
	Filename string `json:"filename" jsonschema_description:"The relative path of a file in the working directory (alternative key: filename)."`
	File     string `json:"file" jsonschema_description:"The relative path of a file in the working directory (alternative key: file)."`
}

var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

func ReadFile(input json.RawMessage) (string, error) {
	//fmt.Println("Reading file...", input)
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		panic(err)
	}

	//fmt.Println("Reading file:", readFileInput.Path+" (alternative key: "+readFileInput.Dollar+")")

	if readFileInput.Dollar != "" {
		readFileInput.Path = readFileInput.Dollar
	} else if readFileInput.Filename != "" {
		readFileInput.Path = readFileInput.Filename
	} else if readFileInput.File != "" {
		readFileInput.Path = readFileInput.File
	}

	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GenerateSchema[T any]() ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	return ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}
