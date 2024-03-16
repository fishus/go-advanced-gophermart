-- +goose Up
create extension pgcrypto;

-- users
create table public.users (
    id uuid primary key not null default gen_random_uuid(),
    username varchar not null,
    password varchar(60) not null,
    created_at timestamp not null default current_timestamp
);
create unique index users_username_uindex on public.users (username);

-- orders
create type order_status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
create table public.orders (
    id uuid primary key not null default gen_random_uuid(),
    user_id uuid not null,
    num varchar not null,
    accrual numeric(15,5) not null default 0,
    status order_status not null default 'NEW'::order_status,
    uploaded_at timestamp not null default current_timestamp,
    CONSTRAINT orders_user_id_fkey foreign key (user_id) references public.users (id) on update cascade on delete cascade
);
create unique index orders_num_uindex on public.orders (num);
create index orders_status_index on public.orders (status);
create index orders_user_id_index on public.orders (user_id);

-- loyalty_history
create table public.loyalty_history (
    user_id uuid not null,
    order_id uuid not null,
    accrual numeric(15,5) not null default 0,
    withdrawal numeric(15,5) not null default 0,
    processed_at timestamp not null default current_timestamp,
    CONSTRAINT loyalty_history_order_id_fkey foreign key (order_id) references public.orders (id) on update cascade on delete cascade,
    CONSTRAINT loyalty_history_user_id_fkey foreign key (user_id) references public.users (id) on update cascade on delete cascade
);

-- loyalty_balance
create table public.loyalty_balance (
    user_id uuid primary key not null,
    current numeric(15,5) not null default 0,
    accrued numeric(15,5) not null default 0,
    withdrawn numeric(15,5) not null default 0,
    CONSTRAINT loyalty_balance_user_id_fkey foreign key (user_id) references public.users (id) on update cascade on delete cascade
);

-- +goose Down
drop table public.loyalty_balance cascade;
drop table public.loyalty_history cascade;
drop index orders_user_id_index;
drop index orders_status_index;
drop index orders_num_uindex;
drop table public.orders cascade;
drop type order_status cascade;
drop index public.users_username_uindex;
drop table public.users cascade;
drop extension if exists pgcrypto;
