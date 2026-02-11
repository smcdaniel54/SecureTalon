import type {
  ApiError,
  Session,
  EffectivePolicy,
  RuleOverride,
  Skill,
  RegisterSkillPayload,
  AuditQueryParams,
  AuditEvent,
  AuditValidateResult,
  ReplayResponse,
} from './types'

const AUTH_KEY = 'securetalon_auth'

export interface Auth {
  apiBase: string
  token: string
}

export interface ApiErrorResult {
  message: string
  code?: string
  details?: Record<string, unknown>
}

export function getStoredAuth(): Auth | null {
  try {
    const raw = localStorage.getItem(AUTH_KEY)
    if (!raw) return null
    const o = JSON.parse(raw) as Auth
    if (!o.apiBase || !o.token) return null
    return o
  } catch {
    return null
  }
}

export function setStoredAuth(auth: Auth): void {
  localStorage.setItem(AUTH_KEY, JSON.stringify(auth))
}

export function clearStoredAuth(): void {
  localStorage.removeItem(AUTH_KEY)
}

async function request<T>(auth: Auth, path: string, options: RequestInit = {}): Promise<T> {
  const url = `${auth.apiBase.replace(/\/$/, '')}${path}`
  const headers: Record<string, string> = {
    Authorization: `Bearer ${auth.token}`,
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  const res = await fetch(url, { ...options, headers })
  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as ApiError
    const err: ApiErrorResult = {
      message: body?.error?.message ?? res.statusText,
      code: body?.error?.code,
      details: body?.error?.details,
    }
    const e = new Error(err.message) as Error & ApiErrorResult
    e.code = err.code
    e.details = err.details
    throw e
  }
  return res.json() as Promise<T>
}

export async function connect(auth: Auth): Promise<void> {
  await request<{ sessions: unknown[] }>(auth, '/v1/sessions?limit=1')
}

// Sessions
export async function createSession(auth: Auth, label: string, metadata?: Record<string, string>): Promise<Session> {
  return request<Session>(auth, '/v1/sessions', {
    method: 'POST',
    body: JSON.stringify({ label, metadata: metadata ?? {} }),
  })
}

export async function listSessions(auth: Auth, limit = 50): Promise<Session[]> {
  const r = await request<{ sessions: Session[] }>(auth, `/v1/sessions?limit=${limit}`)
  return r.sessions ?? []
}

export async function getSession(auth: Auth, sessionId: string): Promise<Session> {
  return request<Session>(auth, `/v1/sessions/${encodeURIComponent(sessionId)}`)
}

export async function listMessages(auth: Auth, sessionId: string, limit = 200) {
  return request<{ messages: { id: string; role: string; content: string; timestamp: string; run_id?: string }[] }>(
    auth,
    `/v1/sessions/${encodeURIComponent(sessionId)}/messages?limit=${limit}`
  )
}

export async function postMessage(
  auth: Auth,
  sessionId: string,
  body: { role?: string; content: string; metadata?: Record<string, string>; intents?: { tool: string; params: Record<string, unknown> }[] }
) {
  return request<{ run_id: string; status: string }>(auth, `/v1/sessions/${encodeURIComponent(sessionId)}/messages`, {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export async function getRun(auth: Auth, runId: string) {
  return request<{ id: string; session_id: string; status: string; started_at: string; ended_at?: string; steps?: { step_id: string; type: string; status: string; tool?: string; details?: Record<string, unknown> }[] }>(
    auth,
    `/v1/runs/${encodeURIComponent(runId)}`
  )
}

// Policy
export async function getEffectivePolicy(auth: Auth, sessionId: string): Promise<EffectivePolicy> {
  return request<EffectivePolicy>(auth, `/v1/policy/effective?session_id=${encodeURIComponent(sessionId)}`)
}

export async function putSessionPolicy(auth: Auth, sessionId: string, overrides: RuleOverride[]): Promise<void> {
  await request(auth, `/v1/sessions/${encodeURIComponent(sessionId)}/policy`, {
    method: 'PUT',
    body: JSON.stringify({ overrides }),
  })
}

// Skills
export async function listSkills(auth: Auth): Promise<Skill[]> {
  const r = await request<{ skills: Skill[] }>(auth, '/v1/skills')
  return r.skills ?? []
}

export async function registerSkill(auth: Auth, payload: RegisterSkillPayload): Promise<void> {
  await request(auth, '/v1/skills', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

// Audit
export async function queryAudit(auth: Auth, params: AuditQueryParams): Promise<AuditEvent[]> {
  const p = new URLSearchParams()
  p.set('limit', String(params.limit ?? 500))
  if (params.session_id) p.set('session_id', params.session_id)
  if (params.run_id) p.set('run_id', params.run_id)
  if (params.type) p.set('type', params.type)
  if (params.since) p.set('since', params.since)
  if (params.until) p.set('until', params.until)
  const r = await request<{ events: AuditEvent[] }>(auth, `/v1/audit?${p}`)
  return r.events ?? []
}

export async function validateAuditChain(auth: Auth, params: { session_id?: string; limit?: number }): Promise<AuditValidateResult> {
  const p = new URLSearchParams()
  p.set('limit', String(params.limit ?? 500))
  if (params.session_id) p.set('session_id', params.session_id)
  return request<AuditValidateResult>(auth, `/v1/audit/validate?${p}`)
}

// Replay (safe only)
export async function safeReplay(auth: Auth, runId: string): Promise<ReplayResponse> {
  return request<ReplayResponse>(auth, `/v1/runs/${encodeURIComponent(runId)}/replay`, {
    method: 'POST',
  })
}
