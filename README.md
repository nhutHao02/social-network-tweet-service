# social-network-tweet-service
## Project Summary
This is project about social network that allows users to share content, images, and emotions, and have real-time communication capabilities, while ensuring high performance, security, and scalability using the microservices architecture.

#### Technologies:
- Back-end:
  - Language: Go.
  - Frameworks/Platforms: Gin-Gonic, gRPC, Swagger, JWT, Google-Wire, SQLX, Redis, RabbitMQ, Zap, WebSocket.
  - Database: MariaDB, MongoDB.
- Front-end:
  - Language: JavaScript.
  - Frameworks/Platforms: React, Tailwind CSS, FireBase.

## The project includes repositories
- [common-service](https://github.com/nhutHao02/social-network-common-service)
- [user-service](https://github.com/nhutHao02/social-network-user-service)
- [tweet-service](https://github.com/nhutHao02/social-network-tweet-service)
- [chat-service](https://github.com/nhutHao02/social-network-chat-service)
- [notification-service](https://github.com/nhutHao02/social-network-notification-service)
- [Front-end-service (in progress)](https://github.com/nhutHao02/)

## This service
This is the service that provides the APIs related to the Post.

## ER Diagram
![ER Diagram](https://github.com/user-attachments/assets/c1b994a5-5592-4e37-9f03-af0822ce453f)

## Project structure
```
.
├── config
│   ├── config.go
│   └── local
│       └── config.yaml
├── database
│   └── database.go
├── __debug_bin3107679626
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── grpc
│   │   ├── http
│   │   │   ├── http_server.go
│   │   │   └── v1
│   │   │       ├── route.go
│   │   │       └── tweet_handler.go
│   │   └── server.go
│   ├── application
│   │   ├── imp
│   │   │   └── tweet_service_imp.go
│   │   └── tweet_service.go
│   ├── domain
│   │   ├── entity
│   │   ├── interface
│   │   │   └── tweet
│   │   │       └── tweet_repository.go
│   │   └── model
│   │       ├── notification.go
│   │       ├── tweet-comment.go
│   │       ├── tweet.go
│   │       ├── user.go
│   │       └── websocket.go
│   ├── infrastructure
│   │   └── tweet
│   │       ├── command_repository.go
│   │       └── query_repository.go
│   ├── wire_gen.go
│   └── wire.go
├── main.go
├── Makefile
├── migrations
│   ├── 000001_tweet.down.sql
│   ├── 000001_tweet.up.sql
│   ├── 000002_tweet_image.down.sql
│   ├── 000002_tweet_image.up.sql
│   ├── 000003_tweet_video.down.sql
│   ├── 000003_tweet_video.up.sql
│   ├── 000004_tweet_comment.down.sql
│   ├── 000004_tweet_comment.up.sql
│   ├── 000005_bookmark_tweet.down.sql
│   ├── 000005_bookmark_tweet.up.sql
│   ├── 000006_love_tweet.down.sql
│   ├── 000006_love_tweet.up.sql
│   ├── 000007_repost_tweet.down.sql
│   └── 000007_repost_tweet.up.sql
├── pkg
│   ├── common
│   │   └── response.go
│   ├── constants
│   │   ├── action_constants.go
│   │   └── constants.go
│   ├── redis
│   │   └── redis.go
│   └── websocket
│       └── websocket.go
├── README.md
└── startup
    └── startup.go
```