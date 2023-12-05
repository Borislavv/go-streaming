document.getElementById('loginForm').addEventListener('submit', function (event) {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const requestBody = {
        email: email,
        password: password
    };

    fetch('http://0.0.0.0:8000/api/v1/authorization', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
    })
        .then(response => response.json())
        .then(data => {
            // Set up the access token into the cookies by x-access-token key.
            document.cookie = `x-access-token=${data.data}`;
            // Redirect to homepage.
            window.location.replace("/");
        })
        .catch(error => console.error('Error:', error));
});