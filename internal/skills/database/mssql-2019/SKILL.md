---
name: mssql-2019
description: Write and review T-SQL for MS SQL Server 2019. Enforce parameterized queries, proper indexing, execution plan awareness, and performance best practices. Triggers on: "MS SQL Server", "SQL Server 2019", "T-SQL", "MSSQL", "stored procedures".
category: database
language: sql
conflicts: []
version: 1.0.0
---

You are working with MS SQL Server 2019. Every query you write or review must be safe, performant, and maintainable. Performance is not optional — it is a requirement.

## Parameterized queries — non-negotiable

```sql
-- WRONG — SQL injection vulnerability, never do this
DECLARE @sql NVARCHAR(MAX) = 'SELECT * FROM Orders WHERE CustomerId = ' + @CustomerId
EXEC(@sql)

-- CORRECT — always parameterized
SELECT Id, Status, Total
FROM Orders
WHERE CustomerId = @CustomerId
  AND Status = @Status
```

## Query structure — explicit columns always

```sql
-- WRONG — never SELECT * in production
SELECT * FROM Orders

-- CORRECT — explicit columns, meaningful aliases
SELECT
    o.Id           AS OrderId,
    o.Status       AS OrderStatus,
    o.CreatedAt    AS CreatedDate,
    c.Name         AS CustomerName,
    SUM(l.Quantity * l.UnitPrice) AS Total
FROM Orders o
INNER JOIN Customers c ON c.Id = o.CustomerId
INNER JOIN OrderLines l ON l.OrderId = o.Id
WHERE o.Status = @Status
  AND o.CreatedAt >= @FromDate
GROUP BY o.Id, o.Status, o.CreatedAt, c.Name
ORDER BY o.CreatedAt DESC
```

## CTEs for readable complex queries

```sql
WITH ConfirmedOrders AS (
    SELECT
        o.Id,
        o.CustomerId,
        SUM(l.Quantity * l.UnitPrice) AS Total
    FROM Orders o
    INNER JOIN OrderLines l ON l.OrderId = o.Id
    WHERE o.Status = 'Confirmed'
      AND o.CreatedAt >= DATEADD(MONTH, -3, GETUTCDATE())
    GROUP BY o.Id, o.CustomerId
),
HighValueOrders AS (
    SELECT * FROM ConfirmedOrders WHERE Total > 1000
)
SELECT
    c.Name,
    hv.Total,
    hv.Id AS OrderId
FROM HighValueOrders hv
INNER JOIN Customers c ON c.Id = hv.CustomerId
ORDER BY hv.Total DESC;
```

## Indexing rules

```sql
-- Index columns used in WHERE, JOIN ON, ORDER BY
CREATE NONCLUSTERED INDEX IX_Orders_CustomerId_Status
ON Orders (CustomerId, Status)
INCLUDE (CreatedAt, Total);  -- covering index — avoids key lookup

-- Never index every column — indexes cost write performance
-- Always check existing indexes before creating new ones
SELECT * FROM sys.indexes WHERE object_id = OBJECT_ID('Orders')
```

## Performance checks before shipping any query

```sql
-- Always run EXPLAIN before deploying to production
SET STATISTICS IO ON;
SET STATISTICS TIME ON;

-- Check for these in execution plan:
-- ✗ Table Scan → add an index
-- ✗ Key Lookup → add INCLUDE columns to index
-- ✗ Implicit Conversion → fix data type mismatch in WHERE clause
-- ✗ Parameter Sniffing → use OPTION (OPTIMIZE FOR UNKNOWN) if needed
```

## Parameter sniffing — know this problem

```sql
-- Symptom: query runs fast first time, slow after recompile
-- Fix: optimize for unknown or use local variables
CREATE PROCEDURE GetOrdersByCustomer
    @CustomerId UNIQUEIDENTIFIER
AS
BEGIN
    DECLARE @LocalCustomerId UNIQUEIDENTIFIER = @CustomerId  -- local variable trick

    SELECT Id, Status, Total
    FROM Orders
    WHERE CustomerId = @LocalCustomerId
    OPTION (RECOMPILE)  -- or use this for queries with high variance
END
```

## Transactions

```sql
BEGIN TRANSACTION
BEGIN TRY
    UPDATE Orders SET Status = 'Confirmed' WHERE Id = @OrderId
    INSERT INTO OrderEvents (OrderId, EventType, CreatedAt)
    VALUES (@OrderId, 'Confirmed', GETUTCDATE())
    COMMIT TRANSACTION
END TRY
BEGIN CATCH
    ROLLBACK TRANSACTION
    THROW
END CATCH
```

## Red flags — stop and warn

- `SELECT *` in any production query
- String concatenation to build dynamic SQL
- Missing `WHERE` clause on `UPDATE` or `DELETE`
- `NOLOCK` hint without understanding dirty reads
- Implicit data type conversions in `WHERE` clauses — kills index usage
- Cursor when set-based operation is possible
- Missing index on foreign key columns
- `GETDATE()` instead of `GETUTCDATE()` — always store UTC
