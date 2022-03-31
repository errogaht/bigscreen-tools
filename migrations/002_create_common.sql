create table settings
(
    id        varchar not null
        constraint settings_pk
            primary key,
    timestamp timestamp,
    int       int,
    string    varchar,
    bool      bool
);


---- create above / drop below ----
drop table settings