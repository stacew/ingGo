
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

function fDie() {
    zbLive = false;
    fLeaveShow(true);
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
        fDie();
    });
    zio.on('sGame', function (msg) {
        fMsgDecoder(msg);
    });
}

//expand
// zio.emit('cReqMsg', function (data) {
// });

function clickCanvas(e) {
    if (zbShot == false || zbLive == false)
        return;
    zbShot = false;
    multi = fGetClickPos(e);
    fDrawArrowLine(znMyX, znMyY, multi[0], multi[1], "black");
    zio.emit('cShot', parseInt(multi[0]).toString() + "," + parseInt(multi[1]).toString());
}
canvas.addEventListener("click", clickCanvas, false);