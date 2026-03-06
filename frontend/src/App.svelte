<script lang="ts">
  import { onMount } from 'svelte'
  import { fly, fade } from 'svelte/transition'
  import { flip } from 'svelte/animate'

  type Task = {
    ID: number
    Task: string
    Description: string
    Status: boolean
    CreatedAt: string
  }

  let tasks: Task[] = []
  let loading = true
  let error = ''
  let filter: 'all' | 'active' | 'done' = 'all'

  onMount(async () => {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/tasks`)
      tasks = await res.json()
    } catch (e) {
      error = 'Failed to fetch tasks'
    } finally {
      loading = false
    }
  })

  $: filtered = tasks.filter(t =>
    filter === 'all' ? true : filter === 'active' ? !t.Status : t.Status
  )
  $: remaining = tasks.filter(t => !t.Status).length
  $: doneCount = tasks.filter(t => t.Status).length

  async function toggleStatus(task: Task) {
    const updated = { ...task, Status: !task.Status }
    tasks = tasks.map(t => t.ID === task.ID ? updated : t)
    try {
      await fetch(`${import.meta.env.VITE_API_URL}/task/${task.ID}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updated),
      })
    } catch {
      tasks = tasks.map(t => t.ID === task.ID ? task : t) // rollback
    }
  }

  let newTask = ''
  let newDesc = ''
  let creating = false

  async function createTask() {
    const text = newTask.trim()
    if (!text || creating) return
    creating = true
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/task/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ Task: text, Description: newDesc.trim(), Status: false }),
      })
      const raw = await res.json()
      console.log('POST /task/create response:', JSON.stringify(raw))
      const taskField = raw.Task ?? raw.task ?? raw.task_name ?? raw.name
      const created: Task = {
        ID:          raw.ID          ?? raw.id          ?? Date.now(),
        Task:        typeof taskField === 'string' ? taskField : text,
        Description: raw.Description ?? raw.description ?? newDesc.trim(),
        Status:      raw.Status      ?? raw.status      ?? false,
        CreatedAt:   raw.CreatedAt   ?? raw.created_at  ?? new Date().toISOString(),
      }
      tasks = [created, ...tasks]
      newTask = ''
      newDesc = ''
    } catch {
      error = 'Failed to create task'
    } finally {
      creating = false
    }
  }

  async function deleteTask(task: Task) {
    tasks = tasks.filter(t => t.ID !== task.ID)
    try {
      await fetch(`${import.meta.env.VITE_API_URL}/task/delete/${task.ID}`, {
        method: 'DELETE',
      })
    } catch {
      tasks = [task, ...tasks] // rollback
    }
  }

  function formatDate(iso: string) {
    if (!iso) return '—'
    const d = new Date(iso)
    if (isNaN(d.getTime())) return '—'
    return d.toLocaleDateString('en-GB', { day: '2-digit', month: 'short' })
  }
</script>

<svelte:head>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;500;600&family=Syne:wght@400;700;800&display=swap" rel="stylesheet" />
</svelte:head>

<div class="app">
  <header>
    <div class="label">// task runner</div>
    <h1>to<span>do</span>s</h1>
    <div class="meta">
      <span>{remaining} remaining</span>
      <span><em>{doneCount}</em> done</span>
    </div>
  </header>

  <div class="input-row" class:focused={false}>
    <div class="prompt">&gt;_</div>
    <div class="input-fields">
      <input
        class="todo-input"
        bind:value={newTask}
        on:keydown={e => e.key === 'Enter' && createTask()}
        placeholder="new task..."
        autocomplete="off"
        disabled={creating}
      />
      <input
        class="todo-input desc-input"
        bind:value={newDesc}
        on:keydown={e => e.key === 'Enter' && createTask()}
        placeholder="description (optional)"
        autocomplete="off"
        disabled={creating}
      />
    </div>
    <button class="add-btn" on:click={createTask} disabled={creating}>
      {#if creating}<span class="btn-spinner" />{:else}+{/if}
    </button>
  </div>

  <div class="filters">
    {#each (['all', 'active', 'done'] as const) as f}
      <button class="filter-btn" class:active={filter === f} on:click={() => filter = f}>{f}</button>
    {/each}
  </div>

  {#if loading}
    <div class="status-msg">
      <span class="spinner" />
      fetching tasks...
    </div>
  {:else if error}
    <div class="status-msg error">{error}</div>
  {:else}
    <ul class="todo-list">
      {#each filtered as task, i (task.ID)}
        <li
          class="todo-item"
          class:done={task.Status}
          animate:flip={{ duration: 200 }}
          in:fly={{ y: -8, duration: 200 }}
          out:fade={{ duration: 150 }}
        >
          <input
            class="check"
            type="checkbox"
            checked={task.Status}
            on:change={() => toggleStatus(task)}
          />
          <div class="content">
            <span class="title">{task.Task}</span>
            {#if task.Description}
              <span class="desc">{task.Description}</span>
            {/if}
          </div>
          <div class="right">
            <span class="date">{formatDate(task.CreatedAt)}</span>
            <span class="index">{String(i + 1).padStart(2, '0')}</span>
          </div>
          {#if task.Status}
            <button class="delete-btn" on:click={() => deleteTask(task)} title="delete">🗑</button>
          {/if}
        </li>
      {:else}
        <div class="empty" transition:fade>
          <div>{filter === 'done' ? '✓' : '○'}</div>
          {filter === 'done' ? 'nothing done yet' : filter === 'active' ? 'all caught up' : 'no tasks'}
        </div>
      {/each}
    </ul>
  {/if}

  <footer>
    <span>{tasks.length} total</span>
    <span class="endpoint">GET localhost:8080/tasks</span>
  </footer>
</div>

<style>
  :global(body) { background: #0d0d0f; }

  .app {
    --surface: #141417;
    --border:  #22222a;
    --muted:   #3a3a48;
    --subtle:  #5a5a72;
    --text:    #e2e2ec;
    --dim:     #888899;
    --accent:  #7c6af7;
    --accent2: #4fd1c5;
    --danger:  #f56565;

    font-family: 'JetBrains Mono', monospace;
    color: var(--text);
    background: #0d0d0f;
    max-width: 600px;
    margin: 60px auto;
    padding: 0 20px;
  }

  header { margin-bottom: 36px; }

  .label {
    font-size: 11px;
    letter-spacing: 0.2em;
    text-transform: uppercase;
    color: var(--accent);
    margin-bottom: 6px;
  }

  h1 {
    font-family: 'Syne', sans-serif;
    font-size: 42px;
    font-weight: 800;
    line-height: 1;
    letter-spacing: -0.02em;
  }
  h1 span { color: var(--accent); }

  .meta {
    margin-top: 10px;
    font-size: 12px;
    color: var(--dim);
    display: flex;
    gap: 20px;
  }
  .meta em { color: var(--accent2); font-style: normal; }

  .filters {
    display: flex;
    gap: 4px;
    margin-bottom: 16px;
    border-bottom: 1px solid var(--border);
    padding-bottom: 16px;
  }

  .filter-btn {
    background: none;
    border: 1px solid transparent;
    color: var(--dim);
    padding: 4px 12px;
    border-radius: 4px;
    font-family: 'JetBrains Mono', monospace;
    font-size: 11px;
    letter-spacing: 0.08em;
    cursor: pointer;
    transition: all 0.15s;
  }
  .filter-btn:hover { color: var(--text); border-color: var(--border); }
  .filter-btn.active {
    color: var(--accent);
    border-color: var(--accent);
    background: rgba(124,106,247,0.08);
  }

  /* Input */
  .input-row {
    display: flex;
    border: 1px solid var(--border);
    border-radius: 6px;
    overflow: hidden;
    background: var(--surface);
    margin-bottom: 28px;
    transition: border-color 0.2s, box-shadow 0.2s;
  }
  .input-row:focus-within {
    border-color: var(--accent);
    box-shadow: 0 0 0 3px rgba(124,106,247,0.12);
  }

  .prompt {
    padding: 0 14px;
    color: var(--accent);
    font-size: 16px;
    display: flex;
    align-items: center;
    border-right: 1px solid var(--border);
    background: rgba(124,106,247,0.05);
    user-select: none;
    flex-shrink: 0;
  }

  .input-fields {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .todo-input {
    background: transparent;
    border: none;
    outline: none;
    padding: 10px 16px;
    font-family: 'JetBrains Mono', monospace;
    font-size: 14px;
    color: var(--text);
    caret-color: var(--accent);
    width: 100%;
  }
  .todo-input::placeholder { color: var(--muted); }
  .todo-input:disabled { opacity: 0.5; cursor: not-allowed; }

  .desc-input {
    font-size: 11px;
    padding-top: 0;
    border-top: 1px solid var(--border);
    color: var(--dim);
  }

  .add-btn {
    background: var(--accent);
    border: none;
    color: #fff;
    padding: 0 14px;
    cursor: pointer;
    font-size: 16px;
    font-weight: 400;
    transition: background 0.15s;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }
  .add-btn:hover:not(:disabled) { background: #6a5ae0; }
  .add-btn:disabled { opacity: 0.6; cursor: not-allowed; }

  .btn-spinner {
    width: 14px;
    height: 14px;
    border: 1.5px solid rgba(255,255,255,0.3);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }

  .status-msg {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 12px;
    color: var(--dim);
    padding: 32px 0;
  }
  .status-msg.error { color: var(--danger); }

  .spinner {
    width: 14px;
    height: 14px;
    border: 1.5px solid var(--muted);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
    flex-shrink: 0;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .todo-list {
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: 2px;
    padding: 0;
  }

  .todo-item {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 13px 16px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 6px;
    transition: border-color 0.2s, background 0.2s;
  }
  .todo-item:hover { border-color: var(--muted); background: #18181e; }
  .todo-item.done { background: #111114; border-color: #1a1a22; }
  .todo-item.done .title {
    color: var(--subtle);
    text-decoration: line-through;
    text-decoration-color: var(--muted);
  }
  .todo-item.done .desc { color: #333340; }

  .check {
    appearance: none;
    -webkit-appearance: none;
    width: 18px;
    height: 18px;
    border: 1.5px solid var(--muted);
    border-radius: 4px;
    cursor: pointer;
    flex-shrink: 0;
    position: relative;
    transition: all 0.15s;
    background: transparent;
  }
  .check:hover { border-color: var(--accent); }
  .check:checked { background: var(--accent); border-color: var(--accent); }
  .check:checked::after {
    content: '';
    position: absolute;
    left: 4px; top: 1px;
    width: 5px; height: 9px;
    border: 2px solid #fff;
    border-top: none;
    border-left: none;
    transform: rotate(45deg);
  }

  .content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 3px;
    overflow: hidden;
  }

  .title {
    font-size: 13.5px;
    color: var(--text);
    transition: color 0.2s;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .desc {
    font-size: 11px;
    color: var(--dim);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    transition: color 0.2s;
  }

  .right {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
    flex-shrink: 0;
  }

  .delete-btn {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 14px;
    padding: 2px 4px;
    border-radius: 3px;
    opacity: 0.3;
    transition: opacity 0.15s, background 0.15s;
    line-height: 1;
    flex-shrink: 0;
  }
  .todo-item:hover .delete-btn { opacity: 0.7; }
  .delete-btn:hover { opacity: 1 !important; background: rgba(245,101,101,0.12); }

  .date { font-size: 10px; color: var(--subtle); }
  .index { font-size: 10px; color: var(--muted); }

  .empty {
    text-align: center;
    color: var(--subtle);
    font-size: 12px;
    padding: 40px 0;
    letter-spacing: 0.05em;
    border: 1px dashed var(--border);
    border-radius: 6px;
  }
  .empty div { font-size: 28px; margin-bottom: 10px; }

  footer {
    margin-top: 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 11px;
    color: var(--dim);
    padding-top: 16px;
    border-top: 1px solid var(--border);
  }

  .endpoint { color: var(--muted); font-size: 10px; }
</style>
