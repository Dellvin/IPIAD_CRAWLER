package main

import (
	"IPIAD_DZ/config"
	"errors"
	"github.com/rai-project/go-fasttext"
	"math"
	"strings"
)

func SentenceVec(sentence string, m *fasttext.Model) ([]float32, error) {
	words := strings.Split(sentence, " ")
	var points []fasttext.Vectors
	for _, w := range words {
		p, err := m.Wordvec(w)
		if err != nil {
			continue
		}
		points = append(points, p)
	}

	var result = make([]float32, points[0].Len())

	for _, p := range points {
		for i := range p {
			result[i] += p[i].Element
		}
	}

	for i := range result {
		result[i] /= float32(len(points))
	}

	return result, nil
}

func Cosine(a []float32, b []float32) (cosine float64, err error) {
	count := 0
	length_a := len(a)
	length_b := len(b)
	if length_a > length_b {
		count = length_a
	} else {
		count = length_b
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= length_a {
			s2 += math.Pow(float64(b[k]), 2)
			continue
		}
		if k >= length_b {
			s1 += math.Pow(float64(a[k]), 2)
			continue
		}
		sumA += float64(a[k]) * float64(b[k])
		s1 += math.Pow(float64(a[k]), 2)
		s2 += math.Pow(float64(b[k]), 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, errors.New("Vectors should not be null (all zeros)")
	}
	return sumA / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}

func calcClusters(m *fasttext.Model) ([]float32, []float32, error) {
	pointA, err := SentenceVec(config.ClusterA, m)
	if err != nil {
		return nil, nil, err
	}

	pointB, err := SentenceVec(config.ClusterB, m)
	if err != nil {
		return nil, nil, err
	}

	return pointA, pointB, nil
}

func chooseCluster(a, b []float32, news string, m *fasttext.Model) (bool, error) { // if true - cluster A, false - cluster B
	pointNews, err := SentenceVec(news, m)
	if err != nil {
		return false, err
	}

	aCluster, err := Cosine(a, pointNews)
	if err != nil {
		return false, err
	}
	bCluster, err := Cosine(b, pointNews)
	if err != nil {
		return false, err
	}

	return aCluster > bCluster, nil
}
