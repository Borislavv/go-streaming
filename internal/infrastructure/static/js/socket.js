const socket = new WebSocket('ws://0.0.0.0:9988/');

const videoPlayer = document.getElementById('videoPlayer');
const nextBtn = document.getElementById('next-btn');
const prevBtn = document.getElementById('prev-btn');

socket.binaryType = 'arraybuffer';

let buffer;
let mediaSource;
let chunks;
let mediaSourceReady;

// Socket events further

socket.onopen = (event) => {
    console.log('WebSocket connection opened');
    socket.send("start")
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

    if (data === 'stop') {
        closeMediaResource()
    }

    if (typeof data === 'string' && data.startsWith('start')) {
        console.log("Starting new video...")
        let dataParts = data.split(':')
        console.log(dataParts)
        makeMediaResource(dataParts[1], dataParts[2])
    }
};

// todo need to make a ticker which will count the awaiting time if it's more than N then send decrease buffer action
videoPlayer.addEventListener('waiting', function() {
    console.warn('Video playback is waiting for data (buffering)');

    // requesting decrease the buffer
    socket.send("decrBuff")
});

// previous video handler
nextBtn.addEventListener('click', function() {
    // requesting new video/audio from server
    socket.send("next")
});

// next video handler
prevBtn.addEventListener('click', function() {
    // requesting the prev video/audio from server
    socket.send("prev")
});

let prevID = '';
document.addEventListener('click', function (event) {
    let ul = document.querySelector('.video-list');
    if (ul !== null) {
        let lis = ul.getElementsByTagName('li');
        Array.from(lis).forEach(function (el) {
            if (prevID === '') {
                prevID = el.id
            }
            if (event.target === el && prevID !== el.id) {
                // requesting the prev video/audio from server
                socket.send("nextID:" + el.id)
                prevID = el.id;
            }
        });
    }
});

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
        console.log(
            "Awaiting data: " +
            "media resource is not ready yet -", mediaSourceReady, ", " +
            "buffer is still updating -", buffer.updating,", " +
            "chunks awaiting -", chunks.length === 0,
            "(len:", chunks.length, ")"
        )
        await new Promise(r => setTimeout(r, 1));
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
                    console.error("Stream closing failed", error, mediaSource.readyState)
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
        console.log(
            "Awaiting closing: " +
            "buffer is still updating -", buffer.updating,", " +
            "chunks is not empty -", chunks.length > 0,
                "(len:", chunks.length, ")"
        )
        await new Promise(r => setTimeout(r, 1));
    }
}

function makeMediaResource(audioCodec, videoCodec) {
    mediaSource = new MediaSource();
    mediaSourceReady = false;
    videoPlayer.src = URL.createObjectURL(mediaSource);
    chunks = [];

    mediaSource.addEventListener('sourceclose', function (e) {
        console.log("MediaSource sourceopen event is closed", e)
    });

    mediaSource.addEventListener('sourceopen', function () {
        console.log("MediaSource sourceopen event is open");

        try {
            let codecsMap = []
            codecsMap['avc1'] = 'avc1.42E01E'
            codecsMap['mp4a'] = 'mp4a.40.2'

            let codecsStr = ''
            if (videoCodec !== '' && codecsMap[videoCodec]) {
                codecsStr += codecsMap[videoCodec]
            }
            if (audioCodec !== '' && codecsMap[audioCodec]) {
                if (videoCodec !== '') {
                    codecsStr += ', '
                }
                codecsStr += codecsMap[audioCodec]
            }
            if (codecsStr === '') {
                console.error('Codecs string a empty! Unable to play video!')
            }

            const codec = 'video/mp4; codecs="' + codecsStr +'"';

            console.log("CODEC: ", codec)

            if (MediaSource.isTypeSupported(codec)) {
                console.log("Type of codec is supported: ", codec);

                if (mediaSource.sourceBuffers.length > 0) {
                    let b = mediaSource.sourceBuffers.length
                    while(b > 1) {
                        mediaSource.removeSourceBuffer(mediaSource.sourceBuffers[b]);
                        b--
                    }
                }

                buffer = mediaSource.addSourceBuffer(codec);
            } else {
                console.log("Type of codec is NOT supported: ", codec);
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
}

