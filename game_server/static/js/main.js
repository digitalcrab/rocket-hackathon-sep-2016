$(document).ready(function() {
    var pressedKeys = {},
            keysMap = {
                37: "left",
                38: "up",
                39: "right",
                40: "down",
                83: "stop"
            },
            buttonMap = {
                37: $('#move-left'),
                38: $('#move-up'),
                39: $('#move-right'),
                40: $('#move-down'),
                83: $('#move-stop')
            };
    
        $(document).keydown(function (e) {
            var k = e.keyCode;
            if (k == 37 || k == 38 || k == 39 || k == 40 || k == 83) {
                pressedKeys[k] = true;
                triggerMovement();
            }
        });
    
        $(document).keyup(function (e) {
            var k = e.keyCode;
            if (k == 37 || k == 38 || k == 39 || k == 40 || k == 83) {
                delete pressedKeys[k];
                buttonMap[k].removeClass('selected');
            }
        });
    
        function triggerMovement() {
            var keys = [];
    
            $.each(pressedKeys, function(key) {
                keys.push(keysMap[key]);
                buttonMap[key].addClass('selected');
            });
    
            move(keys);
        }
    
        $('#car-red').on('click', function () {
            selectCar('red', $(this), $('#car-green'));
        });
    
        $('#car-green').on('click', function () {
            selectCar('green', $(this), $('#car-red'));
        });
    
        var car = null;
        var ws = null;
    
        function wsConnect() {
            var url = (window.location.protocol == 'http:' ? 'ws:' : 'wss:') + '//' +
                window.location.host + '/game';
            ws = new WebSocket(url);
    
            ws.addEventListener('open', function (e) {
                document.getElementById('error').style.display = 'none';
            });
    
            ws.addEventListener('message', function (e) {
                console.log(e.data)
            });
    
            ws.addEventListener('close', function (e) {
                document.getElementById('error').style.display = 'block';
                window.setTimeout(wsConnect, 1000);
            });
    
            ws.addEventListener('error', function (e) {
                document.getElementById('error').style.display = 'block';
            });
    
        }
    
        wsConnect();
    
        $('#car-red').click();
    
        function selectCar(color, me, other) {
            car = color;
    
            me.addClass('selected');
    
            if (other.hasClass('selected')) {
                other.removeClass('selected');
            }
        }
    
        function move(directions) {
            var speed = $('#speed').val(),
              data = JSON.stringify({
                car: car,
                speed: parseInt(speed, 10),
                directions: directions
            });
            console.log(data);
            ws.send(data);
        }
});
