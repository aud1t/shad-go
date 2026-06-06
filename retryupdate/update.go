//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	authErr := new(kvapi.AuthError)
	conflictErr := new(kvapi.ConflictError)

outerLoop:
	for {
		var oldValue *string
		var oldVersion uuid.UUID

		getResp, getErr := c.Get(&kvapi.GetRequest{
			Key: key,
		})
		switch {
		case getErr == nil:
			oldValue = &getResp.Value
			oldVersion = getResp.Version
		case errors.Is(getErr, kvapi.ErrKeyNotFound):
		case errors.As(getErr, &authErr):
			return getErr
		default:
			continue
		}

		newValue, err := updateFn(oldValue)
		if err != nil {
			return err
		}
		newVersion := uuid.Must(uuid.NewV4())

	innerLoop:
		for {
			_, setErr := c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      newValue,
				OldVersion: oldVersion,
				NewVersion: newVersion,
			})
			switch {
			case setErr == nil:
				break outerLoop
			case errors.As(setErr, &authErr):
				return setErr
			case errors.Is(setErr, kvapi.ErrKeyNotFound):
				newValue, err = updateFn(nil)
				if err != nil {
					return err
				}
				oldVersion = uuid.UUID{}
			case errors.As(setErr, &conflictErr):
				if conflictErr.ExpectedVersion == newVersion {
					break outerLoop
				}
				break innerLoop
			default:
				continue
			}
		}
	}
	return nil
}
