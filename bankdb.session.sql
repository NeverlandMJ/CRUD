create table if not exists managers 
(
    id bigserial primary key,
    name	text not null,
    salary integer not null default 0,
    plan    integer not null default 0,
    boss_id bigint references managers,
    departament text,
    login 	text 	,
    password text ,
    active 	boolean not null default true,
    created timestamp not null default current_timestamp 
);