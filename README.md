# Basic Golang API for Kredit-Plus

This project serves as a platform to showcase and demonstrate my skills and capabilities. It involves the development of a Kredit Plus application, incorporating an HTTP API and a Golang Based Application. Through this project, I aim to highlight my proficiency and expertise in the field.

## Prerequisites

Before you begin, please make sure you have the following tools and dependencies installed:

- [Docker](https://www.docker.com/)
- [golangci-lint](https://golangci-lint.run/)
- [Goose](https://github.com/pressly/goose)

## Getting Started

To get started with the project, follow these steps:

1. Clone the repository:

   ```bash
   git clone <repository-url>
   ```

2. Change to the project directory:

   ```bash
   cd <project-directory>
   ```

3. Install project dependencies:

   ```bash 
   go mod download
   ```

4. Create the necessary environment variables or configuration files.

5. Run the linting process to ensure code quality:

   ```bash
   make lint
   ```

6. Start the project using Docker:

   ```bash
   make start
   ```

   This command will build and start the project in detached mode.

7. To start only the PostgreSQL database server:

   ```bash
   make db-start
   ```

8. To create a new migration file:

   ```bash
   make migration
   ```

   Follow the prompts to provide a name for the migration file.

9. To stop the project and remove associated volumes:

   ```bash
   make down
   ```

   This command will stop the project and erase the Docker containers and associated volumes.

Feel free to reach out if you have any questions or need further assistance with the setup. We are here to help you get started with your Kredit-Plus project.