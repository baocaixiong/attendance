package tpls

func init() {
	registerTemplate("download.html", `
{{template "header.html" .}}
<h1>下载</h1>
<form action="/download/do" method="post">
  日期: <input type="text" name="date" />
  <input type="submit" value="下载" />
</form>

{{template "footer.html" .}}
`)
}
