package test

import (
	"testing"

	investBasic "github.com/rz1998/invest-basic"
)

// TestGetSecInfoQmtSuffix locks the QMT-suffix acceptance (2026-10-22,
// best-shm 后缀统一): after the strategy converts every code to the QMT
// form at the signal boundary, GetSecInfo must still classify them.
//
// This is the correctness hard-dependency: GetMinVol / GetTickVol /
// GetTickPrice / GetVolSellableFromAcPos all branch on the ExchangeCD
// this function returns. A .SH that falls to the default branch means
// min vol = 1, tick = 1, wrong sellable formula -- all silent.
func TestGetSecInfoQmtSuffix(t *testing.T) {
	cases := []struct {
		code  string
		want  string
		exch  investBasic.ExchangeCD
	}{
		{"600000.SH", "600000", investBasic.SSE},
		{"000001.SZ", "000001", investBasic.SZSE},
		{"430047.BJ", "430047", investBasic.BSE},
		// The Go form must keep working (both directions, one function).
		{"600000.SSE", "600000", investBasic.SSE},
		{"000001.SZSE", "000001", investBasic.SZSE},
		{"430047.BSE", "430047", investBasic.BSE},
		// Suffix is case-insensitive.
		{"600000.sh", "600000", investBasic.SSE},
		{"600000.Sh", "600000", investBasic.SSE},
		// A bare code keeps the old behaviour: empty exchange.
		{"600000", "600000", ""},
		// An unknown suffix is passed through untouched (old behaviour).
		{"IC2406.CFFEX", "IC2406", "CFFEX"},
	}
	for _, c := range cases {
		code, exch := investBasic.GetSecInfo(c.code)
		if code != c.want || exch != c.exch {
			t.Errorf("GetSecInfo(%q) = (%q, %q), want (%q, %q)",
				c.code, code, exch, c.want, c.exch)
		}
	}
}

// TestGetSecInfoQmtSuffixDrivesTheBranches proves the classification
// actually reaches the suffix-sensitive branches downstream, not just
// the function in isolation.
func TestGetSecInfoQmtSuffixDrivesTheBranches(t *testing.T) {
	// GetSecType branches on the exchange: a .SH that is not recognised
	// would answer its default (unknown type), not STOCK.
	if got := investBasic.GetSecType("600000.SH"); got != investBasic.STOCK {
		t.Errorf("GetSecType(600000.SH) = %q, want STOCK", got)
	}
	if got := investBasic.GetSecType("000001.SZ"); got != investBasic.STOCK {
		t.Errorf("GetSecType(000001.SZ) = %q, want STOCK", got)
	}
	// GetTickPrice branches on the exchange + prefix: the main-board
	// share answers 100. Unrecognised would fall through to 0/default.
	if got := investBasic.GetTickPrice("600519.SH"); got != 100 {
		t.Errorf("GetTickPrice(600519.SH) = %d, want 100", got)
	}
	// 511090 is an ETF: SSE + "5" prefix answers FUND.
	if got := investBasic.GetSecType("511090.SH"); got != investBasic.FUND {
		t.Errorf("GetSecType(511090.SH) = %q, want FUND", got)
	}
}
