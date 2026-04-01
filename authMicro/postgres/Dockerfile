FROM postgres:16.9-bookworm
RUN apt-get update && \
    apt-get install -y postgresql-16-cron && \
    rm -rf /var/lib/apt/lists/*
RUN echo "shared_preload_libraries = 'pg_cron'" >> /usr/share/postgresql/postgresql.conf.sample

RUN mkdir -p /docker-entrypoint-initdb.d
COPY pg_cron.sh ./docker-entrypoint-initdb.d
RUN chmod +x /docker-entrypoint-initdb.d/pg_cron.sh
