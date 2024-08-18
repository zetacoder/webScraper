# Web Scraper

This Go-based web scraper extracts and filters data from the Hacker News website ([https://news.ycombinator.com/](https://news.ycombinator.com/)). It retrieves the first 30 entries, applies specified filters, and stores the results in a PostgreSQL database.

## Prerequisites

Ensure you have the following installed:

- **Go (1.19 or higher)**: [Install Go](https://golang.org/doc/install)
- **PostgreSQL**: [Install PostgreSQL](https://www.postgresql.org/download/)
- **Git**: [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/zetacoder/webScraper.git
   cd webScraper

2. **Set Up Environment Variables:**

Create a .env file in the root directory of the project with the following content, replacing your_database_dsn with your PostgreSQL Data Source Name (DSN):
   
   ```bash
   SQL_DSN=host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable
   Example DSN: postgres://username:password@localhost:5432/mydatabase?sslmode=disable


3. **Install Go Dependencies:**

Ensure you are in the project directory and run:

   ```bash  
   go mod tidy
