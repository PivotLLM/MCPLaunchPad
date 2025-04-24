// Copyright (c) 2025 Tenebris Technologies Inc.
// This software is licensed under the MIT License (see LICENSE for details).

// Package mlogger provides a simple file-based logger with optional debug message
// suppression and logging to stdout.
package mlogger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PivotLLM/MCPLaunchPad/global"
)

type MLogger struct {
	fileHandle *os.File
	logfile    string
	logStdout  bool
	debug      bool
	logLevel   bool
	prefix     string
	dateFormat string
}

// This package implements interfaces.Logger
var _ global.Logger = (*MLogger)(nil)

// Option is a function that configures a MLogger
type Option func(*MLogger) error

// New creates a new instance of MLogger with the provided options
func New(options ...Option) (global.Logger, error) {
	m := &MLogger{
		logLevel:   true,
		dateFormat: "2006-01-02 15:04:05",
	}

	for _, option := range options {
		if err := option(m); err != nil {
			return nil, err
		}
	}

	// Call the OS-specific constructor
	return m.open()
}

// WithPrefix sets a process name or similar short identifier
//
//goland:noinspection GoUnusedExportedFunction
func WithPrefix(prefix string) Option {
	return func(u *MLogger) error {
		if prefix == "" {
			u.prefix = ""
		} else {
			u.prefix = " " + strings.TrimSpace(prefix)
		}
		return nil
	}
}

// WithDateFormat sets the date format for the MLogger
//
//goland:noinspection GoUnusedExportedFunction
func WithDateFormat(dateFormat string) Option {
	return func(u *MLogger) error {
		u.dateFormat = dateFormat
		return nil
	}
}

// WithLogFile sets the log file for the MLogger
//
//goland:noinspection GoUnusedExportedFunction
func WithLogFile(logfile string) Option {
	return func(u *MLogger) error {
		u.logfile = logfile
		return nil
	}
}

// WithLogStdout enables or disables logging to stdout
//
//goland:noinspection GoUnusedExportedFunction
func WithLogStdout(logStdout bool) Option {
	return func(u *MLogger) error {
		u.logStdout = logStdout
		return nil
	}
}

// WithLevel enables or disables logging the level
//
//goland:noinspection GoUnusedExportedFunction
func WithLevel(logLevel bool) Option {
	return func(u *MLogger) error {
		u.logLevel = logLevel
		return nil
	}
}

// WithDebug enables or disables debug logging
//
//goland:noinspection GoUnusedExportedFunction
func WithDebug(debug bool) Option {
	return func(u *MLogger) error {
		u.debug = debug
		return nil
	}
}

// open sets up the logger. This function is not exported, it is called by New
func (m *MLogger) open() (*MLogger, error) {
	var err error
	var fh *os.File

	if m.logfile != "" {

		// Sanitize the file path
		m.logfile = filepath.Clean(m.logfile)

		// Create the directory if it doesn't exist
		dir := filepath.Dir(m.logfile)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open the log file
		fh, err = os.OpenFile(m.logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			m.fileHandle = nil
			// If unable to log to file, force stdout logging
			m.logStdout = true
		} else {
			m.fileHandle = fh

			// Attempt to set the file mode to 0644 on a best-effort basis
			_ = os.Chmod(m.logfile, 0644)
		}
	} else {
		// If no log file is specified, force stdout logging
		m.logStdout = true
	}
	return m, nil
}

// Close closes the logger
func (m *MLogger) Close() {
	if m.fileHandle != nil {
		_ = m.fileHandle.Sync()
		_ = m.fileHandle.Close()
	}
}

// formatMessage formats the log message with a timestamp.
func (m *MLogger) formatMessage(level string, message string) string {
	var levelStr string
	if m.logLevel {
		levelStr = " [" + level + "]"
	} else {
		levelStr = ""
	}
	return fmt.Sprintf("%s%s%s %s",
		time.Now().Format(m.dateFormat),
		m.prefix, levelStr, message)
}

// writeLog writes a log message
func (m *MLogger) writeLog(level string, message string) {

	tmp := m.formatMessage(level, message) + "\n"

	//  Write and flush
	if m.fileHandle != nil {
		_, _ = m.fileHandle.WriteString(tmp)
		_ = m.fileHandle.Sync()
	}

	if m.logStdout {
		_, _ = os.Stdout.Write([]byte(tmp))
	}
}

// Debug logs a debug message.
func (m *MLogger) Debug(message string) {
	if m.debug {
		m.writeLog("DEBUG", message)
	}
}

// Info logs an informational message.
func (m *MLogger) Info(message string) {
	m.writeLog("INFO", message)
}

// Notice logs a notice message.
func (m *MLogger) Notice(message string) {
	m.writeLog("NOTICE", message)
}

// Warning logs a warning message.
func (m *MLogger) Warning(message string) {
	m.writeLog("WARNING", message)
}

// Error logs an error message.
func (m *MLogger) Error(message string) {
	m.writeLog("ERROR", message)
}

// Fatal logs a fatal error message.
func (m *MLogger) Fatal(message string) {
	m.writeLog("FATAL", message)
	m.FatalExit()
}

// Debugf logs a formatted debug message.
func (m *MLogger) Debugf(format string, v ...any) {
	if m.debug {
		m.writeLog("DEBUG", fmt.Sprintf(format, v...))
	}
}

// Infof logs a formatted informational message.
func (m *MLogger) Infof(format string, v ...any) {
	m.writeLog("INFO", fmt.Sprintf(format, v...))
}

// Noticef logs a formatted notice message.
func (m *MLogger) Noticef(format string, v ...any) {
	m.writeLog("NOTICE", fmt.Sprintf(format, v...))
}

// Warningf logs a formatted warning message.
func (m *MLogger) Warningf(format string, v ...any) {
	m.writeLog("WARNING", fmt.Sprintf(format, v...))
}

// Errorf logs a formatted error message.
func (m *MLogger) Errorf(format string, v ...any) {
	m.writeLog("ERROR", fmt.Sprintf(format, v...))
}

// Fatalf logs a formatted fatal message.
func (m *MLogger) Fatalf(format string, v ...any) {
	m.writeLog("FATAL", fmt.Sprintf(format, v...))
	m.FatalExit()
}

// FatalExit attempts to close the log and exits with a status code of 1
func (m *MLogger) FatalExit() {
	m.writeLog("FATAL", "Exiting with code 1 on fatal error")
	m.Close()
	os.Exit(1)
}
