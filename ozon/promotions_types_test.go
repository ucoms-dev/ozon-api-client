package ozon

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

func TestPromotionProductIDExactWire(t *testing.T) {
	t.Parallel()

	valid := []struct {
		name string
		wire string
		want uint64
	}{
		{name: "number above float precision", wire: `9007199254740993`, want: 9007199254740993},
		{name: "string above float precision", wire: `"9007199254740993"`, want: 9007199254740993},
		{name: "maximum number", wire: `18446744073709551615`, want: math.MaxUint64},
		{name: "maximum string with surrounding whitespace", wire: `" 18446744073709551615 "`, want: math.MaxUint64},
	}
	for _, test := range valid {
		t.Run(test.name, func(t *testing.T) {
			var id PromotionProductID
			if err := json.Unmarshal([]byte(test.wire), &id); err != nil {
				t.Fatalf("Unmarshal(%s): %v", test.wire, err)
			}
			if uint64(id) != test.want {
				t.Fatalf("ID = %d, want %d", id, test.want)
			}
			got, err := json.Marshal(id)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if string(got) != strconv.FormatUint(test.want, 10) {
				t.Fatalf("Marshal = %s, want numeric %d", got, test.want)
			}
		})
	}

	for _, wire := range []string{`1.5`, `-1`, `1e3`, `18446744073709551616`, `null`, `""`, `"+1"`, `"1.0"`, `"1e3"`} {
		t.Run("rejects "+wire, func(t *testing.T) {
			var id PromotionProductID
			if err := json.Unmarshal([]byte(wire), &id); err == nil {
				t.Fatalf("Unmarshal(%s) unexpectedly succeeded with %d", wire, id)
			}
		})
	}
}

func TestPromotionMoneyExactWire(t *testing.T) {
	t.Parallel()

	const amount = "12345678901234567890.123456789"
	wire := `{"amount":"` + amount + `","currency":"KZT"}`
	var money PromotionMoney
	if err := json.Unmarshal([]byte(wire), &money); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if money.Amount != amount || money.Currency != "KZT" {
		t.Fatalf("Money = %#v", money)
	}
	got, err := json.Marshal(money)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(got) != wire {
		t.Fatalf("Marshal = %s, want %s", got, wire)
	}

	zero, err := json.Marshal(PromotionMoney{Amount: "0.00", Currency: "RUB"})
	if err != nil || string(zero) != `{"amount":"0.00","currency":"RUB"}` {
		t.Fatalf("zero money = %s, %v", zero, err)
	}

	for _, invalid := range []PromotionMoney{
		{Amount: "", Currency: "RUB"},
		{Amount: "-1", Currency: "RUB"},
		{Amount: "1e3", Currency: "RUB"},
		{Amount: "NaN", Currency: "RUB"},
		{Amount: "1.00", Currency: ""},
	} {
		if _, err := json.Marshal(invalid); err == nil {
			t.Errorf("Marshal(%#v) unexpectedly succeeded", invalid)
		}
	}
}

func FuzzPromotionProductIDJSON(f *testing.F) {
	for _, input := range []string{`0`, `"0"`, `18446744073709551615`, `"9007199254740993"`, `-1`, `1.5`, `1e3`, `null`, `""`} {
		f.Add(input)
	}
	f.Fuzz(func(t *testing.T, input string) {
		var id PromotionProductID
		if err := json.Unmarshal([]byte(input), &id); err == nil {
			if _, err := json.Marshal(id); err != nil {
				t.Fatalf("Marshal after successful Unmarshal: %v", err)
			}
		}
	})
}

func FuzzPromotionMoneyJSON(f *testing.F) {
	for _, input := range []string{
		`{"amount":"0","currency":"RUB"}`,
		`{"amount":"12345678901234567890.123456789","currency":"KZT"}`,
		`{"amount":"-1","currency":"RUB"}`,
		`{"amount":"1e3","currency":"RUB"}`,
		`null`,
	} {
		f.Add(input)
	}
	f.Fuzz(func(t *testing.T, input string) {
		var money PromotionMoney
		if err := json.Unmarshal([]byte(input), &money); err == nil {
			if _, err := json.Marshal(money); err != nil {
				t.Fatalf("Marshal after successful Unmarshal: %v", err)
			}
		}
	})
}
