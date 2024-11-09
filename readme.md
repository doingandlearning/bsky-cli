# BSky CLI v0.1

This is a command-line interface (CLI) tool for interacting with the BSky social platform. It allows you to authenticate and create posts on BSky.

It is very much a prototype right now.

## TODO:

- [ ] Work on build
- [ ] Improve documentation
- [ ] Add unit tests
- [ ] Implement error handling
- [ ] Add other commands

## Prerequisites

- Go 1.19 or later
- A `.env` file with your BSky username and app password

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/doingandlearning/bsky-cli.git
   cd bsky-cli
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Configuration

Create a `.env` file in the root directory of the project with the following content:

```
USERNAME=your_username
APP_PASSWORD=your_password
```

## Usage

To create a post on BSky, run the following command:

```sh
go run cmd/cli/main.go -content "Your post content here"
```

To list the last 10 posts from users in your feed:

```sh
go run cmd/cli/main.go -fetch
```

## License

This project is licensed under the MIT License.
