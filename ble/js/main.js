var bluetoothDevice,
    serialPortCharacteristic,
    pressedKeys = {},
    keysMap = {
        37: 4, // left
        38: 1, // up
        39: 8, // right
        40: 2, // down
        83: 0  // stop
    }
    buttonMap = {
        37: $('#move-left'),
        38: $('#move-up'),
        39: $('#move-right'),
        40: $('#move-down'),
        83: $('#move-stop')
    },
    writing = false;

function onDisconnected(event) {
  let device = event.target;
  console.log('Device ' + device.name + ' is disconnected.');
  bluetoothDevice = null;
  serialPortCharacteristic = null;
}

function triggerMovement() {
    if (!serialPortCharacteristic || writing) {
        return;
    }

    var direction = 0;
    $.each(pressedKeys, function(key) {
        direction += keysMap[key];
        buttonMap[key].addClass('selected');
    });

    var val = new Uint8Array([direction]);
    writing = true;
    serialPortCharacteristic.writeValue(val).then(function() {
        console.log('Direction: ' + direction);
        writing = false;
    });
}

$(document).ready(function() {
    $('#connect').on('click', function () {
        bluetoothDevice = null;
        serialPortCharacteristic = null;
        navigator.bluetooth.requestDevice({
            filters: [{
                services: ['0000dfb0-0000-1000-8000-00805f9b34fb']
            }]
        })
        .then(device => {
            bluetoothDevice = device;
            console.log('Device ' + device.name + ' is connected.');
            device.addEventListener('gattserverdisconnected', onDisconnected);
            return device.gatt.connect();
        })
        .then(server => {
            // DF Robot Service
            return server.getPrimaryService('0000dfb0-0000-1000-8000-00805f9b34fb');
        })
        .then(service => {
            // Serial Port
            return service.getCharacteristic('0000dfb1-0000-1000-8000-00805f9b34fb');
        })
        .then(characteristic => {
            console.log('Characteristic ' + characteristic.uuid + ' is connected.');
            serialPortCharacteristic = characteristic;
        })
        .catch(error => {
            console.log(error);
        });
    });

    $('#disconnect').on('click', function () {
        if (!bluetoothDevice || !bluetoothDevice.gatt.connected) {
            return;
        }
        bluetoothDevice.gatt.disconnect();
    });

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
});
