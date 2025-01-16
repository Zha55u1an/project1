document.addEventListener("DOMContentLoaded", function () {
    const email = document.querySelector("strong").textContent; // Extract email from HTML
    const resendButton = document.getElementById("resendCode");
    const timerSpan = document.getElementById("timer");
    const feedbackMessage = document.getElementById("feedbackMessage");

    let timer = 15;

    // Function to Start Timer
    function startTimer() {
        timer = 15;
        resendButton.disabled = true;
        const interval = setInterval(() => {
            timer--;
            timerSpan.textContent = timer;

            if (timer <= 0) {
                clearInterval(interval);
                resendButton.disabled = false;
                document.getElementById("resendMessage").textContent = "Didn't receive the code? Click the button below to resend.";
            }
        }, 1000);
    }

    // Start Timer on Page Load
    startTimer();

    // Handle Form Submission
    document.getElementById("verifyForm").addEventListener("submit", function (e) {
        e.preventDefault();
        const code = document.getElementById("code").value;

        fetch("/verify-reset-code", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, code }),
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    feedbackMessage.textContent = data.error;
                    feedbackMessage.style.color = "red";
                } else {
                    feedbackMessage.textContent = data.message;
                    feedbackMessage.style.color = "green";
                    setTimeout(() => {
                        window.location.href = data.redirect;
                    }, 2000);
                }
            })
            .catch(() => {
                feedbackMessage.textContent = "Failed to verify the code. Please try again.";
                feedbackMessage.style.color = "red";
            });
    });

    // Handle Resend Button
    resendButton.addEventListener("click", function () {
        fetch("/resend-code", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email }),
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.message) {
                    alert(data.message || "Code resent!");
                    startTimer();
                }
            })
            .catch(() => {
                alert("Failed to resend the code. Please try again.");
            });
    });
});
