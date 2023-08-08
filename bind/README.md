# Bind
Bind binds variables values into string expression in fast and safety way. When the variable pattern is invalid, it return an error.

e.g.: should be `{price}` but written as `{price }` (with space) will return an error.


## Usage
### Bind
```go
    s := "{price} - ({price} * {discount-percentage})"
    v, err := bind.Bind(s,
        "price", 100,
        "discount-percentage", 0.1,
    )
    if err != nil {
        panic(err)
    }

    fmt.Println(v) // "100 - (100 * 0.1)"
```

### SetIdent
Using custom identifier.
```go
    bind.SetIdent(&bind.Ident{
        Prefix: ":", 
        Suffix: "",
    })

    s := ":price - (:price * :discount_percentage)"
    v, err := bind.Bind(s,
        "price", 100,
        "discount_percentage", 0.1,
    )
    if err != nil {
        panic(err)
    }

    fmt.Println(v) // "100 - (100 * 0.1)"
```
