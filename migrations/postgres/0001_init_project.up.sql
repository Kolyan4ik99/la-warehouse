create table if not exists public.warehouse
(
    id        serial primary key,
    name      varchar not null,
    available boolean not null
);

create table if not exists public.products
(
    id     serial PRIMARY KEY,
    name   varchar not null,
    size   int     not null,
    amount int     not null
);

create table if not exists warehouse_product
(
    product_id   int not null,
    warehouse_id int not null,
    amount       int not null,
    unique (product_id, warehouse_id),
    CONSTRAINT fk_product_id
        FOREIGN KEY (product_id)
            REFERENCES public.products (id)
            ON DELETE SET NULL,
    CONSTRAINT fk_warehouse_id
        FOREIGN KEY (warehouse_id)
            REFERENCES public.warehouse (id)
            ON DELETE SET NULL
);