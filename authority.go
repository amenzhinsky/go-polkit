package polkit

import (
	"github.com/godbus/dbus"
)

type PKImplicitAuthorization uint32

const (
	NotAuthorized PKImplicitAuthorization = iota
	AuthenticationRequired
	AdministratorAuthenticationRequired
	AuthenticationRequiredRetained
	AdministratorAuthenticationRequiredRetained
	Authorized
)

func (i PKImplicitAuthorization) String() string {
	switch i {
	case NotAuthorized:
		return "no"
	case AuthenticationRequired:
		return "auth_self"
	case AdministratorAuthenticationRequired:
		return "auth_admin"
	case AuthenticationRequiredRetained:
		return "auth_self_keep"
	case AdministratorAuthenticationRequiredRetained:
		return "auth_admin_keep"
	case Authorized:
		return "yes"
	default:
		panic("unknown flag value")
	}
}

const (
	CheckAuthorizationNone uint32 = iota
	CheckAuthorizationAllowUserInteraction
)

type (
	Authority struct {
		conn    *dbus.Conn
		object  dbus.BusObject
		subject PKSubject
	}

	PKSubject struct {
		Kind    string                  `dbus:"subject_kind"`
		Details map[string]dbus.Variant `dbus:"subject_details"`
	}

	PKAuthorizationResult struct {
		IsAuthorized bool              `dbus:"is_authorized"`
		IsChallenge  bool              `dbus:"is_challenge"`
		Details      map[string]string `dbus:"details"`
	}

	PKActionDescription struct {
		ActionID         string            `dbus:"action_id"`
		Description      string            `dbus:"description"`
		Message          string            `dbus:"message"`
		VendorName       string            `dbus:"vendor_name"`
		VendorURL        string            `dbus:"vendor_url"`
		IconName         string            `dbus:"icon_name"`
		ImplicitAny      uint32            `dbus:"implicit_any"`
		ImplicitInactive uint32            `dbus:"implicit_inactive"`
		ImplicitActive   uint32            `dbus:"implicit_active"`
		Annotations      map[string]string `dbus:"annotations"`
	}
)

func NewAuthority() (*Authority, error) {
	bus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	names := bus.Names()
	if len(names) == 0 {
		panic("empty dbus names")
	}

	return &Authority{
		conn:   bus,
		object: bus.Object("org.freedesktop.PolicyKit1", "/org/freedesktop/PolicyKit1/Authority"),
		subject: PKSubject{
			Kind: "system-bus-name",
			Details: map[string]dbus.Variant{
				"name": dbus.MakeVariant(names[0]),
			},
		},
	}, nil
}

func (a *Authority) EnumerateActions(locale string) ([]PKActionDescription, error) {
	var result []PKActionDescription
	if err := a.call("EnumerateActions", &result, locale); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Authority) CheckAuthorization(
	actionID string,
	details map[string]string,
	flags uint32,
	cancellationID string) (*PKAuthorizationResult, error) {
	result := PKAuthorizationResult{}
	if err := a.call("CheckAuthorization", &result, a.subject, actionID, details, flags, cancellationID); err != nil {
		return nil, err
	}

	return &result, nil
}

func (a *Authority) CancelCheckAuthorization(cancellationID string) error {
	if err := a.call("CancelCheckAuthorization", nil, cancellationID); err != nil {
		return err
	}
	return nil
}

func (a *Authority) call(action string, result interface{}, args ...interface{}) error {
	call := a.object.Call("org.freedesktop.PolicyKit1.Authority."+action, 0, args...)

	if result != nil {
		if err := call.Store(result); err != nil {
			return err
		}
	}

	return nil
}

func (a *Authority) Close() error {
	return a.conn.Close()
}
