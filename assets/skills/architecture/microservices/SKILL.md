---
name: microservices
description: >
  Design and build microservices. Each service owns its bounded context,
  its data, and its deployment. Services communicate via HTTP or events.
  Triggers on: microservices, bounded context, service decomposition,
  API gateway, service mesh.
category: architecture
conflicts: []
version: 1.0.0
license: MIT
---
You are building microservices. Each service is an autonomous unit — it owns its data, its domain, and its deployment. Services never share databases.

## Non-negotiable rules

- One database per service — no cross-service DB queries ever
- Services communicate via well-defined APIs or events — never direct DB access
- Each service must be independently deployable
- Each service exposes `/health` and `/health/ready` endpoints

## Service structure

```
order-service/
  src/
    api/           ← controllers, DTOs, serialization
    application/   ← use cases, commands, queries
    domain/        ← entities, value objects, domain events
    infrastructure/← repository implementations, event publisher
  Dockerfile
  docker-compose.yml
```

## Synchronous communication (HTTP)

```java
// CORRECT — call via HTTP client, not direct DB
@Service
public class OrderService {

    private final CustomerClient customerClient;  // HTTP client to customer-service

    public void placeOrder(PlaceOrderCommand cmd) {
        var customer = customerClient.getById(cmd.customerId())
            .orElseThrow(() -> new CustomerNotFoundException(cmd.customerId()));

        if (!customer.isActive())
            throw new CustomerNotActiveException(cmd.customerId());

        // proceed with order creation
    }
}

// WRONG — querying another service's database directly
@Query("SELECT * FROM customer_service.customers WHERE id = :id")  // never
```

## Asynchronous communication (events)

```java
// Producer — publish after successful operation
public void confirmOrder(UUID orderId) {
    var order = orderRepository.findById(orderId).orElseThrow();
    order.confirm();
    orderRepository.save(order);

    eventPublisher.publish(new OrderConfirmedEvent(
        orderId,
        order.customerId(),
        order.total(),
        Instant.now()
    ));
}

// Consumer — always idempotent
@KafkaListener(topics = "order-confirmed")
public void handleOrderConfirmed(OrderConfirmedEvent event) {
    if (processedEventStore.hasBeenProcessed(event.eventId())) return;

    notificationService.sendConfirmation(event.customerId(), event.orderId());
    processedEventStore.markProcessed(event.eventId());
}
```

## Health checks (required on every service)

```java
@RestController
public class HealthController {

    private final DataSource dataSource;

    @GetMapping("/health")
    public ResponseEntity<Map<String, String>> health() {
        return ResponseEntity.ok(Map.of("status", "UP"));
    }

    @GetMapping("/health/ready")
    public ResponseEntity<Map<String, String>> ready() {
        try {
            dataSource.getConnection().close();
            return ResponseEntity.ok(Map.of("status", "READY"));
        } catch (Exception e) {
            return ResponseEntity.status(503).body(Map.of("status", "NOT_READY"));
        }
    }
}
```

## Circuit breaker (Resilience4j)

```java
@CircuitBreaker(name = "customerService", fallbackMethod = "customerFallback")
public CustomerDto getCustomer(UUID customerId) {
    return customerClient.getById(customerId);
}

public CustomerDto customerFallback(UUID customerId, Exception ex) {
    log.warn("Customer service unavailable for {}", customerId);
    return CustomerDto.unknown(customerId);
}
```

## Red flags — stop immediately

- Cross-service database queries
- Shared database between two services
- Service without health endpoint
- Synchronous chain of 3+ service calls — consider async
- Missing idempotency in event consumers
- Service calling another service's internal API (not its public contract)
