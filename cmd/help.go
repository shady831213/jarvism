package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/shady831213/jarvism/cmd/base"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

func help(args []string) {
	// 'jarvsim help documentation' generates doc.go.
	if len(args) == 1 && args[0] == "documentation" {
		buf := new(bytes.Buffer)
		printUsage(buf, base.Jarvism)
		usage := &base.Command{Long: buf.String()}
		cmds := []*base.Command{usage}
		for _, cmd := range base.Jarvism.Commands {
			cmds = append(cmds, cmd)
			cmds = append(cmds, cmd.Commands...)
		}
		tmpl(&commentWriter{W: os.Stdout}, documentationTemplate, cmds)
		return
	}

	cmd := base.Jarvism
Args:
	for i, arg := range args {
		for _, sub := range cmd.Commands {
			if sub.Name() == arg {
				cmd = sub
				continue Args
			}
		}

		// helpSuccess is the help command using as many args as possible that would succeed.
		helpSuccess := "jarvism help"
		if i > 0 {
			helpSuccess = " " + strings.Join(args[:i], " ")
		}
		fmt.Fprintf(os.Stderr, "javisim help %s: unknown help topic. Run '%s'.\n", strings.Join(args, " "), helpSuccess)
		os.Exit(2)
	}

	if len(cmd.Commands) > 0 {
		printUsage(os.Stdout, cmd)
	} else {
		tmpl(os.Stdout, helpTemplate, cmd)
	}
	return
}

var usageTemplate = `{{.Long | trim}}

Usage:

	{{.UsageLine}} <command> [arguments]

The commands are:
{{range .Commands}}{{if or (.Runnable) .Commands}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "jarvsim help{{with .LongName}} {{.}}{{end}} <command>" for more information about a command.
{{if eq (.UsageLine) "jarvsim"}}
Additional help topics:
{{range .Commands}}{{if and (not .Runnable) (not .Commands)}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "jarvsim help{{with .LongName}} {{.}}{{end}} <topic>" for more information about that topic.
{{end}}
`

var helpTemplate = `{{if .Runnable}}usage: {{.UsageLine}}

{{end}}{{.Long | trim}}
`

var documentationTemplate = `{{range .}}{{if .Short}}{{.Short | capitalize}}

{{end}}{{if .Commands}}` + usageTemplate + `{{else}}{{if .Runnable}}Usage:

	{{.UsageLine}}

{{end}}{{.Long | trim}}


{{end}}{{end}}`

// commentWriter writes a jarvsim comment to the underlying io.Writer,
// using line comment form (//).
type commentWriter struct {
	W            io.Writer
	wroteSlashes bool // Wrote "//" at the beginning of the current line.
}

func (c *commentWriter) Write(p []byte) (int, error) {
	var n int
	for i, b := range p {
		if !c.wroteSlashes {
			s := "//"
			if b != '\n' {
				s = "// "
			}
			if _, err := io.WriteString(c.W, s); err != nil {
				return n, err
			}
			c.wroteSlashes = true
		}
		n0, err := c.W.Write(p[i : i+1])
		n += n0
		if err != nil {
			return n, err
		}
		if b == '\n' {
			c.wroteSlashes = false
		}
	}
	return len(p), nil
}

// An errWriter wraps a writer, recording whether a write error occurred.
type errWriter struct {
	w   io.Writer
	err error
}

func (w *errWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	ew := &errWriter{w: w}
	err := t.Execute(ew, data)
	if ew.err != nil {
		// I/O error writing. Ignore write on closed pipe.
		if strings.Contains(ew.err.Error(), "pipe") {
			os.Exit(1)
		}
		base.Fatalf("writing output: %v", ew.err)
	}
	if err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func printUsage(w io.Writer, cmd *base.Command) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, cmd)
	bw.Flush()
}
