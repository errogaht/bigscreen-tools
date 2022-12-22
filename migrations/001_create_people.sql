create table oculus_profiles
(
    id              varchar primary key not null,
    image_url       varchar             not null,
    small_image_url varchar             not null
);
create table steam_profiles
(
    id                         varchar primary key not null,
    community_visibility_state smallint            not null,
    profile_state              smallint            not null,
    persona_name               varchar             not null,
    profile_url                varchar             not null,
    avatar                     varchar             not null,
    avatar_medium              varchar             not null,
    avatar_full                varchar             not null,
    avatar_hash                varchar             not null,
    persona_state              smallint            not null,
    real_name                  varchar             not null,
    primary_clan_id            varchar             not null,
    created_at                 timestamp           not null,
    persona_state_flags        smallint            not null,
    loc_country_code           varchar             not null
);

create table account_profiles
(
    username          varchar primary key not null,
    created_at        timestamp           not null,
    is_verified       boolean             not null,
    is_banned         boolean             not null,
    is_staff          boolean             not null,
    steam_profile_id  varchar default null,
    oculus_profile_id varchar default null,
    CONSTRAINT fk_steam_profile
        FOREIGN KEY (steam_profile_id)
            REFERENCES steam_profiles (id)
            ON DELETE SET NULL,
    CONSTRAINT fk_oculus_profile
        FOREIGN KEY (oculus_profile_id)
            REFERENCES oculus_profiles (id)
            ON DELETE SET NULL
);
create table rooms
(
    id           varchar primary key default gen_random_uuid() not null,
    created_at   timestamp                                  not null,
    participants int                                        not null,
    status       varchar                                    not null,
    invite_code  varchar                                    not null,
    visibility   varchar                                    not null,
    room_type    varchar                                    not null,
    version      varchar                                    not null,
    size         int                                        not null,
    environment  varchar                                    not null,
    category     varchar                                    not null,
    description  varchar                                    not null,
    name         varchar                                    not null
);

create index rooms_category_index
    on rooms (category);

create table room_users
(
    user_session_id varchar primary key not null,
    seat_index      smallint         not null,
    created_at      timestamp        not null,
    account_profile varchar default null,
    room_id         varchar    default null,
    version         varchar          not null,
    is_staff        boolean          not null,
    is_mod          boolean          not null,
    is_admin        boolean          not null,
    CONSTRAINT fk_account_profile
        FOREIGN KEY (account_profile)
            REFERENCES account_profiles (username)
            ON DELETE SET NULL,
    CONSTRAINT fk_room_id
        FOREIGN KEY (room_id)
            REFERENCES rooms (id)
            ON DELETE CASCADE
);

alter table rooms
    add creator_profile varchar default null;

alter table rooms
    add constraint rooms_account_profiles_username_fk
        foreign key (creator_profile) references account_profiles (username)
            on delete set null;


---- create above / drop below ----
drop table room_users cascade;
drop table account_profiles cascade;
drop table oculus_profiles cascade;
drop table steam_profiles cascade;
drop table rooms cascade;