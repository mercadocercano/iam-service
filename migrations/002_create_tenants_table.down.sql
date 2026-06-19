-- Down: elimina la tabla tenants y todos sus índices, triggers y foreign keys dependientes.
DROP TABLE IF EXISTS tenants CASCADE;
