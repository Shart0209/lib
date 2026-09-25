package stack

import (
	"fmt"
	"net/url"
	"regexp"
	"runtime"
	"strings"
)

const _defaultCallersDepth = 20

type Frame struct {
	Function string
	File     string
	Line     int
}

var vendorRe = regexp.MustCompile("^.*?/vendor/")

func (f Frame) String() string {
	var sb strings.Builder
	sb.WriteString(f.Function)
	if len(f.File) > 0 {
		if sb.Len() > 0 {
			sb.WriteRune(' ')
		}
		fmt.Fprintf(&sb, "(%v", f.File)
		if f.Line > 0 {
			fmt.Fprintf(&sb, ":%d", f.Line)
		}
		sb.WriteRune(')')
	}

	if sb.Len() == 0 {
		return "unknown"
	}

	return sb.String()
}

type Stack []Frame

func (fs Stack) String() string {
	return strings.Join(fs.Strings(), "; ")
}

func (fs Stack) Strings() []string {
	items := make([]string, len(fs))
	for i, f := range fs {
		items[i] = f.String()
	}
	return items
}

func CallerStack(skip, depth int) Stack {
	if depth <= 0 {
		depth = _defaultCallersDepth
	}

	pcs := make([]uintptr, depth)

	// +2 to skip this frame and runtime.Callers.
	n := runtime.Callers(skip+2, pcs)
	pcs = pcs[:n] // truncate to number of frames actually read

	result := make([]Frame, 0, n)
	frames := runtime.CallersFrames(pcs)
	for f, more := frames.Next(); more; f, more = frames.Next() {
		result = append(result, Frame{
			Function: sanitize(f.Function),
			File:     f.File,
			Line:     f.Line,
		})
	}
	return result
}

func sanitize(function string) string {
	if unescaped, err := url.QueryUnescape(function); err == nil {
		function = unescaped
	}
	return vendorRe.ReplaceAllString(function, "vendor/")
}
