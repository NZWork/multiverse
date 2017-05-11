// OT
var Changeset = ot.Changeset
var engine = new diff_match_patch
var his = ""
var seq = 0
var ver = 0
var otLock = false

// WS
var conn = null
var target = "wss://" + location.host + "/sync"
var uid = 0

// Switcher
var send = 0
var lastMsg = 0
var didClose = false
var typing = false

var ChangesetQueueLock = false
var ChangesetQueue = []

const OP_RETAIN = 0
const OP_INSERT = 1
const OP_DELETE = 2

const OT_MSG = 0
const ACK_MSG = 1
const FORCE_SYNC_MSG = 2
const INIT_MSG = 3
const ACTIVE_SYNC_MSG = 4

Object.prototype.getName = function() {
    var funcNameRegex = /function (.{1,})\(/
    var results = (funcNameRegex).exec((this).constructor.toString())
    return (results && results.length > 1) ? results[1] : ""
}

$(document).ready(function() {
    $('#main').on('keyup', function() {
        sync()
    });
})

//  Change 转换为 一个 JSON
function changeToJSON(change) {
    var data = new Object()
    var ops = new Object()
    ops.op = new Array()
    var shouldSend = false
    var last = change.length - 1
    for (var i = 0; i < change.length; i++) {
        var op = change[i].getName()

        if (op == 'Retain') {
            if (i == last) { // 最后一个retain扔掉
                continue
            }
        } else {
            shouldSend = true
        }

        switch (op) {
            case 'Retain':
                change[i]['type'] = OP_RETAIN
                break
            case 'Skip':
                change[i]['type'] = OP_DELETE
                break
            case 'Insert':
                change[i]['type'] = OP_INSERT
                break
            default:
        }
        ops.op.push(change[i])
    }
    if (shouldSend) {
        ops.adden = change.addendum
        ops.inputLength = change.inputLength
        ops.outputLength = change.outputLength
        ops.removen = change.removendum

        seq++
        data['type'] = OT_MSG
        data['seq'] = seq
        data['uid'] = uid
        data['ver'] = ver
        data['ops'] = ops
        return JSON.stringify(data)
    }
}

function JSONToChange(json) {
    var ops = []
    if (json == "") {
        return
    }

    data = JSON.parse(json)

    let obj = $('#main')
    let pos = getCaretPosition('main')
    let cursorDrift = false

    if (data.type == INIT_MSG) {
        console.log('user data init')
        uid = data.uid
        ver = data.ver
    }

    if (data['type'] == ACK_MSG) {
        // console.log('ack')
        ver = data.ver
        ChangesetQueueLock = false // unclock
        return
    }

    if (data['type'] == FORCE_SYNC_MSG) {
        // console.log('force sync')
        // clean all the content
        his = ''
        obj.val('')
        seq = data.seq
    }

    for (var i = 0; i < data.ops.op.length; i++) {
        var current = data.ops.op[i]

        switch (current['type']) {
            case OP_RETAIN:
                ops.push(new ot.Retain(current.length))
                if (current.length < pos) {
                    // console.log('drift')
                    cursorDrift = true
                }
                break
            case OP_INSERT:
                ops.push(new ot.Insert(current.length))
                break;
            case OP_DELETE:
                ops.push(new ot.Skip(current.length))
                break;
        }
    }

    var change = new ot.Changeset(ops)
    change.addendum = data.ops.adden
    change.removendum = data.ops.removen
    change.inputLength = data.ops.inputLength
    change.outputLength = data.ops.outputLength

    let text = obj.val()
    try {
        his = change.apply(text)
    } catch (e) {
        // console.log(e)
        // Empty all
        ChangesetQueue = []
        ChangesetQueueLock = false // unclock
        sendMsg(JSON.stringify({'type': ACTIVE_SYNC_MSG, 'uid': uid}))
        return null
    }

    if (cursorDrift) {
        pos += change.addendum.length
    }

    ver = data.ver
    return pos
}

function sync() {
    let text = $('#main').val()
    let diff = engine.diff_main(his, text)
    if (diff.length == 1 && diff[0][0] == 0) { // 移动，选择
        return
    }
    let change = Changeset.fromDiff(diff)
    his = text
    ChangesetQueue.push(change)
}

// Queue consumer
var changesetQueueConsumer = setInterval(function() {
    if (!ChangesetQueueLock) { // get ACK
        var merged = ChangesetQueue.shift()
        while (ChangesetQueue.length != 0) {
            merged = merged.merge(ChangesetQueue.shift())
        }
        if (merged != null) {
            let msg = changeToJSON(merged)
            if (msg != null) {
                sendMsg(msg)
            } else {
                ChangesetQueueLock = true // lock
            }
        }
    }
}, 1000)

function sendMsg(msg) {
    // console.log("send: " + msg)
    ChangesetQueueLock = true // lock
    conn.send(msg)
}

function connect() {
    let token = $('#token').val().trim()
    let pubkey = $('#pubkey').val().trim()
    if (token.length == 0) {
        alert("token invalid")
        return
    }

    conn = new WebSocket(target + "?token=" + encodeURIComponent(token) + '&pubkey=' + encodeURIComponent(pubkey))
    conn.onopen = function() {
        console.log("connected to sync ")
    }
    conn.onclose = function(e) {
        didClose = true
        console.log("connection closed (" + e.code + ")")
    }

    conn.onmessage = function(e) {
        let data = e.data
        // console.log('received ' + data)
        let modedPos = JSONToChange(e.data)
        if (modedPos == null) {
            return
        }
        $('#main').val(his)
        setCaretPosition('main', modedPos)
    }

    conn.onclose = function(e) {
        console.log('closed ['+e.code+'] '+e.reason)
        $('#main').attr('readonly', true)
    }
    $('#main').attr("readonly", false)
}


function disconnect() {
    didClose = true
    conn.close()
    $('#main').attr("readonly", true)
}
