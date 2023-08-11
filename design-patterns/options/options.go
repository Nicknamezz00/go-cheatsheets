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

package options

type BigClass struct {
	Options
}

func NewBigClass(opts ...Option) *BigClass {
	b := BigClass{}
	for _, apply := range opts {
		apply(&b.Options)
	}
	return &b
}

type Option func(o *Options)

type Options struct {
	name   string
	age    int
	sex    string
	weight float64
	height float64
	width  float64
	fieldA string
	fieldB string
	fieldC string
}

func WithName(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}

func WithAge(age int) Option {
	return func(o *Options) {
		o.age = age
	}
}

func WithSex(sex string) Option {
	return func(o *Options) {
		o.sex = sex
	}
}

func WithWeight(weight float64) Option {
	return func(o *Options) {
		o.weight = weight
	}
}

func WithHeight(height float64) Option {
	return func(o *Options) {
		o.height = height
	}
}

func WithWidth(width float64) Option {
	return func(o *Options) {
		o.width = width
	}
}

func WithFieldA(fieldA string) Option {
	return func(o *Options) {
		o.fieldA = fieldA
	}
}

func WithFieldB(fieldB string) Option {
	return func(o *Options) {
		o.fieldB = fieldB
	}
}

func WithFieldC(fieldC string) Option {
	return func(o *Options) {
		o.fieldC = fieldC
	}
}
