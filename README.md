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
