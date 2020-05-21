//If you touch the js file, the browser may get stuck
//All flow of the game is controlled by the server.
//So if you change this, it's just an optical illusion effect.
//G.G.
let imgBlack = new Image();
imgBlack.src = "/asset/black.png"
let imgWhite = new Image();
imgWhite.src = "/asset/white.png"

const bgCanvas = document.getElementById("bgCanvas");
const bgCtx = bgCanvas.getContext("2d");
function fDrawBoard() {
    bgCtx.fillStyle = "burlywood";
    bgCtx.fillRect(0, 0, bgCanvas.width, bgCanvas.height);
    let nWidthBand = bgCanvas.width / 20;
    let nHeightBand = bgCanvas.height / 20;

    bgCtx.strokeStyle = "black";
    bgCtx.lineWidth = 1;
    for (let i = 1; i < 20; i++) {
        bgCtx.strokeRect(nWidthBand * i, 0, 0, bgCanvas.height);
        bgCtx.strokeRect(0, nHeightBand * i, bgCanvas.width, 0);
    }

    bgCtx.fillStyle = "black";
    for (let i = 0; i < 3; i++) {
        let x = nWidthBand * 3 + i * (nWidthBand * 7);
        for (let j = 0; j < 3; j++) {
            bgCtx.beginPath();
            let y = nHeightBand * 3 + j * (nHeightBand * 7);
            bgCtx.arc(x, y, 5, 0, Math.PI * 2);
            bgCtx.fill();
            bgCtx.closePath();
        }
    }
}
fDrawBoard();

document.addEventListener("contextmenu", handleContextualMenu);
function handleContextualMenu(event) {
    event.preventDefault();
}

const canvas = document.getElementById("fgCanvas");
const ctx = canvas.getContext("2d");
const connectBtn = document.getElementById("connectBtn");
const joinBtn = document.getElementById("joinBtn");
const LeaveBtn = document.getElementById("leaveBtn");


const test2 = document.getElementById("test2");

let zbShot = false;
let zbLive = false;





