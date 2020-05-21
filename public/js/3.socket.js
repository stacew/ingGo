let sio = io('/');
connectBtn.onclick = fConnectClick;
joinBtn.onclick = fJoinClick;
LeaveBtn.onclick = fLeaveClick;
function fConnectShow(bool) { connectBtn.style.display = bool ? "block" : "none"; }
function fJoinShow(bool) { joinBtn.style.display = bool ? "block" : "none"; }
function fLeaveShow(bool) { LeaveBtn.style.display = bool ? "block" : "none"; }
function fConnectClick() {
    if (sio.io.engine.id == null) {
        sio = io('/');
        fConnectShow(true);
        fJoinShow(false);
    } else {
        fConnectShow(false);
        fJoinShow(true);
        //expand
        // sio.emit('cCurrentCon', function (data) {
        // });
    }
}
function fDie() {
    zbLive = false;
    fLeaveShow(true);
}
function fJoinClick() {
    sio.emit('cJoin');
    fJoinShow(false);
}
function fLeaveClick() {
    fClearCanvas();
    sio.emit('cLeave');
    fLeaveShow(false);
    fJoinShow(true);
}
sio.on('sOver', function (msg) {
    fDie();
});
sio.on('sGame', function (msg) {
    fMsgDecoder(msg);
});

function clickCanvas(e) {
    if (zbShot == false || zbLive == false)
        return;
    zbShot = false;
    multi = fGetClickPos(e);
    fDrawArrowLine(znMyX, znMyY, multi[0], multi[1], "sienna");
    sio.emit('cShot', parseInt(multi[0]).toString() + "," + parseInt(multi[1]).toString());
}
canvas.addEventListener("click", clickCanvas, false);