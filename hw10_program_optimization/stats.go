/*
package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
	"sync"
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"Email,nocopy"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int
type PtrString *string

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {

	result := make(DomainStat)
	ch := make(chan string, 10)
	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			user := &User{}
			for {
				s, ok := <-ch
				if !ok {
					break
				}
				if err := json.Unmarshal([]byte(s), user); err == nil {

					if strings.HasSuffix(user.Email, domain) {
						lock.Lock()
						result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
						lock.Unlock()
					}
				}
			}
		}()
	}

	reader := bufio.NewScanner(r)
	for reader.Scan() {
		ch <- reader.Text()
	}

	if err := reader.Err(); err != nil {
		return nil, err
	}

	close(ch)
	wg.Wait()

	return result, nil
}
*/

package hw10programoptimization

import (
	"encoding/json"
	"io"
	"strings"
)

//easyjson:json
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	user := &User{}
	result := make(DomainStat)

	decoder := json.NewDecoder(r)

	for decoder.More() {
		err := decoder.Decode(user)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
