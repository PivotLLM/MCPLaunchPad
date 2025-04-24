// Copyright (c) 2025 Tenebris Technologies Inc.
// Please see LICENSE for details.

package mcpserver

import "encoding/json"

// logInJSON logs data in JSON for debugging
func (s *MCPServer) logInJSON(data any) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err == nil {
		s.logger.Debugf("Failed to marshal type %T to JSON: %v", data, data)
		return
	}
	s.logger.Debugf("JSON DATA:\n%s", string(b))
}
