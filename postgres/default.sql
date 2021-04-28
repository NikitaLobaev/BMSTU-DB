ALTER SYSTEM SET max_wal_senders = 0;
ALTER SYSTEM SET wal_level = minimal;
ALTER SYSTEM SET fsync = OFF;
ALTER SYSTEM SET full_page_writes = OFF;
ALTER SYSTEM SET synchronous_commit = OFF;
ALTER SYSTEM SET archive_mode = OFF;
ALTER SYSTEM SET shared_buffers = '400 MB';
ALTER SYSTEM SET effective_cache_size = '1 GB';
ALTER SYSTEM SET work_mem = '32 MB';

SELECT pg_reload_conf();
