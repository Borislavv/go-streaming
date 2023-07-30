const videoPlayer = document.getElementById('videoPlayer');
const playButton = document.getElementById('playButton');
const socket = new WebSocket('ws://127.0.0.1:9988/'); // Replace with the actual address of your WebSocket server

socket.binaryType = 'arraybuffer';

let mediaSource = new MediaSource();
let buffer;
let chunks = [];
let isPlaying = false;
let mediaSourceReady = false;
videoPlayer.src = URL.createObjectURL(mediaSource);

// MediaSource events further

mediaSource.addEventListener('sourceopen', function () {
    console.log("MediaSource sourceopen event is open!");

    try {
        const codec = 'video/mp4; codecs="avc1.42E01E"';

        if (MediaSource.isTypeSupported(codec)) {
            console.log("Type of codec is supported!");
            buffer = mediaSource.addSourceBuffer(codec);
        } else {
            console.log("Type of codec is NOT supported!");
        }

        console.log("Source buffer added:", buffer);
    } catch (e) {
        console.error("Error adding source buffer:", e);
    }

    buffer.addEventListener('error', function (e) {
        console.error('Buffer error', e)
    });

    buffer.addEventListener('updateend', function () {
        addNextChunk();
    });

    buffer.addEventListener('ended', function () {
        addNextChunk();
    });

    mediaSourceReady = true;
}, false);

mediaSource.addEventListener('sourceclose', function (e) {
    console.log("MediaSource sourceopen event is closed!", e)
});

// Socket events further

socket.onopen = (event) => {
    console.log('WebSocket connection opened');
};

socket.onclose = (event) => {
    console.log('WebSocket connection closed');
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
    if (chunks.length > 0) {

        awaiting()
            .then(
                function () {
                    buffer.appendBuffer(chunks.shift())
                }
            ).catch(
            function (e) {
                console.error("unable to awaiting", e)
            }
        )

        // while ((chunk = chunks.shift()) !== 'undefined') {
        //     awaiting()
        //         .then(
        //             function () {
        //                 buffer.appendBuffer(chunk)
        //             }
        //         ).catch(
        //             function (e) {
        //                 console.error("unable to awaiting", e)
        //             }
        //     )
        // }

        // try {
        //     console.log("Chunks length: ", chunks.length)
        //     buffer.appendBuffer(chunk)
        //     console.log("Chunk successfully added to buffer", buffer)
        // } catch (error) {
        //     console.error("Chunk adding to buffer filed", error)
        // }
    }
}

videoPlayer.addEventListener('timeupdate', function () {
    if (videoPlayer.readyState >= 2) {
        addNextChunk()
    }
});

async function awaiting() {
    while(!mediaSourceReady || buffer.updating || chunks.length === 0) {}
}