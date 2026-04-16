---
name: tdd
description: >
  Apply Test-Driven Development to all feature work. Write a failing test
  first, then write the minimum code to pass it, then refactor.
  No production code without a failing test.
  Triggers on TDD, test-driven development, red green refactor,
  write test first.
category: driven
conflicts: []
version: 1.0.0
license: MIT
---

You are applying Test-Driven Development. The cycle is Red, Green, Refactor. No exceptions.

## The three laws of TDD

1. You may not write production code until you have a failing test
2. You may not write more test code than is sufficient to fail
3. You may not write more production code than is sufficient to pass the failing test

## The cycle in practice

```csharp
// Step 1 — RED: write the smallest failing test
[Fact]
public void Confirm_EmptyOrder_ThrowsDomainException()
{
    var order = Order.CreateDraft(CustomerId.New());

    var act = () => order.Confirm();

    act.Should().Throw<DomainException>()
        .WithMessage("Cannot confirm an empty order.");
}

// Step 2 — GREEN: write minimum code to pass
public void Confirm()
{
    if (!_lines.Any())
        throw new DomainException("Cannot confirm an empty order.");
    Status = OrderStatus.Confirmed;
}

// Step 3 — REFACTOR: clean up without changing behavior
public void Confirm()
{
    Guard.Against.Empty(_lines, "Cannot confirm an empty order.");
    Status = OrderStatus.Confirmed;
    AddDomainEvent(new OrderConfirmedEvent(Id));
}
```

## Test anatomy — Arrange, Act, Assert

```csharp
[Fact]
public void AddLine_ValidProduct_IncreasesLineCount()
{
    // Arrange
    var order = Order.CreateDraft(CustomerId.New());
    var productId = ProductId.New();
    var qty = new Quantity(2);
    var price = Money.Of(50m, "USD");

    // Act
    order.AddLine(productId, qty, price);

    // Assert
    order.Lines.Should().HaveCount(1);
    order.Total.Amount.Should().Be(100m);
}
```

## Test naming — behavior not implementation

```csharp
// CORRECT — describes behavior
void Confirm_WithLines_ChangesStatusToConfirmed()
void Confirm_EmptyOrder_ThrowsDomainException()
void AddLine_ConfirmedOrder_ThrowsDomainException()

// WRONG — describes implementation
void TestConfirmMethod()
void ConfirmWorks()
void Test1()
```

## Start with the degenerate case

```csharp
// First test — the nothing case
[Fact]
public void NewOrder_HasNoLines()
{
    var order = Order.CreateDraft(CustomerId.New());
    order.Lines.Should().BeEmpty();
}
```

## Rules to enforce always

- Never write production code without a failing test requiring it
- One logical assertion per test
- Test names follow Method_Scenario_ExpectedResult format
- Tests are first-class code — same quality standards as production
- A test that always passes is worse than no test

## Red flags — stop and warn

- Writing production code then tests after — that is not TDD
- Tests that test implementation details instead of behavior
- Thread.Sleep or Task.Delay in tests — design flaw
- Mocking everything — if you need 5 mocks, the class has too many dependencies
