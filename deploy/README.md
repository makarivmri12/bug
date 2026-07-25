# HLFA Deployment Guide

## Docker Compose (Development)

```bash
# Start development environment
make docker-up

# Stop environment
make docker-down

# View logs
docker-compose logs -f
```

## Docker Compose (Production)

```bash
# Set environment variables
export DB_USER=hlfa
export DB_PASSWORD=$(openssl rand -base64 32)
export REDIS_PASSWORD=$(openssl rand -base64 32)
export NEO4J_PASSWORD=$(openssl rand -base64 32)

# Start production environment
docker-compose -f deploy/docker-compose.prod.yml up -d

# Check status
docker-compose -f deploy/docker-compose.prod.yml ps
```

## Kubernetes Deployment

### Prerequisites
- Kubernetes cluster 1.20+
- kubectl configured
- Helm 3+

### Deploy with Helm

```bash
# Create namespace
kubectl create namespace hlfa

# Deploy HLFA
helm install hlfa ./deploy/helm \
  -n hlfa \
  -f deploy/helm/values-prod.yaml

# Check deployment
kubectl get pods -n hlfa
kubectl logs -n hlfa deployment/hlfa-backend
```

### Upgrade

```bash
helm upgrade hlfa ./deploy/helm \
  -n hlfa \
  -f deploy/helm/values-prod.yaml
```

### Rollback

```bash
helm rollback hlfa -n hlfa
```

## PostgreSQL Backup

```bash
# Backup database
docker exec hlfa-postgres-1 pg_dump -U hlfa hlfa > backup.sql

# Restore database
docker exec -i hlfa-postgres-1 psql -U hlfa hlfa < backup.sql
```

## Monitoring

### Health Checks

```bash
# Backend
curl http://localhost:8080/health

# Python Engine
curl http://localhost:5000/health

# Frontend
curl http://localhost:3000
```

## Troubleshooting

### Backend Connection Issues

```bash
# Check database connection
docker-compose exec backend psql -h postgres -U hlfa -d hlfa -c "SELECT 1"

# Check Redis connection
docker-compose exec backend redis-cli -h redis ping

# Check Neo4j connection
docker-compose exec backend cypher-shell -a neo4j -p neo4j "RETURN 1"
```

### Python Engine Issues

```bash
# Check logs
docker-compose logs python-engine

# Test API
curl http://localhost:5000/health
```

## Production Checklist

- [ ] Set strong passwords for all services
- [ ] Configure SSL/TLS certificates
- [ ] Set up proper backup strategy
- [ ] Configure monitoring and alerting
- [ ] Set up log aggregation
- [ ] Configure resource limits
- [ ] Set up auto-scaling policies
- [ ] Configure database replication
- [ ] Set up disaster recovery plan
- [ ] Document runbooks
