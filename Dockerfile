# syntax=docker/dockerfile:1
FROM golang:1.21

# Встановлюємо bash та netcat
RUN apt-get update && apt-get install -y bash netcat-openbsd postgresql-client && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копіюємо go.mod та go.sum і завантажуємо залежності
COPY go.mod go.sum ./
RUN go mod download

# Копіюємо весь код у контейнер
COPY . .

# Збираємо бінарник під ім'ям app
RUN go build -o app main.go


# Переконуємося, що app має права на виконання
RUN chmod +x ./app

COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

ENTRYPOINT ["/wait-for-postgres.sh"]
# Запускаємо бінарник
CMD ["./app"]
