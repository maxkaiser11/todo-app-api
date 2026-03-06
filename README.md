# todos

A full-stack task management app built with Go/Gin on the backend and Svelte on the frontend.

## 🔗 Links

- **Frontend:** https://todo-app-api-frontend.vercel.app/
- **Backend API:** https://todo-app-api.fly.dev

## 🛠 Tech Stack

**Backend**
- [Go](https://golang.org/) + [Gin](https://gin-gonic.com/) — HTTP framework
- [GORM](https://gorm.io/) — ORM
- [SQLite](https://www.sqlite.org/) — Database
- [golang-jwt](https://github.com/golang-jwt/jwt) — JWT auth
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — Password hashing
- Deployed on [Fly.io](https://fly.io)

**Frontend**
- [Svelte](https://svelte.dev/) + TypeScript
- [Vite](https://vitejs.dev/)
- Deployed on [Vercel](https://vercel.com)

## 📁 Project Structure

```
├── backend/
│   ├── cmd/api/
│   │   └── main.go
│   ├── internal/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── models/
│   └── Dockerfile
└── frontend/
    └── src/
        ├── App.svelte
        ├── Auth.svelte
        └── TodoList.svelte
```

## 🚀 Running Locally

**Backend**
```bash
cd backend
cp .env.example .env   # set DB_PATH and JWT_SECRET
go mod download
go run ./cmd/api
```

**Frontend**
```bash
cd frontend
cp .env.example .env   # set VITE_API_URL=http://localhost:8080
npm install
npm run dev
```

## 🔑 API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/auth/register` | ❌ | Register a new user |
| POST | `/auth/login` | ❌ | Login and receive JWT |
| GET | `/tasks` | ✅ | Get all tasks for user |
| POST | `/task/create` | ✅ | Create a new task |
| PUT | `/tasks/:id` | ✅ | Update task status |
| DELETE | `/task/delete/:id` | ✅ | Delete a task |

## 🔐 Auth

JWT-based authentication. Token is stored in `localStorage` and sent as `Authorization: Bearer <token>` on every protected request. Tokens expire after 7 days.
