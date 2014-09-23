package tpls

func init() {
	registerTemplate("flush.html", `
    {{if .FlushMessage}} {{.FlushMessage}} {{end}}
`)
}
