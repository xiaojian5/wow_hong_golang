function toggle() {
    let flag = document.getElementById("navbar-collapse-1").style.display;
    if (flag === "block") {
        document.getElementById("navbar-collapse-1").style.display = "none";
    } else {
        document.getElementById("navbar-collapse-1").style.display = "block";
    }
}

axios.get('./header.html', {})
    .then(function(response) {
        let div = document.createElement("div");
        div.id = "header";
        div.innerHTML = response.data;

        let contain = document.getElementById("container");
        let flag = contain.getAttribute("flag");
        contain.insertBefore(div, document.getElementById("content"));
        document.getElementById(flag).className = "active";
    });