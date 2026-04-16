---
name: clean-architecture
description: Enforce Clean Architecture (Robert C. Martin) in all code. Dependency rule is absolute — source code dependencies point inward only. Use when building backend services that must be testable, framework-independent, and maintainable long-term. Triggers on: "clean architecture", "use cases", "entities", "interface adapters", "dependency rule".
category: architecture
conflicts: [hexagonal]
version: 1.0.0
---

You are enforcing Clean Architecture. The Dependency Rule is absolute and non-negotiable: source code dependencies point inward only. Inner layers know nothing about outer layers.

## Layer structure

```
src/
  domain/          ← Entities, value objects, domain events — zero dependencies
  application/     ← Use cases, ports (interfaces), DTOs — depends only on domain
  infrastructure/  ← Repositories, DB, external services — depends on application
  presentation/    ← Controllers, serializers, HTTP — depends on application
```

## The Dependency Rule — enforced always

```
domain        →  depends on nothing
application   →  depends on domain only
infrastructure →  depends on application (implements its ports)
presentation  →  depends on application (calls its use cases)
```

## Domain entity — pure, no framework

```csharp
// CORRECT — pure domain entity, zero dependencies
public class Order
{
    private readonly List<OrderLine> _lines = new();
    private readonly List<IDomainEvent> _events = new();

    public OrderId Id { get; private set; }
    public CustomerId CustomerId { get; private set; }
    public OrderStatus Status { get; private set; }
    public IReadOnlyList<OrderLine> Lines => _lines.AsReadOnly();
    public IReadOnlyList<IDomainEvent> DomainEvents => _events.AsReadOnly();

    public void AddLine(ProductId productId, Quantity qty, Money price)
    {
        if (Status != OrderStatus.Draft)
            throw new DomainException("Cannot modify a confirmed order.");
        _lines.Add(new OrderLine(productId, qty, price));
    }

    public void Confirm()
    {
        if (!_lines.Any())
            throw new DomainException("Cannot confirm an empty order.");
        Status = OrderStatus.Confirmed;
        _events.Add(new OrderConfirmedEvent(Id));
    }
}

// WRONG — entity with infrastructure leak
[Table("orders")]               // ← infrastructure annotation in domain
public class Order
{
    [HttpGet]                   // ← presentation concern in domain
    public IActionResult Get() {}
}
```

## Use case — application layer

```csharp
public class ConfirmOrderUseCase(IOrderRepository orders, IEventBus bus)
{
    public async Task ExecuteAsync(ConfirmOrderCommand cmd, CancellationToken ct = default)
    {
        var order = await orders.GetByIdAsync(cmd.OrderId, ct)
            ?? throw new OrderNotFoundException(cmd.OrderId);

        order.Confirm();

        await orders.SaveAsync(order, ct);
        await bus.PublishAsync(order.DomainEvents, ct);
    }
}
```

## Port — defined in application layer

```csharp
// application/ports/IOrderRepository.cs
public interface IOrderRepository
{
    Task<Order?> GetByIdAsync(OrderId id, CancellationToken ct = default);
    Task SaveAsync(Order order, CancellationToken ct = default);
}
```

## Infrastructure implements the port

```csharp
// infrastructure/persistence/SqlOrderRepository.cs
public class SqlOrderRepository(AppDbContext context) : IOrderRepository
{
    public async Task<Order?> GetByIdAsync(OrderId id, CancellationToken ct = default)
    {
        // EF Core, Dapper, raw SQL — only here, never in domain or application
        var entity = await context.Orders
            .Include(o => o.Lines)
            .FirstOrDefaultAsync(o => o.Id == id.Value, ct);
        return entity?.ToDomain();
    }
}
```

## Controller — presentation layer

```csharp
// presentation/controllers/OrdersController.cs
[ApiController]
[Route("api/[controller]")]
public class OrdersController(ConfirmOrderUseCase confirmOrder) : ControllerBase
{
    [HttpPost("{id:guid}/confirm")]
    public async Task<IActionResult> Confirm(Guid id, CancellationToken ct)
    {
        await confirmOrder.ExecuteAsync(new ConfirmOrderCommand(id), ct);
        return NoContent();
    }
}
```

## Red flags — stop immediately

- Entity importing from `Microsoft.AspNetCore`, `System.Data`, or any ORM namespace
- Controller containing business logic or calculations
- Repository containing domain rules or decisions
- Use case importing from infrastructure layer
- DTOs exposing domain entities directly — always map at boundaries
- `new SqlOrderRepository()` inside a use case — inject via constructor
