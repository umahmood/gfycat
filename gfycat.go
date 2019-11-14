package gfycat

import (
	"bufio"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

// lockedRandSource private random generator
type lockedRandSource struct {
	lock sync.Mutex
	src  rand.Source
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64. to
// satisfy rand.Source interface
func (r *lockedRandSource) Int63() int64 {
	r.lock.Lock()
	ret := r.src.Int63()
	r.lock.Unlock()
	return ret
}

// Seed uses the provided seed value to initialize the random generator. to
// satisfy rand.Source interface
func (r *lockedRandSource) Seed(seed int64) {
	r.lock.Lock()
	r.src.Seed(seed)
	r.lock.Unlock()
}

// Getter HTTP Get request
type Getter interface {
	Get(*http.Request) (*http.Response, error)
}

// retriever instance
type retriever struct{}

// Get HTTP request
func (retriever) Get(req *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Order in which the animal appears in the output
type Order int

// Order in which the animal appears in the output
const (
	AnimalFirst Order = iota + 1
	AnimalSecond
	AnimalThird
)

// Gfycat instance
type Gfycat struct {
	random     *rand.Rand
	adjectives []string
	animals    []string
	assetsPath string
}

// New construct a new Gfycat instance
func New() (*Gfycat, error) {
	adj, err := loadFile("adjectives.txt")
	if err != nil {
		return nil, err
	}
	ani, err := loadFile("animals.txt")
	if err != nil {
		return nil, err
	}
	return &Gfycat{
		random: rand.New(
			&lockedRandSource{
				src: rand.NewSource(time.Now().UnixNano()),
			}),
		adjectives: adj,
		animals:    ani,
	}, nil
}

// GenerateNameOrder generates gfycat name with a specific ordering
func (g *Gfycat) GenerateNameOrder(order Order) string {
	w1 := g.random.Intn(len(g.adjectives))
	w2 := g.random.Intn(len(g.adjectives))
	w3 := g.random.Intn(len(g.animals))

	a := g.adjectives[w1]
	b := g.adjectives[w2]
	c := g.animals[w3]
	switch order {
	case AnimalFirst:
		return c + b + a
	case AnimalSecond:
		return a + c + b
	case AnimalThird:
		fallthrough
	default:
		return a + b + c
	}
}

// GenerateName generates gfycat name
func (g *Gfycat) GenerateName() string {
	return g.GenerateNameOrder(AnimalThird)
}

// loadfile retrieves the requested file
func loadFile(fileName string) ([]string, error) {
	if !fileExists(fileName) {
		fileData, err := fileRetriever(retriever{}, fileName)
		if err != nil {
			return nil, err
		}
		err = createFile(fileName, fileData)
		if err != nil {
			return nil, err
		}
	}
	fileData, err := readFile(fileName)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

// fileRetriever retrieves the requested resource from the network
func fileRetriever(g Getter, fileName string) ([]byte, error) {
	// remove file extension if present
	idx := strings.Index(fileName, ".")
	if idx != -1 {
		fileName = fileName[:idx]
	}
	gfycatURL := "https://assets.gfycat.com/"
	req, err := http.NewRequest("GET", gfycatURL+fileName, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.Get(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	return ioutil.ReadAll(resp.Body)
}

// readFile returns each line in a file
func readFile(fileName string) ([]string, error) {
	path := filepath.Join(basePath(), fileName)
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
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, err
}

// fileExists checks if a files exists on the file system
func fileExists(fileName string) bool {
	path := filepath.Join(basePath(), fileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// createFile creates the directory path and file
func createFile(fileName string, data []byte) error {
	path := basePath()
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	filePath := filepath.Join(path, fileName)
	err = ioutil.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// basePath returns the full path to where the assets are stored
func basePath() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".gfycat")
}
