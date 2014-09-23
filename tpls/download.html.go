package tpls

func init() {
	registerTemplate("download.html", `
{{template "header.html" .}}
<h1>下载</h1>
<form action="/download/do" method="post">

    
</form>

{{template "footer.html" .}}
`)
}
