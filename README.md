# Elitecode CLI

brew formula for the elitecode website

## Installation

### Prerequisites

- Go 1.21 or later
- Docker
- Git

### From Source

1. Clone the repository:
```bash
git clone https://github.com/IshaanNene/EliteCode-brew.git
```

2. Navigate to the cloned directory:
```bash
cd EliteCode-brew
```

3. Build the binary:
```bash
go build -o elitecode
```

4. Install globally:
```bash
sudo mv elitecode /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/IshaanNene/EliteCode-brew@latest
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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Firebase](https://firebase.google.com/) - Backend services
- [Docker](https://www.docker.com/) - Container runtime 