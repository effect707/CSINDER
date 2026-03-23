var wsUrl = (location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + '/ws';
var ws = null;
var reconnectDelay = 1000;

function connectWS() {
    ws = new WebSocket(wsUrl);

    ws.onopen = function() {
        reconnectDelay = 1000;
    };

    ws.onmessage = function(e) {
        try {
            var msg = JSON.parse(e.data);
            if (msg.type === 'like_notification') {
                showToast(msg.data.from_nickname + ' выбрал тебя!', msg.data.from_avatar);
                updateBadge();
            }
        } catch(err) {}
    };

    ws.onclose = function() {
        setTimeout(function() {
            reconnectDelay = Math.min(reconnectDelay * 2, 30000);
            connectWS();
        }, reconnectDelay);
    };

    ws.onerror = function() {
        ws.close();
    };
}

function showToast(text, avatar) {
    var container = document.getElementById('toast-container');
    if (!container) return;

    var toast = document.createElement('div');
    toast.className = 'toast';
    toast.innerHTML =
        (avatar ? '<img src="' + avatar + '" class="toast-avatar">' : '') +
        '<span>' + text + '</span>';

    container.appendChild(toast);
    setTimeout(function() { toast.classList.add('show'); }, 10);
    setTimeout(function() {
        toast.classList.remove('show');
        setTimeout(function() { toast.remove(); }, 300);
    }, 5000);
}

function updateBadge() {
    var badge = document.querySelector('.badge');
    if (badge) {
        var count = parseInt(badge.textContent) || 0;
        badge.textContent = count + 1;
        badge.style.display = 'inline-flex';
    } else {
        var matchesLink = document.querySelector('a[href="/matches"]');
        if (matchesLink) {
            var newBadge = document.createElement('span');
            newBadge.className = 'badge';
            newBadge.textContent = '1';
            matchesLink.appendChild(newBadge);
        }
    }
}

connectWS();
