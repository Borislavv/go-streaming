const videoPlayer = document.getElementById('videoPlayer');
const socket = new WebSocket('ws://127.0.0.1:9988/ws'); // Замените на реальный адрес вашего сокет-сервера

socket.binaryType = 'arraybuffer';

const chunks = []; // Массив для хранения полученных частей видео

socket.onopen = (event) => {
    console.log('WebSocket connection opened');
};

socket.onmessage = (event) => {
    const data = event.data;
    if (data instanceof ArrayBuffer) {
        console.log("Chunk received");
        // Добавляем полученную часть видео в массив
        chunks.push(data);
    }
};

socket.onclose = (event) => {
    console.log('WebSocket connection closed: ' + event.reason +', '+ event.code +', '+ event.type +', '+ event.eventPhase);
};

socket.onerror = (event) => {
    console.error('WebSocket error:', event);
};

function playVideo() {
    if (chunks.length === 0) {
        console.log("No video chunks received");
        return;
    }

    // Склеиваем все части видео в один ArrayBuffer
    const fullVideoBuffer = new Uint8Array(chunks.reduce((acc, chunk) => acc + chunk.byteLength, 0));
    let offset = 0;
    chunks.forEach(chunk => {
        fullVideoBuffer.set(new Uint8Array(chunk), offset);
        offset += chunk.byteLength;
    });

    const videoBlob = new Blob([fullVideoBuffer], { type: 'video/mp4' }); // Предполагается, что данные в формате MP4
    const videoURL = URL.createObjectURL(videoBlob);
    videoPlayer.src = videoURL;

    // Воспроизводим видео только после действия пользователя (например, клика)
    videoPlayer.addEventListener("click", function() {
        console.log("Video playback started on click");
        videoPlayer.play();
    });
}

// Воспроизведение видео начинается после клика по видеоплееру
videoPlayer.addEventListener("click", function() {
    playVideo();
});
