# Ambassador Protocol (v2)

**Role**: Multi-gastown coordination with capability and policy-based work routing
**Applies to**: ALL agents working on lirt

## Purpose

The Ambassador ensures that:
1. **Work is routed** to gastowns capable of performing it (capability matching)
2. **IP is protected** through domain and client isolation policies
3. **Beads are qualified** with policy metadata at creation time
4. **Approval workflows** gate sensitive or out-of-scope work

This enables multiple gastowns (and multiple users) to collaborate on the same project with proper work routing and intellectual property isolation.

## Key Concepts

### Town Identity (v2)

Each gastown has:
- **UUID**: Globally unique identifier (e.g., `550e8400-e29b-41d4-a716-446655440000`)
- **Alias**: Human-readable name (`hostname.overseer`, e.g., `mbp-m5.james`)

The UUID is used for bead assignment (collision-proof). The alias is used for display and logs.

```bash
lirt-ambassador id      # Returns UUID
lirt-ambassador alias   # Returns hostname.overseer
lirt-ambassador whoami  # Full identity display
```

### Capabilities

Gastowns declare what they can do:

| Type | Examples | How Set |
|------|----------|---------|
| OS | darwin, linux, windows | Auto-detected |
| Arch | arm64, amd64 | Auto-detected |
| Tools | go, node, docker, kubectl | Auto-detected |
| Custom | gpu-access, vpn-access, admin | Manual |

### Capability Matching

Beads can specify both **required** and **excluded** capabilities:

| Label | Meaning |
|-------|---------|
| `requires:darwin` | Gastown MUST have darwin |
| `requires:gpu` | Gastown MUST have gpu |
| `excludes:windows` | Gastown MUST NOT have windows |
| `excludes:gpu` | Gastown MUST NOT have gpu |

Example: Deterministic tests might use `excludes:gpu` to ensure reproducibility.

### Policy System

Gastowns can configure policies for domain and IP isolation:

```bash
# Domain policies
lirt-ambassador policy set-domain allowed general-programming
lirt-ambassador policy set-domain restricted competitor-code
lirt-ambassador policy set-domain requires_approval external-api-access

# IP isolation
lirt-ambassador policy set-ip-client acme-corp --exclusive
```

**Domain Categories:**
- `allowed`: Work in these domains proceeds automatically
- `restricted`: Work in these domains is rejected (routed elsewhere)
- `requires_approval`: Work requires overseer approval before proceeding

**IP Isolation:**
- `client_projects`: List of client identifiers this gastown can work on
- `exclusive`: If true, gastown ONLY works on listed clients (and rejects others)

## Bead Creation Enforcement

**CRITICAL**: All beads must have policy classification. Use `ambassador create` instead of `bd create`.

### Creating Policy-Qualified Beads

```bash
# Explicit classification
lirt-ambassador create \
  --title "Fix authentication bug" \
  --type bug \
  --domain identity-management \
  --ip public

# Auto-classification from title
lirt-ambassador create \
  --title "Fix authentication bug" \
  --auto-classify
```

### Required Labels

Every bead MUST have:
- `domain:<topic>` - What topic area (e.g., `domain:identity-management`)
- `ip:<scope>` - IP context (e.g., `ip:public` or `ip:acme-corp`)

### Classifying Existing Beads

```bash
# Manual classification
lirt-ambassador classify <issue-id> --domain testing --ip public

# Auto-classification
lirt-ambassador classify <issue-id> --auto
```

### Quarantine

Beads without policy classification are quarantined and cannot be claimed:

```bash
# List quarantined beads
lirt-ambassador quarantine

# Only policy-qualified beads appear in ready list
lirt-ambassador ready
```

## Assessment Flow

Before starting ANY task, run assessment:

```bash
lirt-ambassador assess <issue-id>
```

### Assessment Phases

```
1. POLICY CLASSIFICATION CHECK
   └─ Does bead have domain: label?
      ├─ No → QUARANTINED (exit 2)
      └─ Yes → Continue

2. CAPABILITY CHECK (requires:)
   └─ Does gastown have all required capabilities?
      ├─ No → ROUTE (exit 1)
      └─ Yes → Continue

3. CAPABILITY CHECK (excludes:)
   └─ Does gastown have any excluded capabilities?
      ├─ Yes → ROUTE (exit 1)
      └─ No → Continue

4. DOMAIN POLICY CHECK
   └─ Is task domain in gastown's policy?
      ├─ Restricted → REJECT (exit 4)
      ├─ Requires approval (not granted) → HOLD (exit 3)
      └─ Allowed or empty policy → Continue

5. IP ISOLATION CHECK
   └─ Does task IP scope match gastown's client list?
      ├─ Exclusive mode + no match → ROUTE (exit 1)
      └─ Match or non-exclusive → PROCEED (exit 0)
```

### Exit Codes

| Code | Meaning | Action |
|------|---------|--------|
| 0 | CAN HANDLE | Proceed with work |
| 1 | CANNOT HANDLE | Route to other gastown |
| 2 | QUARANTINED | Classify bead first |
| 3 | HOLD | Await overseer approval |
| 4 | REJECTED | Policy violation, cannot run here |

## Routing Work

When this gastown cannot handle a task:

```bash
# Route to any capable gastown
lirt-ambassador route <issue-id> --capability darwin

# Route to specific gastown (by UUID)
lirt-ambassador route <issue-id> --to <uuid>

# Add reason for audit trail
lirt-ambassador route <issue-id> --capability linux \
  --reason "IP isolation: acme-corp"
```

## Approval Workflow

For tasks requiring approval:

```bash
# List pending approvals
lirt-ambassador approvals

# Grant approval
lirt-ambassador approve <issue-id>

# Deny and route elsewhere
lirt-ambassador deny <issue-id> --route

# Deny and ban from this gastown
lirt-ambassador deny <issue-id> --ban
```

## Agent Protocol

### Before Creating Any Bead

1. Determine domain classification from task context
2. Determine IP scope (public or client-specific)
3. Use `lirt-ambassador create` with classification

### Before Starting Any Task

1. Run `lirt-ambassador assess <issue-id>`
2. Check exit code:
   - 0: Proceed with `bd update <id> --claim && bd sync`
   - 1: Route with `lirt-ambassador route <id>`
   - 2: Classify with `lirt-ambassador classify <id>`
   - 3: Wait for approval or escalate to overseer
   - 4: Do not attempt, work is prohibited here

### Bead-Creating Roles

Roles that create beads (chronicler, archivist, etc.) MUST:
- Use `lirt-ambassador create` instead of `bd create`
- Include domain classification based on detected work
- Default to `ip:public` unless task is client-specific

## Example Scenarios

### Scenario 1: Capability Exclusion

```
Task: "Run deterministic ML tests"
Labels: requires:python3, excludes:gpu

Gastown A (has gpu):
  + python3
  ! gpu (PRESENT - violates exclusion)
  → ROUTE

Gastown B (no gpu):
  + python3
  - gpu (absent, ok)
  → CAN HANDLE
```

### Scenario 2: Domain Policy

```
Task: "Integrate Stripe payment API"
Labels: domain:external-api-access, ip:public

Gastown policy:
  allowed: [general-programming, devops]
  requires_approval: [external-api-access]

Assessment:
  Domain requires approval
  → HOLD (exit 3)

Overseer runs: lirt-ambassador approve <id>
  → CAN HANDLE
```

### Scenario 3: IP Isolation

```
Task: "Fix client auth bug"
Labels: domain:identity-management, ip:acme-corp

Gastown A (exclusive, clients: [acme-corp]):
  IP acme-corp in client list
  → CAN HANDLE

Gastown B (exclusive, clients: [initech]):
  IP acme-corp NOT in client list
  → ROUTE (IP isolation)
```

### Scenario 4: Quarantined Bead

```
Task: "Fix bug in login flow"
Labels: (none)

Assessment:
  Missing domain: label
  → QUARANTINED (exit 2)

Action:
  lirt-ambassador classify <id> --domain identity-management
  → Re-assess
```

## Tools Reference

| Command | Purpose |
|---------|---------|
| `lirt-ambassador init` | Initialize gastown with UUID |
| `lirt-ambassador migrate` | Migrate v1 config to v2 |
| `lirt-ambassador id` | Get town UUID |
| `lirt-ambassador alias` | Get town alias |
| `lirt-ambassador whoami` | Full identity display |
| `lirt-ambassador caps list` | List capabilities |
| `lirt-ambassador caps add <cap>` | Add capability |
| `lirt-ambassador policy show` | Show policies |
| `lirt-ambassador policy set-domain` | Configure domain policy |
| `lirt-ambassador policy set-ip-client` | Configure IP isolation |
| `lirt-ambassador create` | Create policy-qualified bead |
| `lirt-ambassador classify <id>` | Add policy classification |
| `lirt-ambassador assess <id>` | Check if can handle |
| `lirt-ambassador route <id>` | Route to capable gastown |
| `lirt-ambassador ready` | List qualified ready work |
| `lirt-ambassador quarantine` | List unclassified beads |
| `lirt-ambassador approvals` | List pending approvals |
| `lirt-ambassador approve <id>` | Grant approval |
| `lirt-ambassador deny <id>` | Deny approval |
| `lirt-ambassador status` | Show gastown status |
| `lirt-ambassador sync-setup` | Set up sync branch |

## Success Metrics

Track ambassador effectiveness:
- Beads created with vs without policy classification
- Tasks routed due to capability vs policy reasons
- Approval requests granted vs denied
- IP isolation violations caught
- Time from task creation to capable gastown claiming
