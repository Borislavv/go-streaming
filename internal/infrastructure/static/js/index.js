const listBtn = document.getElementById('list-btn');
const listDdwn = document.getElementById('video-list');

listBtn.addEventListener('click', function () {
    // Загрузка списка видео через HTTP запрос
    loadVideoList();
});

// Функция для загрузки списка видео
function loadVideoList() {
    const limit = 25; // Количество видео на странице
    const page = 1; // Номер страницы (начнем с первой)

    // Очищаем предыдущий список видео
    listDdwn.innerHTML = '';

    // Отправляем HTTP GET запрос для получения списка видео
    fetch(`http://0.0.0.0:8000/api/v1/video?limit=${limit}&page=${page}`)
        .then(response => response.json())
        .then(data => {
            data = data.data

            if (data.list && data.list.length > 0) {
                const videoList = document.createElement('ul');
                videoList.className = 'video-list';

                data.list.forEach(video => {
                    const listItem = document.createElement('li');
                    listItem.textContent = video.name;
                    videoList.appendChild(listItem);
                });

                listDdwn.appendChild(videoList);

                const paginationInfo = document.createElement('div');
                paginationInfo.textContent = `Страница ${data.pagination.page} из ${data.pagination.total}`;
                listDdwn.appendChild(paginationInfo);
            } else {
                listDdwn.textContent = 'Нет доступных видео.';
            }
        })
        .catch(error => {
            console.error('Ошибка при загрузке списка видео:', error);
            listDdwn.textContent = 'Ошибка при загрузке списка видео.';
        });
}