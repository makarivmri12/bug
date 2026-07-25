package heuristic

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// FlowNode represents a node in the application flow graph
type FlowNode struct {
	ID    string
	URL   string
	Method string
	Params []string
}

// FlowEdge represents a connection between nodes
type FlowEdge struct {
	From string
	To   string
	Relation string
}

// FlowGraph represents the complete application flow
type FlowGraph struct {
	Nodes map[string]*FlowNode
	Edges []*FlowEdge
	mu    sync.RWMutex
}

// FlowMapper maps application flows and dependencies
type FlowMapper struct {
	graph  *FlowGraph
	mu     sync.RWMutex
	logger *zap.Logger
}

// NewFlowMapper creates a new flow mapper
func NewFlowMapper(logger *zap.Logger) *FlowMapper {
	return &FlowMapper{
		graph: &FlowGraph{
			Nodes: make(map[string]*FlowNode),
			Edges: make([]*FlowEdge, 0),
		},
		logger: logger,
	}
}

// AddNode adds a node to the flow graph
func (fm *FlowMapper) AddNode(ctx context.Context, id, url, method string, params []string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.graph.mu.Lock()
	defer fm.graph.mu.Unlock()

	if _, ok := fm.graph.Nodes[id]; ok {
		return fmt.Errorf("node already exists: %s", id)
	}

	fm.graph.Nodes[id] = &FlowNode{
		ID:     id,
		URL:    url,
		Method: method,
		Params: params,
	}

	fm.logger.Info("flow node added",
		zap.String("node_id", id),
		zap.String("url", url),
		zap.String("method", method),
	)

	return nil
}

// AddEdge adds an edge connecting two nodes
func (fm *FlowMapper) AddEdge(ctx context.Context, from, to, relation string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.graph.mu.Lock()
	defer fm.graph.mu.Unlock()

	// Verify both nodes exist
	if _, ok := fm.graph.Nodes[from]; !ok {
		return fmt.Errorf("source node not found: %s", from)
	}
	if _, ok := fm.graph.Nodes[to]; !ok {
		return fmt.Errorf("destination node not found: %s", to)
	}

	fm.graph.Edges = append(fm.graph.Edges, &FlowEdge{
		From:     from,
		To:       to,
		Relation: relation,
	})

	fm.logger.Info("flow edge added",
		zap.String("from", from),
		zap.String("to", to),
		zap.String("relation", relation),
	)

	return nil
}

// GetGraph returns the complete flow graph
func (fm *FlowMapper) GetGraph(ctx context.Context) *FlowGraph {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	fm.graph.mu.RLock()
	defer fm.graph.mu.RUnlock()

	return fm.graph
}
