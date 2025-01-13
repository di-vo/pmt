# (P)roject (M)anagement (T)erminal

A simple and lightweight project management app inside your terminal.

Uses the [Bubbletea](https://github.com/charmbracelet/bubbletea) framework.

## Installation

### Using docker

Clone the repository:

```bash
    git clone https://github.com/di-vo/pmt && cd pmt
```

Build the docker image:

```bash
    docker build -t pmt .
```

Run a container:

```bash
    docker run -it -e "TERM=xterm-256color" pmt-test
```

If you want to run a stopped container, you need to do:

```bash
    docker start -ai <name>
```

### From source

Clone the repository:

```bash
    git clone https://github.com/di-vo/pmt && cd pmt
```

Build the project:

```bash
    go build -o pmt .
```
Adjust the output path if you want to build it somewhere else

Run the binary:

```bash
    ./pmt
```
