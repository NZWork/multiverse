// OT
var Changeset = ot.Changeset
var engine = new diff_match_patch
var his = ""
var seq = 0
var ver = 0

// WS
var conn = null
var target = "ws://" + location.host + "/sync"
var uid = 0

// Switcher
var send = 0
var lastMsg = 0
var didClose = false
var typing = false

const OP_RETAIN = 0
const OP_INSERT = 1
const OP_DELETE = 2

Object.prototype.getName = function() {
    var funcNameRegex = /function (.{1,})\(/
    var results = (funcNameRegex).exec((this).constructor.toString())
    return (results && results.length > 1) ? results[1] : ""
}

//  Change 转换为 一个 JSON
function changeToJSON(change) {
    var data = new Object()
    data['op'] = new Array()

    var shouldSend = false

    var last = change.length - 1

    var addCursor, removeCursor = 0

    for (var i = 0; i < change.length; i++) {
        var op = change[i].getName()

        if (op == 'Retain') {
            if (i == last) { // 最后一个retain扔掉
                continue
            }

            data['op'].push({'type': OP_RETAIN, 'op': change[i].length})
        } else if (op == 'Insert') {
            data['op'].push({'type': OP_INSERT, 'op': change.addendum.substring(addCursor, change[i]['length'])})
            addCursor += change[i]['length'] // move cursor

            shouldSend = true
        } else if (op == 'Skip') {
            data['op'].push({'type': OP_DELETE, 'op': change.removendum.substring(removeCursor, change[i]['length'])})
            removeCursor += change[i]['length']

            shouldSend = true
        }
    }
    if (shouldSend) {
        seq++
        data['seq'] = seq
        data['uid'] = uid
        data['ver'] = ver
        return JSON.stringify(data)
    }
}

function JSONToChange(json) {
    var ops = [],
        removendum = '',
        addendum = ''
    if (json == "") {
        return
    }

    data = JSON.parse(json)
    if (data.uid == uid) {
        return
    }

    let obj = $('#main')
    let pos = getCaretPosition('main')
    let cursorDrift = false

    console.log('current ' + pos, obj.selectionStart)

    for (var i = 0; i < data['op'].length; i++) {
        var current = data['op'][i]

        switch (current['type']) {
            case OP_RETAIN:
                ops.push(new ot.Retain(current.op))
                if (current.op < pos) {
                    console.log('drift')
                    cursorDrift = true
                }
                break;
            case OP_INSERT:
                ops.push(new ot.Insert(current.op.length))
                addendum += current.op
                break;
            case OP_DELETE:
                ops.push(new ot.Skip(current.op.length))
                break;
            default:
        }
    }

    // if (data.op != undefined) {
    //     if (data.op.retain != undefined) {
    //         ops.push(new ot.Retain(data.op.retain))
    //         if (data.op.retain < pos) {
    //             console.log('drift')
    //             cursorDrift = true
    //         }
    //     }
    //     if (data.op.insert != undefined) {
    //         ops.push(new ot.Insert(data.op.insert.length))
    //         addendum += data.op.insert
    //     }
    //     if (data.op.delete != undefined) {
    //         ops.push(new ot.Skip(data.op.delete))
    //     }
    // }
    var change = new ot.Changeset(ops)
    change.addendum = addendum


    let text = obj.val()
    change.inputLength = text.length
    his = change.apply(text)

    if (cursorDrift) {
        pos += addendum.length
    }

    let newOne = new Object
    newOne.content = his
    newOne.pos = pos
    return newOne
}

function sync() {
    let text = $('#main').val()
    let diff = engine.diff_main(his, text)
    if (diff.length == 1 && diff[0][0] == 0) { // 移动，选择
        return
    }
    let change = Changeset.fromDiff(diff)
    console.log(change)
    let s = changeToJSON(change)
    if (s != null) {
        sendMsg(s)
        his = text
    }
}

function sendMsg(msg) {
    console.log("send: " + msg)
    send = Date.now()
    conn.send(msg)
}

function connect() {
    let token = $('#token').val()
    if (token.length == 0) {
        alert("token could not be empty")
        return
    }

    conn = new WebSocket(target + "?token=" + token)
    console.log("connect with sync")
    uid = Math.random()
    conn.onopen = function() {
        console.log("connected to sync ")
    }
    conn.onclose = function(e) {
        didClose = true
        console.log("connection closed (" + e.code + ")")
    }

    conn.onmessage = function(e) {
        let data = e.data
        console.log('received ' + data)
        let moded = JSONToChange(e.data)
        if (moded == null) {
            return
        }
        $('#main').val(moded.content)
        setCaretPosition('main', moded.pos)
    }

    $('#cb').attr("disabled", true)
    $('#db').attr("disabled", false)
    $('#main').attr("readonly", false)
}


function disconnect() {
    didClose = true
    conn.close()
    $('#cb').attr("disabled", false)
    $('#db').attr("disabled", true)
    $('#main').attr("readonly", true)
}
