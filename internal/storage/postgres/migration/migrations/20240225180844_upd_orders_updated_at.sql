-- +goose Up
ALTER TABLE public.orders ADD updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE public.loyalty_history DROP CONSTRAINT loyalty_history_order_id_fkey;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION mdt_orders_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
        NEW.updated_at = now();
        RETURN NEW;
    ELSE
        RETURN OLD;
    END IF;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

DROP TRIGGER IF EXISTS mdt_orders_updated_at ON orders;

-- +goose StatementBegin
CREATE TRIGGER mdt_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
EXECUTE PROCEDURE mdt_orders_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS mdt_orders_updated_at ON orders CASCADE;
DROP FUNCTION IF EXISTS mdt_orders_updated_at() CASCADE;
ALTER TABLE public.loyalty_history ADD CONSTRAINT loyalty_history_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders (id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE public.orders DROP COLUMN updated_at;