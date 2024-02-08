DROP TRIGGER IF EXISTS update_time_trigger ON "user-info";

DROP FUNCTION IF EXISTS update_time();

DROP TABLE IF EXISTS "user-info";
DROP TABLE IF EXISTS "users";

DROP EXTENSION IF EXISTS "uuid-ossp";
