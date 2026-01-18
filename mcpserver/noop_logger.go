/******************************************************************************
 * Copyright (c) 2025 Tenebris Technologies Inc.                              *
 * Please see LICENSE file for details.                                       *
 ******************************************************************************/

package mcpserver

// noopLogger implements the Logger interface but does nothing
// Used when no logger is provided to mcpserver
type noopLogger struct{}

func (n *noopLogger) Debug(string)              {}
func (n *noopLogger) Info(string)               {}
func (n *noopLogger) Notice(string)             {}
func (n *noopLogger) Warning(string)            {}
func (n *noopLogger) Error(string)              {}
func (n *noopLogger) Fatal(string)              {}
func (n *noopLogger) Debugf(string, ...any)     {}
func (n *noopLogger) Infof(string, ...any)      {}
func (n *noopLogger) Noticef(string, ...any)    {}
func (n *noopLogger) Warningf(string, ...any)   {}
func (n *noopLogger) Errorf(string, ...any)     {}
func (n *noopLogger) Fatalf(string, ...any)     {}
func (n *noopLogger) Close()                    {}
