function debug() {
}


function getCaretPosition(id) {
    var txt = document.getElementById(id);
    var startPos = txt.selectionStart;
    var endPos = txt.selectionEnd;
    return endPos;
}

function setCaretPosition(id, pos) {
    var txt = document.getElementById(id);
    if(txt.createTextRange) {
      var range = txt.createTextRange();
      range.collapse(true);
      range.moveEnd('character', pos);
      range.moveStart('character', pos);
      range.select();
      return;
    }

    if(txt.selectionStart) {
      txt.focus();
      txt.setSelectionRange(pos, pos);
    }
}
