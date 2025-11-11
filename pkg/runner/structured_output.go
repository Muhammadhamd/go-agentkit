package runner

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// parseStructuredOutput parses and validates structured output from model response
// Uses the existing prepareOutputSchema infrastructure
func (r *Runner) parseStructuredOutput(
	content string,
	outputType reflect.Type,
) (interface{}, error) {
	if outputType == nil {
		return content, nil
	}

	// Parse JSON from content
	var jsonData interface{}
	if err := json.Unmarshal([]byte(content), &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	// Get the schema for validation
	schema := r.prepareOutputSchema(outputType)
	if schema == nil {
		// No schema, just return parsed JSON
		return jsonData, nil
	}

	// Validate against schema (basic validation - check required fields)
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if err := r.validateAgainstSchema(jsonData, schemaMap); err != nil {
			return nil, fmt.Errorf("output validation failed: %w", err)
		}
	}

	// Convert JSON to Go struct
	resultValue := reflect.New(outputType)
	if err := r.unmarshalToStruct(jsonData, resultValue.Interface()); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to struct: %w", err)
	}

	return resultValue.Elem().Interface(), nil
}

// validateAgainstSchema performs basic schema validation
func (r *Runner) validateAgainstSchema(data interface{}, schema map[string]interface{}) error {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected object, got %T", data)
	}

	// Check required fields
	if required, ok := schema["required"].([]string); ok {
		for _, field := range required {
			if _, exists := dataMap[field]; !exists {
				return fmt.Errorf("missing required field: %s", field)
			}
		}
	}

	// Check properties
	if properties, ok := schema["properties"].(map[string]interface{}); ok {
		for fieldName, fieldSchema := range properties {
			if fieldValue, exists := dataMap[fieldName]; exists {
				if fieldSchemaMap, ok := fieldSchema.(map[string]interface{}); ok {
					if err := r.validateFieldType(fieldValue, fieldSchemaMap); err != nil {
						return fmt.Errorf("field %s: %w", fieldName, err)
					}
				}
			}
		}
	}

	return nil
}

// validateFieldType validates a field value against its schema type
func (r *Runner) validateFieldType(value interface{}, schema map[string]interface{}) error {
	schemaType, ok := schema["type"].(string)
	if !ok {
		return nil // No type specified, skip validation
	}

	switch schemaType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case "integer":
		// JSON numbers can be float64, so check if it's a whole number
		switch v := value.(type) {
		case float64:
			if v != float64(int64(v)) {
				return fmt.Errorf("expected integer, got float")
			}
		case int, int32, int64:
			// OK
		default:
			return fmt.Errorf("expected integer, got %T", value)
		}
	case "number":
		switch value.(type) {
		case float64, int, int32, int64:
			// OK
		default:
			return fmt.Errorf("expected number, got %T", value)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("expected array, got %T", value)
		}
	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
	}

	return nil
}

// unmarshalToStruct unmarshals JSON data to a Go struct
func (r *Runner) unmarshalToStruct(data interface{}, target interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, target)
}
