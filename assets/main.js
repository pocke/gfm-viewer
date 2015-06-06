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

  wsConn();
})();
