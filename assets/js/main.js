(function () {
  "use strict";

  var wsConn = function () {
    var host = location.host;
    var path = location.pathname.replace(/^\/files\//, "");

    var conn = new WebSocket('ws://' + host + '/ws/' + path);

    conn.onmessage = function (e) {
      console.log("message!");
      location.reload();
    };
  };

  $(document).on('click', '.page', function (e) {
    e.preventDefault();
    var a = $(this);
    var url = a.attr('href');
    console.log(url);

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
