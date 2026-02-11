# UI Templates (Reusable Components)

The UI should be template-driven so the design improves as templates evolve.

## Layout templates
1. `AppShell`
   - Left nav, top bar, content area
   - Environment badge (DEV/PROD)
2. `PageHeader`
   - Title, subtitle, primary actions
3. `Card`
   - Standard container with header/body/footer
4. `DataTable`
   - sortable columns, pagination, empty state
5. `Timeline`
   - vertical event list, icons, expandable details
6. `FormPanel`
   - label/help text, validation, save/cancel

## UX patterns
- Inline toasts for success/errors
- “Explain this decision” drawer on denies:
  - shows policy rule matched/failed
  - shows recommended least-privilege fix

## Security-first patterns
- **Danger zone** section for any high-risk policy (e.g., shell.exec enable)
- Confirmation modal with explicit consequences
- Default constraints filled automatically:
  - short TTL
  - max bytes
  - network none

## Design tokens (recommended)
- spacing scale: 4/8/12/16/24/32
- radius: 12px
- typography: simple (system fonts OK)
- consistent badges: ALLOW (green), DENY (red), APPROVAL (amber)

