# chain-health

Lightweight multi-chain node health monitoring daemon. Tracks block height, peer count, sync status, and RPC latency across multiple networks. Fires alerts via Slack, PagerDuty, or webhook when nodes fall behind or become unreachable.

## Features

- **Block height tracking** — compares against public explorers to detect sync lag
- **RPC health checks** — configurable method calls with latency thresholds
- **Peer monitoring** — alerts on peer count drops
- **Multi-chain** — single binary monitors Ethereum, Polygon, Arbitrum, Optimism, BSC, Avalanche, Solana
- **Alert channels** — Slack, PagerDuty, OpsGenie, generic webhook
- **Prometheus metrics** — exposes `/metrics` endpoint for Grafana dashboards

## Quick Start

```bash
go install github.com/arctoslabs-io/chain-health@latest

chain-health --config config.yaml
```

## Configuration

```yaml
chains:
  - name: ethereum
    rpc: ${ETH_RPC_URL}
    type: evm
    checks:
      block_lag: 5
      min_peers: 10
      rpc_timeout: 3s
      methods:
        - eth_blockNumber
        - eth_syncing
        - net_peerCount

  - name: polygon
    rpc: ${POLYGON_RPC_URL}
    type: evm
    checks:
      block_lag: 20
      min_peers: 5
      rpc_timeout: 5s

  - name: solana
    rpc: ${SOLANA_RPC_URL}
    type: solana
    checks:
      slot_lag: 50
      rpc_timeout: 5s
      methods:
        - getHealth
        - getSlot

alerts:
  slack:
    webhook: ${SLACK_WEBHOOK}
    channel: "#infra-alerts"
  pagerduty:
    routing_key: ${PD_KEY}
    severity: critical

monitoring:
  interval: 30s
  prometheus:
    enabled: true
    port: 9090
```

## Architecture

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Chain Poller │────▶│  Health Eval  │────▶│  Alert Sink  │
│  (per chain)  │     │  (thresholds) │     │  (fanout)    │
└──────────────┘     └──────────────┘     └──────────────┘
       │                                         │
       ▼                                         ▼
┌──────────────┐                          ┌──────────────┐
│  Prometheus  │                          │ Slack / PD / │
│  /metrics    │                          │ Webhook      │
└──────────────┘                          └──────────────┘
```

## Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `chain_health_block_height` | Gauge | Current block height per chain |
| `chain_health_block_lag` | Gauge | Blocks behind public explorer |
| `chain_health_rpc_latency_ms` | Histogram | RPC response time |
| `chain_health_peer_count` | Gauge | Connected peer count |
| `chain_health_checks_total` | Counter | Total health checks run |
| `chain_health_alerts_total` | Counter | Total alerts fired |

## Deployment

### Docker

```bash
docker run -d \
  -v $(pwd)/config.yaml:/etc/chain-health/config.yaml \
  -p 9090:9090 \
  ghcr.io/arctoslabs-io/chain-health:latest
```

### Kubernetes

```bash
helm install chain-health ./charts/chain-health \
  --set config.eth_rpc=$ETH_RPC_URL \
  --namespace monitoring
```

## Development

```bash
git clone https://github.com/arctoslabs-io/chain-health.git
cd chain-health
make build
make test
make lint
```

## License

MIT
