create or replace function public.insert_data(data json) returns text
    language sql
as
$$
INSERT INTO orders VALUES (data->>'order_uid', data) RETURNING order_id as id;
$$;

alter function public.insert_data(json) owner to postgres;