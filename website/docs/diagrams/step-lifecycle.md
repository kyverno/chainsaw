# Step lifecycle

A test step is made of 3 main blocks used to determine the actions Chainsaw will perform, depending on operations outcome:

- The `try` block *(required)*
- The `catch` block *(optional)*
- The `finally` block *(optional)*

Each block can be represented as an ordered sequence of operations.

## Try, Catch, Finally flow

Operations defined in the `try` block are executed first, then:

- If an operation fails to execute, Chainsaw won't execute the remaining operations and will execute **all** operations defined in the `catch` block instead (if any).
- If all operations succeed, Chainsaw will NOT execute operations defined in the `catch` block (if any).
- Regardless of the step outcome (success or failure), Chainsaw will execute **all** operations defined in the `finally` block (if any).

!!! tip

    Note that all operations coming from the `catch` or `finally` blocks are executed. If one operation fails, Chainsaw will mark the test as failed and continue executing with the next operations.

## Sequence

### Without failure

<div style="text-align: center;">

```mermaid
sequenceDiagram
    autonumber
    participant S1 as Step N
    create participant T as try ...
        S1 ->>  T  : execute
        T  -->> S1 : success

    create participant F as finally ...
        S1 ->>  F  : execute
        F  -->> S1 : done

    participant S2 as Step N+1

    S1  ->> S2 : next step
```

```mermaid
sequenceDiagram
    autonumber
    participant T as Try

    create participant T1 as Op 1
        T ->>  T1  : execute
    create participant T2 as Op 2
        T1 ->>  T2  : execute

    participant C as Catch

    participant F as Finally

    T2 -->> F: done

    create participant F1 as Op 1
        F ->>  F1  : execute
    create participant F2 as Op 2
        F1 ->>  F2  : execute

    F2 -->> T: done
```

```mermaid
sequenceDiagram
    autonumber

    participant S as Step N

    box Try block
    participant T1 as Op 1
    participant T2 as Op N
    end
    box Catch block
    end
    box Finally block
    participant F1 as Op 1
    participant F2 as Op N
    end
    participant S1 as Step N+1

    S  -->> T1 : try
    T1 ->>  T2 : success
    T2 -->> S  : done
    S  -->> F1 : finally
    F1 ->>  F2 : done
    F2 -->> S  : done
    S  -->> S1 : next step
```

```mermaid
graph TD
    start --> t1
    start ~~~ c1
    start ~~~ f1

    subgraph try
        direction LR
        t1 --> t2 --> t3
    end
    subgraph catch
        direction LR
        c1 ~~~ c2 ~~~ c3
    end
    subgraph finally
        direction LR
        f1 --> f2 --> f3
    end

    t3 --> f1
    f3 --> finish
```

</div>

!!! info ""
    1. Test starts by executing Step 1
    1. Step 1 terminates -> Step 2 starts executing
    1. Step 2 terminates -> Step 3 starts executing
    1. Step 3 terminates -> Cleanup for Step 3 starts
    1. Cleanup for Step 3 terminates -> Cleanup for Step 2 starts
    1. Cleanup for Step 2 terminates -> Cleanup for Step 1 is executed

### With failure

<div style="text-align: center;">

```mermaid
sequenceDiagram
    autonumber
    participant S0 as Step N-1
    participant S1 as Step N
    create participant T as try ...
        S1 ->>  T  : execute
        T  -->> S1 : error

    create participant C as catch ...
        S1 ->>  C  : execute
        C  -->> S1 : done

    create participant F as finally ...
        S1 ->>  F  : execute
        F  -->> S1 : done

    S1  -->> S0 : error

```

</div>
