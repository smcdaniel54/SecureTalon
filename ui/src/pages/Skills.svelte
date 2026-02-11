<script lang="ts">
  import { onMount } from 'svelte'
  import { authStore, toastStore } from '../app/store'
  import { listSkills, registerSkill } from '../lib/api'
  import type { Skill, RegisterSkillPayload } from '../lib/types'
  import PageHeader from '../components/PageHeader.svelte'
  import Card from '../components/Card.svelte'

  let skills: Skill[] = []
  let loading = true
  let showRegister = false
  let regName = ''
  let regVersion = '1.0.0'
  let regImage = ''
  let regSignature = ''
  let regKeyId = ''
  let regManifest = '{}'
  let registering = false
  let regErr = ''

  function load() {
    const auth = $authStore
    if (!auth) return
    loading = true
    listSkills(auth)
      .then((r) => { skills = r })
      .catch((e: Error & { details?: unknown }) => {
        toastStore.push('error', e.message, e.details)
      })
      .finally(() => { loading = false })
  }

  function parseManifest(): Record<string, unknown> {
    try {
      const s = regManifest.trim()
      return s ? JSON.parse(s) : {}
    } catch {
      return {}
    }
  }

  async function register() {
    const auth = $authStore
    if (!auth) return
    regErr = ''
    if (!regName.trim()) {
      regErr = 'Name is required.'
      return
    }
    if (!regImage.trim()) {
      regErr = 'Image (repo/name@sha256:...) is required.'
      return
    }
    if (!regImage.includes('@sha256:')) {
      regErr = 'Image must use digest format: repo/name@sha256:...'
      return
    }
    const manifest = parseManifest()
    if (regManifest.trim() && typeof manifest !== 'object') {
      regErr = 'Manifest must be valid JSON object.'
      return
    }
    registering = true
    try {
      const payload: RegisterSkillPayload = {
        name: regName.trim(),
        version: regVersion.trim() || '1.0.0',
        image: regImage.trim(),
        manifest: Object.keys(manifest).length ? manifest : undefined,
      }
      if (regSignature.trim()) payload.signature = regSignature.trim()
      if (regKeyId.trim()) payload.public_key_id = regKeyId.trim()
      await registerSkill(auth, payload)
      showRegister = false
      regName = regVersion = '1.0.0'
      regImage = regSignature = regKeyId = ''
      regManifest = '{}'
      toastStore.push('success', 'Skill registered')
      load()
    } catch (e: unknown) {
      const err = e as Error & { details?: unknown }
      regErr = err.message ?? 'Register failed'
      toastStore.push('error', regErr, err.details)
    } finally {
      registering = false
    }
  }

  onMount(load)
</script>

<div class="page">
  <PageHeader title="Skills" subtitle="Registered agent skills (digest-only images)." />
  {#if !$authStore}
    <p class="muted">Not connected. <a href="#/login">Login</a></p>
  {:else}
    <button class="primary" on:click={() => showRegister = true}>Register skill</button>
    {#if loading}
      <p class="muted">Loading…</p>
    {:else}
      {#if skills.length === 0}
        <p class="muted">No skills registered. Register a skill (image must use digest, e.g. <code>repo/name@sha256:...</code>).</p>
      {:else}
        <table>
          <thead><tr><th>Name</th><th>Version</th><th>Image</th><th>Signed</th></tr></thead>
          <tbody>
            {#each skills as s}
              <tr>
                <td>{s.name}</td>
                <td>{s.version}</td>
                <td><code>{s.image}</code></td>
                <td>{#if s.signed}<span class="badge signed">✓ Signed</span>{:else}<span class="badge">—</span>{/if}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    {/if}
    {#if showRegister}
      <Card title="Register skill">
        <p class="hint">Only images by digest (<code>@sha256:...</code>) are allowed to run. Signing is recommended.</p>
        <label for="reg-name">Name</label>
        <input id="reg-name" type="text" bind:value={regName} placeholder="hello-world" />
        <label for="reg-version">Version</label>
        <input id="reg-version" type="text" bind:value={regVersion} placeholder="1.0.0" />
        <label for="reg-image">Image (digest required)</label>
        <input id="reg-image" type="text" bind:value={regImage} placeholder="registry/hello-world@sha256:..." />
        <label for="reg-sig">Signature (optional)</label>
        <input id="reg-sig" type="text" bind:value={regSignature} placeholder="base64..." />
        <label for="reg-keyid">Public key ID (optional)</label>
        <input id="reg-keyid" type="text" bind:value={regKeyId} placeholder="key_001" />
        <label for="reg-manifest">Manifest JSON (optional)</label>
        <textarea id="reg-manifest" bind:value={regManifest} rows="4" placeholder={'{"key": "value"}'}></textarea>
        {#if regErr}<p class="error">{regErr}</p>{/if}
        <div class="actions">
          <button class="primary" on:click={register} disabled={registering}>{registering ? 'Registering…' : 'Register'}</button>
          <button class="secondary" on:click={() => { showRegister = false; regErr = '' }}>Cancel</button>
        </div>
      </Card>
    {/if}
  {/if}
</div>

<style>
  .page { max-width: 720px; }
  .hint { font-size: 0.85rem; color: var(--text-muted); margin-bottom: var(--space-3); }
  .page label { display: block; margin-top: var(--space-2); }
  .page label:first-of-type { margin-top: 0; }
  .page input, .page textarea { width: 100%; margin-top: var(--space-1); box-sizing: border-box; }
  .page textarea { font-family: ui-monospace, monospace; font-size: 0.9rem; }
  .actions { margin-top: var(--space-4); display: flex; gap: var(--space-2); }
  .error { color: var(--error); font-size: 0.875rem; }
  .muted { color: var(--text-muted); }
  .badge { font-size: 0.85rem; color: var(--text-muted); }
  .badge.signed { color: var(--success); font-weight: 500; }
  .page table { margin-top: var(--space-4); }
</style>
