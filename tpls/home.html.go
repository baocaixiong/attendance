package tpls

func init() {
	registerTemplate("home.html", `
{{template "header.html" .}}
<div class="container">
  <div class="page-header">
    <h1>Sticky footer with fixed navbar</h1>
</div>

<form action="/upload/do" method="post" class="form-horizontal" enctype="multipart/form-data">
<div class="form-group">
    <label for="date" class="col-sm-2 control-label">日期</label>
    <div class="col-sm-10">
      <input type="text" value="" id="datetimepicker" name="date">
    </div>
    
  </div>

  <div class="form-group">
    <label for="attendanceFile"  class="col-sm-2 control-label">文件</label>
    <input type="file" name="uploadFile" id="attendanceFile" style="margin-left: 200px;">
  </div>

  <div class="form-group">
    <div class="col-sm-offset-2 col-sm-10">
      <div class="checkbox">
        <label>
          <input type="checkbox" name="isContainHead" value="on"> 包含首部
        </label>
      </div>
    </div>
  </div>
  <div class="form-group">
    <div class="col-sm-offset-2 col-sm-10">
      <input type="submit" class="btn btn-default" value="上传"/>
    </div>
  </div>
</form>
</div>

{{template "footer.html" .}}
`)
}
