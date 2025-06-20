let np = document.getElementById("nav-portrait");
np.style.position = "fixed";
np.style.right = "-" + np.offsetWidth + 50 + "px";
function showNavPortrait() {
    np.style.right = 0;
    setTimeout(function () {
        document.addEventListener('click', tf, false);
    }, 50);
}

function tf() {
    np.style.right = "-" + np.offsetWidth + "px";
    document.removeEventListener('click', tf);
}
