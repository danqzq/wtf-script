# 🤯 WTFScript

**WTFScript** is a minimal, experimental scripting language with **randomization as a key feature**.  

> For the record, **WTFScript** stands for "Wild Type Factory Script" – a chaotic playground for exploring types, randomness, and scripting.

Using it will make you go "WTF?" and make you question your sanity, but in a fun way!

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
├── docs/           # Language specification and roadmap
├── examples/       # Sample WTFScript programs
├── interpreter/    # AST, parser, lexer, interpreter logic
├── types/          # Type definitions and variable system
├── go.mod
├── main.go         # CLI entry point
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
go build -o wtf ./main.go
```

### 3. Run an example

```bash
./wtf examples/intro.wtf
```

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
* [ ] Proper lexer and AST implementation
* [ ] Branching: `if`, `else` (+ random branching with `ifrand`)
* [ ] Loops: `while`, `for` (+ random loops)
* [ ] Functions with parameters and returns
* [ ] Arrays and maps
* [ ] REPL mode
* [ ] Syntax highlighting plugin for VSCode
* [ ] Web playground

---

## 🤝 Contributing

Pull requests are welcome. Feel free to file issues for bugs, suggestions, or chaotic ideas.

---

## 📜 License

This project is under the [MIT License](LICENSE)

---

> *WTFScript: Because determinism is overrated. Enjoy the chaos!* 🎉