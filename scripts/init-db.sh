#!/bin/bash

echo "Executing migrations"

# Ejecutar las migraciones
for file in /migrations/*.up.sql; do
  echo "Executing migration: $file"
  PGPASSWORD=$POSTGRES_PASSWORD psql -h "$DB_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$file"
done

echo "Migrations completed"
