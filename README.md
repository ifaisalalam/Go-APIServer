## [Template] GO API Server with Swagger

### Prerequisites

The API server files and models are generated using [Swagger Codegen](https://swagger.io/tools/swagger-codegen/).
Please install the Swagger CLI to generate the server and model code files.

Visit the official page for instructions on [installing the Go Swagger CLI](https://goswagger.io/install.html) tool.

### Setup Instructions

1. Clone the repository

    ```shell
    git clone git@github.com:/ifaisalalam/Go-awesome-service
    ```

2. Run Swagger Codegen to generate the server and model files

    ```shell
    make build
    ```

3. To start the dev server

    ```shell
    make server
    ```
