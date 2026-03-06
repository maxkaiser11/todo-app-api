<script lang="ts">
  import Auth from './Auth.svelte'
  import TodoList from './TodoList.svelte'

  let token: string | null = localStorage.getItem('token')

  function handleLogin(t: string) {
    token = t
    localStorage.setItem('token', t)
  }

  function handleLogout() {
    token = null
    localStorage.removeItem('token')
  }
</script>

{#if token}
  <TodoList {token} onLogout={handleLogout} />
{:else}
  <Auth onLogin={handleLogin} />
{/if}
