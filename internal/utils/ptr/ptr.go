package ptr

import (
	"time"
)

func ToPtr[T any](v T) *T {
	return &v
}

func DerefOrNil[T any](v *T) any {
	if v == nil {
		return nil
	}

	return *v
}

func DerefDefault[T any](v *T, d T) T {
	if v == nil {
		return d
	}

	return *v
}

func UnsafeDerefOrNil[T any](v *T) any {
	if v == nil {
		return nil
	}

	return *v
}

// PtrMilisToTime converts a millisecond timestamp to a time.Time pointer
func PtrMilisToTime(v *int64) *time.Time {
	if v == nil {
		return nil
	}
	return ToPtr(time.UnixMilli(*v))
}

// PtrTimeToMilis converts a time.Time pointer to a millisecond timestamp
func PtrTimeToMilis(v *time.Time) *int64 {
	if v == nil {
		return nil
	}

	return ToPtr(v.UnixMilli())
}

func BrandedToStringPtr[T ~string](b *T) *string {
	if b == nil {
		return nil
	}
	str := string(*b)
	return &str
}
