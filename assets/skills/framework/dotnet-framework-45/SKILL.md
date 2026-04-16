---
name: dotnet-framework-45
description: >
  Write all code targeting .NET Framework 4.5 with C#. Use only APIs,
  patterns, and NuGet packages compatible with this version. Apply when
  working on legacy enterprise applications that cannot be migrated.
  Triggers on: .NET Framework 4.5, legacy .NET, full framework, net45.
category: framework
language: csharp
conflicts: [dotnet-core-8, aspnet-core]
version: 1.0.0
license: MIT
---
You are working on a .NET Framework 4.5 codebase written in C#. This is a legacy environment. Your job is to make it better without breaking it.

## Non-negotiable constraints

- Target framework is `net45` — never suggest APIs introduced after .NET Framework 4.5
- Language version is C# 5 — async/await is available but records, pattern matching, nullable reference types, and primary constructors are NOT
- NuGet packages must explicitly support `net45` — always verify before suggesting one
- Never suggest migrating to .NET Core or .NET 8 unless explicitly asked

## Project structure

```
src/
  MyApp/
    Controllers/     ← ASP.NET MVC or Web API controllers
    Models/          ← View models and domain models
    Services/        ← Business logic
    Repositories/    ← Data access
    App_Start/       ← Route and filter configuration
  MyApp.Tests/
    Services/
    Repositories/
```

## C# 5 patterns to use

```csharp
// Async/await — available in C# 5
public async Task<Order> GetOrderAsync(int id)
{
    return await _repository.GetByIdAsync(id);
}

// Generic collections
public IEnumerable<Order> GetActiveOrders()
{
    return _orders.Where(o => o.IsActive).ToList();
}

// Extension methods
public static class OrderExtensions
{
    public static decimal CalculateTotal(this Order order)
    {
        return order.Lines.Sum(l => l.Quantity * l.UnitPrice);
    }
}
```

## What NOT to use

```csharp
// WRONG — C# 6+ features not available in net45
public string Name { get; } = "default";        // auto-property initializer
var msg = $"Hello {name}";                       // string interpolation
order?.Lines?.FirstOrDefault();                  // null conditional operator

// Use these instead
public string Name { get; set; }
var msg = string.Format("Hello {0}", name);
order != null ? order.Lines.FirstOrDefault() : null;
```

## Dependency Injection in net45

Use Unity, Autofac, or Ninject — not `Microsoft.Extensions.DependencyInjection`:

```csharp
// Unity container setup in Global.asax or UnityConfig.cs
var container = new UnityContainer();
container.RegisterType<IOrderRepository, SqlOrderRepository>();
container.RegisterType<IOrderService, OrderService>();
GlobalConfiguration.Configuration.DependencyResolver =
    new UnityDependencyResolver(container);
```

## Data access

Entity Framework 6 is the ORM for net45:

```csharp
public class OrderRepository : IOrderRepository
{
    private readonly AppDbContext _context;

    public OrderRepository(AppDbContext context)
    {
        _context = context;
    }

    public async Task<Order> GetByIdAsync(int id)
    {
        return await _context.Orders
            .Include(o => o.Lines)
            .FirstOrDefaultAsync(o => o.Id == id);
    }
}
```

## ConfigureAwait rule

Always use `ConfigureAwait(false)` in library and repository code to avoid deadlocks in classic ASP.NET:

```csharp
var result = await _context.Orders.ToListAsync().ConfigureAwait(false);
```

## Red flags — stop and warn the user

- Any import from `Microsoft.Extensions.*` — not compatible with net45
- Use of `System.Text.Json` — not available, use `Newtonsoft.Json`
- `IHttpClientFactory` — not available, use `new HttpClient()` carefully
- `ValueTask<T>` — not available in net45
- Any C# 6+ syntax in a net45 project
- Suggesting Docker or containerization for this stack
