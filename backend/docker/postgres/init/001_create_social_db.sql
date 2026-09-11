SELECT 'CREATE DATABASE gav_social'
WHERE NOT EXISTS (
    SELECT FROM pg_database WHERE datname = 'gav_social'
)\gexec
