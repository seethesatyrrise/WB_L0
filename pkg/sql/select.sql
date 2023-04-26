create function public.select_data(id text) returns json
    language plpgsql
as
$$
DECLARE
    _data json;
BEGIN

    SELECT data INTO _data FROM orders as O WHERE O.order_id = id ORDER BY O.add_time DESC limit 1;
    RETURN _data;
END
$$;

alter function public.select_data(text) owner to postgres;


select select_data('b563feb7b2b84b6test');