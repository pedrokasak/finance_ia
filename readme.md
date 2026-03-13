# Finance IA

This project is a web application for personal finance management, that allows users to create an account, login and see their financial reports. The application is built using Next.js, a popular React framework, and uses a PostgreSQL database to store user data.

investment-portfolio-platform
│
├ backend-api (NestJS)
├ worker-service (RabbitMQ)
├ frontend (React)
├ docker-compose
└ architecture.md

## Stack

The stack used in this project is:

### Frontend

- React: A best library for frontend
- TypeORM: A TypeScript ORM that allows you to interact with the database using a simple and intuitive API.
- Lucide: A set of React components that provide a simple and intuitive way to create UI components.
- Tailwind CSS: A popular CSS framework that allows you to write more concise and maintainable code.

### Backend

- Golang: A use golang how this backend.
- Gorm: A popular ORM for Golang that provides a simple and intuitive way to interact with the database.
- PostgreSQL: A powerful, open source object-relational database system.

## How to run

To run this project, you need to have Node.js installed on your machine. Then, you can run the following commands:

To run the project, execute the following commands:

- In the `web` folder:

  ```bash
  npm install
  npm run dev
  ```

- In the `server` folder:

  ```bash
  docker-compose up -d postgres
  go run main.go
  ```

## 1. Configure envs

- `cp server/.env.example server/.env`
- `edite server/.env com suas credenciais`

## 2. Suba a infraestrutura

- `cd server && make docker-up`

## Serviços disponíveis

- `App:       http://localhost:8080`
- `Prometheus: http://localhost:9090`
- `Grafana:   http://localhost:3001  (admin/grafana_password)`
- `RabbitMQ:  http://localhost:15672 (finance_user/rabbitmq_password)`
- `Mailhog:   http://localhost:8025`

## 3. Rodar testes unitários

- `make coverage`

## 4. Testar rotas no Postman

- `make postman-test`

## 5. Rodar E2E

- `cd web && npm run dev`
- `cd web && npx playwright test`

## Dashboard Preview

Below is a preview of the application's dashboard:

### Dashboard

![Dashboard Screenshot](./web/assets/images/image.png)

### Reports

![Dashboard Screenshot](./web/assets/images/image2.png)

### Profile

![Dashboard Screenshot](./web/assets/images/profile.png)
