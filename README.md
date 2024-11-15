# Messages API that uses Postgresql developed in Go

A API that gets messages from a PostgreSQL database table with attachments.

---

🔥 Cool Features
* Fetch Mastodon posts and store them in PostgreSQL
* Supports media attachments
* Has user authentication
--- 


🚀 Installation

**Initial setup**<br>
1. **Install Go:** https://golang.org/doc/install
2. **Clone repository:** the repository and navigate to the project folder.
3. **Install dependencies:** Run `go mod tidy` to install dependencies.

**Setup PostgreSQL**<br>
* Create a database `msg-dev`.
* Username: `postgres`
* Password: `root`

Add the connection strings for PostgreSQL to your **environment variables**.
In Goland editor this is done by
Edit Configurations > Run/Debug Configurations + Go Build<br>
Files: `main.go`<br>
Environment: `DB_HOST=localhost;DB_NAME=msg-dev;DB_PASSWORD=root;DB_PORT=5432;DB_USER=postgres;`<br>

**Database Tables**<br>
The API gets its data from the following two database tables that has to exists:

```
CREATE TABLE IF NOT EXISTS messages_index (
    msg_id SERIAL PRIMARY KEY,
    msg_platform VARCHAR(200) DEFAULT NULL,
    msg_external_id VARCHAR(200) DEFAULT NULL,
    msg_created_at TIMESTAMP DEFAULT NULL,
    msg_language VARCHAR(200) DEFAULT NULL,
    msg_url VARCHAR(200) DEFAULT NULL,
    msg_content TEXT DEFAULT NULL,
    msg_external_account_id VARCHAR(200) DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS messages_attachments (
    attachment_id SERIAL PRIMARY KEY,
    attachment_msg_id INT DEFAULT NULL,
    attachment_external_id VARCHAR(200) DEFAULT NULL,
    attachment_url VARCHAR(200) DEFAULT NULL,
    attachment_type VARCHAR(200) DEFAULT NULL,
    attachment_meta_description TEXT DEFAULT NULL
);

```

**Running Locally**<br>
* **Run locally:** `go run main.go`.
* **Access API:** The API will be available at `http://localhost:8080`.

**Running Locally with Docker**<br>
* **Build:** `docker build -t mastodon-statuses-to-postgres-go .`.
* **Run:** `docker run -p 8080:8080 mastodon-statuses-to-postgres-go`

**Setup service account**<br>

* POST http://localhost:8080/api/v1/users/signup<br>
  Body > Raw:
```
{
    "email": "service.account@email.com",
    "password": "GolangMsg1#"
}
```

Set the new users as `approved` with the following SQL:

```sql
UPDATE users SET approved=1;
```

Then login
* POST http://localhost:8080/api/v1/users/login<br>
  Body > Raw:
```
{
    "email": "service.account@email.com",
    "password": "GolangMsg1#"
}
```

Next you can use the different routes with the authorization token:

* GET http://localhost:8080/api/v1/messages<br>
  Headers:
    - Authorization Bearer TOKEN



**Running Locally with Docker**<br>
* **Build:** `docker build -t messages-api-postgresql-go .`.
* **Run:** `docker run -p 8080:8080 messages-api-postgresql-go`


---

## 📖 License

This project is licensed under the
[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

```
Copyright 2024 github.com/ditlef9

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```