package utils

import (
	"database/sql"
	"math/big"
	"strconv"
	"strings"
)

func NullStringScan(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NullStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

func ComputeCash(totalMoney int, grandTotal *big.Int) (string, string) {
	s := strconv.Itoa(totalMoney)
	raw := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
	paidBI := new(big.Int)
	paidBI.SetString(raw, 10)
	changeBI := new(big.Int).Sub(paidBI, grandTotal)
	return paidBI.String(), changeBI.String()
}
