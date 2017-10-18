$('.deleteContainer').on('click', function() {
  var ID = $(this).find("input").attr('value');
  
  var formData = new FormData();
  formData.append("id", ID);
  var request = new XMLHttpRequest();
  request.open("POST", "/v1/staging/environment/remove");
  request.send(formData);
  
  setTimeout(location.reload.bind(location), 1000);
  
  $(this)
  .addClass('loading')
  .siblings()
  .removeClass('loading');
});

$('.ui.positive.button').on('click',function () {
    $('.ui.mini.modal')
        .modal('setting', 'transition', 'browse')
        .modal('setting', 'closable', false)
        .modal('show');
});

function modalSubmitForm() {
  if (document.getElementsByClassName("javaPR")[0].value  || document.getElementsByClassName("htmlPR")[0].value || document.getElementsByClassName("slackUser")[0].value ){
  
    var javaPR = document.getElementsByClassName("javaPR")[0].value;
    var htmlPR = document.getElementsByClassName("htmlPR")[0].value;
    var slackUser = document.getElementsByClassName("slackUser")[0].value;
  
    var formData = new FormData();
  
    formData.append("javaPR", javaPR);
    formData.append("htmlPR", htmlPR);
    formData.append("slackUser", slackUser);
  
    var request = new XMLHttpRequest();
  
    request.open("POST", "/v1/staging/environment/create");
    request.send(formData);
    
    if (request.statusText == 200) {
      alert("Job iniciado!");
      setTimeout(location.reload.bind(location), 1000);
    }
    
  } else {
    alert("Field empty.")
  }
};