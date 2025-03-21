# go-feedo

An RSS feed aggregator microservice written in Go.  
The command-line interface (CLI) is homemade, without using frameworks like Bubble Tea or Cobra.  
Interactions with the PostgreSQL database are implemented without an ORM like GORM.

# Requirements

Before anything, you need to have [Go](https://go.dev/doc/install) and [PostgreSQL](https://www.postgresql.org/docs/current/installation.html) installed on your system.

Configure the database. In the `~/.config/go-feedo/config.json` file, add the following line:

```json
{
  "db_url": "postgres://<db_user>:<db_password>@localhost:5432/<db_name>?sslmode=disable"
}
```

# Installation

The straightforward option is to install the CLI globally using Go tooling:

```bash
go install github.com/alnah/go-feedo@latest
```

Or, clone the repo from GitHub and compile it using the Makefile:

```bash
git clone https://github.com/alnah/go-feedo && cd go-feedo && make
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
