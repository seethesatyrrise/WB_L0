create or replace function public.select_all()
    returns TABLE(id text, _data json)
    language plpgsql
as
$$
BEGIN
    return query SELECT order_id, data FROM orders as O WHERE O.add_time =
                                                    (SELECT max(O1.add_time) FROM orders as O1 WHERE O1.order_id = O.order_id);
END
$$;

alter function public.select_all() owner to postgres;

select * from select_all();
