---
name: angular-17
description: >
  Build Angular 17+ applications using standalone components, signals, new
  control flow syntax, and modern Angular patterns. No NgModules for new code.
  Triggers on Angular 17, Angular 18, standalone components, signals,
  Angular new control flow.
category: framework
language: typescript
conflicts: [angular-14]
version: 1.0.0
license: MIT
---

You are working on an Angular 17+ application. Standalone components only. Signals for state. New control flow syntax. No NgModules for new features.

## Standalone component

```typescript
@Component({
  selector: 'app-order-list',
  standalone: true,
  imports: [RouterModule, CurrencyPipe],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    @if (orders().length > 0) {
      <ul>
        @for (order of orders(); track order.id) {
          <li>
            {{ order.id }} — {{ order.total | currency }}
            @if (order.status === 'draft') {
              <button (click)="confirm(order.id)">Confirm</button>
            }
          </li>
        }
      </ul>
    } @else {
      <p>No orders found.</p>
    }
  `
})
export class OrderListComponent {
  private readonly orderService = inject(OrderService);
  readonly orders = toSignal(this.orderService.getAll(), { initialValue: [] });

  confirm(id: string): void {
    this.orderService.confirm(id).subscribe();
  }
}
```

## Signal-based state

```typescript
export class OrderDetailComponent {
  private readonly route = inject(ActivatedRoute);
  private readonly orderService = inject(OrderService);

  readonly orderId = toSignal(
    this.route.paramMap.pipe(map(p => p.get('id') ?? ''))
  );

  readonly order = toSignal(
    toObservable(this.orderId).pipe(
      switchMap(id => this.orderService.getById(id))
    )
  );

  readonly isConfirmed = computed(() => this.order()?.status === 'confirmed');
}
```

## App config (replaces AppModule)

```typescript
export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(withInterceptors([errorInterceptor])),
    provideAnimations(),
  ]
};
```

## Rules to enforce

- inject() over constructor injection in components
- @if / @for over *ngIf / *ngFor — new control flow is faster
- toSignal() to wrap observables — avoid manual subscriptions
- computed() for derived state
- OnPush on every component

## Red flags — stop and warn

- NgModule in new code — standalone is the standard from Angular 17
- *ngIf / *ngFor directives — use new control flow syntax
- Manual subscriptions without cleanup
- any type — always type explicitly
- Constructor injection in components — use inject() function
