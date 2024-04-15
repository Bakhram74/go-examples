REVOKE CREATE ON DATABASE postgres FROM create_obj;
REASSIGN OWNED BY sw_migration TO "user";
DROP ROLE sw_service, sw_migration, service, admin, create_obj;
