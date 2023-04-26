CREATE or replace FUNCTION select_data(id text, inout _data json)
    LANGUAGE plpgsql
AS $$
BEGIN
    SELECT data FROM orders as O WHERE O.order_id = id ORDER BY O.add_time DESC limit 1
                                 INTO _data;
END
$$;

select select_data('b563feb7b2b84b6test');