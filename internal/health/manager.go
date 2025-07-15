package health

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/database"
	"github.com/google/uuid"
)

type HealthResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}

type HealthManager struct {
	lockKey        string
	instanceID     string
	processors     map[string]string
	healthTicker   *time.Ticker
	stopHealth     chan struct{}
	healthStartOnce sync.Once
}

func NewHealthManager() *HealthManager {
	return &HealthManager{
		instanceID: uuid.NewString(),
		lockKey:    "health-manager-leader",
		processors: map[string]string{
			"default":  config.ProcessorDefaultURL + "/payments/service-health",
			"fallback": config.ProcessorFallbackURL + "/payments/service-health",
		},
	}
}

func (m *HealthManager) Start() {
	m.ensureInitialHealth()

	electionTicker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-electionTicker.C:
				m.evaluateLeadership()
			case <-database.RedisContext.Done():
				electionTicker.Stop()
				m.stopHealthCheck()
				return
			}
		}
	}()
}

func (m *HealthManager) ensureInitialHealth() {
	for name := range m.processors {
		key := getKey(name)
		exists, err := database.RedisClient.Exists(database.RedisContext, key).Result()
		if err != nil {
			log.Printf("Error checking existence of key %s: %v", key, err)
			continue
		}
		if exists == 0 {
			defaultHealth := &HealthResponse{Failing: false, MinResponseTime: 0}
			if err := saveHealthResponse(name, defaultHealth); err != nil {
				log.Printf("Error saving default health for %s: %v", name, err)
			} else {
				log.Printf("Default health status initialized for %s: %+v", name, defaultHealth)
			}
		}
	}
}

func (m *HealthManager) evaluateLeadership() {
	isLeader, err := m.tryToBecomeLeader()
	if err != nil {
		log.Printf("Error during leader election: %v", err)
		return
	}

	if isLeader {
		log.Println("Acting as leader. Ensuring health checks are running.")
		m.startHealthCheck()
	} else {
		log.Println("Not the leader.")
		m.stopHealthCheck()
	}
}

func (m *HealthManager) startHealthCheck() {
	m.healthStartOnce.Do(func() {
		m.healthTicker = time.NewTicker(5 * time.Second)
		m.stopHealth = make(chan struct{})

		go func() {
			for {
				select {
				case <-m.healthTicker.C:
					m.runHealthCheck()
				case <-m.stopHealth:
					m.healthTicker.Stop()
					m.healthTicker = nil
					m.healthStartOnce = sync.Once{}
					return
				}
			}
		}()
	})
}

func (m *HealthManager) stopHealthCheck() {
	if m.stopHealth != nil {
		close(m.stopHealth)
		m.stopHealth = nil
	}
}

func (m *HealthManager) runHealthCheck() {
	client := &http.Client{Timeout: 1 * time.Second}
	for name, url := range m.processors {
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("Health check failed for %s: %v", name, err)
			continue
		}

		var healthResponse HealthResponse
		if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
			log.Printf("Error decoding health for %s: %v", name, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		if err := saveHealthResponse(name, &healthResponse); err != nil {
			log.Printf("Error saving health for %s: %v", name, err)
		} else {
			log.Printf("Health check for %s OK: %+v", name, healthResponse)
		}
	}
}

func (m *HealthManager) tryToBecomeLeader() (bool, error) {
	ok, err := database.RedisClient.SetNX(database.RedisContext, m.lockKey, m.instanceID, 12*time.Second).Result()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	val, err := database.RedisClient.Get(database.RedisContext, m.lockKey).Result()
	if err != nil {
		return false, err
	}
	if val == m.instanceID {
		if err := database.RedisClient.Expire(database.RedisContext, m.lockKey, 12*time.Second).Err(); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func saveHealthResponse(name string, health *HealthResponse) error {
	value, err := json.Marshal(health)
	if err != nil {
		return err
	}
	return database.RedisClient.Set(database.RedisContext, getKey(name), value, 0).Err()
}

func CheckHealth(name string) (*HealthResponse, error) {
	raw, err := database.RedisClient.Get(database.RedisContext, getKey(name)).Result()
	if err != nil {
		log.Printf("Error retrieving health for %s: %v", name, err)
		return nil, err
	}

	var health HealthResponse
	if err := json.Unmarshal([]byte(raw), &health); err != nil {
		log.Printf("Error parsing health for %s: %v", name, err)
		return nil, err
	}
	return &health, nil
}

func SetAsFailing(name string) error {
	current, err := CheckHealth(name)
	isAlreadyFailing := err == nil && current.Failing
	if isAlreadyFailing {
		return nil
	}
	log.Printf("Setting %s as failing", name)
	return saveHealthResponse(name, &HealthResponse{Failing: true})
}

func getKey(name string) string {
	return "health:" + name
}
