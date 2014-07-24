var canvas = document.getElementById('canvas');
var c = canvas.getContext('2d');

var SCALE = 8

var offscreenCanvas = document.createElement('canvas');
offscreenCanvas.width = canvas.width / SCALE;
offscreenCanvas.height = canvas.height / SCALE;
var offscreenCtx = offscreenCanvas.getContext('2d');
var buffer = offscreenCtx.getImageData(0, 0, offscreenCanvas.width, offscreenCanvas.height);

// Define pix stuff.
var pix = {};

pix.setPixel = function(x, y, r, g, b, a) {
	index = (x + y * buffer.width) * 4;
	buffer.data[index + 0] = r;
	buffer.data[index + 1] = g;
	buffer.data[index + 2] = b;
	buffer.data[index + 3] = a;
}

pix.width = buffer.width;
pix.height = buffer.height;

var lastT = null;
var step = function(ts) {
	var dt = ts - (lastT || ts);
	lastT = ts;

	render(dt);

	for (var x = 0, x2 = 0; x < buffer.width; ++x, x2+=SCALE) {
		for (var y = 0, y2 = 0; y < buffer.height; ++y, y2+=SCALE) {
			var i = (y * buffer.width + x) * 4;
			var r = buffer.data[i+0];
			var g = buffer.data[i+1];
			var b = buffer.data[i+2];
			var a = buffer.data[i+3];
			c.fillStyle = 'rgba(' + r + ',' + g + ',' + b + ',' + (a/255) + ')';
			c.fillRect(x2, y2, SCALE, SCALE);
		}
	}

	window.requestAnimationFrame(step);
}

window.requestAnimationFrame(step);
