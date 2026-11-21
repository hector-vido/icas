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
        req.addEventListener("load", function () {
            loginButton.className = "btn btn-secondary";
            loginButton.disabled = true;
            loginError.style.display = "none";
            if (req.status >= 200 && req.status < 300) {
                loginSuccess.style.display = "block";
                setTimeout(function () {
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

if (document.getElementById("map")) {

    const map = L.map('map').setView([-14.2350, -51.9253], 4);
    const markerGroup = L.layerGroup().addTo(map);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);
    updateMarkers = function(locations) {
        markerGroup.clearLayers();
        locations.forEach((loc, index) => {
            L.marker([loc.lat, loc.lng])
                .bindPopup(`Ponto ${index + 1}<br>${loc.registered_at}<br>(${loc.lat}, ${loc.lng})`)
                .addTo(markerGroup);
        });
    };

    let tamanduasList = document.getElementById("tamanduas-list") || null;
    tamanduasList.querySelectorAll("a").forEach(a => {
        a.addEventListener("click", function () {
            const req = new XMLHttpRequest();
            let animal = a.textContent.replace("#", "");
            req.open('GET', `/animal/${animal}`, true)
            req.addEventListener("load", function () {
                if (req.status >= 200 && req.status < 300) {
                    let animal = JSON.parse(req.responseText);
                    document.getElementById("animal-nome").textContent = animal.name;
                    document.getElementById("link-coordenadas").href = `downloadCoords/${animal.name}`;
                    document.getElementById("animal-peso").textContent = animal.weight;
                    document.getElementById("animal-ultima-verificacao").textContent = animal.lastCheck;
                    document.getElementById("animal-data-equipamento").textContent = animal.equipmentDate;
                    document.getElementById("link-coordenadas").style.display = 'inline';
                    updateMarkers(animal.coordinates);
                    if (animal.coordinates.length > 0) {
                        map.setView([animal.coordinates[0].lat, animal.coordinates[0].lng], 15);
                    }
                } else {
                    alert("Algo deu errado, tente novamente mais tarde.")
                }
            });
            req.send();
        });
    });
}