# 📜 WTFScript Language Specification

## 🎯 Overview

WTFScript is a minimal, chaos-inspired scripting language where **randomization is a key feature**. Declaring a variable without an explicit value assigns it a random one within its type’s range.

---

## 🔤 Supported Types

| Type                      | Description                            | Default Range / Behavior                                           |
|---------------------------|----------------------------------------|--------------------------------------------------------------------|
| `int`                     | 64-bit signed integer                  | Random between -1000 and 1000<br>[-1000; 1000]                     |
| [`uint`](uint.md)         | 64-bit unsigned integer                | Random between 0 and 2000<br>[0; 2000]                             |
| `float`                   | 64-bit float                           | Random between -1000.0 and 1000.0 (exclusive)<br>[-1000.0; 1000.0) |
| [`unofloat`](unofloat.md) | Unit float (strictly between 0 and 1)  | Random between 0.0 and 1.0 (exclusive)<br>[0.0; 1.0)               |
| `string`                  | Random alphanumeric string             | Random 10 characters long, consisting of alphanumeric characters   |
| `bool`                    | Random true or false                   | Random true or false (50% chance)                                  |

> Note: The default range is configurable by creating a `config.json` file in the working directory (see [here](../README.md#configuration-options) for details).

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

## � Comparison Operators

Supported comparison operators:

* Equal: `==`
* Not equal: `!=`
* Less than: `<`
* Less than or equal: `<=`
* Greater than: `>`
* Greater than or equal: `>=`
* Logical NOT: `!`

Comparisons return a `bool` value and work with all numeric types (`int`, `uint`, `float`, `unofloat`), `string`, and `bool`.

Example:

```wtf
int x = 10;
int y = 5;

bool isGreater = x > y;  // true
bool isEqual = x == y;    // false

if (x > y) {
    print("x is greater than y");
}
```

---

## 🧠 Logical Operators

WTFScript supports logical operators for combining boolean expressions:

* Logical AND: `&&`
* Logical OR: `||`
* Logical NOT: `!`

### Short-Circuit Evaluation

Both `&&` and `||` use **short-circuit evaluation**:
- `&&`: If the left operand is false, the right operand is **not evaluated**
- `||`: If the left operand is true, the right operand is **not evaluated**

This prevents unnecessary computation and potential errors.

### Operator Precedence

From highest to lowest:
1. `!` (NOT)
2. Comparison operators (`==`, `!=`, `<`, `<=`, `>`, `>=`)
3. `&&` (AND)
4. `||` (OR)

### Truthiness

When non-boolean values are used with logical operators, they are converted to boolean:
- **Truthy**: Non-zero numbers, non-empty strings
- **Falsy**: `0`, `0.0`, empty string `""`, `false`

Example:

```wtf
bool a = true;
bool b = false;

// Basic logical operations
bool and_result = a && b;  // false
bool or_result = a || b;   // true
bool not_result = !a;      // false

// With comparisons
int x = 10;
int y = 20;
int z = 15;

bool inRange = x < z && z < y;  // true (15 is between 10 and 20)

// Complex expressions with precedence
bool complex = true || false && false;  // true (evaluated as: true || (false && false))

// Short-circuit prevents errors
bool safe = false && x / 0 == 0;  // false (division never happens)

// Using numeric truthiness
int num = 5;
bool truthy = num && true;  // true (5 is truthy)
```

---

## 🔀 Control Flow

### If Statements

WTFScript supports standard conditional statements with C-like syntax:

```wtf
if (condition) {
    // code block
}
```

### Else If and Else

Chain multiple conditions:

```wtf
int score = 85;

if (score >= 90) {
    print("Grade: A");
} else if (score >= 80) {
    print("Grade: B");
} else if (score >= 70) {
    print("Grade: C");
} else {
    print("Grade: F");
}
```

### 🎲 Random Conditionals: `ifrand`

WTFScript introduces `ifrand` - a conditional statement with **probabilistic execution**:

**Default 50% probability:**
```wtf
ifrand {
    print("This has a 50% chance of executing");
}
```

**Custom probability (0.0 to 1.0):**
```wtf
ifrand(0.8) {
    print("This has an 80% chance of executing");
}

ifrand(0.1) {
    print("This has a 10% chance of executing");
}
```

**With else blocks:**
```wtf
ifrand(0.5) {
    print("True branch");
} else {
    print("False branch");
}
```

**Chaining ifrand:**
```wtf
ifrand(0.3) {
    print("30% chance");
} else ifrand(0.6) {
    print("~42% chance (60% of remaining 70%)");
} else {
    print("~28% chance");
}
```

**Mixing regular if and ifrand:**
```wtf
int x = 10;

if (x > 5) {
    print("x is greater than 5");
    ifrand(0.7) {
        print("And we got lucky! (70% chance)");
    }
}
```

---

## �🚫 Error Handling

* Division by zero produces a runtime error.
* Assigning incompatible types (e.g. `uint x = -5;`) produces a parse or evaluation error.
* Undefined variables raise runtime errors.

> See [`examples/errors.wtf`](../examples/errors.wtf) for examples of possible errors.

---

## 🔮 Future Planned Features

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