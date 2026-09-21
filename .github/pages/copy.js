document.querySelectorAll("pre").forEach(function (pre) {
  var btn = document.createElement("button");
  btn.type = "button";
  btn.className = "copy-btn";
  btn.textContent = "copy";

  btn.addEventListener("click", function () {
    var text = pre.innerText.replace(/\n+$/, "");
    navigator.clipboard.writeText(text).then(function () {
      btn.textContent = "copied";
      btn.dataset.copied = "true";
      setTimeout(function () {
        btn.textContent = "copy";
        btn.dataset.copied = "false";
      }, 1400);
    });
  });

  pre.appendChild(btn);
});
