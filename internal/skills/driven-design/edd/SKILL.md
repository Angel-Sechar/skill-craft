---
name: edd
description: Apply Event-Driven Design to the system. Domain events capture what happened. Services react to events asynchronously and independently. Triggers on: "event-driven", "domain events", "event bus", "event sourcing", "pub/sub", "async messaging".
category: driven-design
conflicts: []
version: 1.0.0
---

You are applying Event-Driven Design. Things that happen in the domain are captured as events. Other parts of the system react to those events — independently, asynchronously, without tight coupling.

## What is a domain event

A domain event is a fact. Something that happened. It is immutable, past-tense, and named after what occurred in the domain.

```csharp
// CORRECT — past tense, immutable, carries what matters
public record OrderConfirmedEvent(
    OrderId OrderId,
    CustomerId CustomerId,
    Money Total,
    DateTime ConfirmedAt
) : IDomainEvent
{
    public Guid EventId { get; } = Guid.NewGuid();
    public DateTime OccurredAt { get; } = DateTime.UtcNow;
}

// WRONG — command disguised as event
public record ConfirmOrderEvent(...) // ← commands are not events
public record OrderEvent(...)        // ← too generic, meaningless
```

## Raising events inside the aggregate

```csharp
public class Order
{
    private readonly List<IDomainEvent> _events = new();
    public IReadOnlyList<IDomainEvent> DomainEvents => _events.AsReadOnly();

    public void Confirm()
    {
        if (!_lines.Any())
            throw new DomainException("Cannot confirm empty order.");

        Status = OrderStatus.Confirmed;

        // Raise event AFTER state change — fact is recorded after it happens
        _events.Add(new OrderConfirmedEvent(Id, CustomerId, Total, DateTime.UtcNow));
    }

    public void ClearEvents() => _events.Clear();
}
```

## Publishing events — after persistence

```csharp
public class ConfirmOrderUseCase(IOrderRepository orders, IEventBus bus)
{
    public async Task ExecuteAsync(ConfirmOrderCommand cmd, CancellationToken ct)
    {
        var order = await orders.GetByIdAsync(cmd.OrderId, ct)
            ?? throw new OrderNotFoundException(cmd.OrderId);

        order.Confirm();

        // Save first — then publish. Order matters.
        await orders.SaveAsync(order, ct);
        await bus.PublishAsync(order.DomainEvents, ct);
        order.ClearEvents();
    }
}
```

## Consumer — always idempotent

```csharp
public class OrderConfirmedConsumer(
    IProcessedEventStore processedEvents,
    INotificationService notifications)
{
    public async Task HandleAsync(OrderConfirmedEvent evt, CancellationToken ct)
    {
        // Idempotency check — processing same event twice = same outcome
        if (await processedEvents.HasBeenProcessedAsync(evt.EventId, ct))
            return;

        await notifications.SendConfirmationAsync(evt.CustomerId, evt.OrderId, ct);
        await processedEvents.MarkProcessedAsync(evt.EventId, ct);
    }
}
```

## Rules to enforce always

- Events are facts — immutable, past tense, never modified after creation
- Always include `EventId` for idempotency and `OccurredAt` for ordering
- Publish events AFTER successful persistence — never before
- Consumers must be idempotent — same event twice = same result
- Never call another aggregate's method inside an event handler — publish a new event instead
- Dead letter queue for all consumers — failed events must not be silently dropped

## Red flags — stop and warn

- Event published before data is saved — inconsistency risk
- Event consumer that is not idempotent
- Event with mutable properties
- Using events as commands — "ProcessOrderEvent" is a command, not an event
- Synchronous chain of event handlers — defeats the purpose of async decoupling
- Missing correlation ID for tracing events across services
