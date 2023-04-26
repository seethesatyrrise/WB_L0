CREATE or replace FUNCTION insert_data(data json) RETURNS VOID
    LANGUAGE SQL
    AS $$
INSERT INTO orders VALUES (data->>'order_uid', data);
$$;