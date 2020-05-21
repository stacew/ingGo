function fClearCanvas() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
}

function fDrawCircle(multi) {
    console.log(multi);
    ctx.beginPath();
    ctx.arc(multi[0], multi[1], multi[2], 0, Math.PI * 2);
    ctx.strokeStyle = "red";
    ctx.lineWidth = 7;
    ctx.stroke();
    ctx.closePath();
}
function fDrawPlayer(img, x, y, nRadius) {
    ctx.drawImage(img,
        x - nRadius, y - nRadius,
        nRadius * 2, nRadius * 2);
}
//////
function fInfoRoom(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strPlayerCount = multi[0];
    multi = msgToken([msg, multi[1]]);
    let strRoomCapacity = multi[0];

    test2.textContent = strPlayerCount + "/" + strRoomCapacity;

    return multi[1];
}
function fStarting(msg, nFindIndex) {
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
    fDrawPlayer((charBW == 'b') ? imgBlack : imgWhite, nX, nY, nR);
    if (strCID == sio.io.engine.id) { fDrawCircle([nX, nY, nR]); }
    return multi[1];
}
function fOneshotStartEnd(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let charSE = multi[0];
    if (zbLive) {
        zbShot = (charSE == 's' ? true : false);
    }
    else {
        zbShot = false;
    }
    return multi[1];
}
function fClientTimer(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strTime = multi[0];
    test2.textContent = strTime;
    return multi[1];
}
function fDieMessage(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strCID = multi[0];
    if (strCID == sio.io.engine.id) { fDie(); }
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
    fDrawPlayer((charBW == 'b') ? imgBlack : imgWhite, nX, nY, nR);
    if (strCID == sio.io.engine.id) { fDrawCircle([nX, nY, nR]); }
    return multi[1];
}
