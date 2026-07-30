CREATE TYPE notification_type AS ENUM ('cve', 'system');

create table notifications (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    type notification_type not null default 'cve',
    cve_id uuid,
    user_id uuid,
    created_at timestamp not null default now(),
    viewed boolean not null default false,
    viewd_at timestamp,
    updated_at timestamp  not null default now(),
    deleted_at timestamp,
    CONSTRAINT fk_cve_id FOREIGN KEY (cve_id)  REFERENCES cves(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)  REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL
)