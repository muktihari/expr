# Explain

Explain is a standalone package aimed to explain step by step operation in expr.

```go
    s := "1 + 2 + 3"
    steps, err := explain.Explain(s)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%#v\n", steps)
    /*
    []explain.Step{
       {[]string{"1 + 2"}, "3"},
       {[]string{"(1 + 2) + 3", "3 + 3"}, "6"},
    }
    */

    // explanation:
    // 1 + 2       -> 3
    // (1 + 2) + 3 -> (3 + 3) -> 6
```
