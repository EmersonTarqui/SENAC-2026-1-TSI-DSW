const API_BASE_URL = 'http://localhost:8080';

var api = {
    async login(email, password) {
        const res = await fetch(`${API_BASE_URL}/api/v1/login`,{
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        if (!res.ok) throw new Error('Credenciais inválidas');
        const userData = await res.json();
        SessionManager.saveSession(userData);
        return userData;
    },

    async register(email, password) {
        const res = await fetch(`${API_BASE_URL}/api/v1/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        if (!res.ok) throw new Error('Erro ao cadastrar. Tente outro e-mail.');
        return await res.json();
    },
    
    // Busca tarefas
    async getTasks() {
        const userId = SessionManager.getUserId();
        const res = await fetch(`${API_BASE_URL}/api/v1/tasks?user_id=${userId}`);
        if (!res.ok) throw new Error('Erro ao buscar tarefas');
        const response = await res.json();
        return response.data || [];
    },

    // Cria tarefa
    async createTask(title) {
        const userId = SessionManager.getUserId();
        const res = await fetch(`${API_BASE_URL}/api/v1/tasks`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                title: title, 
                user_id: Number(userId), 
                done: false 
            })
        });
        if (!res.ok) throw new Error('Erro ao criar');
        const response = await res.json();
        return response.data; // Retorna a tarefa com ID
    },

    // Atualiza status
    async updateTask(taskData) {
        const taskId = taskData.id || taskData.ID;
        const res = await fetch(`${API_BASE_URL}/api/v1/tasks/${taskId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                id: Number(taskId),
                user_id: Number(taskData.user_id || taskData.UserID),
                title: taskData.title || taskData.Title,
                done: Boolean(taskData.done || taskData.Done) 
            })
        });
        if (!res.ok) throw new Error('Erro ao atualizar');
        return await res.json();
    },

    // Deleta tarefa
    async deleteTask(taskId) {
        const userId = SessionManager.getUserId();
        const res = await fetch(`${API_BASE_URL}/api/v1/tasks/${taskId}`, {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ user_id: Number(userId) })
        });
        if (!res.ok) throw new Error('Erro ao excluir');
    }
};