SELECT cron.schedule('cleanup_registration_sessions', '* */3 * * *',
                     $$ DELETE FROM registration_session WHERE code_expires < NOW() $$);
SELECT cron.schedule('cleanup_login_session', '* */5 * * *',
                     $$ DELETE FROM login_session WHERE code_expires < NOW() $$);
SELECT cron.schedule('cleanup_token_session', '30 */2 * * *',
                     $$ DELETE FROM refresh_token_session WHERE expires_at < NOW() $$);

SELECT * FROM cron.job;