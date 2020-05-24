//G.G.
document.addEventListener("contextmenu", handleContextualMenu);
function handleContextualMenu(event) {
    event.preventDefault();
}

let imgBlack = new Image();
imgBlack.src = "/game/asset/black.png"
let imgWhite = new Image();
imgWhite.src = "/game//asset/white.png"

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
const msgCanvas = document.getElementById("msgCanvas");
const msgCtx = msgCanvas.getContext("2d");
const canvas = document.getElementById("fgCanvas");
function fGetClickPos(e) {
    let rect = canvas.getBoundingClientRect();
    return [e.layerX / rect.width * canvas.width, e.layerY / rect.height * canvas.height];
}

const ctx = canvas.getContext("2d");
const connectBtn = document.getElementById("connectBtn");
const joinBtn = document.getElementById("joinBtn");
const LeaveBtn = document.getElementById("leaveBtn");
const top2 = document.getElementById("top2");
//io
let zio = null;
//render attack color
let zbAttackTeamBlack = true;
//for click arrow
let zbLive = false;
let zbShot = false;
let znMyX = 0;
let znMyY = 0;







