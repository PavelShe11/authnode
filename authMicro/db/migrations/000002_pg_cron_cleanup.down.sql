SELECT cron.unschedule('cleanup_registration_sessions');
SELECT cron.unschedule('cleanup_login_session');
SELECT cron.unschedule('cleanup_token_session');

SELECT * FROM cron.job;
