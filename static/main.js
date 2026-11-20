let loginButton = document.getElementById("login-button") || null
if (loginButton) {
    let loginError = document.getElementById("login-error");
    let loginErrorMessage = document.getElementById("login-error-message");
    let loginSuccess = document.getElementById("login-success");
    document.getElementById("login-button").addEventListener("click", function () {
        let loginForm = document.getElementById("login-form")
        loginForm.reportValidity()
        let email = document.getElementById("email");
        let password = document.getElementById("password");

        const req = new XMLHttpRequest();
        let reqListener = function() {
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
        }
        req.open('POST', '/login', true)
        req.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        req.addEventListener("load", reqListener);
        req.send(JSON.stringify({email: email.value, password: password.value}));
    });
}