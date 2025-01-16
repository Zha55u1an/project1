document.getElementById('updatePasswordForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const newPassword = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const feedback = document.getElementById('feedbackMessage');

    // Validate if passwords match
    if (newPassword !== confirmPassword) {
        feedback.textContent = 'Passwords do not match.';
        feedback.style.color = 'red';
        return;
    }

    try {
        const response = await fetch('/update-password', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, new_password: newPassword, confirm_password: confirmPassword }),
        });

        const result = await response.json();

        if (response.ok) {
            feedback.textContent = result.message || 'Password updated successfully.';
            feedback.style.color = 'green';
            setTimeout(() => {
                window.location.href = '/login';
            }, 2000);
        } else {
            feedback.textContent = result.error || 'Failed to update password.';
            feedback.style.color = 'red';
        }
    } catch (error) {
        console.error('Error:', error);
        feedback.textContent = 'Something went wrong. Please try again.';
        feedback.style.color = 'red';
    }
});
