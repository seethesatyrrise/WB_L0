create function public.select_all()
    returns TABLE(_data json)
    language plpgsql
as
$$
BEGIN
    return query SELECT data FROM orders as O WHERE O.add_time =
                                                    (SELECT max(O1.add_time) FROM orders as O1 WHERE O1.order_id = O.order_id);
END
$$;

alter function public.select_all() owner to postgres;

