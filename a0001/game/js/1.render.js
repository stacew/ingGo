function fClearCanvas() {
    msgCtx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.clearRect(0, 0, canvas.width, canvas.height);
}

function fDrawCircle(x, y, r, color) {
    ctx.beginPath();
    ctx.arc(x, y, r, 0, Math.PI * 2);
    ctx.strokeStyle = color;
    ctx.lineWidth = 7;
    ctx.stroke();
    ctx.closePath();
}
function fDrawPlayer(img, x, y, nRadius) {
    ctx.drawImage(img, x - nRadius, y - nRadius, nRadius * 2, nRadius * 2);
}
function fDrawArrowLine(x1, y1, x2, y2, color) {
    ctx.beginPath();
    ctx.moveTo(x1, y1);
    ctx.lineTo(x2, y2);
    ctx.strokeStyle = color;
    ctx.lineWidth = 7;
    ctx.stroke();
    ctx.closePath();
}

function fDrawText(text, color) {
    msgCtx.clearRect(0, 0, canvas.width, canvas.height);
    msgCtx.font = '240px serif';
    msgCtx.textBaseline = 'middle';
    msgCtx.textAlign = 'center';
    msgCtx.fillStyle = color;
    msgCtx.fillText(text, 500, 500);
}
//////
function fInfoRoom(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strPlayerCount = multi[0];
    multi = msgToken([msg, multi[1]]);
    let strRoomCapacity = multi[0];

    fDrawText(strPlayerCount + " / " + strRoomCapacity, "white");

    return multi[1];
}

function fOneshotStartEnd(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let charSE = multi[0];

    zbShot = zbLive ? ((charSE == 's') ? true : false) : false;

    return multi[1];
}
function fClientTimer(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strTime = multi[0];

    fDrawText(strTime, "red");

    return multi[1];
}
function fDieMessage(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strCID = multi[0];

    if (strCID == zio.io.engine.id) { fDie(); }

    return multi[1];
}
function fPlaying(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strCID = multi[0];
    multi = msgToken([msg, multi[1]]);
    let nX = parseInt(multi[0]);
    multi = msgToken([msg, multi[1]]);
    let nY = parseInt(multi[0]);
    multi = msgToken([msg, multi[1]]);
    let nR = parseInt(multi[0]);
    multi = msgToken([msg, multi[1]]);
    let charBW = multi[0];

    let circleColor;
    if (charBW == 'b') {
        fDrawPlayer(imgBlack, nX, nY, nR);
        fDrawCircle(nX, nY, nR, "red");
    } else {
        fDrawPlayer(imgWhite, nX, nY, nR);
    }

    if (strCID == zio.io.engine.id) {
        znMyX = nX;
        znMyY = nY;
        fDrawCircle(nX, nY, nR + 10, "black");
    }

    return multi[1];
}
function fAttackTeamTurn(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let charTeam = multi[0];

    zbAttackTeamBlack = (charTeam == 'b') ? true : false;

    return multi[1];
}
