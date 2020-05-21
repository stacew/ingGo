function msgToken(multi) {
    let nStart = multi[1] + 1;
    multi[1] = multi[0].indexOf(",", nStart);
    multi[0] = multi[0].substring(nStart, multi[1]);
    return [multi[0], multi[1]];
}

function fMsgDecoder(msg) {
    let bNeedCanvasClear = true;
    let nFindIndex = 0;
    while ((nFindIndex = msg.indexOf(".", nFindIndex)) > -1) {
        let charType = msg.charAt(nFindIndex + 1);
        switch (charType) {
            case "i":
                nFindIndex = fInfoRoom(msg, nFindIndex + 1);
                break;
            case "s":
                nFindIndex = fStarting(msg, nFindIndex + 1);
                zbLive = true;
                break;
            case "o":
                nFindIndex = fOneshotStartEnd(msg, nFindIndex + 1);
                break;
            case "t":
                nFindIndex = fClientTimer(msg, nFindIndex + 1);
                break;
            case "d":
                nFindIndex = fDieMessage(msg, nFindIndex + 1);
                break;
            case "p":
                if (bNeedCanvasClear) { fClearCanvas(); bNeedCanvasClear = false; }
                nFindIndex = fPlaying(msg, nFindIndex + 1);
                break;
            default:
                break;
        }
    }
}