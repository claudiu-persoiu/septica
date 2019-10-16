var communicationHandler = function (url, callback) {
    var socket, socketOpened = false;

    var connect = function (uuid) {
        socket = new WebSocket(url);

        socket.onopen = function () {
            socketOpened = true;
            send('identify', uuid || getUUID());
        };

        socket.onmessage = function (event) {
            var messages = event.data.split("\n");
            for (var message of messages) {
                callback(JSON.parse(message));
            }
        };

        socket.onclose = function () {
            socketOpened = false;
            throw "Disconnected";
        };
    };

    var getUUID = function () {
        if (!window.localStorage.getItem('uuid')) {
            var uuid = ([1e7]+-1e3+-1e3+-1e3+-1e11).replace(/[018]/g, c =>
                (c ^ crypto.getRandomValues(new Uint8Array(1))[0]  & 15 >> c ).toString(16)
              );
            window.localStorage.setItem('uuid', uuid);
        }

        return window.localStorage.getItem('uuid');
    };

    var send = function (action, data) {
        if (!data) {
            data = "";
        }

        if (socketOpened) {
            return socket.send(JSON.stringify({"action": action, "data": data}));
        }
        return false;
    }

    return {
        send: function (action, data) {
            return send(action, data);
        },
        isConnected: function () {
            return socketOpened;
        },
        connect: function (uuid) {
            connect(uuid);
        }
    };
};