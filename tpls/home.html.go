package tpls

func init() {
	registerTemplate("home.html", `{{template "header.html" .}}
<div id='form'>
    <form action="/upload/do" method="POST" enctype="multipart/form-data">
    <table>
        <tr>
            <label>
                <td width='90' align='left'>时间</td>
                <td align='left'><input class= 'sl' type="text" name="date" placeholder="2000-01-01"></td>
            </label>
        </tr>
        <tr>
            <label>
                <td width='90' align='left'>文件</td>
                <td align='left'><input name="uploadfile" type="file"></td>
            </label>
        </tr>
        <tr>
            <label>
                <td width='90' align='left'>包含首部?</td>
                <td align='left'><input  type="checkbox" name="isContainHead" checked=""></td>
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
