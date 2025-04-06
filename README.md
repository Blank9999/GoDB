# GoDB

GoDB is a lightweight, JSON-based database built with Go, inspired by MongoDB. It supports fast and concurrent data operations by leveraging Go's powerful concurrency model. Ideal for projects that need a simple, file-based NoSQL database without the overhead of external dependencies.

## 🌟💡 Inspiration

GoDB was born out of the curiosity to explore low-level database design while taking advantage of Go’s excellent support for concurrency. The goal was to build a simplified version of a document-based NoSQL system like MongoDB that stores data in `.json` files and supports basic operations such as create, read, update, and delete (CRUD).

## 🛠️📋 How It Works

### ⚙️ Server (Go)

- `main.go`: Launches the Go server and listens for HTTP requests.
- Routes are available to:
  - ➕ Insert new documents
  - 🔍 Retrieve documents
  - 📝 Update existing data
  - ❌ Delete entries
- All data is stored as `.json` files in the filesystem, organized by collection.

### 🌐 Frontend (Optional)

- You can build your own HTML/JS UI or use the provided HTML to send HTTP requests to the GoDB server.
- Open the HTML file with **Live Server** (VSCode extension or similar) to interact with GoDB via forms or fetch requests.

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/Blank9999/GoDB.git
cd GoDB
```

### 2. Run the Go Server

```bash
go run main.go
```
### 3. Open the Frontend (Optional)

- Use Live Server or any static server to host the provided HTML.
- Interact with the Go backend via the web interface to store and retrieve data from .json files.

### 🧠 Key Features

- ⚡ Built with Go for performance and concurrency.
- 💾 JSON-based, document-style data storage.
- 🔄 Basic CRUD functionality through HTTP requests.
- 🧩 Easy to integrate into small tools or projects.

### 🔮 Future Enhancements

- 🗃️ Support for indexing and querying on specific fields.
- 🧩 Schema validation for collections.
- 🔐 Authentication and access control for endpoints.
