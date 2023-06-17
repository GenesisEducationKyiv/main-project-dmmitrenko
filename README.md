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
│   └── models/
│       └── Rate.go
├── repositories/
│   └── emailRepository.go
└── services/
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
docker run -p 8080:8080 --env APIKEY=<API key for sending messages> <your-image-name>
```
