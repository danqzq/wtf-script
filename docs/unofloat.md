# Unit Float (`unofloat`)

The `unofloat` type represents a floating-point number constrained to the range **[0.0, 1.0]**. It is designed for normalized values such as probabilities, ratios, and percentages represented as decimals.

## 1. First-Come-First-Served (FCFS) Coercion with Clamping

In mixed-type arithmetic, the **Left Operand** determines the result type. The Right operand is coerced to match the Left operand. When the result is stored in a `unofloat`, it is **clamped** to the valid range [0.0, 1.0].

| Left Operand | Right Operand | Result Type | Behavior |
| :--- | :--- | :--- | :--- |
| `unofloat` | `int` | `unofloat` | `int` is converted to `float`, then result is clamped to [0.0, 1.0]. |
| `unofloat` | `uint` | `unofloat` | `uint` is converted to `float`, then result is clamped to [0.0, 1.0]. |
| `unofloat` | `float` | `unofloat` | Result is clamped to [0.0, 1.0]. |
| `int` | `unofloat` | `int` | `unofloat` is converted to `int` (truncated). |
| `uint` | `unofloat` | `uint` | `unofloat` is converted to `uint` (truncated). |
| `float` | `unofloat` | `float` | `unofloat` is converted to `float` (no clamping). |

**Note on Clamping**: When arithmetic results exceed the valid range, they are automatically clamped:
- Values < 0.0 become 0.0
- Values > 1.0 become 1.0
- Example: `unofloat x = 0.2; int y = 1; unofloat z = x + y;` results in `z = 1.0`

## 2. Strict Assignment Rules

To prevent accidental errors, strict rules apply when assigning values to `unofloat` variables from literals and existing variables.

- **Out-of-Range Literals**: Assigning a literal value outside [0.0, 1.0] (e.g., `unofloat x = 1.1;` or `unofloat x = -0.5;`) causes a **Runtime Error**.
- **Out-of-Range Variables**: Assigning a variable holding a value outside [0.0, 1.0] causes a **Runtime Error**.
  - `int` variables with value < 0 or > 1 → **Runtime Error**
  - `uint` variables with value > 1 → **Runtime Error**
  - `float` variables with value < 0.0 or > 1.0 → **Runtime Error**
- **Computed Values with Clamping**: Calculated values are allowed and automatically clamped to [0.0, 1.0].
  - Example: `unofloat x = 0.5 + 0.8;` results in `x = 1.0` (clamped)
  - Example: `unofloat y = 0.2 - 0.5;` results in `y = 0.0` (clamped)

## 3. Implicit Casting & Clamping

Variables implicitly cast assigned values to their declared type. For `unofloat`, this includes clamping computed values.

| Target Variable | Value Type | Behavior |
| :--- | :--- | :--- |
| `int` | `unofloat` | Casts to `int` (truncates to 0 or 1). |
| `uint` | `unofloat` | Casts to `uint` (truncates to 0 or 1). |
| `float` | `unofloat` | Casts to `float` (widens without clamping). |
| `unofloat` | `int` | **Runtime Error** if value < 0 or > 1 (variable/literal). Clamps if computed. |
| `unofloat` | `uint` | **Runtime Error** if value > 1 (variable/literal). Clamps if computed. |
| `unofloat` | `float` | **Runtime Error** if value < 0.0 or > 1.0 (variable/literal). Clamps if computed. |

## 4. Common Use Cases

The `unofloat` type is ideal for:
- **Probabilities**: Representing values that must be between 0 and 1
- **Ratios and Percentages**: `unofloat ratio = completed / total;`
- **Normalized Values**: Alpha values, confidence scores, completion percentages
- **Safe Arithmetic**: Automatic clamping prevents invalid normalized values

## 5. Examples

```wtf
// Valid assignments
unofloat probability = 0.75;
unofloat ratio = 0.5;
unofloat zero = 0.0;
unofloat one = 1.0;

// Computed values with clamping
unofloat a = 0.2 + 0.7;     // Result: 0.9
unofloat b = 0.1 - 0.2;     // Result: 0.0 (clamped)
unofloat c = 0.6 + 0.8;     // Result: 1.0 (clamped)

// Ratio calculation (common use case)
float completed = 15.0;
float total = 20.0;
unofloat progress = completed / total;  // Result: 0.75

// Mixed-type arithmetic with clamping
unofloat x = 0.3;
int increment = 1;
unofloat result = x + increment;  // Result: 1.0 (clamped)

// Runtime Errors (commented out)
// unofloat invalid1 = 1.5;      // Error: out of range
// unofloat invalid2 = -0.1;     // Error: out of range
// int negative = -1;
// unofloat invalid3 = negative; // Error: negative variable
// uint large = 100;
// unofloat invalid4 = large;    // Error: value > 1
```

## 6. Type Conversion Summary

| Conversion | Behavior |
| :--- | :--- |
| `unofloat` → `int` / `uint` | Truncates to 0 or 1 |
| `int` / `uint` → `unofloat` | Runtime Error if out of [0, 1] range for literal/variable; clamps if computed |
| `float` → `unofloat` | Runtime Error if out of [0.0, 1.0] range for literal/variable; clamps if computed |
| `unofloat` → `float` | Direct cast (no restrictions) |
