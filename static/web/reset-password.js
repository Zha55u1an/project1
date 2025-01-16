document.addEventListener("DOMContentLoaded", function () {
    const forgotPasswordForm = document.getElementById("forgotPasswordForm");
    const submitButton = document.getElementById("submitButton");
    const btnText = submitButton.querySelector(".btn-text");
    const btnLoader = submitButton.querySelector(".btn-loader");
    const responseMessage = document.getElementById("responseMessage");

    forgotPasswordForm.addEventListener("submit", async function (e) {
        e.preventDefault();
        const email = document.getElementById("email").value;

        // Show loader and disable button
        btnText.style.display = "none";
        btnLoader.style.display = "inline-block";
        submitButton.disabled = true;

        try {
            const response = await fetch('/reset-password', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email }),
            });

            if (response.ok) {
                const result = await response.json();
                responseMessage.textContent = result.message || 'Reset code sent successfully.';
                responseMessage.style.color = 'green';

                // Redirect the user
                setTimeout(() => {
                    window.location.href = result.redirect;
                }, 2000);
            } else {
                const errorResult = await response.json();
                responseMessage.textContent = errorResult.error || 'Error occurred while sending reset code.';
                responseMessage.style.color = 'red';
            }
        } catch (error) {
            console.error('Error:', error);
            responseMessage.textContent = 'Something went wrong. Please try again.';
            responseMessage.style.color = 'red';
        } finally {
            // Hide loader and re-enable button
            btnText.style.display = "inline-block";
            btnLoader.style.display = "none";
            submitButton.disabled = false;
        }
    });
});
