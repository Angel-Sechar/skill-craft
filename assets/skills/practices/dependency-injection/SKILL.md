---
name: dependency-injection
description: >
  Apply Dependency Injection correctly. Constructor injection for required
  dependencies, no service locator, no new keyword for infrastructure objects
  inside business classes. Triggers on: dependency injection, DI, IoC,
  inversion of control, constructor injection.
category: practices
conflicts: []
version: 1.0.0
license: MIT
---
You are enforcing Dependency Injection. Dependencies are declared, not created. Business classes never instantiate their own infrastructure.

## Constructor injection — the only acceptable pattern

```csharp
// CORRECT — dependencies declared, injected externally
public class OrderService(
    IOrderRepository orderRepository,
    IEventBus eventBus,
    ILogger<OrderService> logger)
{
    public async Task ConfirmAsync(Guid id, CancellationToken ct)
    {
        var order = await orderRepository.GetByIdAsync(id, ct)
            ?? throw new NotFoundException(id);

        order.Confirm();
        await orderRepository.SaveAsync(order, ct);
        await eventBus.PublishAsync(order.DomainEvents, ct);
    }
}

// WRONG — creating dependencies internally
public class OrderService
{
    private readonly IOrderRepository _repo = new SqlOrderRepository();  // ← wrong
    private readonly IEventBus _bus = new RabbitMqBus("amqp://...");    // ← wrong
}
```

## Registration — .NET

```csharp
// Lifetimes matter — use the right one
builder.Services.AddScoped<IOrderRepository, SqlOrderRepository>();   // per request
builder.Services.AddSingleton<IEventBus, RabbitMqBus>();              // one instance
builder.Services.AddTransient<IEmailSender, SmtpEmailSender>();       // new each time

// WRONG — registering as wrong lifetime
builder.Services.AddSingleton<IOrderRepository, SqlOrderRepository>(); // DbContext is scoped!
```

## Registration — Spring Boot

```java
// CORRECT — Spring manages lifecycle
@Configuration
public class ApplicationConfig {

    @Bean
    public OrderService orderService(
        OrderRepository orderRepository,
        EventPublisher eventPublisher) {
        return new OrderService(orderRepository, eventPublisher);
    }
}

// WRONG — field injection
@Service
public class OrderService {
    @Autowired
    private OrderRepository orderRepository;  // never use field injection
}
```

## No service locator

```csharp
// WRONG — service locator anti-pattern
public class OrderService(IServiceProvider provider)
{
    public void Confirm(Guid id)
    {
        var repo = provider.GetService<IOrderRepository>();  // ← anti-pattern
    }
}

// CORRECT — declare what you need upfront
public class OrderService(IOrderRepository repo) {}
```

## Lifetime rules (.NET)

- `Scoped` — one instance per HTTP request. Use for DbContext, repositories, use cases.
- `Singleton` — one instance for app lifetime. Use for caches, HTTP clients, configuration.
- `Transient` — new instance every time. Use for lightweight stateless services.
- Never inject `Scoped` into `Singleton` — causes captive dependency bug.

## Red flags — stop and warn

- `new ConcreteClass()` inside a business service or use case
- `IServiceProvider` injected into a business class — service locator
- Static methods that create dependencies: `DatabaseFactory.Create()`
- `@Autowired` on fields in Spring — always constructor injection
- Scoped service injected into singleton — lifetime mismatch
- More than 4-5 constructor parameters — class has too many responsibilities
