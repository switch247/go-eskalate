task-manager/
├── Delivery/
│   ├── main.go
│   ├── controllers/
│   │   └── controller.go
│   └── routers/
│       └── router.go
├── Domain/
│   └── Domain.go
├── Infrastructure/
│   ├── auth_middleWare.go
│   ├── jwt_service.go
│   └── password_service.go
├── Repositories/
│   ├── task_repository.go
│   └── user_repository.go
└── Usecases/
    ├── task_usecases.go
    └── user_usecases.go


client
client        Router -> Controller-> UseCase -> Repository   ====>DataBase
client                ----------------------------------
client                               Domain



old : [controller]-> [repository]  
new : [controller] -> [usecase] -> [repository] 
old : [model]
new : [Domain]