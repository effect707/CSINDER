var candidates = [];
var currentIndex = 0;
var cardStack = document.getElementById('card-stack');
var emptyState = document.getElementById('empty-state');
var btnSkip = document.getElementById('btn-skip');
var btnLike = document.getElementById('btn-like');

function fetchCandidates() {
    fetch('/api/candidates')
        .then(function(r) { return r.json(); })
        .then(function(data) {
            candidates = data || [];
            currentIndex = 0;
            if (candidates.length === 0) {
                showEmpty();
            } else {
                renderCard();
            }
        });
}

function showEmpty() {
    cardStack.innerHTML = '';
    emptyState.style.display = 'block';
    btnSkip.disabled = true;
    btnLike.disabled = true;
}

function renderCard() {
    if (currentIndex >= candidates.length) {
        fetchCandidates();
        return;
    }
    emptyState.style.display = 'none';
    btnSkip.disabled = false;
    btnLike.disabled = false;

    var c = candidates[currentIndex];
    var card = document.createElement('div');
    card.className = 'player-card';
    card.innerHTML =
        '<div class="card-header">' +
            (c.avatar ? '<img src="' + c.avatar + '" alt="' + c.nickname + '" class="card-avatar">' :
                '<div class="card-avatar-placeholder">' + c.nickname.charAt(0) + '</div>') +
            '<div class="card-title">' +
                '<h3>' + c.nickname + '</h3>' +
                (c.country ? '<span class="card-country">' + c.country + '</span>' : '') +
            '</div>' +
        '</div>' +
        '<div class="card-stats">' +
            '<div class="card-stat"><span class="stat-val">' + c.faceit_elo + '</span><span class="stat-lbl">ELO</span></div>' +
            '<div class="card-stat"><span class="stat-val">' + c.skill_level + '</span><span class="stat-lbl">Уровень</span></div>' +
            '<div class="card-stat"><span class="stat-val">' + c.kd_ratio.toFixed(2) + '</span><span class="stat-lbl">K/D</span></div>' +
            '<div class="card-stat"><span class="stat-val">' + Math.round(c.win_rate) + '%</span><span class="stat-lbl">Винрейт</span></div>' +
            '<div class="card-stat"><span class="stat-val">' + c.matches + '</span><span class="stat-lbl">Матчи</span></div>' +
            '<div class="card-stat"><span class="stat-val">' + Math.round(c.headshot_pct) + '%</span><span class="stat-lbl">HS</span></div>' +
        '</div>';

    card.dataset.targetId = c.id;

    var startX = 0;
    var currentX = 0;
    var dragging = false;

    card.addEventListener('mousedown', function(e) {
        dragging = true;
        startX = e.clientX;
        card.style.transition = 'none';
    });

    card.addEventListener('touchstart', function(e) {
        dragging = true;
        startX = e.touches[0].clientX;
        card.style.transition = 'none';
    });

    document.addEventListener('mousemove', onMove);
    document.addEventListener('touchmove', onMove);

    function onMove(e) {
        if (!dragging) return;
        var x = e.clientX || (e.touches && e.touches[0].clientX);
        currentX = x - startX;
        card.style.transform = 'translateX(' + currentX + 'px) rotate(' + (currentX * 0.1) + 'deg)';
        if (currentX > 50) {
            card.classList.add('swiping-right');
            card.classList.remove('swiping-left');
        } else if (currentX < -50) {
            card.classList.add('swiping-left');
            card.classList.remove('swiping-right');
        } else {
            card.classList.remove('swiping-right', 'swiping-left');
        }
    }

    document.addEventListener('mouseup', onEnd);
    document.addEventListener('touchend', onEnd);

    function onEnd() {
        if (!dragging) return;
        dragging = false;
        document.removeEventListener('mousemove', onMove);
        document.removeEventListener('touchmove', onMove);
        document.removeEventListener('mouseup', onEnd);
        document.removeEventListener('touchend', onEnd);

        if (currentX > 100) {
            swipe(true);
        } else if (currentX < -100) {
            swipe(false);
        } else {
            card.style.transition = 'transform 0.3s';
            card.style.transform = '';
            card.classList.remove('swiping-right', 'swiping-left');
        }
    }

    cardStack.innerHTML = '';
    cardStack.appendChild(card);
}

function swipe(liked) {
    var c = candidates[currentIndex];
    if (!c) return;

    var card = cardStack.querySelector('.player-card');
    if (card) {
        card.style.transition = 'transform 0.4s, opacity 0.4s';
        card.style.transform = 'translateX(' + (liked ? 500 : -500) + 'px) rotate(' + (liked ? 30 : -30) + 'deg)';
        card.style.opacity = '0';
    }

    fetch('/api/swipe', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({target_id: c.id, liked: liked})
    });

    setTimeout(function() {
        currentIndex++;
        renderCard();
    }, 400);
}

btnSkip.addEventListener('click', function() { swipe(false); });
btnLike.addEventListener('click', function() { swipe(true); });

fetchCandidates();
