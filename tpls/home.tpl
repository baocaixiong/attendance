{{define "body"}}
<!-- Begin page content -->
<div class="container">
  <div class="page-header">
    <h1>Sticky footer with fixed navbar</h1>
</div>

<form action="/update/do" method="post" class="form-horizontal" role="form">
<div class="form-group">
    <label for="date" class="col-sm-2 control-label">日期</label>
    <div class="col-sm-10">
      <input type="text" value="" id="datetimepicker" name="date">
    </div>
    
  </div>

  <div class="form-group">
    <div class="col-sm-offset-2 col-sm-10">
      <div class="checkbox">
        <label>
          <input type="radio" name="isContainHead" value="on"> 包含首部
        </label>
      </div>
    </div>
  </div>
  <div class="form-group">
    <div class="col-sm-offset-2 col-sm-10">
      <button type="submit" class="btn btn-default">上传</button>
    </div>
  </div>
</form>

<script type="text/javascript">
    $('#datetimepicker').datetimepicker({
        format: 'yyyy-mm',
        language:  'zh-CN',
        weekStart: 1,
        todayBtn:  1,
        autoclose: 1,
        todayHighlight: 1,
        startView: 2,
        minView: 3,
        forceParse: 0
    });
</script>

</div>

{{end}}