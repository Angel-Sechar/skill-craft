---
name: ddd
description: >
  Apply Domain-Driven Design tactical patterns. Aggregates, entities, value
  objects, domain services, bounded contexts, repositories, and ubiquitous
  language. Triggers on DDD, domain-driven design, aggregate, value object,
  bounded context, ubiquitous language.
category: driven
conflicts: []
version: 1.0.0
license: MIT
---

You are applying Domain-Driven Design. The domain model is the heart of the application. Code must speak the ubiquitous language. Every name must match the language the domain expert uses.

## Ubiquitous language — non-negotiable

If the domain expert says "confirm an order" — the method is Confirm(), not UpdateStatus() or SetStatusToConfirmed().
If the domain expert says "line item" — the class is OrderLine, not OrderItem or LineEntry.

## Value Objects — identity-less, immutable, self-validating

```csharp
public record Money
{
    public decimal Amount { get; }
    public string Currency { get; }

    public Money(decimal amount, string currency)
    {
        if (amount < 0)
            throw new DomainException("Money amount cannot be negative.");
        if (string.IsNullOrWhiteSpace(currency) || currency.Length != 3)
            throw new DomainException("Currency must be a 3-letter ISO code.");

        Amount = amount;
        Currency = currency.ToUpperInvariant();
    }

    public Money Add(Money other)
    {
        if (Currency != other.Currency)
            throw new DomainException($"Cannot add {Currency} and {other.Currency}.");
        return new Money(Amount + other.Amount, Currency);
    }

    public static Money Of(decimal amount, string currency) => new(amount, currency);
}
```

## Aggregate root — consistency boundary

```csharp
public class Order : Entity<OrderId>
{
    private readonly List<OrderLine> _lines = new();
    private readonly List<IDomainEvent> _events = new();

    public CustomerId CustomerId { get; private set; }
    public OrderStatus Status { get; private set; }
    public Money Total => _lines.Aggregate(
        Money.Of(0, "USD"), (sum, l) => sum.Add(l.Subtotal));

    public static Order Place(CustomerId customerId)
    {
        var order = new Order(OrderId.New(), customerId);
        order._events.Add(new OrderPlacedEvent(order.Id, customerId));
        return order;
    }

    public void AddLine(ProductId productId, Quantity quantity, Money price)
    {
        if (Status != OrderStatus.Draft)
            throw new DomainException("Cannot add lines to a non-draft order.");
        if (_lines.Any(l => l.ProductId == productId))
            throw new DomainException("Product already in order. Update quantity instead.");
        _lines.Add(new OrderLine(productId, quantity, price));
    }

    public void Confirm()
    {
        if (Status != OrderStatus.Draft)
            throw new DomainException("Only draft orders can be confirmed.");
        if (!_lines.Any())
            throw new DomainException("Cannot confirm an empty order.");
        Status = OrderStatus.Confirmed;
        _events.Add(new OrderConfirmedEvent(Id, CustomerId, Total));
    }
}
```

## Repository — one per aggregate root

```csharp
// One repository per AGGREGATE ROOT — not per entity
public interface IOrderRepository
{
    Task<Order?> GetByIdAsync(OrderId id, CancellationToken ct = default);
    Task SaveAsync(Order order, CancellationToken ct = default);
}

// WRONG — repository for a child entity
public interface IOrderLineRepository  // never do this
```

## Rules to enforce always

- Aggregate roots are the only public entry point
- All invariants enforced inside the aggregate — not in services or controllers
- Value objects are immutable
- Use the ubiquitous language in all names
- Domain events use past tense — OrderConfirmedEvent, never ConfirmOrderEvent
- One repository per aggregate root — never per child entity

## Red flags — stop and warn

- Setting properties directly on aggregate from outside
- Repository for a child entity
- Domain service with infrastructure dependencies
- Method names that do not match ubiquitous language
- Anemic domain model — entities with only getters and setters
