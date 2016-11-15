# go-polkit

PolKit client library and CLI for golang.

## CLI Installation
```
$ cd $GOPATH/src/github.com/amenzhinsky/go-polkit
$ go install 
```

## CLI Usage

List all PolKit actions:
```
$ pkquery -verbose
```

View particular actions:
```
$ pkquery -verbose org.freedesktop.udisks2.filesystem-fstab org.freedesktop.udisks2.filesystem-mount
```

Check access to an action:
```
$ pkquery -check-access org.freedesktop.udisks2.filesystem-fstab
```

Check access and allow user to interact by typing his password when PolKit requires it:
```
$ pkquery -check-access -allow-password org.freedesktop.udisks2.filesystem-fstab
```

## Library Usage

### Authorization checking
```go
authority, err := polkit.NewAuthority()
if err != nil {
	panic(err)
}

result, err := authority.CheckAuthorization(
	"org.freedesktop.udisks2.filesystem-fstab",
	nil,
	polkit.CheckAuthorizationAllowUserInteraction, "",
)

if err != nil {
	panic(err)
}

fmt.Printf("Is authorized: %t\n", result.IsAuthorized)
fmt.Printf("Is challenge:  %t\n", result.IsChallenge)
fmt.Printf("Details:       %v\n", result.Details)
```

### Actions Enumerating
```go
authority, err := polkit.NewAuthority()
if err != nil {
	panic(err)
}

actions, err := authority.EnumerateActions("")
if err != nil {
	panic(err)
}

for _, action := range actions {
	fmt.Println(action.ActionID)
}
````
