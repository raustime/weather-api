# syntax=docker/dockerfile:1
FROM golang:1.21

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

# Запускаємо бінарник
CMD ["./app"]