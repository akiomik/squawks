// Copyright 2022 Akiomi Kamakura
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package json

import (
	"strconv"
)

// https://developer.twitter.com/en/docs/twitter-api/v1/data-dictionary/object-model/geo

type LatLong [2]float64

func (l LatLong) Latitude() float64 {
	return l[0]
}

func (l LatLong) Longitude() float64 {
	return l[1]
}

func (l LatLong) String() string {
	return strconv.FormatFloat(l.Latitude(), 'f', -1, 64) + "," + strconv.FormatFloat(l.Longitude(), 'f', -1, 64)
}

type LongLat [2]float64

func (l LongLat) Longitude() float64 {
	return l[0]
}

func (l LongLat) Latitude() float64 {
	return l[1]
}

func (l LongLat) String() string {
	return strconv.FormatFloat(l.Longitude(), 'f', -1, 64) + "," + strconv.FormatFloat(l.Latitude(), 'f', -1, 64)
}

type Geo struct {
	Type        string  `json:"type"`
	Coordinates LatLong `json:"coordinates"`
}

type Coordinates struct {
	Type        string  `json:"type"`
	Coordinates LongLat `json:"coordinates"`
}

type BoundingBox struct {
	Type        string      `json:"type"`
	Coordinates [][]LongLat `json:"coordinates"`
}

type Place struct {
	Id          string      `json:"id"`
	Url         string      `json:"url"`
	PlaceType   string      `json:"place_type"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	CountryCode string      `json:"country_code"`
	Country     string      `json:"country"`
	BoundingBox BoundingBox `json:"bounding_box"`
}
