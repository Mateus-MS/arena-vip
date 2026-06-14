(function () {
  const canvas = document.getElementById("hero-canvas");
  if (!canvas) return;
  const ctx = canvas.getContext("2d");

  fetch("/static/particles.json")
    .then(r => r.json())
    .then(init)
    .catch(e => console.warn("particles.json not found:", e));

  function init(data) {
    const fps = data.fps;
    const srcW = data.width;
    const srcH = data.height;
    const frames = data.frames;
    const frameDuration = 1000 / fps;
    let frameIdx = 0;
    let lastTime = 0;

    function resize() {
      canvas.width  = canvas.clientWidth;
      canvas.height = canvas.clientHeight;
    }
    new ResizeObserver(resize).observe(canvas);
    resize();

    function render(now) {
      requestAnimationFrame(render);
      if (now - lastTime < frameDuration) return;
      lastTime = now;

      const cw = canvas.width;
      const ch = canvas.height;
      const aspect = srcW / srcH;
      var rw = cw, rh = ch, rx = 0, ry = 0;
      if (cw / ch > aspect) {
        rw = ch * aspect;
        rx = (cw - rw) / 2;
      } else {
        rh = cw / aspect;
        ry = (ch - rh) / 2;
      }

      ctx.clearRect(0, 0, cw, ch);

      var pts = frames[frameIdx] || [];
      for (var i = 0; i < pts.length; i++) {
        var pt = pts[i];
        var x  = rx + pt[0] * rw;
        var y  = ry + pt[1] * rh;
        var b  = pt[2] !== undefined ? pt[2] : 0.5;
        var r  = Math.max(1, 1.0 + b * 1.4);
        var a  = (0.4 + b * 0.6).toFixed(2);
        ctx.beginPath();
        ctx.arc(x, y, r, 0, 6.2832);
        ctx.fillStyle = "rgba(201,168,76," + a + ")";
        ctx.fill();
      }

      frameIdx = (frameIdx + 1) % frames.length;
    }

    requestAnimationFrame(render);
  }
})();
