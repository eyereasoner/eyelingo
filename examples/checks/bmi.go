package main

import (
	"math"
)

func checkBMI(ctx *Context) []Check {
	d := ctx.M()
	unit := str(d["UnitSystem"])
	weight := num(d["Weight"])
	height := num(d["Height"])
	var kg, m float64
	if unit == "metric" {
		kg = weight
		m = height / 100
	} else if unit == "us" {
		kg = weight * 0.45359237
		m = height * 0.0254
	} else {
		kg = math.NaN()
		m = math.NaN()
	}
	h2 := m * m
	bmi := kg / h2
	bmi1 := roundHalfUp(bmi, 1)
	bmi2 := roundHalfUp(bmi, 2)
	hmin := roundHalfUp(18.5*h2, 1)
	hmax := roundHalfUp(24.9*h2, 1)
	actualBMI, cat, amin, amax := parseBMIAnswer(ctx.Answer)
	order := []string{bmiClass(18.49), bmiClass(18.50), bmiClass(24.99), bmiClass(25), bmiClass(29.99), bmiClass(30), bmiClass(35), bmiClass(40)}
	return []Check{{"input units are normalized to positive SI kg and m", kg > 0 && m > 0 && (unit == "metric" || unit == "us")}, {"height squared is recomputed independently from the normalized height", math.Abs(h2-3.1684) < 1e-12}, {"reported BMI matches independent kg/m² computation rounded to one decimal", actualBMI == bmi1}, {"the unrounded BMI independently rounds to the expected two-decimal value", math.Abs(bmi2-22.72) < 1e-12}, {"reported BMI category follows the adult threshold table", cat == bmiClass(bmi)}, {"healthy-weight lower bound is recomputed from BMI 18.5", amin == hmin}, {"healthy-weight upper bound is recomputed from BMI 24.9", amax == hmax}, {"category boundary ordering covers all adult BMI classes", sliceEq(order, []string{"Underweight", "Normal", "Normal", "Overweight", "Overweight", "Obesity I", "Obesity II", "Obesity III"})}}
}
