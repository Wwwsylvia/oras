/*
Copyright The ORAS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trace

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type contextKey int

// loggerKey is the associated key type for logger entry in context.
const loggerKey contextKey = iota

// CustomTextFormatter wraps the existing TextFormatter
type CustomTextFormatter struct {
	logrus.TextFormatter
}

// Format overrides the TextFormatter's Format method
// TODO: problem: this does not work for non-tty outputs (key-value pairs)
func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Use the TextFormatter to format the log entry
	var buf bytes.Buffer
	f.TextFormatter.DisableTimestamp = true // Disable the default timestamp
	f.TextFormatter.DisableQuote = true
	f.TextFormatter.DisableLevelTruncation = true

	formattedLog, err := f.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	// Get the timestamp and level
	timestamp := entry.Time.Format(time.RFC3339)

	// Prepend the timestamp and log level to the formatted message
	buf.WriteString(fmt.Sprintf("[%s] ", timestamp))
	buf.Write(formattedLog)
	buf.WriteString("\n\n")

	return buf.Bytes(), nil
}

// NewLogger returns a logger.
func NewLogger(ctx context.Context, debug bool, verbose bool) (context.Context, logrus.FieldLogger) {
	var logLevel logrus.Level
	if debug {
		logLevel = logrus.DebugLevel
	} else if verbose {
		logLevel = logrus.InfoLevel
	} else {
		logLevel = logrus.WarnLevel
	}

	logger := logrus.New()
	// logger.SetFormatter(&logrus.TextFormatter{
	// 	DisableQuote:           true,
	// 	FullTimestamp:          true,
	// 	DisableLevelTruncation: true,
	// })
	logger.SetFormatter(&CustomTextFormatter{})

	logger.SetLevel(logLevel)
	entry := logger.WithContext(ctx)
	return context.WithValue(ctx, loggerKey, entry), entry
}

// Logger return the logger attached to context or the standard one.
func Logger(ctx context.Context) logrus.FieldLogger {
	logger, ok := ctx.Value(loggerKey).(logrus.FieldLogger)
	if !ok {
		return logrus.StandardLogger()
	}
	return logger
}
