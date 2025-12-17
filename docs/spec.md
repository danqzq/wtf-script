# 📜 WTFScript Language Specification

## 🎯 Overview

WTFScript is a minimal, chaos-inspired scripting language where **randomization is a key feature**. Declaring a variable without an explicit value assigns it a random one within its type’s range.

---

## 🔤 Supported Types

| Type        | Description                              | Default Range / Behavior                  |
|-------------|------------------------------------------|-------------------------------------------|
| `int`       | 64-bit signed integer                    | Random between -1000 and 1000 (inclusive)            |
| [`uint`](uint.md)      | 64-bit unsigned integer                  | Random between 0 and 2000 (inclusive)            |
| `float`     | 64-bit float                             | Random between -1000.0 and 1000.0 (inclusive)            |
| [`unofloat`](unofloat.md)  | Unit float (strictly between 0 and 1)            | Random between 0.0 and 1.0 (inclusive)                |
| `string`    | Random alphanumeric string               | Random 10 characters long, alphanumeric |
| `bool`      | Random true or false                     | Random true or false (50% chance)         |

> Note: The default range is configurable by adjusting constant limits set in [`config/config.go`](../config/config.go).

---

### 🔢 Ranged Type Declarations

Specify min and max for numeric types:

```wtf
int(0, 100) x;
uint(10, 500) y;
````

If omitted, defaults to the default type range as specified above.

---

## 🔧 Built-in Functions

### 📤 `print(args...)`

Prints variables, numbers, strings, or booleans.

Example:

```wtf
print(42, 69);
print("Hello World!");
print(x);
```

### 🎲 `seed(int)`

Sets the random seed for reproducible runs.

Example:

```wtf
seed(12345);
```

---

## ➗ Arithmetic Operations

Supported:

* Addition: `+`
* Subtraction: `-`
* Multiplication: `*`
* Division: `/` (integer division for ints and uints)
* Parentheses: `()`

Example:

```wtf
int result = (10 + 5) * 2;
print(result);
```
### 🔄 Type Coercion & Strictness

WTFScript enforces **Foundational Type Strictness** with specific coercion rules:

**1. No Implicit Conversions in Assignment:**
Variables are strictly typed. You cannot assign a `string` to an `int` or a `bool` to a `float`.

**2. First Come First Served (FCFS) Arithmetic:**
For mixed numeric types (`int` and `float`), the **left-hand operand determines the result type**.
* `int + float` → `int` (the float is treated as an int)
* `float + int` → `float` (the int is treated as a float)

**3. Boolean Isolation:**
Booleans cannot be added to integers or strings. Logical operations are strict.

**4. String Concatenation:**
Strings can only be concatenated with other strings. `string + int` is invalid.

Example:

```wtf
float a = 10.5;
int b = 5;

print(a + b); // Result: 15.5 (float)
print(b + a); // Result: 15 (int)
```
---

## 🚫 Error Handling

* Division by zero produces a runtime error.
* Assigning incompatible types (e.g. `uint x = -5;`) produces a parse or evaluation error.
* Undefined variables raise runtime errors.

> See [`examples/errors.wtf`](../examples/errors.wtf) for examples of possible errors.

---

## 🔮 Future Planned Features

* **Conditionals:** `if`, `else`, `ifrand` (randomized conditionals)
* **Loops:** `while`, `for` (with random iterations)
* **User-defined functions:** with parameters and return values
* **Arrays and maps**
* **Modules and imports**
* **REPL mode**

---

## ✨ Philosophy

WTFScript is built on the idea that **randomness is a core programming concept**, encouraging creativity, chaos testing, and philosophical questioning of determinism in code.

---

*This spec is evolving as the language develops.*