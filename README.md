# u-chess - Terminal Chess Variant

_u-chess_ is a terminal-based chess variant written in Go, featuring original pieces and gameplay mechanics designed to challenge traditional strategy. Run directly from the command line, it offers an intuitive interface for local turn-based matches.

It’s played on a user-defined square board and features custom units like the King, Tower, and Knight — each with unique movement rules and tactical implications.

## Features

- Interactive Shell: Type commands like `select` and `move` to control the game in real time.
- Custom Ruleset: Includes unique mechanics such as sweep captures for Towers and leaping Knights.
- Highlighted Moves: Shows all legal destinations on the board using `x` markers.
- Game State Feedback: Informs players about invalid moves, captures, and turn progression.
- Smart Endgame Detection: Automatically declares victory when a King is taken.
- Clear Board Visualization: Redraws the board after each action for a clean view.

## Installation and Usage

### Prerequisites

- Requires Go (version 1.18 or newer)

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
├── out/
└── README.md
```

### Running Locally

To simply run the code just run:

```bash
go run ./cmd
```

If you'd like to compile a standalone binary:

```bash
go build -o uchess ./cmd
./uchess
```

### Running with Docker

You can also run the game inside a container using the provided Dockerfile.

**1. Build the Image**

```bash
docker build -t uchess .
```

**2. Run the Container**

```bash
docker run -it uchess
```

This will launch the game in interactive mode right in your terminal.

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

### 🧑‍💼 King (♔)

- Moves exactly one square in any direction (vertical, horizontal, or diagonal), like the king in standard chess.
- If defeated, the game ends immediately and the opponent wins.

### 👩‍💻 Tower (♜)

- Can run up to three squares in any direction, skipping over enemies but never landing on an occupied square.
- If enemies lie between origin and destination, they are captured automatically during the move.

### 🎨 Knight (♘)

- Moves in an L-shape: two squares in one direction and one square perpendicular, like a standard knight.
- Can jump over any piece and land on an empty or enemy square.

Each piece uses Unicode symbols for intuitive display:

- `♔` = King  
- `♜` = Tower  
- `♘` = Knight

**Objective**: Capture the opponent's King before yours is captured!



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
rm uchess
```

