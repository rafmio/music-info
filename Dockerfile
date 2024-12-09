FROM postgres:latest

# Копируем SQL-скрипт для создания таблицы в контейнер
COPY migrations/create_song_table.sql /docker-entrypoint-initdb.d/

# Указываем пользовательский порт для доступа к базе данных
EXPOSE 5432

# Указываем команду для запуска PostgreSQL при старте контейнера
# CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]
