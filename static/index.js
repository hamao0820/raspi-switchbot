main = () => {
    $button = document.getElementById('turnOnButton');
    $button.addEventListener('click', async () => {
        try {
            const response = await fetch('/api/turn_on', {method: 'POST'});
            if (response.ok) {
                alert('Turned on the SwitchBot!');
            } else {
                alert('Failed to turn on the SwitchBot.');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('An error occurred while trying to turn on the SwitchBot.');
        }
    });
}

main()
