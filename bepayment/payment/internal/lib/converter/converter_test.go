package converter_test

import (
	"testing"

	"github.com/chonlatee/payment/internal/lib/converter"
	"github.com/jackc/pgx/v5/pgtype"
)

func Test_PgTypeNumericsToFloat64s(t *testing.T) {
	t.Run("all valid", func(t *testing.T) {

		var n pgtype.Numeric
		err := n.Scan("10.20")
		if err != nil {
			t.Error(err)
		}

		var n1 pgtype.Numeric
		err = n1.Scan("222.2222")
		if err != nil {
			t.Error(err)
		}

		in := []pgtype.Numeric{
			n,
			n1,
		}

		r, err := converter.PgtypeNumericsToFloat64s(in)
		if err != nil {
			t.Fatal(err)
		}

		if len(r) != 2 {
			t.Errorf("want len 2 got %d", len(r))
		}

		if r[0] != 10.20 {
			t.Errorf("want 10.20 got %f", r[0])
		}

		if r[1] != 222.2222 {
			t.Errorf("want 222.2222 got %f", r[1])
		}

	})
}
