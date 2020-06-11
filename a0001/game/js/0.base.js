//G.G.
document.addEventListener("contextmenu", handleContextualMenu);
function handleContextualMenu(event) {
    event.preventDefault();
}

const imgBlack = new Image();
const imgWhite = new Image();
const imgBlackDie = new Image();
const imgWhiteDie = new Image();
imgBlack.src = "/game/asset/black.png";
imgWhite.src = "/game/asset/white.png";
imgBlackDie.src = "/game/asset/black_die.png";
imgWhiteDie.src = "/game/asset/white_die.png";
const audioTurn = new Audio("/game/asset/turn.mp3");
const audioTong = new Audio("/game/asset/tong.mp3");

const zCanvasWH = 1000;
const baseCanvas = document.getElementById("baseCanvas");
const bgCanvas = document.getElementById("bgCanvas");
const msgCanvas = document.getElementById("msgCanvas");
const fgCanvas = document.getElementById("fgCanvas");
const animCanvas = document.getElementById("animCanvas");
const baseCtx = baseCanvas.getContext("2d");
const bgCtx = bgCanvas.getContext("2d");
const msgCtx = msgCanvas.getContext("2d");
const fgCtx = fgCanvas.getContext("2d");
const animCtx = animCanvas.getContext("2d");
function fDrawBoard() {
    baseCtx.fillStyle = "burlywood";
    baseCtx.fillRect(0, 0, zCanvasWH, zCanvasWH);
    let nWidthBand = zCanvasWH / 20;
    let nHeightBand = zCanvasWH / 20;

    baseCtx.strokeStyle = "black";
    baseCtx.lineWidth = 1;
    for (let i = 1; i < 20; i++) {
        baseCtx.strokeRect(nWidthBand * i, 0, 0, zCanvasWH);
        baseCtx.strokeRect(0, nHeightBand * i, zCanvasWH, 0);
    }

    baseCtx.fillStyle = "black";
    for (let i = 0; i < 3; i++) {
        let x = nWidthBand * 3 + i * (nWidthBand * 7);
        for (let j = 0; j < 3; j++) {
            baseCtx.beginPath();
            let y = nHeightBand * 3 + j * (nHeightBand * 7);
            baseCtx.arc(x, y, 5, 0, Math.PI * 2);
            baseCtx.fill();
            baseCtx.closePath();
        }
    }
}
fDrawBoard();

function fGetClickPos(e) {
    let rect = fgCanvas.getBoundingClientRect();
    return [e.layerX / rect.width * zCanvasWH, e.layerY / rect.height * zCanvasWH];
}
const connectBtn = document.getElementById("connectBtn");
const joinBtn = document.getElementById("joinBtn");
const LeaveBtn = document.getElementById("leaveBtn");
//io
let zio = null;
//render attack color
let zbAttackTeamBlack = true;
//for click arrow
let zbPlaying = false;
let zbLive = false;
let zbShotChance = false;
let znMyX = 0;
let znMyY = 0;
let znR = 0;
let znShotX = 0;
let znShotY = 0;