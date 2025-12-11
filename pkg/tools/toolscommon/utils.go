package toolscommon

import (
	"encoding/json"
	"github.com/invopop/jsonschema"
)

func TypeToJsonSchema[T any]() string {
	var zero T

	reflector := jsonschema.Reflector{
		DoNotReference:            true, // Removes $defs map, outputs entire structure inline
		Anonymous:                 true, // Hides auto-generated Schema IDs
		AllowAdditionalProperties: true, // Removes additionalProperties: false
	}
	schema := reflector.Reflect(zero)
	schema.Version = ""

	//data, err := json.MarshalIndent(schema, "", "")
	data, err := json.Marshal(schema)
	if err != nil {
		return ""
	}

	return string(data)
}
