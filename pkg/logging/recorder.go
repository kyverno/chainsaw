package logging

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
)

// WithRecorder returns a context whose logger writes every log entry to w, as plain text
// without colors, and then passes the entry to the logger of ctx.
func WithRecorder(ctx context.Context, w io.Writer) context.Context {
	next := getLogger(ctx)
	return WithLogger(ctx, LoggerFunc(func(ctx context.Context, operation Operation, status Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
		_, _ = fmt.Fprintln(w, format(operation, status, obj, args...))
		if next != nil {
			next.Log(ctx, operation, status, obj, color, args...)
		}
	}))
}

func format(operation Operation, status Status, obj client.Object, args ...fmt.Stringer) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s | %s", operation, status)
	if obj != nil {
		gvk := obj.GetObjectKind().GroupVersionKind()
		fmt.Fprintf(&b, " | %s/%s @ %s", gvk.GroupVersion(), gvk.Kind, client.Name(client.Key(obj)))
	}
	for _, arg := range args {
		b.WriteString("\n")
		b.WriteString(arg.String())
	}
	return b.String()
}
