# ADR-001: Web Crawler Architecture

**Date:** August 18, 2024

**Status:** Proposed

## Context
We are tasked with developing a web crawler to extract and filter data from the Hacker News website (https://news.ycombinator.com/). The requirements are to:

1. Extract the first 30 entries, including number, title, points, and number of comments.
2. Implement two filters:
   - Filter entries with more than five words in the title, ordered by the number of comments.
   - Filter entries with five words or fewer in the title, ordered by points.
3. Store usage data including request timestamp and applied filter.

## Decision
We have chosen to implement the solution in Golang using the following stack:

- **Language:** Golang
- **Web Scraping Library:** Colly
- **Database:** PostgreSQL
- **ORM:** GORM

## Architecture Overview

### Scraping
- We use the Colly library to perform web scraping. This library allows us to create a web collector and specify how to process HTML elements.
- The `Scrape()` method fetches data from the target website and extracts required information (title, source, points, and comments) using CSS selectors.

### Filtering
- After scraping, the `FilterPosts()` method processes the collected posts to apply the required filters:
  - **More than Five Words:** Filters posts with titles containing more than five words, ordered by the number of comments.
  - **Five Words or Fewer:** Filters posts with titles containing five words or fewer, ordered by points.
- The filtering methods (`filterMoreThanFiveWords` and `filterFiveWordsOrLess`) use regular expressions to count words and sort the results accordingly.

### Storage
- The data is stored in a PostgreSQL database using the GORM ORM for ease of database interaction and schema management.
- Two main entities are managed:
  - **Post:** Stores information about each post.
  - **UsageData:** Stores metadata about the scraping operation, including timestamps and filter details.

### Usage Tracking
- The `UsageData` entity captures metadata about each scraping operation, such as start and end times, total bytes scraped, and filtering details.

## Design Decisions

### Choice of Technology
- **Golang:** Chosen for its performance, strong concurrency support, and ease of integration with the Colly library.
- **Colly:** Selected for its simplicity and efficiency in web scraping tasks.
- **GORM:** Used for its robust support for PostgreSQL and ease of use in managing database models.
- **PostgreSQL:** Opted for its reliability and support for complex queries.

### Data Structure and Storage
- **Post Struct:** Represents a post entry with attributes for title, points, comments, source, and filter applied.
- **UsageData Struct:** Captures metadata for auditing and analysis purposes.
- **Database Schema:** Defined to accommodate the necessary data attributes and relationships.

### Error Handling
- Implemented comprehensive error handling to ensure robustness during scraping, filtering, and database operations.

### Performance Considerations
- **Concurrency:** Colly handles concurrency internally, which aids in efficient data fetching.
- **Data Volume:** Limited to the first 30 entries to manage memory and processing time effectively.

### Testing
- Automated testing is recommended to ensure the correctness of scraping and filtering logic. Testing should cover edge cases like empty responses, network errors, and invalid data formats.

## Consequences

### Pros
- Efficient data scraping and filtering due to Colly's capabilities and Golang's performance.
- Clear separation of concerns with distinct methods for scraping, filtering, and storage.
- Robust data storage and usage tracking using PostgreSQL and GORM.

### Cons
- Dependence on external libraries may require monitoring for updates or changes.
- PostgreSQL requires additional setup and management, which may increase deployment complexity.

## Next Steps
1. Implement automated tests for the scraper and filters.
2. Review and optimize code for performance and maintainability.
3. Document and review the implementation with stakeholders to ensure all requirements are met.

## References
- [Colly Documentation](https://pkg.go.dev/github.com/gocolly/colly/v2)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

This ADR will be updated as necessary to reflect changes in the implementation or technology choices.
