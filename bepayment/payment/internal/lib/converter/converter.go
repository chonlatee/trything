package converter

import (
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
