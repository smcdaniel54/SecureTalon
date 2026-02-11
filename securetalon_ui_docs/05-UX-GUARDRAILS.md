# UX Guardrails for Security

The UI must prevent users from accidentally weakening security.

## Guardrail rules
1. Do not allow blank allowlists:
   - If enabling http.fetch, require at least one domain
   - If enabling file.read/write, require at least one allowed root
2. Auto-suggest least privilege:
   - when user types a URL, extract domain and propose exact domain allowlist
3. Explicit warnings:
   - enabling network for skills
   - enabling shell.exec
   - allowing host mounts
4. Provide “why this is risky” in plain language.
5. Provide a safe alternative:
   - “Use brokered http.fetch instead of enabling container network”
6. Default TTL short:
   - 60s; user can increase but show warning above 5 min.

## Advisor panel (MVP-lite)
A right-side panel on Policies page:
- lists “Issues found”:
  - overly broad domain patterns
  - large max_bytes
  - long TTL
- offers one-click “tighten” suggestions

Phase 2:
- true “Security Advisor agent” consuming audit and recommending changes.

