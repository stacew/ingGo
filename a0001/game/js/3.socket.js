
connectBtn.onclick = fConnectClick;
joinBtn.onclick = fJoinClick;
LeaveBtn.onclick = fLeaveClick;
function fConnectShow(bool) { connectBtn.style.display = bool ? "block" : "none"; }
function fJoinShow(bool) { joinBtn.style.display = bool ? "block" : "none"; }
function fLeaveShow(bool) { LeaveBtn.style.display = bool ? "block" : "none"; }
function fConnectClick() {
    if (zio == null || zio.io.engine.id == null) {
        fRegisterReceiver();
    }

    fConnectShow(false);
    fJoinShow(true);
}

function fJoinClick() {
    if (zio.io.engine.id == null) {
        fConnectShow(true);
        fJoinShow(false);
    }
    zio.emit('cJoin');
    fJoinShow(false);
}
function fLeaveClick() {
    fClearCanvas();
    zio.emit('cLeave');
    fLeaveShow(false);
    fJoinShow(true);
}

function fRegisterReceiver() {
    zio = io('/');
    zio.on('sOver', function (msg) {
        zbPlaying = false;
        fLeaveShow(true);
    });
    zio.on('sDecoder', function (msg) {
        fMsgDecoder(msg);
    });
}

//expand
// zio.emit('cReqMsg', function (data) {
// });

function clickCanvas(e) {
    if (zbShotChance == false || zbLive == false)
        return;
    zbShotChance = false;
    multi = fGetClickPos(e);
    znShotX = multi[0];
    znShotY = multi[1];
    zio.emit('cShot', parseInt(multi[0]).toString() + "," + parseInt(multi[1]).toString());
}
fgCanvas.addEventListener("click", clickCanvas, false);