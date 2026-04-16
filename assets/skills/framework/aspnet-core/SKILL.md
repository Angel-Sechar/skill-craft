---
name: aspnet-core
description: >
  Build HTTP APIs and web applications with ASP.NET Core. Apply controller
  patterns, middleware, filters, model binding, and proper error handling.
  Triggers on ASP.NET Core, Web API, controllers, middleware, REST API .NET.
category: framework
language: csharp
conflicts: [dotnet-framework-45]
version: 1.0.0
license: MIT
---

You are building an ASP.NET Core Web API. Controllers are thin. Business logic never lives here.

## Controller rules

Controllers do exactly three things:
1. Parse and validate the incoming request
2. Call the appropriate service or use case
3. Map the result to an HTTP response

```csharp
[ApiController]
[Route("api/[controller]")]
public class OrdersController(OrderService orderService) : ControllerBase
{
    [HttpPost("{id:guid}/confirm")]
    [ProducesResponseType(StatusCodes.Status204NoContent)]
    [ProducesResponseType(typeof(ProblemDetails), StatusCodes.Status404NotFound)]
    public async Task<IActionResult> Confirm(Guid id, CancellationToken ct)
    {
        await orderService.ConfirmAsync(id, ct);
        return NoContent();
    }

    [HttpGet("{id:guid}")]
    [ProducesResponseType(typeof(OrderResponse), StatusCodes.Status200OK)]
    public async Task<ActionResult<OrderResponse>> GetById(Guid id, CancellationToken ct)
    {
        var order = await orderService.GetByIdAsync(id, ct);
        return Ok(order);
    }
}
```

## Global exception handling — Problem Details (RFC 9457)

```csharp
builder.Services.AddProblemDetails();
builder.Services.AddExceptionHandler<GlobalExceptionHandler>();

public class GlobalExceptionHandler(IProblemDetailsService problemDetails)
    : IExceptionHandler
{
    public async ValueTask<bool> TryHandleAsync(
        HttpContext context, Exception exception, CancellationToken ct)
    {
        var (status, title) = exception switch
        {
            NotFoundException   => (404, "Resource not found"),
            ConflictException   => (409, "Conflict"),
            ValidationException => (400, "Validation failed"),
            _                   => (500, "Internal server error")
        };

        context.Response.StatusCode = status;
        return await problemDetails.TryWriteAsync(new()
        {
            HttpContext = context,
            ProblemDetails = { Title = title, Status = status }
        });
    }
}
```

## Model validation

```csharp
public record CreateOrderRequest(
    [Required] Guid CustomerId,
    [Required, MinLength(1)] List<OrderLineRequest> Lines
);
```

## Middleware order — this matters

```csharp
app.UseExceptionHandler();
app.UseHttpsRedirection();
app.UseAuthentication();
app.UseAuthorization();
app.MapControllers();
```

## Red flags — stop and warn

- Business logic inside a controller action
- try/catch in controllers instead of global exception handler
- Missing CancellationToken parameter on async actions
- Returning hardcoded status integers instead of named methods
- Missing ApiController attribute — disables automatic model validation
