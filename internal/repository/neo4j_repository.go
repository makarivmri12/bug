package repository

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

// Neo4jRepository handles graph database operations
type Neo4jRepository struct {
	driver neo4j.Driver
	logger *zap.Logger
}

// NewNeo4jRepository creates a new Neo4j repository
func NewNeo4jRepository(driver neo4j.Driver, logger *zap.Logger) *Neo4jRepository {
	return &Neo4jRepository{
		driver: driver,
		logger: logger,
	}
}

// CreateEndpointNode creates an endpoint node in the graph
func (nr *Neo4jRepository) CreateEndpointNode(ctx context.Context, url, method string, params []string) error {
	session := nr.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
		CREATE (e:Endpoint {url: $url, method: $method, params: $params, created_at: timestamp()})
		RETURN e.url
	`

	_, err := session.Run(ctx, query, map[string]interface{}{
		"url":    url,
		"method": method,
		"params": params,
	})

	if err != nil {
		nr.logger.Error("failed to create endpoint node", zap.Error(err))
		return err
	}

	nr.logger.Info("endpoint node created", zap.String("url", url), zap.String("method", method))
	return nil
}

// CreateSessionNode creates a session node in the graph
func (nr *Neo4jRepository) CreateSessionNode(ctx context.Context, sessionID, userID string) error {
	session := nr.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
		CREATE (s:Session {session_id: $sessionID, user_id: $userID, created_at: timestamp()})
		RETURN s.session_id
	`

	_, err := session.Run(ctx, query, map[string]interface{}{
		"sessionID": sessionID,
		"userID":   userID,
	})

	if err != nil {
		nr.logger.Error("failed to create session node", zap.Error(err))
		return err
	}

	return nil
}

// CreateRelationship creates a relationship between two nodes
func (nr *Neo4jRepository) CreateRelationship(ctx context.Context, fromType, fromID, relType, toType, toID string) error {
	session := nr.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := fmt.Sprintf(`
		MATCH (a:%s {id: $fromID})
		MATCH (b:%s {id: $toID})
		CREATE (a)-[r:%s]->(b)
		RETURN r
	`, fromType, toType, relType)

	_, err := session.Run(ctx, query, map[string]interface{}{
		"fromID": fromID,
		"toID":   toID,
	})

	if err != nil {
		nr.logger.Error("failed to create relationship", zap.Error(err))
		return err
	}

	return nil
}

// GetEndpointNodes retrieves all endpoint nodes
func (nr *Neo4jRepository) GetEndpointNodes(ctx context.Context) ([]map[string]interface{}, error) {
	session := nr.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := "MATCH (e:Endpoint) RETURN e.url as url, e.method as method, e.params as params"

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		nr.logger.Error("failed to get endpoint nodes", zap.Error(err))
		return nil, err
	}

	var endpoints []map[string]interface{}
	for result.Next(ctx) {
		endpoints = append(endpoints, result.Record().AsMap())
	}

	return endpoints, result.Err()
}

// GetFlowPaths retrieves flow paths between two endpoints
func (nr *Neo4jRepository) GetFlowPaths(ctx context.Context, fromURL, toURL string) ([]map[string]interface{}, error) {
	session := nr.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `
		MATCH p = (from:Endpoint {url: $fromURL})-[*]->(to:Endpoint {url: $toURL})
		RETURN p
	`

	result, err := session.Run(ctx, query, map[string]interface{}{
		"fromURL": fromURL,
		"toURL":   toURL,
	})

	if err != nil {
		nr.logger.Error("failed to get flow paths", zap.Error(err))
		return nil, err
	}

	var paths []map[string]interface{}
	for result.Next(ctx) {
		paths = append(paths, result.Record().AsMap())
	}

	return paths, result.Err()
}
