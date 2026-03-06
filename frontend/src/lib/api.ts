const BASE = "http://localhost:8080"

export async function getTasks() {
    const res = await fetch(`${BASE}/tasks`)
    return res.json()
}


export async function createTask(task: any) {
    const res = await fetch(`${BASE}/task/create`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(task),
    })
    return res.json()
}
