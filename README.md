<div align="center">
    <img src="https://github.com/codern-org/codern/assets/17198802/43ae20dc-74ba-42d3-9a5d-c77dde27217c" width="96" />
    <h2>Codern</h2>
    <p><b>An open-source Coding Platform</b></p>
</div>

**Codern** is an open-source coding platform that focuses on providing a seamless setup for coding event organizers and end-users, enabling them to enhance their programming skills through interactive coding experiences.

> [!NOTE]
> This repository contains the API Server of Codern only. If you are interested in exploring how the Codern user interface was created, please visit the link [here](https://github.com/codern-org/ui).

> This project also serves as an experimental playground for the founders of Codern, allowing them to learn and develop a robust system as part of their Computer Engineering senior project.

## Proposal

### Architecture

<details><summary>Click to expand!</summary>

```mermaid
flowchart TB
    subgraph route[Route]
        swaggerRoute[Swagger route]
        apiRoute[Public API route]
        fallbackRoute[Fallback route]

        apiRoute -.->|Not match any routes| fallbackRoute
        swaggerRoute
    end

    subgraph middleware[Middleware]
        loggerMiddleware[Logger middleware]
    end

    subgraph controller[Controller]
        authController[Auth Controller]
    end

    subgraph usecase[Usecase]
        authUsecase[Auth usecase]
        googleUsecase[Google usecase]
        sessionUsecase[Session usecase]
        userUsecase[User usecase]
    end

    subgraph repository[Repository]
        sessionRepository[Session Repository]
        userRepository[User Repository]
    end

    subgraph fiber[Fiber]
        route
        middleware
        controller
        usecase
        repository

        apiRoute --> loggerMiddleware
        loggerMiddleware --> controller

        authController --> authUsecase
        authController --> googleUsecase

        authUsecase --> googleUsecase
        authUsecase --> sessionUsecase
        authUsecase --> userUsecase

        %% For alignment
        authController --> userUsecase

        sessionUsecase --> sessionRepository
        userUsecase --> userRepository
    end

    style fiber fill:#dfe6e9
    style route fill:#b2bec3
    style controller fill:#b2bec3
    style usecase fill:#b2bec3
    style repository fill:#b2bec3

    style platform fill:#dfe6e9
    style database fill:#b2bec3

    loggerMiddleware -->|Measurement executation time| influxdb
    grafana -.->|Visualization| influxdb
    repository ---> mysql

    config[YAML file]
    fiber -->|Load configuration| config
    vault -->|Render template| config

    subgraph platform[Other platform]
        grafana[Grafana]
        vault[Vault]

        subgraph database[Database]
            mysql[(MySQL)]
            influxdb[(InfluxDB)]
        end
    end
```

#### We don't need Microservice

Behind the scenes, the Codern API server was built using lighting-fast web framework, [Fiber](https://docs.gofiber.io/). Our codebase was designed with a monolithic architecture. Previosuly, we adopted a Microservice architecture for building everything (see [legacy version](https://github.com/codern-org/legacy)), but we eventually made the decision to switch back to a monolith. We found that the Microservice architecture didn't provide us any significant advantages, only introducing development difficulties. As the result, **we opted for the more streamlined monolithic approach.**

#### Clean architecture

This project follows the Clean Architecture principles, ensuring a modular and maintainable codebase. With clear separation of core business logic from external dependencies, it promotes flixibility and scalability. This approach also facilitates easier testing.

_Thanks to Uncle Bob, for the [article](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) of Clean Architecture._

#### Logging & Measurement

This project utilizes [InfluxDB](https://www.influxdata.com/), [Prometheus](https://prometheus.io/), [Grafana](https://grafana.com/) and to achieve a robust log and metric management system. With InfluxDB for time-series storage, Prometheus for monitoring, and Grafana for visualization, we ensure optimal performance and proactive issue detection.

</details>

## Contribution

### Prerequisite

We highly recommended [Visual Studio Code](https://code.visualstudio.com/) for development with a [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension.

## License

Copyright (c) Vectier. All rights reserved.
