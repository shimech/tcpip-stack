package log

import (
	"fmt"
	"os"
	"runtime"
	"time"
	"unicode"

	"golang.org/x/exp/utf8string"
)

func Errorf(format string, ap ...any) (int, error) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return 0, fmt.Errorf("runtime error")
	}
	function := runtime.FuncForPC(pc)
	return lprintf(os.Stderr, 'E', file, line, function.Name(), format, ap...)
}

func Infof(format string, ap ...any) (int, error) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return 0, fmt.Errorf("runtime error")
	}
	function := runtime.FuncForPC(pc)
	return lprintf(os.Stderr, 'I', file, line, function.Name(), format, ap...)
}

func Debugf(format string, ap ...any) (int, error) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return 0, fmt.Errorf("runtime error")
	}
	function := runtime.FuncForPC(pc)
	return lprintf(os.Stderr, 'D', file, line, function.Name(), format, ap...)
}

func Debugdump(data []byte) {
	hexdump(os.Stderr, data)
}

func lprintf(
	fp *os.File,
	level int,
	file string,
	line int,
	funcName string,
	format string,
	ap ...any) (int, error) {
	result := 0
	n, err := fmt.Fprintf(fp, "%s [%c] %s: ", time.Now().Format("15:04:05"), level, funcName)
	if err != nil {
		return result, err
	}
	result += n

	n, err = fmt.Fprintf(fp, format, ap...)
	if err != nil {
		return result, err
	}
	result += n

	n, err = fmt.Fprintf(fp, " (%s:%d)\n", file, line)
	if err != nil {
		return result, err
	}
	result += n

	return result, nil
}

func hexdump(fp *os.File, data []byte) {
	size := len(data)
	fmt.Fprintf(fp, "+------+-------------------------------------------------+------------------+\n")
	for offset := 0; offset < size; offset += 16 {
		fmt.Fprintf(fp, "| %04x | ", offset)
		for index := 0; index < 16; index++ {
			if offset+index < size {
				fmt.Fprintf(fp, "%02x ", 0xff&data[offset+index])
			} else {
				fmt.Fprintf(fp, "   ")
			}
		}
		fmt.Fprintf(fp, "| ")
		for index := 0; index < 16; index++ {
			d := data[offset+index]
			if offset+index < size {
				if utf8string.NewString(string(d)).IsASCII() && unicode.IsPrint(rune(d)) {
					fmt.Fprintf(fp, "%c", data[offset+index])
				} else {
					fmt.Fprintf(fp, ".")
				}
			} else {
				fmt.Fprintf(fp, " ")
			}
		}
		fmt.Fprintf(fp, " |\n")
	}
	fmt.Fprintf(fp, "+------+-------------------------------------------------+------------------+\n")
}
