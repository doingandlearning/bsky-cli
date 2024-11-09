# BSky CLI v0.1

This is a command-line interface (CLI) tool for interacting with the BSky social platform. It allows you to authenticate and create posts on BSky.

It is very much a prototype right now.

## TODO:

- [ ] Work on build
- [ ] Improve documentation
- [ ] Add unit tests
- [ ] Implement error handling
- [x] Add other commands (fetch, stream)
- [ ] Anything else?

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

## Build 

```sh
go build -o bsky ./cmd/cli
```

## Usage 

To create a post on BSky, run the following command:

```sh
./bsky -content "Your post content here"
```

To list the last 10 posts from users in your feed:

```sh
./bsky -fetch
```

To stream posts use the following command.

```sh
./bsky -stream
```

It defaults to 10 seconds but you can pass and optional interval flag if you'd like more or less frequent updates.

```sh
./bsky -stream -interval=1s
```

## Usage with go run (dev mode)

To create a post on BSky, run the following command:

```sh
go run cmd/cli/*.go -content "Your post content here"
```

To list the last 10 posts from users in your feed:

```sh
go run cmd/cli/*.go -fetch
```

## License

This project is licensed under the MIT License.
