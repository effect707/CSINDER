CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    faceit_id     TEXT UNIQUE NOT NULL,
    nickname      TEXT NOT NULL,
    avatar        TEXT NOT NULL DEFAULT '',
    country       TEXT NOT NULL DEFAULT '',
    faceit_url    TEXT NOT NULL DEFAULT '',
    faceit_elo    INT NOT NULL DEFAULT 0,
    skill_level   INT NOT NULL DEFAULT 0,
    win_rate      REAL NOT NULL DEFAULT 0,
    kd_ratio      REAL NOT NULL DEFAULT 0,
    matches       INT NOT NULL DEFAULT 0,
    headshot_pct  REAL NOT NULL DEFAULT 0,
    longest_streak INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id         TEXT PRIMARY KEY,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

CREATE TABLE swipes (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    swiper_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    liked      BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(swiper_id, target_id)
);
CREATE INDEX idx_swipes_target_liked ON swipes(target_id, liked) WHERE liked = true;

CREATE TABLE notifications (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seen         BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_notifications_user_unseen ON notifications(user_id, seen) WHERE seen = false;
