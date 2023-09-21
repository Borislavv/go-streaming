const container = document.querySelector('.container');
const listBtn = document.getElementById('list-btn');
const videoList = document.getElementById('video-list');

videoList.style.width = `${videoPlayer.clientWidth}px`;

listBtn.addEventListener('click', function () {
    loadVideoList();
});

listBtn.addEventListener('click', function (event) {
    event.stopPropagation(); // Предотвращение всплытия события до document
    videoList.style.display = (videoList.style.display === 'block') ? 'none' : 'block';
    videoList.classList.toggle('active');
    container.classList.toggle('with-list');
});

document.addEventListener('click', function (event) {
    if (event.target !== listBtn) {
        videoList.style.display = 'none';
    }
});

function loadVideoList(page = 1, limit = 25) {
    // clear the previous video list
    const ul = document.querySelector('.dropdown-content ul');
    ul.innerHTML = '';

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
                    listItem.id = video.id
                    videoList.appendChild(listItem);
                });

                ul.appendChild(videoList);

                const paginationInfo = document.querySelector('.dropdown-content .pagination-info');

                // Displaying pagination info
                const currentPage = data.pagination.page;
                const totalPages = data.pagination.total;

                paginationInfo.textContent = `Page ${currentPage} of ${totalPages}`;
            } else {
                ul.textContent = 'There are no available videos';
            }
        })
        .catch(error => {
            console.error('Error occurred while loading a video list:', error);
            ul.textContent = 'Sorry, there is an error occurred while loading a video list';
        });
}

// Обработчик выбора лимита элементов
const limitSelect = document.getElementById('limit-select');
limitSelect.addEventListener('change', () => {
    const selectedLimit = parseInt(limitSelect.value, 10);
    loadVideoList(1, selectedLimit);
});

// Загрузка списка с лимитом по умолчанию при загрузке страницы
loadVideoList();