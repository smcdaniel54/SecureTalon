<script lang="ts">
  import Router from 'svelte-spa-router'
  import AppShell from './components/AppShell.svelte'
  import Login from './pages/Login.svelte'
  import Dashboard from './pages/Dashboard.svelte'
  import Sessions from './pages/Sessions.svelte'
  import SessionDetail from './pages/SessionDetail.svelte'
  import RunDetail from './pages/RunDetail.svelte'
  import Audit from './pages/Audit.svelte'
  import Policies from './pages/Policies.svelte'
  import Skills from './pages/Skills.svelte'
  import Replay from './pages/Replay.svelte'
  import NotFound from './pages/NotFound.svelte'

  const routes = {
    '/': Dashboard,
    '/login': Login,
    '/sessions': Sessions,
    '/sessions/:id': SessionDetail,
    '/runs/:id': RunDetail,
    '/policies': Policies,
    '/skills': Skills,
    '/replay': Replay,
    '/audit': Audit,
    '*': NotFound,
  }
</script>

<AppShell>
  <Router {routes} />
</AppShell>
