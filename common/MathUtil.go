package common

import (
	"code.google.com/p/godec/dec"
)

type MathUtil struct{}

func (c MathUtil) Add(str1 string, str2 string) string {
	dec1 := dec.Dec{}
	dec1.SetString(str1)
	dec2 := dec.Dec{}
	dec2.SetString(str2)
	
	result := dec.Dec{}
	result.Add(&dec1, &dec2)
	return result.String()
}

/**
	return x - y
*/
func (c MathUtil) Sub(str1 string, str2 string) string {
	dec1 := dec.Dec{}
	dec1.SetString(str1)
	dec2 := dec.Dec{}
	dec2.SetString(str2)
	
	result := dec.Dec{}
	result.Sub(&dec1, &dec2)
	return result.String()
}

func (c MathUtil) Mul(str1 string, str2 string) string {
	dec1 := dec.Dec{}
	dec1.SetString(str1)
	dec2 := dec.Dec{}
	dec2.SetString(str2)
	
	result := dec.Dec{}
	result.Mul(&dec1, &dec2)
	return result.String()
}

/**
	return x / y
*/
func (c MathUtil) Quo(str1 string, str2 string, scaler dec.Scaler, rounder dec.Rounder) string {
	dec1 := dec.Dec{}
	dec1.SetString(str1)
	dec2 := dec.Dec{}
	dec2.SetString(str2)
	
	result := dec.Dec{}
	result.Quo(&dec1, &dec2, scaler, rounder)
	return result.String()
}
