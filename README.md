# Basic kredit-plusger User, Admin RBAC with Posts

This project is a create a kredit-plus consisting of an HTTP API and a database using GoLang. You should be able to demonstrate working code and document how to run the application.

## Prerequisites

Before getting started, ensure that you have the following installed:

- [Docker](https://www.docker.com/)
- [golangci-lint](https://golangci-lint.run/)
- [Goose](https://github.com/pressly/goose)

## Getting Started

To get started with the project, follow the steps below:

1. Clone the repository:

   ```bash
   git clone <repository-url>

2. Change into the project directory:
   ```bash
    cd <project-directory>

3. Install project dependencies:
    ```bash 
    go mod download

4. Create the required environment variables or configuration files.

5. Run the linting process to ensure code quality:
    ```bash
    make lint

6. Start the project using Docker:
    ```bash
    make start

    This command will build and start the project in detached mode.

7. To start only the PostgreSQL database server:
    ```bash
    make db-start

8. To create a new migration file:
    ```bash
    make migration
    
    Follow the prompts to provide a name for the migration file.

9. To stop the project and remove associated volumes:
    ```bash
    make down
    
    This command will stop the project and erase the Docker containers and associated volumes.

## Contributing
Please feel free to contribute to this project by forking the repository and creating a pull request.

## License
This project is licensed under the [License Name] - see the LICENSE file for details.

Make sure to replace `<repository-url>`, `<project-directory>`, and `[License Name]` with the appropriate values for your project.

This `README.md` file provides an overview of the project, lists the prerequisites, and guides users through the steps to get started with the project using the provided Makefile commands. It also includes sections for contributing and license information.

Feel free to modify the content and structure of the `README.md` file according to your project's specific requirements.

