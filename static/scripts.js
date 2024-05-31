document.getElementById('getWeatherBtn').addEventListener('click', function() {
    const city = document.getElementById('cityInput').value;
    if (city) {
        fetch(`/weather/${city}`)
            .then(response => response.json())
            .then(data => {
                document.getElementById('cityName').innerText = data.name;
                document.getElementById('temperature').innerText = data.main.temp;
                document.getElementById('humidity').innerText = data.main.humidity;
                document.getElementById('weatherResult').classList.remove('hidden');
            })
            .catch(error => {
                alert('Error fetching weather data');
                console.error('Error:', error);
            });
    } else {
        alert('Please enter a city name');
    }
});
