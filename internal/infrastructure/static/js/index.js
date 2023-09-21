const container = document.querySelector('.container');
const listBtn = document.getElementById('list-btn');
const videoList = document.getElementById('video-list');

let prevPage = 1;
let prevLimit = 25;

let currentPage = 1;
let currentLimit = 25;

videoList.style.width = `${videoPlayer.clientWidth}px`;

listBtn.addEventListener('click', function (event) {
    event.stopPropagation();
    videoList.style.display = (videoList.style.display === 'block') ? 'none' : 'block';
    videoList.classList.toggle('active');
    container.classList.toggle('with-list');
});

const paginationInfo = document.querySelector('.pagination-info')
const paginationControl = document.querySelector('.pagination-control')
document.addEventListener('click', function (event) {
    let matched = false
    let ul = document.querySelector('.video-list');
    if (ul !== null) {
        let lis = ul.getElementsByTagName('li');
        Array.from(lis).forEach(function (el) {
            if (event.target === el) {
                matched = true
            }
        });
    }

    if (
        (event.target !== videoList &&
        event.target !== paginationInfo &&
        event.target !== paginationControl &&
        event.target !== limitSelect &&
        event.target !== pageSelect &&
        event.target !== reqBtn) &&
        matched === false
    ) {
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
                    listItem.className = 'list-item';
                    listItem.textContent = video.name;
                    listItem.id = video.id
                    videoList.appendChild(listItem);
                });

                ul.appendChild(videoList);

                const paginationInfo = document.querySelector('.dropdown-content .pagination-info');
                paginationInfo.textContent = `Page ${currentPage} of ${Math.ceil(data.pagination.total / currentLimit)}`;
            } else {
                ul.textContent = 'There are no available videos';
            }
        })
        .catch(error => {
            console.error('Error occurred while loading a video list:', error);
            ul.textContent = 'Sorry, there is an error occurred while loading a video list';
        });
}

// handling limit 'select' box
const limitSelect = document.getElementById('limit-select');
limitSelect.addEventListener('change', () => {
    prevLimit    = currentLimit
    currentLimit = parseInt(limitSelect.value, 10);
});

// handling page 'select' box
const pageSelect = document.getElementById('page-select');
pageSelect.addEventListener('change', () => {
    prevPage    = currentPage
    currentPage = parseInt(pageSelect.value, 10);
});

// handling list request btn
const reqBtn = document.getElementById('request-btn');
reqBtn.addEventListener('click', function () {
    if (currentPage !== prevPage || currentLimit !== prevLimit) {
        loadVideoList(currentPage, currentLimit);
    }
});

loadVideoList(currentPage, currentLimit);