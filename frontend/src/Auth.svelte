<script lang="ts">
  import { fade, fly } from 'svelte/transition'

  const API = import.meta.env.VITE_API_URL

  export let onLogin: (token: string) => void

  let mode: 'login' | 'register' = 'login'
  let email = ''
  let password = ''
  let loading = false
  let error = ''

  async function submit() {
    error = ''
    if (!email || !password) {
      error = 'email and password required'
      return
    }
    loading = true
    try {
      const res = await fetch(`${API}/auth/${mode}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      })
      const data = await res.json()
      if (!res.ok) {
        error = data.error ?? 'something went wrong'
        return
      }
      onLogin(data.token)
    } catch {
      error = 'could not reach server'
    } finally {
      loading = false
    }
  }

  function switchMode() {
    mode = mode === 'login' ? 'register' : 'login'
    error = ''
  }
</script>

<svelte:head>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;500;600&family=Syne:wght@400;700;800&display=swap" rel="stylesheet" />
</svelte:head>

<div class="auth" in:fade={{ duration: 200 }}>
  <div class="card" in:fly={{ y: 16, duration: 300 }}>
    <div class="label">// task runner</div>
    <h1>to<span>do</span>s</h1>
    <p class="sub">{mode === 'login' ? 'sign in to continue' : 'create an account'}</p>

    <div class="fields">
      <div class="field">
        <div class="field-label">email</div>
        <input
          type="email"
          bind:value={email}
          on:keydown={e => e.key === 'Enter' && submit()}
          placeholder="you@example.com"
          autocomplete="email"
          disabled={loading}
        />
      </div>

      <div class="field">
        <div class="field-label">password</div>
        <input
          type="password"
          bind:value={password}
          on:keydown={e => e.key === 'Enter' && submit()}
          placeholder="••••••••"
          autocomplete={mode === 'login' ? 'current-password' : 'new-password'}
          disabled={loading}
        />
      </div>
    </div>

    {#if error}
      <div class="error" in:fly={{ y: -4, duration: 150 }}>{error}</div>
    {/if}

    <button class="submit-btn" on:click={submit} disabled={loading}>
      {#if loading}
        <span class="spinner" />
      {:else}
        {mode === 'login' ? 'sign in' : 'create account'}
      {/if}
    </button>

    <button class="switch-btn" on:click={switchMode}>
      {mode === 'login' ? "don't have an account? register" : 'already have an account? sign in'}
    </button>
  </div>
</div>

<style>
  .auth {
    --surface: #141417;
    --border:  #22222a;
    --muted:   #3a3a48;
    --text:    #e2e2ec;
    --dim:     #888899;
    --accent:  #7c6af7;
    --danger:  #f56565;

    font-family: 'JetBrains Mono', monospace;
    min-height: 100vh;
    background: #0d0d0f;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
  }

  .card {
    width: 100%;
    max-width: 360px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 40px 36px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .label {
    font-size: 11px;
    letter-spacing: 0.2em;
    text-transform: uppercase;
    color: var(--accent);
    margin-bottom: 6px;
    text-align: center;
    width: 100%;
  }

  h1 {
    font-family: 'Syne', sans-serif;
    font-size: 42px;
    font-weight: 800;
    line-height: 1;
    letter-spacing: -0.02em;
    color: var(--text);
    text-align: center;
    width: 100%;
  }
  h1 span { color: var(--accent); }

  .sub {
    margin-top: 8px;
    font-size: 12px;
    color: var(--dim);
    margin-bottom: 32px;
    text-align: center;
    width: 100%;
  }

  .fields {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 20px;
    width: 100%;
  }

  .field {
    width: 100%;
  }

  .field-label {
    font-size: 10px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--dim);
    margin-bottom: 6px;
    text-align: left;
  }

  input {
    width: 100%;
    background: #0d0d0f;
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 12px 14px;
    font-family: 'JetBrains Mono', monospace;
    font-size: 13px;
    color: var(--text);
    caret-color: var(--accent);
    outline: none;
    transition: border-color 0.2s, box-shadow 0.2s;
    box-sizing: border-box;
  }
  input:focus {
    border-color: var(--accent);
    box-shadow: 0 0 0 3px rgba(124,106,247,0.12);
  }
  input::placeholder { color: var(--muted); }
  input:disabled { opacity: 0.5; cursor: not-allowed; }

  .error {
    font-size: 11px;
    color: var(--danger);
    margin-bottom: 12px;
    padding: 8px 12px;
    background: rgba(245,101,101,0.08);
    border: 1px solid rgba(245,101,101,0.2);
    border-radius: 4px;
    width: 100%;
  }

  .submit-btn {
    width: 100%;
    background: var(--accent);
    border: none;
    border-radius: 6px;
    color: #fff;
    font-family: 'JetBrains Mono', monospace;
    font-size: 13px;
    padding: 13px;
    cursor: pointer;
    transition: background 0.15s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    margin-bottom: 16px;
  }
  .submit-btn:hover:not(:disabled) { background: #6a5ae0; }
  .submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }

  .spinner {
    width: 14px;
    height: 14px;
    border: 1.5px solid rgba(255,255,255,0.3);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .switch-btn {
    background: none;
    border: none;
    color: var(--dim);
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px;
    cursor: pointer;
    width: 100%;
    text-align: center;
    transition: color 0.15s;
    padding: 4px;
  }
  .switch-btn:hover { color: var(--accent); }
</style>
