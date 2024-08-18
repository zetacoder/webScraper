# Web Scraper

This Go-based web scraper extracts and filters data from the Hacker News website ([https://news.ycombinator.com/](https://news.ycombinator.com/)). It retrieves the first 30 entries, applies specified filters, and stores the results in a PostgreSQL database.

## Prerequisites

Ensure you have the following installed:

- **Go (1.19 or higher)**: [Install Go](https://golang.org/doc/install)
- **PostgreSQL**: [Install PostgreSQL](https://www.postgresql.org/download/)
- **Git**: [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Installation

1. **Clone the Repository:**

   ```
   git clone https://github.com/zetacoder/webScraper.git
   cd webScraper
   ```

2. **Set Up Environment Variables:**

Create a .env file in the root directory of the project with the following content, replacing your_database_dsn with your PostgreSQL Data Source Name (DSN):
   
   ```
   SQL_DSN=your_database_dsn
   ```

   Example DSN: `host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable`



3. **Install Go Dependencies:**

Ensure you are in the project directory and run:

```
go mod tidy
```

This command will download and install the required Go modules.


4. **Set Up PostgreSQL Database:**

Ensure you have a PostgreSQL database running. Create a database if you haven’t already. The DSN provided in the .env file should point to this database.


## Running the Application

To run the web scraper, use the following command (in the root dir):

```
go run main.go
```
    
This will start the scraper, which performs the following tasks:

1. **Fetches the first 30 entries** from Hacker News.
2. Applies the **filtering logic**.
3. **Saves the filtered results** and **usage data** to the PostgreSQL database.


## Testing

To run all the tests of the project, mainly hosted in **./scraper/scraper_test.go** just run in root dir:

   ```
go test ./...
   ```
   

## Error Handling

The application includes comprehensive error handling for:

**Environment Variables**: Ensures the .env file is loaded correctly.
**Database Initialization**: Handles errors during database connection and table migration.
**Web Scraping**: Captures and logs errors encountered while scraping data.
**Data Storage**: Manages errors during data saving to the database.


## Additional Notes

* Database Schema: The schema is managed with GORM’s AutoMigrate function, creating necessary tables and fields in PostgreSQL.
* Environment Configuration: Add your .env file to .gitignore to prevent it from being version-controlled.
* Dependencies: Monitor updates for external libraries (Colly, GORM) for potential changes or improvements.


## References
* [Go Documentation](https://go.dev/)
* [Colly Documentation](https://github.com/gocolly/colly)
* [GORM Documentation](https://gorm.io/)
* [PostgreSQL Documentation](https://www.postgresql.org/docs/current/intro-whatis.html)

Feel free to open issues or submit pull requests if you encounter problems or have suggestions for improvements.

Happy coding!
