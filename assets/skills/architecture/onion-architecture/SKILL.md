---
name: onion-architecture
description: >
  Enforce Onion Architecture in all code. Domain model at the center,
  dependencies point inward through concentric layers.
  Triggers on onion architecture, concentric layers, domain model center.
category: architecture
conflicts: []
version: 1.0.0
license: MIT
---

You are enforcing Onion Architecture. The domain model is at the center. Every layer depends only on layers closer to the center — never outward.

## Layer structure (inside out)

```
Domain Model         ← center — entities, value objects, domain events
Domain Services      ← domain logic spanning multiple entities
Application Services ← orchestration, use cases, DTOs
Infrastructure       ← outermost — DB, HTTP, messaging, file system
```

## Layer dependency rules

```
Infrastructure    → Application Services
Application Srvcs → Domain Services
Domain Services   → Domain Model
Domain Model      → nothing
```

## Domain model — innermost, pure

```csharp
public class Order
{
    public OrderId Id { get; private set; }
    public Money Total => Lines.Sum(l => l.Subtotal);
    private readonly List<OrderLine> _lines = new();
    public IReadOnlyList<OrderLine> Lines => _lines.AsReadOnly();

    public void AddLine(ProductId productId, Quantity qty, Money price)
    {
        Guard.Against.Null(productId);
        Guard.Against.NegativeOrZero(qty.Value);
        _lines.Add(new OrderLine(productId, qty, price));
    }
}
```

## Domain service — spans multiple entities

```csharp
public class OrderPricingService
{
    public Money CalculateTotal(Order order, IEnumerable<Discount> discounts)
    {
        var subtotal = order.Lines.Sum(l => l.Subtotal);
        var discount = discounts
            .Where(d => d.AppliesTo(order))
            .Sum(d => d.Calculate(subtotal));
        return subtotal - discount;
    }
}
```

## Application service — orchestration layer

```csharp
public class OrderApplicationService(
    IOrderRepository orderRepo,
    OrderPricingService pricingService,
    IEventPublisher eventPublisher)
{
    public async Task<OrderDto> CreateAsync(CreateOrderCommand cmd, CancellationToken ct)
    {
        var order = Order.Create(new CustomerId(cmd.CustomerId));
        foreach (var line in cmd.Lines)
            order.AddLine(new ProductId(line.ProductId), new Quantity(line.Qty), Money.Of(line.Price));

        var total = pricingService.CalculateTotal(order, await orderRepo.GetDiscountsAsync(ct));
        await orderRepo.SaveAsync(order, ct);
        await eventPublisher.PublishAsync(new OrderCreatedEvent(order.Id, total), ct);

        return OrderDto.From(order);
    }
}
```

## Red flags — stop immediately

- Framework annotations in domain model or domain services
- Application service importing directly from infrastructure
- Domain model calling repositories or external services
- Skipping layers — infrastructure calling domain directly
