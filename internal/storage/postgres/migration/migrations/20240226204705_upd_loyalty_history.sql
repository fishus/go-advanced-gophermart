-- +goose Up
alter table public.loyalty_history add order_num varchar not null default '';

UPDATE public.loyalty_history lh SET order_num = (SELECT num FROM public.orders o WHERE o.id = lh.order_id);

alter table public.loyalty_history alter column order_num drop default;

alter table public.loyalty_history drop column order_id;

-- +goose Down
alter table public.loyalty_history add order_id uuid not null default gen_random_uuid();

UPDATE public.loyalty_history lh SET order_id = (SELECT id FROM public.orders o WHERE o.num = lh.order_num);

alter table public.loyalty_history alter column order_id drop default;

alter table public.loyalty_history drop column order_num;