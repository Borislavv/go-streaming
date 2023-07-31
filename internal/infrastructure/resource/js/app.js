const videoPlayer = document.getElementById('videoPlayer');
const playButton = document.getElementById('playButton');
const socket = new WebSocket('ws://127.0.0.1:9988/'); // Replace with the actual address of your WebSocket server

socket.binaryType = 'arraybuffer';

let mediaSource = new MediaSource();
let buffer;
let chunks = [];
let mediaSourceReady = false;
videoPlayer.src = URL.createObjectURL(mediaSource);

// MediaSource events further

mediaSource.addEventListener('sourceopen', function () {
    console.log("MediaSource sourceopen event is open");

    try {
        const codec = 'video/mp4; codecs="avc1.42E01E"';

        if (MediaSource.isTypeSupported(codec)) {
            console.log("Type of codec is supported");
            buffer = mediaSource.addSourceBuffer(codec);
        } else {
            console.log("Type of codec is NOT supported");
        }

        console.log("Source buffer added");
    } catch (e) {
        console.error("Error adding source buffer:", e);
    }

    buffer.addEventListener('error', function (e) {
        console.error('Buffer error', e)
    });

    mediaSourceReady = true;
}, false);

mediaSource.addEventListener('sourceclose', function (e) {
    console.log("MediaSource sourceopen event is closed", e)
});

// Socket events further

socket.onopen = (event) => {
    console.log('WebSocket connection opened');
};

socket.onclose = (event) => {
    console.log('WebSocket connection closed');
    closeMediaResource()
};

socket.onerror = (event) => {
    console.error('WebSocket error: ', event);
};

socket.onmessage = (event) => {
    const data = event.data;
    if (data instanceof ArrayBuffer) {
        console.log("Chunk received");
        chunks.push(data)
        addNextChunk()
    }
};

function addNextChunk() {
    awaiting()
        .then(
            function () {
                try {
                    buffer.appendBuffer(chunks.shift())
                    console.log("Chunk successfully added to buffer")
                } catch (error) {
                    console.error("Chunk adding to buffer filed", error)
                }
            }
        ).catch(
        function (e) {
            console.error("unable to awaiting", e)
        }
    )
}

async function awaiting() {
    while(!mediaSourceReady || buffer.updating || chunks.length === 0) {
        console.log("awaiting...", mediaSourceReady, !buffer.updating, chunks.length !== 0)
        await new Promise(r => setTimeout(r, 250));
    }
}

function closeMediaResource() {
    awaitingClose()
        .then(
            function () {
                try {
                    mediaSource.endOfStream()
                    console.log("Stream is successfully closed")
                } catch (error) {
                    console.error("Stream closing failed")
                }
            }
        ).catch(
        function (e) {
            console.error("unable to awaiting closing", e)
        }
    )
}

async function awaitingClose() {
    while(buffer.updating || chunks.length > 0) {
        console.log("awaiting closing...", mediaSourceReady, !buffer.updating, chunks.length !== 0)
        await new Promise(r => setTimeout(r, 250));
    }
}