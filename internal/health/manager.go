// TODO improve this by rotating the leader
package health

import (
	"encoding/json"
	"log"
	"net/http"
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
	lockKey    string
	instanceID string
	processors map[string]string
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

	isLeader, err := m.tryToBecomeLeader()
	if err != nil {
		log.Printf("Error trying to become leader: %v", err)
		return
	}

	if !isLeader {
		log.Println("Not the leader, exiting health check manager.")
		return
	}

	log.Println("I'm the leader, starting health checks.")

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.runHealthCheck()
			case <-database.RedisContext.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func CheckHealth(name string) (*HealthResponse, error) {
	health, err := database.RedisClient.Get(database.RedisContext, getKey(name)).Result()
	if err != nil {
		log.Printf("Error retrieving health check for %s: %v", name, err)
		return nil, err
	}

	var healthResponse HealthResponse
	if err := json.Unmarshal([]byte(health), &healthResponse); err != nil {
		log.Printf("Error unmarshalling health response for %s: %v", name, err)
		return nil, err
	}

	return &healthResponse, nil
}

func SetAsFailing(name string) error {
	health := &HealthResponse{Failing: true, MinResponseTime: 0}
	value, err := json.Marshal(health)
	if err != nil {
		log.Printf("Error marshalling updated health status for %s: %v", name, err)
		return err
	}

	if err := database.RedisClient.Set(database.RedisContext, getKey(name), value, 0).Err(); err != nil {
		log.Printf("Error setting updated health status for %s: %v", name, err)
		return err
	}

	log.Printf("Set failing status for %s to true", name)
	return nil
}

func (m *HealthManager) runHealthCheck() {
	client := &http.Client{Timeout: 4 * time.Second}
	for name, url := range m.processors {
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("Error checking health for %s: %v", name, err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("Health check for %s failed with status code: %d", name, resp.StatusCode)
			continue
		}
		var healthResponse HealthResponse
		if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
			log.Printf("Error decoding health response for %s: %v", name, err)
			continue
		}
		value, err := json.Marshal(healthResponse)
		if err != nil {
			log.Printf("Error marshalling health response for %s: %v", name, err)
			continue
		}
		if err := database.RedisClient.Set(database.RedisContext, getKey(name), value, 0).Err(); err != nil {
			log.Printf("Error saving health check result for %s: %v", name, err)
			continue
		}
		log.Printf("Health check for %s successful: %+v", name, healthResponse)
	}
}

func (m *HealthManager) tryToBecomeLeader() (bool, error) {
	ok, err := database.RedisClient.SetNX(database.RedisContext, m.lockKey, m.instanceID, 0).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func getKey(name string) string {
	return "health:" + name
}
