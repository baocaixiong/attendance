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

    <!-- Fixed navbar -->
    <nav class="navbar navbar-default navbar-fixed-top" role="navigation">
      <div class="container">
        <div class="navbar-header">
          <a class="navbar-brand" href="/">Project name</a>
        </div>
        <div id="navbar" class="collapse navbar-collapse">
          <ul class="nav navbar-nav">
            <li {{if .IsUpload}}class="active" {{end}}><a href="/">上传</a></li>
            <li {{if .IsDownload}}class="active" {{end}}><a href="/download">下载</a></li>
          </ul>
        </div><!--/.nav-collapse -->
      </div>
    </nav>
`)
}
