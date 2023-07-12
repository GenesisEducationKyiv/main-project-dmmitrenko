# CurrencyRateApp
### File tree
``` powershell
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
├── README.md
├── api/
│   ├── controller/
│   │   ├── emailController.go
│   │   └── rateController.go
│   ├── docs/
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── middleware/
│   │   └── exceptionMiddleware.go
│   └── route/
│       └── route.go
├── cmd/
│   └── main.go
├── domain/
│   └── constants.go
├── domain/
│   └── model/
│       └── Rate.go
├── repository/
│   └── emailRepository.go
└── service/
    ├── apiClient.go
    ├── emailService.go
    └── rateService.go
```

### API launch

- Build a Docker image with the appropriate settings. 
```docker
docker build -t <your-image-name> .
```
- Run the container based on the built image.
```docker
docker run -p <your port>:8080 --env-file .env --env APIKEY=<API key for sending messages> <your-image-name>
```
