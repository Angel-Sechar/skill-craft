---
name: dotnet-core-8
description: >
  Write all backend code targeting .NET 8 LTS with C# 12. Use minimal APIs,
  primary constructors, nullable reference types, and modern patterns.
  Triggers on: .NET Core 8, .NET 8, C# 12, net8, modern .NET.
category: framework
language: csharp
conflicts: [dotnet-framework-45]
version: 1.0.0
license: MIT
---
You are working on a .NET 8 backend written in C# 12. Use all modern language features deliberately. No legacy patterns.

## Project setup

```xml
<Project Sdk="Microsoft.NET.Sdk.Web">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
    <Nullable>enable</Nullable>
    <ImplicitUsings>enable</ImplicitUsings>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
  </PropertyGroup>
</Project>
```

## C# 12 features — use these

```csharp
// Primary constructors on classes
public class OrderService(IOrderRepository repo, IEventBus bus)
{
    public async Task ConfirmAsync(Guid id, CancellationToken ct = default)
    {
        var order = await repo.GetByIdAsync(id, ct)
            ?? throw new NotFoundException($"Order {id} not found");
        order.Confirm();
        await repo.SaveAsync(order, ct);
    }
}

// Collection expressions
List<string> tags = ["backend", "api", "orders"];
int[] ids = [1, 2, 3, 4];

// Records for DTOs
public record CreateOrderRequest(Guid CustomerId, List<OrderLineDto> Lines);
public record OrderResponse(Guid Id, string Status, decimal Total);

// Pattern matching
string Describe(Order order) => order.Status switch
{
    OrderStatus.Draft     => "Pending confirmation",
    OrderStatus.Confirmed => $"Confirmed — {order.Total:C}",
    OrderStatus.Cancelled => "Cancelled",
    _                     => throw new UnreachableException()
};
```

## Minimal API (preferred for new services)

```csharp
var builder = WebApplication.CreateBuilder(args);

builder.Services.AddScoped<IOrderRepository, SqlOrderRepository>();
builder.Services.AddScoped<OrderService>();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

var orders = app.MapGroup("/api/orders").WithOpenApi();

orders.MapPost("/{id:guid}/confirm", async (
    Guid id,
    OrderService service,
    CancellationToken ct) =>
{
    await service.ConfirmAsync(id, ct);
    return Results.NoContent();
})
.WithName("ConfirmOrder")
.Produces(204)
.Produces<ProblemDetails>(404);

app.Run();
```

## Async rules

- All I/O methods return `Task<T>` and end with `Async`
- Always accept and forward `CancellationToken` on public async methods
- Never use `.Result` or `.Wait()` — always `await`
- No `ConfigureAwait(false)` needed in ASP.NET Core — it's safe by default

## Nullable reference types

```csharp
// CORRECT — be explicit about nullability
public async Task<Order?> FindAsync(Guid id, CancellationToken ct = default)
{
    return await _context.Orders.FirstOrDefaultAsync(o => o.Id == id, ct);
}

// CORRECT — use null-forgiving only when you've already checked
var order = await FindAsync(id, ct) ?? throw new NotFoundException(id);
```

## Configuration

```csharp
// CORRECT — Options pattern
builder.Services.Configure<DatabaseOptions>(
    builder.Configuration.GetSection("Database"));

// WRONG — never inject IConfiguration into business services
public class OrderService(IConfiguration config) // ← wrong
```

## Red flags — stop and warn

- `.Result` or `.Wait()` on async code — deadlock risk
- `new HttpClient()` in a service — use `IHttpClientFactory`
- Reading `IConfiguration` directly in business logic
- Missing `CancellationToken` on public async methods
- Nullable warnings suppressed with `!` without explanation
