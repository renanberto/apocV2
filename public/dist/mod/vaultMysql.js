$('#mysqlForm')
.form({
  on: 'blur',
  fields: {
    empty: {
      identifier: 'githubToken',
      rules: [
        {
          type: 'empty',
          prompt: 'Please enter github token.'
        }
      ]
    }
  }
});

$('#submitMysqlForm').on('click', function() {
  
  if ($('#mysqlForm').form('is valid')) {
    
    $('#submitMysqlForm')
      .addClass('loading')
      .siblings()
      .removeClass('loading');
      
    var options = {
      dataType: 'json',
      success: function(response) {
        var myData = JSON.stringify(response);
        var myDataObject = JSON.parse(myData);
        var info = 'Username: <code>' + myDataObject.username + '</code></br> Password: <code>' + myDataObject.password + '</code>';
        
        $("#mysqlCredsInformation").append(info);
        $('.ui.modal')
        .modal('setting', 'transition', 'browse')
        .modal('setting', 'closable', false)
        .modal('show');
      }
    };
    
    $('#mysqlForm').ajaxForm(options);
  }
});

function quitWithReload(){location.reload();}
$('.ui.radio.checkbox').checkbox();
$('.ui.cancel.button');
