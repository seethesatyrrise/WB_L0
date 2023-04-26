create table public.orders
(
    order_id text not null,
    data     json,
    add_time timestamp default now()
);

alter table public.orders
    owner to postgres;

create index orders_order_id_index
    on public.orders (order_id);