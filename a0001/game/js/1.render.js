function fClearCanvas() {
    msgCtx.clearRect(0, 0, zCanvasWH, zCanvasWH);
    fgCtx.clearRect(0, 0, zCanvasWH, zCanvasWH);
}
function fDrawPlayer(img, x, y, nRadius) {
    fgCtx.drawImage(img, x - nRadius, y - nRadius, nRadius * 2, nRadius * 2);
}
function fDrawCircle(anyCtx, x, y, r, color, style = []) {
    anyCtx.beginPath();
    anyCtx.setLineDash(style);
    anyCtx.arc(x, y, r, 0, Math.PI * 2);
    anyCtx.strokeStyle = color;
    anyCtx.lineWidth = 7;
    anyCtx.stroke();
    anyCtx.setLineDash([]);
    anyCtx.closePath();
}
function fDrawArrowLine(anyCtx, x1, y1, x2, y2, color) {
    anyCtx.beginPath();
    anyCtx.setLineDash([50, 10]);
    anyCtx.moveTo(x1, y1);
    anyCtx.lineTo(x2, y2);
    anyCtx.strokeStyle = color;
    anyCtx.lineWidth = 7;
    anyCtx.stroke();
    anyCtx.setLineDash([]);
    anyCtx.closePath();
}
function fDrawText(text, color) {
    msgCtx.clearRect(0, 0, zCanvasWH, zCanvasWH);
    msgCtx.font = '240px serif';
    msgCtx.textBaseline = 'middle';
    msgCtx.textAlign = 'center';
    msgCtx.fillStyle = color;
    msgCtx.fillText(text, 500, 500);
}
//////
var zDashList = [40, 10];
var zDashOffset = 0;
function fDrawAnim() {
    animCtx.clearRect(0, 0, zCanvasWH, zCanvasWH);
    zDashOffset--;
    animCtx.lineDashOffset = zDashOffset;
    if (zbPlaying) {
        fDrawCircle(animCtx, znMyX, znMyY, znMyR + 10, "gray", zDashList);
        if (zbLive) {
            fDrawArrowLine(animCtx, znMyX, znMyY, znShotX, znShotY, "gray");
        }
    }
    setTimeout(fDrawAnim, 10);
}
fDrawAnim();
//////
function fInfoRoom(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strPlayerCount = multi[0]; multi = msgToken([msg, multi[1]]);
    let strRoomCapacity = multi[0];

    fDrawText(strPlayerCount + " / " + strRoomCapacity, "white");
    return multi[1];
}

function fOneshotStartEnd(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let charSE = multi[0];

    if (charSE == 's') {
        audioTurn.play();
        zbShotChance = zbLive ? true : false;
    }
    else {
        audioTong.play();
        zbShotChance = false;
    }
    return multi[1];
}
function fClientTimer(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strTime = multi[0];

    fDrawText(strTime, "red");
    return multi[1];
}
function fPlaying(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let strCID = multi[0]; multi = msgToken([msg, multi[1]]);
    let nX = parseInt(multi[0]); multi = msgToken([msg, multi[1]]);
    let nY = parseInt(multi[0]); multi = msgToken([msg, multi[1]]);
    let nR = parseInt(multi[0]); multi = msgToken([msg, multi[1]]);
    let charLive = multi[0]; multi = msgToken([msg, multi[1]]);
    let charBlack = multi[0];

    let bLive = (charLive == 'l');
    if (strCID == zio.io.engine.id) {
        znMyX = nX; znMyY = nY; znMyR = nR;
        znShotX = nX; znShotY = nY;
        zbLive = bLive;
    }

    if (charBlack == 'b') {
        fDrawPlayer(bLive ? imgBlack : imgBlackDie, nX, nY, nR);
        if (bLive && zbAttackTeamBlack) fDrawCircle(fgCtx, nX, nY, nR, "red");
    } else {
        fDrawPlayer(bLive ? imgWhite : imgWhiteDie, nX, nY, nR);
        if (bLive && zbAttackTeamBlack == false) fDrawCircle(fgCtx, nX, nY, nR, "red");
    }
    return multi[1];
}
function fAttackTeamTurn(msg, nFindIndex) {
    let multi = msgToken([msg, nFindIndex]);
    let charTeam = multi[0];

    zbAttackTeamBlack = (charTeam == 'b') ? true : false;
    return multi[1];
}
