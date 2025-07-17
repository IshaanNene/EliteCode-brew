# Elitecode CLI

A comprehensive coding practice platform that helps you improve your programming skills through interactive problem-solving.

## Features

- 🔍 Browse and search coding problems
- 📝 Submit solutions in multiple programming languages
- 🐳 Secure code execution in Docker containers
- 📊 Track your progress and statistics
- 🏆 Compete on the leaderboard
- 🔄 Sync solutions with GitHub
- 📱 Cross-platform support

## Installation

### Prerequisites

- Go 1.21 or later
- Docker
- Git

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/elitecode.git
cd elitecode

# Build the binary
go build -o elitecode

# Install globally
sudo mv elitecode /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/yourusername/elitecode@latest
```

## Quick Start

1. Log in to your account:
   ```bash
   elitecode login
   ```

2. Browse available problems:
   ```bash
   elitecode problems
   ```

3. Set up a problem:
   ```bash
   elitecode set_problem two-sum
   ```

4. Submit your solution:
   ```bash
   elitecode submit
   ```

## Commands

### Authentication
- `login` - Log in to your account
- `logout` - Log out from your account
- `github login` - Log in with GitHub

### Problem Management
- `problems` - List available problems
- `search` - Search for problems
- `set_problem` - Set up a problem for solving
- `bookmark` - Manage bookmarked problems

### Solution Management
- `submit` - Submit your solution
- `run` - Run your solution locally
- `push` - Push solution to GitHub

### Statistics
- `stats` - View problem statistics
- `profile` - View your profile
- `leaderboard` - View rankings

## Configuration

Configuration is stored in `~/.elitecode/config.json`. You can override settings using command-line flags or environment variables.

### Environment Variables

- `ELITECODE_CONFIG` - Path to config file
- `ELITECODE_VERBOSE` - Enable verbose output
- `FIREBASE_EMULATOR_HOST` - Firebase emulator host (for development)

## Development

### Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run Firebase emulators:
   ```bash
   firebase emulators:start
   ```

3. Build and run:
   ```bash
   go run main.go
   ```

### Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Firebase](https://firebase.google.com/) - Backend services
- [Docker](https://www.docker.com/) - Container runtime 