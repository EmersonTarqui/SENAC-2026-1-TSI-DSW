const API_BASE_URL = 'http://localhost:8080';

const api = {
    async login(email, password) {
        const res = await fetch(`${API_BASE_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        if (!res.ok) throw new Error('Credenciais inválidas ou rota /login não existe no Go.');
        return await res.json();
    },

    async getTasks() {
        const res = await fetch(`${API_BASE_URL}/tasks`);
        if (!res.ok) throw new Error('Erro ao buscar tarefas. A API está rodando?');
        return await res.json();
    },

    async createTask(title, userId) {
        const res = await fetch(`${API_BASE_URL}/tasks`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title, user_id: userId, done: false })
        });
        if (!res.ok) throw new Error('Erro ao criar tarefa no banco.');
        return await res.json();
    },

    async updateTask(id, done) {
        const res = await fetch(`${API_BASE_URL}/tasks/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ done })
        });
        if (!res.ok) throw new Error('Aviso pra equipe: O Go não tem a rota PUT /tasks/{id}');
        return await res.json();
    },

    async deleteTask(id) {
        const res = await fetch(`${API_BASE_URL}/tasks/${id}`, {
            method: 'DELETE'
        });
        if (!res.ok) throw new Error('Aviso pra equipe: O Go não tem a rota DELETE /tasks/{id}');
    }
};