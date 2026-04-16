---
name: solid
description: Apply all five SOLID principles to every class and module. Enforce single responsibility, open/closed, Liskov substitution, interface segregation, and dependency inversion. Triggers on: "SOLID", "single responsibility", "open closed principle", "Liskov", "interface segregation", "dependency inversion".
category: practices
conflicts: []
version: 1.0.0
---

You are enforcing SOLID principles. Every class you write or review must satisfy all five. Call out violations by name before fixing them.

## S — Single Responsibility

One class, one reason to change, one actor it serves.

```csharp
// WRONG — three responsibilities in one class
public class UserService
{
    public void Register(string email, string password)
    {
        if (!email.Contains('@')) throw new Exception("Invalid email"); // validation
        var hash = BCrypt.HashPassword(password);                       // security
        _db.Execute("INSERT INTO users...", email, hash);               // persistence
        _smtp.Send(email, "Welcome!");                                  // notification
    }
}

// CORRECT — each class has exactly one job
public class UserRegistrationValidator { public void Validate(RegisterRequest r) {} }
public class PasswordHasher            { public string Hash(string raw) {} }
public class UserRepository            { public void Save(User user) {} }
public class WelcomeEmailSender        { public void Send(string email) {} }
public class UserRegistrationService   { /* orchestrates the four above */ }
```

## O — Open/Closed

Open for extension, closed for modification. Add behavior without changing existing code.

```csharp
// WRONG — every new payment method requires modifying this class
public decimal ProcessPayment(string method, decimal amount)
{
    if (method == "credit")  return amount * 0.98m;
    if (method == "paypal")  return amount * 0.97m;
    if (method == "crypto")  return amount * 0.99m;  // ← modify to add
    return amount;
}

// CORRECT — new payment methods extend without modifying
public interface IPaymentProcessor { decimal Process(decimal amount); }
public class CreditCardProcessor : IPaymentProcessor { ... }
public class PayPalProcessor      : IPaymentProcessor { ... }
// Adding crypto = new class, no existing code touched
```

## L — Liskov Substitution

A subtype must be usable wherever the base type is expected, without breaking behavior.

```csharp
// WRONG — Square breaks Rectangle's behavioral contract
public class Rectangle { public virtual int Width { get; set; } public virtual int Height { get; set; } }
public class Square : Rectangle
{
    public override int Width  { set { base.Width = base.Height = value; } }  // breaks LSP
    public override int Height { set { base.Width = base.Height = value; } }
}

// CORRECT — separate hierarchy
public interface IShape { int Area(); }
public class Rectangle : IShape { public int Width; public int Height; public int Area() => Width * Height; }
public class Square    : IShape { public int Side; public int Area() => Side * Side; }
```

## I — Interface Segregation

Clients should not depend on methods they don't use. Split fat interfaces.

```csharp
// WRONG — fat interface forces implementors to stub unused methods
public interface IWorker { void Work(); void Eat(); void Sleep(); void Charge(); }

// CORRECT — small focused interfaces
public interface IWorkable  { void Work(); }
public interface IFeedable  { void Eat(); void Sleep(); }
public interface IChargeable{ void Charge(); }

public class Human : IWorkable, IFeedable {}
public class Robot : IWorkable, IChargeable {}
```

## D — Dependency Inversion

High-level modules must not depend on low-level modules. Both depend on abstractions.

```csharp
// WRONG — OrderService depends directly on concrete SqlOrderRepository
public class OrderService
{
    private readonly SqlOrderRepository _repo = new SqlOrderRepository();
    public void Place(Order order) => _repo.Insert(order);
}

// CORRECT — depends on abstraction, injected externally
public class OrderService(IOrderRepository repo)
{
    public void Place(Order order) => repo.Save(order);
}
```

## Review checklist — run before finalizing any class

1. Can you state this class's single responsibility in one sentence?
2. To add new behavior — do you modify this class or extend it?
3. Can you substitute any subclass and have all tests pass unchanged?
4. Does this interface have methods the current consumer doesn't use?
5. Does this class instantiate its own infrastructure dependencies?

## Red flags — stop and warn

- Class with more than one reason to change
- Switch/if-else that grows every time a new type is added
- Subclass that throws `NotImplementedException` for inherited methods
- Interface with 8+ methods
- `new ConcreteRepository()` inside a business service
