-- создание роли администратора
CREATE ROLE admin PASSWORD 'admin123' LOGIN SUPERUSER CONNECTION LIMIT 5;

-- создание групповых ролей
CREATE ROLE create_obj;
GRANT CREATE ON DATABASE postgres TO create_obj;
CREATE ROLE service IN ROLE pg_read_all_data, pg_write_all_data;

-- создание ролей внутри каждой групповой роли
CREATE ROLE sw_migration PASSWORD 'test123' CONNECTION LIMIT 10 LOGIN IN ROLE create_obj;
CREATE ROLE sw_service PASSWORD 'test123' CONNECTION LIMIT 10 LOGIN IN ROLE service;
