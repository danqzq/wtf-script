# Unsigned Integers (`uint`)

The `uint` type represents an unsigned 64-bit integer. It supports standard arithmetic operations with specific rules for underflow and type mixing.

## 1. First-Come-First-Served (FCFS) Coercion

In mixed-type arithmetic, the **Left Operand** determines the result type. The Right operand is coerced to match the Left operand.

| Left Operand | Right Operand | Result Type | Behavior |
| :--- | :--- | :--- | :--- |
| `int` | `uint` | `int` | `uint` is constrained to `int`. Large `uint` values may overflow or wrap negatively. |
| `uint` | `int` | `uint` | `int` is converted to `uint`. Negative `int` values **underflow** (wrap to large positive). |
| `uint` | `float` | `uint` | `float` is truncated to `uint`. |
| `float` | `uint` | `float` | `uint` is converted to `float`. |

## 2. Strict Assignment Rules

To prevent accidental errors, strict rules apply when assigning values to `uint` variables.

- **Negative Literals**: Assigning a negative integer or float literal directly (e.g., `uint x = -1;` or `uint x = -10.5;`) causes a **Runtime Error**.
- **Negative Variables**: Assigning a variable holding a negative value (e.g., `uint x = neg_int;`) causes a **Runtime Error**.
- **Computed Underflow**: Calculated values are allowed to underflow (e.g., `uint x = 0 - 1;` results in `MAX_UINT`).

## 3. Implicit Casting & Truncation

Variables implicitly cast assigned values to their declared type.

| Target Variable | Value Type | Behavior |
| :--- | :--- | :--- |
| `int` | `uint` / `float` | Casts to `int` (Potential overflow/truncation). |
| `uint` | `int` / `float` | Casts to `uint`. Negative values cause **Runtime Error** (if literal/variable) or Underflow (if computed). |
| `float` | `int` / `uint` | Casts to `float` (Implicit widening). |
