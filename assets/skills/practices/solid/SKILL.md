---
name: solid
description: >
  Apply all five SOLID principles to every class and module. Enforce single
  responsibility, open/closed, Liskov substitution, interface segregation,
  and dependency inversion. Triggers on SOLID, single responsibility,
  open closed principle, Liskov, interface segregation, dependency inversion.
category: practices
conflicts: []
version: 1.0.0
license: MIT
---

You are enforcing SOLID principles. Every class you write or review must satisfy all five. Call out violations by name before fixing them.

## S — Single Responsibility

One class, one reason to change.

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

Open for extension, closed for modification.

```csharp
// WRONG — every new payment method modifies this class
public decimal ProcessPayment(string method, decimal amount)
{
    if (method == "credit")  return amount * 0.98m;
    if (method == "paypal")  return amount * 0.97m;
    return amount;
}

// CORRECT — new methods extend without modifying
public interface IPaymentProcessor { decimal Process(decimal amount); }
public class CreditCardProcessor : IPaymentProcessor { ... }
public class PayPalProcessor      : IPaymentProcessor { ... }
```

## L — Liskov Substitution

A subtype must be usable wherever the base type is expected.

```csharp
// WRONG — Square breaks Rectangle behavioral contract
public class Square : Rectangle
{
    public override int Width  { set { base.Width = base.Height = value; } }
    public override int Height { set { base.Width = base.Height = value; } }
}

// CORRECT — separate hierarchy
public interface IShape { int Area(); }
public class Rectangle : IShape { public int Width; public int Height; public int Area() => Width * Height; }
public class Square    : IShape { public int Side; public int Area() => Side * Side; }
```

## I — Interface Segregation

Clients should not depend on methods they do not use.

```csharp
// WRONG — fat interface
public interface IWorker { void Work(); void Eat(); void Sleep(); void Charge(); }

// CORRECT — small focused interfaces
public interface IWorkable   { void Work(); }
public interface IFeedable   { void Eat(); void Sleep(); }
public interface IChargeable { void Charge(); }

public class Human : IWorkable, IFeedable {}
public class Robot : IWorkable, IChargeable {}
```

## D — Dependency Inversion

High-level modules must not depend on low-level modules. Both depend on abstractions.

```csharp
// WRONG — depends on concrete implementation
public class OrderService
{
    private readonly SqlOrderRepository _repo = new SqlOrderRepository();
    public void Place(Order order) => _repo.Insert(order);
}

// CORRECT — depends on abstraction, injected
public class OrderService(IOrderRepository repo)
{
    public void Place(Order order) => repo.Save(order);
}
```

## Review checklist

1. Can you state this class single responsibility in one sentence?
2. To add new behavior — do you modify or extend this class?
3. Can you substitute any subclass and have all tests pass?
4. Does this interface have methods the consumer does not use?
5. Does this class instantiate its own infrastructure dependencies?

## Red flags — stop and warn

- Class with more than one reason to change
- Switch or if-else that grows every time a new type is added
- Subclass that throws NotImplementedException for inherited methods
- Interface with 8+ methods
- new ConcreteRepository() inside a business service
