create table notification_rules(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid,
    get_immediate_notification boolean not null default false,
    get_filtered_notification boolean not null default false,
    get_immediate_notification_for_selected_domains boolean not null default true,
    selected_domains varchar[],
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
    CONSTRAINT fk_user_id_notification_rules FOREIGN KEY (user_id)  REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL
);