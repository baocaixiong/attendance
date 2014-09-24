package tpls

func init() {
	registerTemplate("download.html", `
{{template "header.html" .}}
<div id='form'>
    <form action="/download/do" method="POST" enctype="multipart/form-data">
    <table>
        <tr>
            <label>
                <td width='90' align='left'>时间</td>
                <td align='left'><input class= 'sl' type="text" name="date" placeholder="2000-01"></td>
            </label>
        </tr>
        <tr>
            <td width='90' align='left'></td>
            <td align='left'><input id='sub' type="submit" value='提交'></td>
        </tr>

    </table>

    </form>
</div>
{{template "footer.html" .}}
`)
}
