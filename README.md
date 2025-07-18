# unvoid-chess - Terminal Chess Variant

_unvoid-chess_ is a terminal-based chess variant written in Go, featuring original pieces and gameplay mechanics designed to challenge traditional strategy. Run directly from the command line, it offers an intuitive interface for local turn-based matches.

It’s played on a user-defined square board and features custom units like the Product Owner, Developer, and Designer—each with unique movement rules and tactical implications.

> _This project was developed as a technical test for the company Unvoid._

## Features

- Interactive Shell: Type commands like `select` and `move` to control the game in real time.
- Custom Ruleset: Includes unique mechanics such as sweep captures for Developers and leaping Designers.
- Highlighted Moves: Shows all legal destinations on the board using `x` markers.
- Game State Feedback: Informs players about invalid moves, captures, and turn progression.
- Smart Endgame Detection: Automatically declares victory when a Product Owner is taken.
- Clear Board Visualization: Redraws the board after each action for a clean view.

## Installation and Usage

### Prerequisites

- Requires Go (version 1.18 or newer)
- You may use `make` to simplify building

## Build and Run

Your project structure looks like this:

```
.
├── cmd
│   ├── main.go
│   ├── pieces.go
│   └── utils.go
├── Dockerfile
├── go.mod
├── main.sh
├── out/
├── prepare-to-upload.sh
└── README.md
```

### Running Locally

To simply run the code just run:

```bash
go run ./cmd
```

If you'd like to compile a standalone binary:

```bash
go build -o unvoid ./cmd
./unvoid
```

### Running with Docker

You can also run the game inside a container using the provided Dockerfile.

**1. Build the Image**

```bash
docker build -t unvoid-chess .
```

**2. Run the Container**

```bash
docker run -it unvoid-chess
```

This will launch the game in interactive mode right in your terminal.

Alternatively you can run the shell `main.sh`.

`main.sh` code:

```bash
docker build -t candidate-test . && docker run --rm -it candidate-testa
```

And if you want to make a `.tag.gz` of all the project simply run the `prepare-to-upload.sh`.

---

Let me know if you'd like instructions for installing globally, compiling for 

## Commands

- `select <square>` — Highlights available moves for the piece in that square
- `move <from> <to>` — Attempts a move and prints the result
- `restart` — Restarts with a new board size
- `exit` — Quits the game
- `help` — Shows command list

## Custom Pieces and Movement Rules

The game features three custom pieces per team, each with unique behavior:

### 🧑‍💼 Product Owner (♔)

- Moves exactly one square in any direction (vertical, horizontal, or diagonal), like the king in standard chess.
- If defeated, the game ends immediately and the opponent wins.

### 👩‍💻 Developer (♜)

- Can run up to three squares in any direction, skipping over enemies but never landing on an occupied square.
- If enemies lie between origin and destination, they are captured automatically during the move.

### 🎨 Designer (♘)

- Moves in an L-shape: two squares in one direction and one square perpendicular, like a standard knight.
- Can jump over any piece and land on an empty or enemy square.

Each piece uses Unicode symbols for intuitive display:

- `♔` = Product Owner  
- `♜` = Developer  
- `♘` = Designer

**Objective**: Capture the opponent's Product Owner before yours is captured!



## Example Turn

```
> select b1
    A  B  C D  E  F
 6  .  .  . ♞ ♜ ♚
 5  .  .  . .  .  .
 4  .  x  . .  x  .
 3  .  x  . x  .  .
 2  x  x  x .  .  .
 1  ♔ ♖ ♘ .  .  .
Valid moves for ♖ at B1: [B2 B3 B4 C2 D3 E4 A2]

Turn: White
> 
```

## Acknowledgements

This project was made possible with:

- Go’s standard library
- Terminal ANSI codes for board redraws

## Uninstall

To remove the binary:

```bash
rm unvoid
```

