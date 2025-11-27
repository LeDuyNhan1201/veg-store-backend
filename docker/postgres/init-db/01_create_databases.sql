CREATE DATABASE "veg-store" WITH OWNER = ldnhan ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8' TABLESPACE = pg_default CONNECTION
LIMIT = -1;

\c "veg-store"
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
