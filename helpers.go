package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
)

func randChoice([]string) string {
	rand.Seed(time.Now().Unix())
	agent := googleAgents[rand.Intn(len(googleAgents))]
	return agent
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func randString(length int) string {
	const (
		letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		size    = len(letters)
	)
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(size)]
	}
	return string(b)
}

func strInSlice(str string, slice []string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func getUUID() string {
	return uuid.New().String()
}

func SHA1hash(str string) string {
	hash := sha1.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func SHA256hash(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

func sliceToIntSet(intSlice []int) []int {
	return []int{}
}
