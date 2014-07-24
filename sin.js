var a = 0;
function render(dt) {
	var xoff = pix.width / 2;
	var yoff = pix.height / 2;

	for (var y = 0; y < pix.height; y++) {
		for (var x = 0; x < pix.width; x++) {
			var x2 = x - xoff;
			var y2 = y - yoff;
			var d = Math.sqrt(x2*x2 + y2*y2);
			var t = Math.sin(d);

			var r = Math.max(0, Math.min(255, Math.sin(a) * t * 200));
			var g = Math.max(0, Math.min(255, Math.sin(a/3) * 125 + t * 80));
			var b = Math.max(0, Math.min(255, Math.sin(a/2) * 235 + t * 20));

			pix.setPixel(x, y, r, g, b, 255);
		}
	}
	a += dt / 1000.0;
}

