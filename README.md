# 🤯 WTFScript

**WTFScript** is a minimal, experimental scripting language with **randomization as a key feature**.

> For the record, **WTFScript** stands for "Wild Type Factory Script" – a chaotic playground for exploring types, randomness, and scripting.

Using it will make you go "WTF?" and make you question your sanity, but in a fun way!

![wtf](https://github.com/user-attachments/assets/632d31fe-cccf-44b0-bd2c-e658507cdb13)

---

## 🚀 Features

- Variable declarations with random initialization
- Type support: `int`, `uint`, `float`, `unofloat`, `bool`, `string`
- Arithmetic operations: `+ - * /` with parentheses
- Built-in functions:
    - `print(args)` – prints arguments (variables or literals)
    - `seed(int)` – sets the randomness seed
- Range-based type declarations, e.g. `int(0, 1000) x;`

---

## 🔧 Project Structure

```
wtf-script/
├── builtins/       # Built-in functions implementation
├── cmd/            # CLI entry point
├── config/         # Configuration system
├── docs/           # Language specification and roadmap
├── examples/       # Sample WTFScript programs
├── interpreter/    # AST, parser, lexer, interpreter logic (including test files)
├── types/          # Type definitions and variable system
├── go.mod
└── README.md
```

---

## 📦 Installation

### 1. Clone the repo

```bash
git clone https://github.com/danqzq/wtf-script.git
cd wtf-script
```

### 2. Build

Requires **Go 1.22+**:

```bash
go build -o wtf ./cmd/wtf/main.go
```

### 3. Run an example

```bash
./wtf examples/intro.wtf
```

---

## ⚙️ Configuration

WTFScript supports custom configuration via JSON files to define default random ranges for all types.

### Usage

```bash
./wtf --config config.json script.wtf
```

### Configuration Options

Create a `config.json` file:

```json
{
  "int": {
    "min": -100,
    "max": 100
  },
  "float": {
    "min": -50.0,
    "max": 50.0
  },
}
```

**Default values** (when no config is provided):
- `int`: -1000 to 1000
- `uint`: 0 to 2000
- `float`: -1000.0 to 1000.0
- `unofloat`: 0.0 to 1.0
- String length: 10 characters

> See [config.json](config.json) for a complete example configuration file.

---

## ✨ Example

```wtf
seed(42);
int(1, 100) x;
float pi = 3.14;
string name;
bool flag;

print(x);
print(pi);
print(name);
print(flag);

print("Seeded randomness!");
```

Possible output:

```
27
3.14
YvH4wqPj
true
Seeded randomness!
```

> See more example scripts under [`examples/`](examples).

---

## 📝 Specification

See [`docs/spec.md`](docs/spec.md) for the full language specification.

---

## 🌐 Roadmap

* [x] MVP with variable declarations and print
* [x] Arithmetic operations with operator precedence
* [x] Proper lexer and AST implementation
* [x] Branching: `if`, `else` (+ random branching with `ifrand`)
* [ ] Loops: `while`, `for` (+ random loops)
* [ ] Functions with parameters and returns
* [ ] Arrays and maps
* [ ] REPL mode
* [ ] Syntax highlighting plugin for VSCode
* [x] Web playground

---

## 🤝 Contributing

Pull requests are welcome. Feel free to file issues for bugs, suggestions, or chaotic ideas.

---

## 📜 License

This project is under the [MIT License](LICENSE)

---

> *WTFScript: Because determinism is overrated. Enjoy the chaos!* 🎉
