---
name: angular-14
description: Build Angular 14 applications using TypeScript, NgModules, RxJS, and Angular CLI conventions. Use Angular 14 patterns — no standalone components, no signals. Triggers on: "Angular 14", "Angular NgModule", "Angular RxJS".
category: framework
language: typescript
conflicts: [angular-17]
version: 1.0.0
---

You are working on an Angular 14 application with TypeScript. Use NgModules. No standalone components — that's Angular 15+. RxJS is your state and async tool.

## Module structure

```
src/app/
  core/
    core.module.ts          ← singleton services, guards, interceptors
    services/
    guards/
    interceptors/
  shared/
    shared.module.ts        ← reusable components, pipes, directives
    components/
    pipes/
  features/
    orders/
      orders.module.ts
      orders-routing.module.ts
      components/
      services/
      models/
```

## Component

```typescript
@Component({
  selector: 'app-order-list',
  templateUrl: './order-list.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class OrderListComponent implements OnInit, OnDestroy {
  orders$: Observable<Order[]>;
  private destroy$ = new Subject<void>();

  constructor(private orderService: OrderService) {}

  ngOnInit(): void {
    this.orders$ = this.orderService.getAll().pipe(
      takeUntil(this.destroy$)
    );
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  confirm(id: string): void {
    this.orderService.confirm(id)
      .pipe(takeUntil(this.destroy$))
      .subscribe();
  }
}
```

## Service with HttpClient

```typescript
@Injectable({ providedIn: 'root' })
export class OrderService {
  private readonly apiUrl = '/api/orders';

  constructor(private http: HttpClient) {}

  getAll(): Observable<Order[]> {
    return this.http.get<Order[]>(this.apiUrl);
  }

  getById(id: string): Observable<Order> {
    return this.http.get<Order>(`${this.apiUrl}/${id}`);
  }

  confirm(id: string): Observable<void> {
    return this.http.post<void>(`${this.apiUrl}/${id}/confirm`, {});
  }
}
```

## HTTP interceptor (error handling)

```typescript
@Injectable()
export class ErrorInterceptor implements HttpInterceptor {
  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return next.handle(req).pipe(
      catchError((error: HttpErrorResponse) => {
        if (error.status === 401) {
          // handle unauthorized
        }
        return throwError(() => error);
      })
    );
  }
}
```

## RxJS rules

- Always unsubscribe — use `takeUntil(destroy$)` or `async` pipe
- Never subscribe inside a subscribe — use `switchMap`, `mergeMap`, `concatMap`
- `OnPush` change detection on all components
- Use `async` pipe in templates over manual subscription

## Red flags — stop and warn

- Standalone components — not available in Angular 14
- Signals — not available in Angular 14
- Missing `takeUntil` or `async` pipe — memory leak
- Business logic in components — belongs in services
- `any` type in TypeScript — always type explicitly
- Subscribing inside `ngOnInit` without unsubscribing
