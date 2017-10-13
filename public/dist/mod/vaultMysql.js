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
    },
    checkbox: {
      identifier: 'accessMode',
      rules: [
        {
          type   : 'checked',
          prompt : 'Please enter access mode.'
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
        var Data = JSON.stringify(response);
        var DataObject = JSON.parse(Data);
        
        if (DataObject.respError != null) {
          var myErrorResponse = JSON.stringify(DataObject.respError);
          var myErrorResponseObject = JSON.parse(myErrorResponse);
          var info = '<h3>Error: </h3>' + myErrorResponseObject
          
        } else {
          var info =
            '<div class="header"> Username: </div>\n' +
            '    <li class="list">' + DataObject.username + '</li>' +
            '</div> </br>\n' +
  
            '<div class="header"> Password: </div>\n' +
            '    <li class="list">' + DataObject.password + '</li>' +
            '</div> </br>\n' +
  
            '<div class="header"> Tempo de duração: </div>\n' +
            '    <li class="list">' + DataObject.lease_time + ' seconds </li>' +
            '</div>\n'
          ;
        }
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
