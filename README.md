# 🚢 Battleship Game

A modern implementation of the classic Battleship game, built with Go and React. Challenge your friends in this strategic naval warfare game!

## 🌟 Features

- **Real-time Gameplay**: Play against other players in real-time
- **Modern UI**: Clean and intuitive React-based interface
- **Scalable Backend**: Built with Go for high performance
- **MongoDB Storage**: Persistent game state and player data
- **Docker Support**: Easy deployment and scaling
- **RESTful API**: Well-documented endpoints for game operations

## 🎮 How to Play

1. **Setup**: Place your ships on the grid
2. **Battle**: Take turns firing at your opponent's grid
3. **Win**: Sink all your opponent's ships to claim victory!

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.18 or later (for development)
- Node.js 16 or later (for frontend development)

### Running with Docker

```bash
# Clone the repository
git clone https://github.com/Jagreen1970/battleship.git
cd battleship

# Start the application
docker-compose up -d
```

The application will be available at:
- Game: http://localhost:3000
- MongoDB Express: http://localhost:8081

### Development Setup

```bash
# Start MongoDB
docker-compose up -d mongo-db mongo-express

# Backend
cd backend
go mod download
go run cmd/battleship/main.go

# Frontend
cd frontend
npm install
npm start
```

## 🏗️ Architecture

- **Frontend**: React with TypeScript
- **Backend**: Go with Gin framework
- **Database**: MongoDB
- **Containerization**: Docker
- **API Documentation**: Swagger/OpenAPI

## 📚 API Documentation

The API documentation is available at `/api/docs` when running the application.

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📝 License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the classic Battleship board game
- Built with modern web technologies
- Thanks to all contributors!

## 📞 Support

If you encounter any issues or have questions, please open an issue in the GitHub repository.

---

Made with ❤️ by [Jagreen1970](https://github.com/Jagreen1970) 