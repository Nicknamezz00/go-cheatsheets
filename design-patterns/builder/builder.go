/*
 * MIT License
 *
 * Copyright (c) 2023 Runze Wu
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package builder

import "errors"

type Food struct {
	name   string
	weight float64
	price  float64
	brand  string
}

func NewFood(name string, weight float64, price float64, brand string) *Food {
	return &Food{
		name:   name,
		weight: weight,
		price:  price,
		brand:  brand,
	}
}

type FoodBuilder struct {
	Food
}

func NewFoodBuilder() *FoodBuilder {
	return &FoodBuilder{}
}

func (f *FoodBuilder) Build() (*Food, error) {
	if f.name == "" || f.weight == 0 {
		return nil, errors.New("missing field")
	}
	return &Food{
		name:   f.name,
		weight: f.weight,
		price:  f.price,
		brand:  f.brand,
	}, nil
}

func (f *FoodBuilder) Name(name string) *FoodBuilder {
	f.name = name
	return f
}

func (f *FoodBuilder) Weight(weight float64) *FoodBuilder {
	f.weight = weight
	return f
}

func (f *FoodBuilder) Price(price float64) *FoodBuilder {
	f.price = price
	return f
}

func (f *FoodBuilder) Brand(brand string) *FoodBuilder {
	f.brand = brand
	return f
}
