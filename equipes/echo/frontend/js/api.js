const API_BASE_URL = 'http://localhost:8080';

const api = {
    // Autenticação de usuário
    async login(email, password) {
        const res = await fetch(${API_BASE_URL}/login, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        
        if (!res.ok) throw new Error('Credenciais inválidas ou erro no servidor.');
        
        const userData = await res.json();
        
        // Salva os dados na sessão através do SessionManager
        SessionManager.saveSession(userData);
        
        return userData;
    },

    // Busca tarefas filtradas pelo ID do usuário logado
    async getTasks() {
        const userId = SessionManager.getUserId();
        
        // Envia o user_id na URL para o backend Go filtrar os resultados
        const res = await fetch(${API_BASE_URL}/tasks?user_id=${userId});
        
        if (!res.ok) throw new Error('Erro ao buscar tarefas. Verifique a conexão com a API.');
        return await res.json();
    },

    // Cria uma nova tarefa vinculada ao usuário atual
    async createTask(title) {
        const userId = SessionManager.getUserId();
        
        const res = await fetch(${API_BASE_URL}/tasks, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                title: title, 
                user_id: parseInt(userId), 
                done: false 
            })
        });
        
        if (!res.ok) throw new Error('Erro ao criar tarefa no banco de dados.');
        return await res.json();
    },

    // Atualiza o status de uma tarefa existente
    async updateTask(id, done) {
        const res = await fetch(${API_BASE_URL}/tasks/${id}, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ done: done })
        });
        
        if (!res.ok) throw new Error('Erro ao atualizar a tarefa. Rota PUT não encontrada ou inválida.');
        return await res.json();
    },

    // Remove uma tarefa do banco
    async deleteTask(id) {
        const res = await fetch(${API_BASE_URL}/tasks/${id}, {
            method: 'DELETE'
        });
        
        if (!res.ok) throw new Error('Erro ao excluir a tarefa. Rota DELETE não encontrada ou inválida.');
    }
};