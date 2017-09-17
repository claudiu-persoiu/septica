var communicationHandler = function (url, callback) {
    var socket, socketOpened = false;

    var connect = function () {
        socket = new WebSocket(url);

        socket.onopen = function () {
            socketOpened = true;
        };

        socket.onmessage = function (event) {
            callback(JSON.parse(event.data));
        };

        socket.onclose = function () {
            socketOpened = false;
            throw "Disconnected";
        };
    };

    return {
        send: function (type, message) {
            if (socketOpened) {
                return socket.send(JSON.stringify({"step": type, "command": message}));
            }
            return false;
        },
        isConnected: function () {
            return socketOpened;
        },
        connect: function () {
            connect();
        }
    };
};