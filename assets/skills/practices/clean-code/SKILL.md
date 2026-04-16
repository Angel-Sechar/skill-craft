---
name: clean-code
description: >
  Apply Clean Code principles to all generated and reviewed code. Enforce
  naming, function size, comment quality, error handling, and readability.
  Triggers on clean code, readable code, naming conventions, code quality,
  refactor for readability.
category: practices
conflicts: []
version: 1.0.0
license: MIT
---

You are enforcing Clean Code. Code is read far more than it is written. Every name, every function, every comment must earn its place.

## Naming — reveal intent

Names must answer what this is, why it exists, and how it is used.

```csharp
// WRONG — names that require a comment to understand
int d;           // elapsed time in days
bool flag;       // whether user is active
void Process();

// CORRECT — names that explain themselves
int elapsedDays;
bool isUserActive;
void ConfirmOrder();
void SendWelcomeEmail(string recipientEmail);
```

## Functions — one level of abstraction, one job

```python
# WRONG — three levels of abstraction in one function
def process_order(order_data):
    conn = psycopg2.connect(...)
    raw = conn.execute("SELECT ...")
    total = sum(i['price'] for i in raw)
    if total > 100: discount = total * 0.1
    requests.post("https://notify/", ...)

# CORRECT — each function at one level of abstraction
def process_order(order: Order) -> ProcessedOrder:
    validate_order(order)
    pricing = calculate_pricing(order)
    result = finalize_order(order, pricing)
    notify_customer(result)
    return result
```

## Function size — max 20 lines

If a function exceeds 20 lines it is doing more than one thing. Extract.

## Comments — explain why, never what

```csharp
// WRONG — comment restates the code
// increment counter by 1
counter++;

// WRONG — commented-out dead code
// var old = OldCalculation(x);
var result = NewCalculation(x);

// CORRECT — comment explains non-obvious reasoning
// Delay 50ms to allow the hardware write buffer to flush.
// Removing this delay causes intermittent data loss on slow disks.
await Task.Delay(50);
```

## Error handling

```typescript
// WRONG — swallowed exception
function findUser(id: string): User | null {
  try {
    return db.query(id);
  } catch (e) {
    return null;
  }
}

// CORRECT — explicit and informative
function findUser(id: string): User {
  const user = db.query(id);
  if (!user) throw new UserNotFoundException(id);
  return user;
}
```

## No magic numbers

```csharp
// WRONG
if (order.Lines.Count > 10)
    ApplyBulkDiscount(order, 0.15m);

// CORRECT
const int BulkOrderThreshold = 10;
const decimal BulkDiscountRate = 0.15m;

if (order.Lines.Count > BulkOrderThreshold)
    ApplyBulkDiscount(order, BulkDiscountRate);
```

## Rules to enforce always

- No abbreviations — ord becomes order, custId becomes customerId
- No double negatives — isEnabled not isNotDisabled
- Class names are nouns, method names are verbs
- Max 3 parameters per function — use parameter object for more
- Delete dead code — version control is the history keeper
- Boolean parameters are a design smell — split into two functions

## Red flags — stop and warn

- Functions longer than 20 lines
- Names that require a comment to explain
- Commented-out code blocks
- Magic numbers without named constants
- Nested conditionals deeper than 3 levels
- catch block that swallows exceptions silently
