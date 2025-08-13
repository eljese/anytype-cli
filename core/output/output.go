package output

import (
	"fmt"
	"os"
)

func Success(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, "✓ "+format+"\n", args...)
}

func Info(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

func Warning(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "⚠ "+format+"\n", args...)
}

func Error(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func Debug(format string, args ...interface{}) {
	// Debug messages are hidden unless debug flag is set
}

func Print(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}
