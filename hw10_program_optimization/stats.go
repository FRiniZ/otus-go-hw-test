package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
	"sync"

	easyjson "github.com/mailru/easyjson"
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"Email,nocopy"` //nolint
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

const nWorkers = 10

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	ch := make(chan []byte, 1000)

	wg := sync.WaitGroup{}
	lock := sync.Mutex{}

	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			user := &User{}
			for b := range ch {
				if err := easyjson.Unmarshal(b, user); err == nil {
					if strings.HasSuffix(user.Email, domain) {
						key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
						lock.Lock()
						result[key]++
						lock.Unlock()
					}
				}
			}
		}()
	}

	reader := bufio.NewScanner(r)

	for reader.Scan() {
		b := reader.Bytes()
		bc := make([]byte, len(b))
		copy(bc, b)
		ch <- bc
	}

	if err := reader.Err(); err != nil {
		return nil, err
	}

	close(ch)
	wg.Wait()

	return result, nil
}
