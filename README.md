# \[form3-go\] Form3 Take Home Exercise (Name: Kris Tunprasert)

## Introduction

Hi I'm Kris, a PHP developer. I am new to Golang, have never used it commercially, only for hobby projects and LeetCode. This project was created using OOP design which is what I am most familiar with. Since this project encapsulates only a small subsection of the `api` I made an executive decision to not further split `domain` and `repository` into a more complex structure.

**Current**
- `/domain/account.go`
- `/repository/account.go`

**End game**
- `/domain/organisation/account.go`
- `/domain/organisation/account_identification.go`
- `/domain/.../...`
- `/repository/organisation/account.go`
- `/repository/organisation/account_identification.go`
- `/repository/.../...`

## Structure

The client library is laid out in the following structure:

- `/sdk/`
    - `/client/`
        - `client.go`
        - `error.go`
    - `/domain/`
        - `account.go`
        - `request.go`
        - `response.go`
    - `/repository/`
        - `account.go`

Folder | Explanation 
------ | -----------
`/client/` | A reusable `Client` entity that exposes the base methods to be used by the `Repository` entity 
`/domain/` | Application level entities to be used by all services 
`/repository/` | The main application interface that other application can quickly reuse. Each repository exposes a meaningful `Create`, `Fetch` and `Delete` method that contains exactly what the user needs without any extras

## Tests

Tests will be executed as soon as the containers are up and running

```bash
$ docker-compose up
```

Alternatively, you can launch the containers in detached mode and inspect the stdout logs

```bash
$ docker-compose up -d
$ docker logs form3-client
```

## Example Usage

```go
package main

import (
    "github.com/ktunprasert/form3-go/sdk/client"
    "github.com/ktunprasert/form3-go/sdk/domain"
    "github.com/ktunprasert/form3-go/sdk/repository"
)

func main() {
    client := client.New("http://localhost:8080")
    accountRepo := repository.New(client)

    // Creating an Organisation Account
    created, err := accountRepo.Create(&domain.Account{...})

    // Fetching an Organisation Account
    account, err := accountRepo.Fetch("7f249544-8841-4a1e-9576-8c6b32dce797")

    // Deleting an Organisation Account
    err := accountRepo.Delete(
        id: "7f249544-8841-4a1e-9576-8c6b32dce797",
        version: 1,
    )
}
```

## Further Improvements

I think I still could improve the library by using the `/internal/` to hide the `Client` implementation to only be used internally. This would include moving the errors out into another section and exposing it so that it may be reused by the user.
Domains could also be generated via a factory pattern to enforce client-side validation.
