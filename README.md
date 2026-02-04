# 💰 Budget Tracker

A robust, secure personal finance management application built with **Go** (Buffalo Framework) and **PostgreSQL**. Track your income, expenses, and net margins with ease.

## ✨ Features

### 🔐 Authentication & Security
- **Secure Signup/Login**: Full user registration and authentication loop.
- **OTP Verification**: Email-based One-Time Password verification for new accounts.
- **Forgot Password**: Secure password reset flow using OTPs.
- **Bcrypt Hashing**: Industry-standard password hashing for security.
- **Session Management**: Secure session handling with Buffalo.

### 💸 Budget Management
- **Transaction Tracking**: Log Income and Expense transactions.
- **Smart Dashboard**: View real-time Net Margin and financial health indicators.
- **Date Filtering**: Filter transactions by date range (coming soon).
- **Responsive Design**: Mobile-friendly UI built with Bootstrap.

## 🛠️ Tech Stack

- **Language**: [Go (Golang)](https://golang.org/)
- **Framework**: [Buffalo](https://gobuffalo.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Frontend**: Plush Templates (HTML), Bootstrap 4, CSS
- **Containerization**: Docker, Docker Compose

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (v1.16+)
- [Docker](https://www.docker.com/products/docker-desktop)
- [PostgreSQL](https://www.postgresql.org/) (if running locally without Docker)
- [Buffalo CLI](https://gobuffalo.io/docs/installation) (`go install github.com/gobuffalo/cli/cmd/buffalo@latest`)

### 🐳 Run with Docker (Recommended)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/budget_tracker.git
   cd budget_tracker
   ```

2. **Start the services:**
   ```bash
   docker-compose up --build
   ```

3. **Access the App:**
   Open [http://localhost:3000](http://localhost:3000) in your browser.

### 💻 Run Locally

1. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

2. **Setup Database:**
   Update `database.yml` with your local Postgres credentials if necessary.
   ```bash
   buffalo pop create
   buffalo pop migrate
   ```

3. **Start the Server:**
   ```bash
   buffalo dev
   ```

## 🧪 Testing

Run the test suite to ensure everything is working correctly:

```bash
buffalo test
```

## 📂 Project Structure

- `actions/`: Application handlers and business logic.
- `models/`: Database models and struct definitions.
- `templates/`: Plush HTML templates for the UI.
- `migrations/`: Database migration files.
- `public/`: Static assets (CSS, JS, Images).

---

Built with ❤️ using Go & Buffalo.
