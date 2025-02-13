# MDCLI

**mdcli** is a simple command-line tool for processing Markdown files.

##  Features
- Accepts a Markdown file as an argument
- Reads and processes the file
- Supports flag-based and positional arguments

##  Installation
Clone the repository and build the project:

```sh
git clone https://github.com/tacheraSasi/mdcli.git
cd mdcli
make build_(your os) #e.g build_linux
go build -o mdcli
```

##  Usage

Run `mdcli` with a Markdown file as an argument:

```sh
./mdcli filename.md
```

or using a flag:

```sh
./mdcli -file=filename.md
```

## ⚡ Example Output
```sh
Processing file: filename.md
```

##  License
This project is licensed under the MIT License.

---
Made with ❤️ in Go by ***Tachera SASI***.