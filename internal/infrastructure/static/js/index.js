const listBtn = document.getElementById('list-btn');
const videoList = document.getElementById('video-list');

videoList.style.width = `${videoPlayer.clientWidth}px`;

listBtn.addEventListener('click', function () {
    loadVideoList();
});

listBtn.addEventListener('click', function (event) {
    event.stopPropagation(); // Предотвращение всплытия события до document
    videoList.style.display = (videoList.style.display === 'block') ? 'none' : 'block';
});

document.addEventListener('click', function (event) {
    if (event.target !== listBtn) {
        videoList.style.display = 'none';
    }
});

function loadVideoList() {
    const page = 1;
    const limit = 25;

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
                    videoList.appendChild(listItem);
                });

                ul.appendChild(videoList);

                const paginationInfo = document.querySelector('.dropdown-content .pagination-info');
                paginationInfo.textContent = `Page ${data.pagination.page} of ${data.pagination.total}`;
            } else {
                ul.textContent = 'There are not available videos';
            }
        })
        .catch(error => {
            console.error('Ошибка при загрузке списка видео:', error);
            ul.textContent = 'Sorry, there is error occurred while loading a video list';
        });
}
