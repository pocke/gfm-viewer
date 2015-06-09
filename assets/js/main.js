(function () {
  "use strict";

  var current_file = "";

  var wsConn = function () {
    var host = location.host;

    var conn = new WebSocket('ws://' + host + '/ws');

    conn.onmessage = function (e) {
      if (e.data === current_file) {
        $('a.page[data-filepath="' + e.data + '"]').click();
      }
    };
  };

  wsConn();

  $(document).on('click', '.page', function (e) {
    e.preventDefault();
    var a = $(this);
    var url = a.attr('href');
    current_file = a.attr('data-filepath');

    $.ajax({url: url}).done(function (data) {
      $('#md-body').html(data);
      $('#md-title').text(url.replace(/^.+\/([^\/]+)$/, "$1"));
      $("#md").removeClass("hidden");
    }).fail(function (xhr) {
      $('#modal-body').text(xhr.responseText);
      $('#modal').modal('show');
    });
  });

})();
