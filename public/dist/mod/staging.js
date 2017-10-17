$('.ui.red.button').on('click', function() {
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