---
name: spring-boot-3-java21
description: >
  Write Java backend code using Spring Boot 3.x with Java 21. Use jakarta.*
  namespace, virtual threads, records, and modern Spring patterns.
  Triggers on: Spring Boot 3, Spring Boot 3.x, Java 21 Spring, jakarta,
  virtual threads.
category: framework
language: java
conflicts: [spring-boot-2-java17]
version: 1.0.0
license: MIT
---
You are working on a Spring Boot 3.x backend with Java 21. Use `jakarta.*` namespace exclusively. Virtual threads are available — use them. Java 21 features are fully supported.

## Critical — jakarta.* not javax.*

```java
// CORRECT — Spring Boot 3
import jakarta.persistence.Entity;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotNull;
import jakarta.servlet.http.HttpServletRequest;

// WRONG — this is Spring Boot 2
import javax.persistence.Entity;
import javax.validation.constraints.NotNull;
```

## Project setup (pom.xml)

```xml
<parent>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-parent</artifactId>
    <version>3.2.0</version>
</parent>
<properties>
    <java.version>21</java.version>
</properties>
```

## Enable virtual threads (Java 21 + Boot 3.2)

```yaml
# application.yml
spring:
  threads:
    virtual:
      enabled: true
```

## Java 21 features to use

```java
// Records for DTOs
public record CreateOrderRequest(
    @NotNull UUID customerId,
    @NotNull @Size(min = 1) List<OrderLineRequest> lines
) {}

// Sequenced collections
List<OrderLine> lines = new ArrayList<>();
var first = lines.getFirst();
var last  = lines.getLast();

// Pattern matching in switch (finalized in Java 21)
String describe(Object obj) {
    return switch (obj) {
        case Order o when o.status() == CONFIRMED -> "Confirmed: " + o.id();
        case Order o                              -> "Draft: " + o.id();
        case null                                 -> "null";
        default                                   -> obj.toString();
    };
}

// Virtual thread-aware — avoid blocking calls
// WRONG — blocks a virtual thread unnecessarily
Thread.sleep(1000);

// CORRECT — use async or reactive where blocking is needed
CompletableFuture.runAsync(() -> expensiveOperation());
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

## Problem Details — built into Boot 3

```java
// application.yml
spring:
  mvc:
    problemdetails:
      enabled: true

// Exception handler
@ControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(OrderNotFoundException.class)
    public ProblemDetail handleNotFound(OrderNotFoundException ex) {
        var problem = ProblemDetail.forStatusAndDetail(
            HttpStatus.NOT_FOUND, ex.getMessage());
        problem.setTitle("Order not found");
        return problem;
    }
}
```

## Red flags — stop and warn

- Any `javax.*` import — wrong namespace for Boot 3
- `@Autowired` on fields — use constructor injection
- Calling `Thread.sleep()` or blocking I/O on virtual threads unnecessarily
- `@Transactional` on domain services — only on application/repository layer
- Not enabling virtual threads when on Java 21 + Boot 3.2
