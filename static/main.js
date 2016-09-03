
function init() {
    var pressedKeys = {},
        keysMap = {
            37: "left",
            38: "up",
            39: "right",
            40: "down"
        };

    $(document).keydown(function (e) {
        var k = e.keyCode;
        if (k == 37 || k == 38 || k == 39 || k == 40) {
            pressedKeys[k] = true;
            triggerMovement()
        }
    });

    $(document).keyup(function (e) {
        var k = e.keyCode;
        if (k == 37 || k == 38 || k == 39 || k == 40) {
            delete pressedKeys[k];
        }
    });

    function triggerMovement() {
        var keys = [];
        $.each(pressedKeys, function(key) {
            keys.push(keysMap[key])
        })
        move(keys);
    }

    console.log($('#car-red'));

    $('#car-red').on('click', function () {
        selectCar('red');
        if ($('#car-green').hasClass('active')) {
            $('#car-green').button('toggle');
        }
    });

    $('#car-green').on('click', function () {
        selectCar('green')
        if ($('#car-red').hasClass('active')) {
            $('#car-red').button('toggle');
        }
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

	function selectCar(color) {
	    car = color;
	    console.log(car);
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
}

$(document).ready(function() {
    init()
});
