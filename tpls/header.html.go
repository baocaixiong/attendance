package tpls

func init() {
	registerTemplate("header.html", `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    {{template "css.html" .}}
</head>
<body>
   <body>
        <div>
          <a href="/">Project name</a>
        </div>
        <div>
          <ul>
            <li {{if .IsUpload}}class="active" {{end}}><a href="/">上传</a></li>
            <li {{if .IsDownload}}class="active" {{end}}><a href="/download">下载</a></li>
          </ul>
        </div> {{template "flush.html" .}}
`)
}
