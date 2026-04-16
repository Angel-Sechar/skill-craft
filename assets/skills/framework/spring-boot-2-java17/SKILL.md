---
name: spring-boot-2-java17
description: >
  Write Java backend code using Spring Boot 2.x with Java 17. Use javax.*
  namespace, constructor injection, and Java 17 features like records and
  sealed classes. Triggers on Spring Boot 2, Spring Boot 2.x, Java 17
  Spring, javax.
category: framework
language: java
conflicts: [spring-boot-3-java21]
version: 1.0.0
license: MIT
---

You are working on a Spring Boot 2.x backend with Java 17. Use javax.* namespace — never jakarta.*. Java 17 features are available but Spring Boot 3 features are not.

## Critical — javax.* not jakarta.*

```java
// CORRECT — Spring Boot 2
import javax.persistence.Entity;
import javax.persistence.Table;
import javax.validation.constraints.NotNull;
import javax.servlet.http.HttpServletRequest;

// WRONG — this is Spring Boot 3
import jakarta.persistence.Entity;
import jakarta.validation.constraints.NotNull;
```

## Project setup (pom.xml)

```xml
<parent>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-parent</artifactId>
    <version>2.7.18</version>
</parent>
<properties>
    <java.version>17</java.version>
</properties>
```

## Java 17 features to use

```java
// Records for DTOs
public record CreateOrderRequest(
    @NotNull UUID customerId,
    @NotNull @Size(min = 1) List<OrderLineRequest> lines
) {}

// Sealed classes for domain results
public sealed interface OrderResult
    permits OrderResult.Success, OrderResult.NotFound {
    record Success(OrderDto order) implements OrderResult {}
    record NotFound(UUID orderId) implements OrderResult {}
}

// Pattern matching
if (result instanceof OrderResult.NotFound notFound) {
    return ResponseEntity.notFound().build();
}

// Switch expressions
String describe(OrderStatus status) {
    return switch (status) {
        case DRAFT     -> "Pending confirmation";
        case CONFIRMED -> "Confirmed";
        case CANCELLED -> "Cancelled";
    };
}
```

## REST controller

```java
@RestController
@RequestMapping("/api/orders")
@RequiredArgsConstructor
public class OrderController {

    private final OrderService orderService;

    @PostMapping("/{id}/confirm")
    public ResponseEntity<Void> confirm(@PathVariable UUID id) {
        orderService.confirm(id);
        return ResponseEntity.noContent().build();
    }
}
```

## Dependency injection — constructor only

```java
// CORRECT — constructor injection via @RequiredArgsConstructor
@Service
@RequiredArgsConstructor
public class OrderService {
    private final OrderRepository orderRepository;
    private final EventPublisher eventPublisher;
}

// WRONG — field injection, never do this
@Service
public class OrderService {
    @Autowired
    private OrderRepository orderRepository;
}
```

## Red flags — stop and warn

- Any jakarta.* import — wrong namespace for Boot 2
- @Autowired on fields — use constructor injection
- Business logic in @RestController methods
- @Transactional on domain services — only on application services
- Missing final on injected fields
