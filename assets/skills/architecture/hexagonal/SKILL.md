---
name: hexagonal
description: Enforce Hexagonal Architecture (Ports and Adapters) in all code. The application core has zero framework dependencies. Driving adapters call in, driven adapters are called out. Triggers on: "hexagonal architecture", "ports and adapters", "driving adapter", "driven adapter", "domain isolation".
category: architecture
conflicts: [clean-architecture]
version: 1.0.0
---

You are enforcing Hexagonal Architecture. The core (domain + application) is completely isolated. No framework imports inside the hexagon — ever.

## Structure

```
src/
  core/
    domain/          ← Entities, value objects, domain events
    ports/
      driving/       ← Interfaces the outside world calls (input ports)
      driven/        ← Interfaces the app needs from outside (output ports)
    application/     ← Use cases implementing driving ports
  adapters/
    driving/
      http/          ← REST controllers
      cli/           ← CLI commands
      messaging/     ← Message consumers
    driven/
      persistence/   ← Repository implementations
      messaging/     ← Event publishers
      external/      ← Third-party API clients
```

## Driving port (input) — defined in core

```java
// core/ports/driving/ConfirmOrderPort.java
public interface ConfirmOrderPort {
    void confirm(ConfirmOrderCommand command);
}
```

## Driven port (output) — defined in core

```java
// core/ports/driven/OrderRepositoryPort.java
public interface OrderRepositoryPort {
    Optional<Order> findById(OrderId id);
    void save(Order order);
}
```

## Application service — implements driving port, uses driven ports

```java
// core/application/ConfirmOrderService.java
public class ConfirmOrderService implements ConfirmOrderPort {

    private final OrderRepositoryPort orderRepository;
    private final EventPublisherPort eventPublisher;

    public ConfirmOrderService(
        OrderRepositoryPort orderRepository,
        EventPublisherPort eventPublisher
    ) {
        this.orderRepository = orderRepository;
        this.eventPublisher = eventPublisher;
    }

    @Override
    public void confirm(ConfirmOrderCommand command) {
        var order = orderRepository.findById(new OrderId(command.orderId()))
            .orElseThrow(() -> new OrderNotFoundException(command.orderId()));

        order.confirm();

        orderRepository.save(order);
        eventPublisher.publish(new OrderConfirmedEvent(order.id()));
    }
}
```

## Driving adapter (HTTP) — calls into core via port

```java
// adapters/driving/http/OrderController.java
@RestController
@RequestMapping("/api/orders")
public class OrderController {

    private final ConfirmOrderPort confirmOrder;  // depends on PORT, not implementation

    public OrderController(ConfirmOrderPort confirmOrder) {
        this.confirmOrder = confirmOrder;
    }

    @PostMapping("/{id}/confirm")
    public ResponseEntity<Void> confirm(@PathVariable UUID id) {
        confirmOrder.confirm(new ConfirmOrderCommand(id));
        return ResponseEntity.noContent().build();
    }
}
```

## Driven adapter (persistence) — implements driven port

```java
// adapters/driven/persistence/JpaOrderRepository.java
@Repository
public class JpaOrderRepository implements OrderRepositoryPort {

    private final OrderJpaRepository jpa;

    public JpaOrderRepository(OrderJpaRepository jpa) {
        this.jpa = jpa;
    }

    @Override
    public Optional<Order> findById(OrderId id) {
        return jpa.findById(id.value()).map(OrderEntity::toDomain);
    }

    @Override
    public void save(Order order) {
        jpa.save(OrderEntity.fromDomain(order));
    }
}
```

## The test that matters

You can test the entire core without Spring, without a database, without HTTP:

```java
@Test
void confirm_existingOrder_publishesEvent() {
    var repo = new InMemoryOrderRepository();
    var publisher = new InMemoryEventPublisher();
    var service = new ConfirmOrderService(repo, publisher);

    var order = Order.createDraft(new CustomerId(UUID.randomUUID()));
    order.addLine(new ProductId(UUID.randomUUID()), new Quantity(1), Money.of(100));
    repo.save(order);

    service.confirm(new ConfirmOrderCommand(order.id().value()));

    assertThat(publisher.publishedEvents()).hasSize(1);
    assertThat(publisher.publishedEvents().get(0)).isInstanceOf(OrderConfirmedEvent.class);
}
```

## Red flags — stop immediately

- Any Spring annotation (`@Service`, `@Component`, `@Repository`) inside `core/`
- HTTP status codes or JSON annotations in application services
- SQL or ORM calls in domain classes
- Application service calling another adapter directly — only via ports
- Domain entities serialized directly in HTTP responses
