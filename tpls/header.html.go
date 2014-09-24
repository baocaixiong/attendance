package tpls

func init() {
	registerTemplate("header.html", `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv="Content-Type" content="text/html;charset=UTF-8">
    <title>xlxs处理工具</title>
    {{template "css.html" .}}
</head>
<body>
    <div id='box'>
        <div id='box_heard'>wellcome to here ...</div>
        <div id='box_body'>
            <div id='nav'>
                <a class="btn {{if .IsUpload}}active {{end}}" href="/">上传</a>
                <a class="btn {{if .IsDownload}}active {{end}}" href="/download">下载</a>
            </div>
            <div id='text'>
               {{if .FlushMessage}} {{.FlushMessage}} {{end}} 
            </div>
`)
}
