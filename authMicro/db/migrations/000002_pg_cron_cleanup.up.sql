CREATE EXTENSION IF NOT EXISTS pg_cron;
SELECT cron.schedule('cleanup_registration_sessions', '*/1 * * * *',
                     $$ DELETE FROM registration_session WHERE code_expires < NOW() $$);
SELECT cron.schedule('cleanup_login_session', '*/1 * * * *',
                     $$ DELETE FROM login_session WHERE code_expires < NOW() $$);
SELECT cron.schedule('cleanup_token_session', '*/1 * * * *',
                     $$ DELETE FROM refresh_token_session WHERE expires_at < NOW() $$);
