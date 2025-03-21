# go-feedo

An RSS feed aggregator microservice written in Go.  
The command-line interface (CLI) is homemade, without using frameworks like Bubble Tea or Cobra.  
Interactions with the PostgreSQL database are implemented without an ORM like GORM.

# Requirements

Before anything, you need to have [Go](https://go.dev/doc/install),
[PostgreSQL](https://www.postgresql.org/docs/current/installation.html),
and [Goose](https://github.com/pressly/goose?tab=readme-ov-file#install) installed on your system.

Configure the database. In the `~/.config/go-feedo/config.json` file, add the following line:

```json
{
  "db_url": "postgres://<db_user>:<db_password>@localhost:5432/<db_name>?sslmode=disable"
}
```

Migrate up:

```bash
git clone https://github.com/alnah/go-feedo
cd go-feedo
goose postgres -dir="sql/schema" "postgres://<db_user>:<db_password>@localhost:5432/<db_name> up"
```

# Installation

The straightforward option is to install the CLI globally using Go tooling:

```bash
go install github.com/alnah/go-feedo@latest
```

Or, after cloning the repo from GitHub, compile it using the Makefile from your repo directory:

```bash
make
```

Otherwise, download the appropriate [artificat](https://github.com/alnah/go-feedo/releases)
extract the downloaded file, and install the CLI according to your custom setup.

# Usage

Ask for help inside the CLI:

```bash
go-feedo help
```

# Licence

This project is distributed under the Apache License.
See the [LICENCE](https://github.com/alnah/go-feedo/blob/main/LICENCE) file for more details.
