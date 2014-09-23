package tpls

func init() {
	registerTemplate("home.html", `{{template "header.html" .}}
<form enctype="multipart/form-data" action="/upload/do" method="post">
  日期: <input type="text" name="date" />
  文件: <input type="file" name="uploadfile" />
  <input type="submit" value="upload" />
</form>
{{template "footer.html" .}}
`)
}
