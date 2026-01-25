package converter

import "github.com/jackc/pgx/v5/pgtype"

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
