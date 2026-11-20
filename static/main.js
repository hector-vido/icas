let loginForm = document.getElementById("login-form") || null;
if (loginForm) {
    let loginButton = document.getElementById("login-button")
    let loginError = document.getElementById("login-error");
    let loginErrorMessage = document.getElementById("login-error-message");
    let loginSuccess = document.getElementById("login-success");
    loginForm.addEventListener("submit", function (e) {
        e.preventDefault();
        if (!loginForm.reportValidity()) {
            return
        }
        let email = document.getElementById("email");
        let password = document.getElementById("password");

        const req = new XMLHttpRequest();
        req.open('POST', '/login', true)
        req.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        req.addEventListener("load", function() {
            loginButton.className = "btn btn-secondary";
            loginButton.disabled = true;
            loginError.style.display = "none";
            if (req.status >= 200 && req.status < 300) {
                loginSuccess.style.display = "block";
                setTimeout(function() {
                    window.location.replace("/");
                }, 1000)
            } else {
                loginErrorMessage.textContent = req.responseText
                loginError.style.display = "block";
                loginButton.className = "btn btn-primary";
                loginButton.disabled = false;
            }
        });
        req.send(JSON.stringify({email: email.value, password: password.value}));
    });
}

let tamanduasList = document.getElementById("tamanduas-list") || null;
if (tamanduasList) {
    tamanduasList.querySelectorAll("a").forEach((a) => {
        a.addEventListener("click", function () {
            const req = new XMLHttpRequest();
            let animal = a.textContent.replace("#","");
            req.open('GET', `/animal/${animal}`, true)
            req.addEventListener("load", function() {
                if (req.status >= 200 && req.status < 300) {
                    console.log(req.responseText);
                } else {
                    console.log(req.responseText);
                }
            });
            req.send();
        });
    });
}

if (document.getElementById("map")) {
    const map = L.map('map').setView([-23.56380, -46.73389], 15);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    const locations = [
        {lat: -23.55980, lng: -46.73600},
        {lat: -23.55980, lng: -46.73404},
        {lat: -23.55982, lng: -46.73208},
        {lat: -23.55986, lng: -46.73012},
        {lat: -23.55992, lng: -46.72816},
        {lat: -23.56000, lng: -46.72620},
        {lat: -23.56012, lng: -46.72424},
        {lat: -23.56026, lng: -46.72228},
        {lat: -23.56050, lng: -46.72080},
        {lat: -23.56090, lng: -46.71980},
        {lat: -23.56140, lng: -46.71940},
        {lat: -23.56190, lng: -46.71980},
        {lat: -23.56220, lng: -46.72080},
        {lat: -23.56240, lng: -46.72176},
        {lat: -23.56240, lng: -46.72372},
        {lat: -23.56240, lng: -46.72568},
        {lat: -23.56240, lng: -46.72764},
        {lat: -23.56240, lng: -46.72960},
        {lat: -23.56240, lng: -46.73156},
        {lat: -23.56240, lng: -46.73352},
        {lat: -23.56240, lng: -46.73548},
        {lat: -23.56250, lng: -46.73744},
        {lat: -23.56280, lng: -46.73900},
        {lat: -23.56320, lng: -46.73980},
        {lat: -23.56360, lng: -46.74020},
        {lat: -23.56400, lng: -46.73980},
        {lat: -23.56430, lng: -46.73900},
        {lat: -23.56440, lng: -46.73752},
        {lat: -23.56440, lng: -46.73556},
        {lat: -23.56440, lng: -46.73360},
    ];

    locations.forEach((loc, index) => {
        // O Leaflet usa [latitude, longitude], a ordem é oposta à do Google!
        const marker = L.marker([loc.lat, loc.lng]).addTo(map);
        marker.bindPopup(`Ponto ${index + 1}<br>(${loc.lat}, ${loc.lng})`);
    });
}