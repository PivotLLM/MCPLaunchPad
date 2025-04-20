package gavin

import (
	"fmt"

	"github.com/PivotLLM/MCPLaunchPad/global"
)

// ValidateAndBuildQueryParams validates the options and builds query parameters.
func (c *Config) ValidateAndBuildQueryParams(toolName string, options map[string]any) (map[string]string, error) {

	// Find the tool definition from the registration
	var toolDef *global.ToolDefinition
	for _, def := range c.Register() {
		if def.Name == toolName {
			toolDef = &def
			break
		}
	}

	if toolDef == nil {
		return nil, fmt.Errorf("tool '%s' not found in registration", toolName)
	}

	// Validate and build query parameters
	queryParams := make(map[string]string)
	for _, param := range toolDef.Parameters {
		value, exists := options[param.Name]
		if !exists {
			if param.Required {
				return nil, fmt.Errorf("missing required parameter: %s", param.Name)
			}
			continue
		}

		// Convert the value to a string or handle numbers
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int, int8, int16, int32, int64, float32, float64:
			strValue = fmt.Sprintf("%v", v)
		default:
			return nil, fmt.Errorf("parameter '%s' must be a string or a number", param.Name)
		}

		queryParams[param.Name] = strValue
	}

	return queryParams, nil
}
