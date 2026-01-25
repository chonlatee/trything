package converter

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func StringsToPgTypeUUIDs(in []string) ([]pgtype.UUID, error) {
	out := make([]pgtype.UUID, len(in))

	for i, s := range in {
		var u pgtype.UUID
		if err := u.Scan(s); err != nil {
			return nil, err
		}
		out[i] = u
	}

	return out, nil
}

func Float64sToPgTypeNumerics(in []float64) ([]pgtype.Numeric, error) {
	out := make([]pgtype.Numeric, len(in))

	for i, s := range in {
		var u pgtype.Numeric
		if err := u.Scan(s); err != nil {
			return nil, err
		}
		out[i] = u
	}

	return out, nil
}

func TimesToPgTypeTimes(in []time.Time) ([]pgtype.Timestamptz, error) {
	out := make([]pgtype.Timestamptz, len(in))

	for i, s := range in {
		var u pgtype.Timestamptz
		if err := u.Scan(s); err != nil {
			return nil, err
		}

		out[i] = u
	}

	return out, nil
}

func StringToPgtypeUUID(in string) (pgtype.UUID, error) {
	var out pgtype.UUID
	if err := out.Scan(in); err != nil {
		return pgtype.UUID{}, err
	}

	return out, nil
}

func Float64ToPgtypeNumeric(in float64) (pgtype.Numeric, error) {
	var out pgtype.Numeric
	if err := out.Scan(in); err != nil {
		return pgtype.Numeric{}, err
	}

	return out, nil
}

func StringToPgtypeText(in string) (pgtype.Text, error) {
	var out pgtype.Text
	if err := out.Scan(in); err != nil {
		return pgtype.Text{}, err
	}

	return out, nil
}

func TimeToPgtypeTimestamptz(in time.Time) (pgtype.Timestamptz, error) {
	var out pgtype.Timestamptz
	if err := out.Scan(in); err != nil {
		return pgtype.Timestamptz{}, nil
	}

	return out, nil
}

func PgtypeUUIDsToStrings(in []pgtype.UUID) ([]string, error) {

	out := make([]string, len(in))

	for i, o := range in {
		if !o.Valid {
			return nil, fmt.Errorf("%s not valid uuid", in[i].Bytes)
		}
		out[i] = o.String()
	}

	return out, nil
}

func PgtypeNumericsToFloat64s(in []pgtype.Numeric) ([]float64, error) {
	out := make([]float64, len(in))

	for i, o := range in {
		if !o.Valid {
			return nil, fmt.Errorf("index %d not valid", i)
		}

		r, err := o.Float64Value()
		if err != nil {
			return nil, err
		}

		out[i] = r.Float64
	}

	return out, nil
}

func PgtypeTimestamptzsToTimes(in []pgtype.Timestamptz) ([]time.Time, error) {
	out := make([]time.Time, len(in))
	for i, o := range in {
		if !o.Valid {
			return nil, fmt.Errorf("index %d not valid", i)
		}

		out[i] = o.Time

	}

	return out, nil
}

func PgtypeUUIDToString(in pgtype.UUID) (string, error) {
	if !in.Valid {
		return "", fmt.Errorf("uuid not valid")
	}

	return in.String(), nil
}

func PgtypeNumericToFloat64(in pgtype.Numeric) (float64, error) {
	if !in.Valid {
		return 0, fmt.Errorf("value not valid")
	}

	r, err := in.Float64Value()
	if err != nil {
		return 0, fmt.Errorf("cat not get float64 value")
	}

	return r.Float64, nil
}

func PgtypeTextToString(in pgtype.Text) (string, error) {
	if !in.Valid {
		return "", fmt.Errorf("value not valid")
	}

	return in.String, nil
}

func PgtypeTimestamptzToTime(in pgtype.Timestamptz) (time.Time, error) {
	if !in.Valid {
		return time.Time{}, fmt.Errorf("value not valid")
	}

	return in.Time, nil
}
