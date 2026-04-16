---
name: oop
description: >
  Apply Object-Oriented Programming principles correctly. Encapsulation,
  inheritance, polymorphism, and abstraction used with intention.
  No anemic models, no public setters on domain objects.
  Triggers on OOP, object-oriented, encapsulation, polymorphism,
  inheritance, abstraction.
category: practices
conflicts: []
version: 1.0.0
license: MIT
---

You are enforcing proper Object-Oriented Design. Objects are not data bags — they have behavior and they protect their own state.

## Encapsulation — objects protect their state

```csharp
// WRONG — anemic model, just a data bag
public class Order
{
    public Guid Id { get; set; }
    public string Status { get; set; }
    public List<OrderLine> Lines { get; set; }
}

// CORRECT — encapsulated, behavior-rich object
public class Order
{
    public OrderId Id { get; private set; }
    public OrderStatus Status { get; private set; }
    private readonly List<OrderLine> _lines = new();
    public IReadOnlyList<OrderLine> Lines => _lines.AsReadOnly();

    public void Confirm()
    {
        if (Status != OrderStatus.Draft)
            throw new DomainException("Only draft orders can be confirmed.");
        Status = OrderStatus.Confirmed;
    }
}
```

## Inheritance — for IS-A relationships only

```csharp
// WRONG — inheritance just for code reuse
public class BaseRepository
{
    protected void LogQuery(string sql) { }
}
public class OrderRepository : BaseRepository {}  // not really "is a" BaseRepository

// CORRECT — composition for code reuse
public class OrderRepository(ILogger<OrderRepository> logger) {}

// CORRECT — inheritance for IS-A
public abstract class Animal { public abstract string Speak(); }
public class Dog : Animal   { public override string Speak() => "Woof"; }
public class Cat : Animal   { public override string Speak() => "Meow"; }
```

## Polymorphism — replace conditionals with objects

```csharp
// WRONG — type-checking is a missed polymorphism opportunity
public decimal Calculate(string discountType, decimal price)
{
    if (discountType == "seasonal") return price * 0.9m;
    if (discountType == "loyalty")  return price * 0.85m;
    return price;
}

// CORRECT — polymorphism replaces the switch
public interface IDiscountStrategy { decimal Apply(decimal price); }
public class SeasonalDiscount : IDiscountStrategy { public decimal Apply(decimal p) => p * 0.9m; }
public class LoyaltyDiscount  : IDiscountStrategy { public decimal Apply(decimal p) => p * 0.85m; }

public class PricingService(IDiscountStrategy discount)
{
    public decimal Calculate(decimal price) => discount.Apply(price);
}
```

## Abstraction — hide complexity, expose intent

```csharp
// WRONG — exposes implementation details
public class EmailService
{
    public void SendViaSmtp(string host, int port, string user,
        string pass, string to, string subject, string body) {}
}

// CORRECT — hides complexity behind meaningful abstraction
public interface IEmailService
{
    Task SendAsync(Email email, CancellationToken ct = default);
}

public record Email(string To, string Subject, string Body);
```

## Tell, don't ask

```csharp
// WRONG — asking for data then acting on it
if (order.Status == OrderStatus.Draft && order.Lines.Any())
{
    order.Status = OrderStatus.Confirmed;
    order.ConfirmedAt = DateTime.UtcNow;
}

// CORRECT — telling the object to do something
order.Confirm();
```

## Rules to enforce always

- No public setters on domain objects
- Collections exposed as IReadOnlyList — never expose internal list
- Prefer composition over inheritance for code reuse
- If you are checking the type of an object to decide what to do — use polymorphism
- Objects tell other objects what to do — they do not ask for data and act on it

## Red flags — stop and warn

- if (obj is TypeA) or if (obj is TypeB) — use polymorphism
- Public setters on domain entity properties
- Class with only getters and setters and no behavior — anemic model
- Inheritance chain deeper than 2 levels
- static methods for business logic — breaks testability
