document
    .getElementById('loginForm')
    .addEventListener('submit', function (event)
    {
        event.preventDefault();

        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const errorMessage = document.getElementById('error-message');
        const errorText = document.getElementById('error-text');
        const closeError = document.getElementById('close-error');

        // Функция для очистки и скрытия блока с ошибкой
        function clearErrorMessage() {
            errorText.textContent = '';
            errorMessage.style.display = 'none';
        }

        // Закрытие уведомления об ошибке при клике на '×'
        closeError.onclick = clearErrorMessage;


        const requestBody = {
            email: email,
            password: password
        };

        fetch(
            'http://0.0.0.0:8000/api/v1/authorization',
            {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(requestBody)
            }
        )
            .then(response => response.json())
            .then(data => {
                if (data.error !== undefined && data.error.message !== undefined) {
                    errorMessage.textContent = data.error.message || "A server error occurred";
                    errorMessage.style.display = 'block';
                    setTimeout(clearErrorMessage, 5000);
                    return
                }

                // Set up the access token into the cookies by x-access-token key.
                document.cookie = `x-access-token=${data.data}`;
                // Redirect to homepage.
                window.location.replace("/");
            })
            .catch(error => console.error('Error:', error));
    }
);