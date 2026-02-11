export interface Session {
  id: string
  created_at: string
  label: string
  status: string
  metadata?: Record<string, string>
}

export interface Message {
  id: string
  role: string
  content: string
  timestamp: string
  metadata?: Record<string, string>
  run_id?: string
}

export interface Run {
  id: string
  session_id: string
  status: string
  started_at: string
  ended_at?: string
  steps?: Step[]
}

export interface Step {
  step_id: string
  type: string
  status: string
  tool?: string
  details?: Record<string, unknown>
}

export interface ToolIntent {
  tool: string
  params: Record<string, unknown>
  subject?: string
}

export interface AuditEvent {
  event_id: string
  ts: string
  session_id: string
  run_id?: string
  type: string
  data: Record<string, unknown>
  prev_hash: string
  hash: string
}

export interface ApiError {
  error: { code: string; message: string; details?: Record<string, unknown> }
}

// Policy
export interface RuleOverride {
  tool: string
  allow: boolean
  constraints: Record<string, unknown>
}

export interface EffectivePolicy {
  default: string
  overrides: RuleOverride[]
}

// Skills
export interface Skill {
  name: string
  version: string
  image: string
  signed?: boolean
  manifest?: Record<string, unknown>
}

export interface RegisterSkillPayload {
  name: string
  version: string
  image: string
  signature?: string
  public_key_id?: string
  manifest?: Record<string, unknown>
}

// Audit
export interface AuditQueryParams {
  session_id?: string
  run_id?: string
  type?: string
  since?: string
  until?: string
  limit?: number
}

export interface AuditValidateResult {
  valid: boolean
  invalid_index: number
  event_count: number
}

// Replay
export interface ReplayEvent {
  ts: string
  type: string
  data: Record<string, unknown>
  hash: string
  prev_hash: string
}

export interface ReplayResponse {
  run_id: string
  mode: string
  valid: boolean
  events: ReplayEvent[]
}
