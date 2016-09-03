
function init() {
	$(document).keydown(function(e) {
        switch (e.keyCode) {
            case 38: // up
                moveUp();
            break;
            case 40: // down
                moveDown();
            break;
            case 37: // left
                moveLeft();
            break;
            case 39: // right
                moveRight();
            break;
        }
    });

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

	function moveUp() {
	    move('up');
	}

	function moveDown() {
	    move('down')
	}

	function moveLeft() {
	    move('left');
	}

	function moveRight() {
	    move('right');
	}

    function move(direction) {
        var speed = $('#speed').val(),
          data = JSON.stringify({
            car: car,
            speed: parseInt(speed, 10),
            direction: direction
        });
        console.log(data);
        ws.send(data);
    }
}

$(document).ready(function() {
    init()
});
