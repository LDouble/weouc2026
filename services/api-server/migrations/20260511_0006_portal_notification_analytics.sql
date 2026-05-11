CREATE TABLE IF NOT EXISTS portal_id_sequences (
    name TEXT PRIMARY KEY,
    value BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS portal_banners (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    action_url TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS portal_notices (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    content TEXT NOT NULL,
    audience TEXT NOT NULL,
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    pinned BOOLEAN NOT NULL DEFAULT FALSE,
    publisher_user_id TEXT NOT NULL,
    publisher TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_portal_banners_sort ON portal_banners (sort_order, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_portal_notices_published_at ON portal_notices (published_at DESC);

INSERT INTO portal_id_sequences (name, value) VALUES
    ('notice', 300)
ON CONFLICT (name) DO UPDATE SET value = GREATEST(portal_id_sequences.value, EXCLUDED.value);

INSERT INTO portal_banners (id, title, description, image_url, action_url, sort_order, created_at) VALUES
    ('banner-101', '新学期校园生活入口升级', '跑腿、组局、二手、资料与失物招领统一接入新首页。', 'https://example.com/portal/banner-101.png', '/pages/home/index', 1, '2026-05-11T09:00:00Z'),
    ('banner-102', '教务绑定后可查看联系方式', '涉及联系方式的内容均以后端绑定状态裁剪结果为准。', 'https://example.com/portal/banner-102.png', '/pages/profile/academic/index', 2, '2026-05-11T09:10:00Z')
ON CONFLICT (id) DO NOTHING;

INSERT INTO portal_notices (id, title, summary, content, audience, tags, pinned, publisher_user_id, publisher, published_at, created_at) VALUES
    ('notice-101', '校园综合应用内测启动', '微信小程序已开放跑腿、组局、二手、资料与失物招领基础能力。', '本周开放第一轮内测，欢迎同学体验并通过站内消息反馈问题。', 'all', '["内测","公告"]'::jsonb, TRUE, 'admin-001', '校园运营中心', '2026-05-11T09:00:00Z', '2026-05-11T09:00:00Z'),
    ('notice-102', '二手与跑腿发布规范', '新发布内容默认进入审核，违规内容会被驳回或下线。', '请勿发布违法违规、交易风险高或联系方式异常的信息，审核员会依据规则处理。', 'student', '["审核","发布规范"]'::jsonb, FALSE, 'admin-001', '校园运营中心', '2026-05-11T09:30:00Z', '2026-05-11T09:30:00Z')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS notification_id_sequences (
    name TEXT PRIMARY KEY,
    value BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS notification_messages (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT NOT NULL,
    target_scope TEXT NOT NULL,
    target_user_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
    action_url TEXT NOT NULL DEFAULT '',
    publisher_user_id TEXT NOT NULL,
    publisher TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    read_by_user_ids JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_notification_messages_created_at ON notification_messages (created_at DESC);

INSERT INTO notification_id_sequences (name, value) VALUES
    ('notification', 500)
ON CONFLICT (name) DO UPDATE SET value = GREATEST(notification_id_sequences.value, EXCLUDED.value);

INSERT INTO notification_messages (
    id, title, content, category, target_scope, target_user_ids, action_url, publisher_user_id, publisher, created_at, read_by_user_ids
) VALUES
    ('notification-101', '欢迎使用校园综合应用', '你可以通过首页动态快速查看跑腿、组局、二手、资料和失物招领。', 'system', 'all', '[]'::jsonb, '/pages/home/index', 'system', '系统助手', '2026-05-11T10:00:00Z', '{}'::jsonb),
    ('notification-102', '教务绑定后可查看联系方式', '若你需要联系发布者，请先完成教务绑定以解锁联系方式查看权限。', 'reminder', 'users', '["u-1"]'::jsonb, '/pages/profile/academic/index', 'admin-001', '校园运营中心', '2026-05-11T10:15:00Z', '{}'::jsonb)
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS audit_id_sequences (
    name TEXT PRIMARY KEY,
    value BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY,
    actor_id TEXT NOT NULL,
    actor_name TEXT NOT NULL,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    result TEXT NOT NULL,
    message TEXT NOT NULL,
    details JSONB,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_action_created_at ON audit_logs (action, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_created_at ON audit_logs (actor_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs (created_at DESC);

INSERT INTO audit_id_sequences (name, value) VALUES
    ('audit', 900)
ON CONFLICT (name) DO UPDATE SET value = GREATEST(audit_id_sequences.value, EXCLUDED.value);
