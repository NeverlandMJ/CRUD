create table customers
(
      id bigserial primary key,
      name text not null,
      phone text not null,
      phone text not null unique,
      password text not null,
      active boolean not null default true,
      created timestamp not null default current_timestamp
);

create table customer_tokens
(
      token text not null unique,
      customer_id bigint not null references customers,
      expire timestamp not null default current_timestamp + interval '1 hour',
      created timestamp not null defautl current_timestamp
);

create table if not exists managers 
(
    id bigserial primary key,
    name	text not null,
    salary integer not null default 0,
    plan    integer not null default 0,
    boss_id bigint references managers,
    departament text,
    login 	text unique,
    password text ,
    active 	boolean not null default true,
    created timestamp not null default current_timestamp 
);

