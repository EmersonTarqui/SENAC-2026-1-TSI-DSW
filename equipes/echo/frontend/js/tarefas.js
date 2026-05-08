// Seleção de elementos do DOM
const taskInput = document.getElementById('task-input');
const addBtn = document.getElementById('add-btn');
const taskList = document.getElementById('task-list');
const filterBtns = document.querySelectorAll('.filter-btn');

let allTasks = [];
let currentFilter = 'all';

// Carregar tarefas iniciais
window.loadTasks = async function() {
    try {
        const data = await api.getTasks();
        allTasks = data || [];
        render();
    } catch (error) {
        taskList.innerHTML = <li>Erro ao carregar tarefas: ${error.message}</li>;
    }
};

// Adicionar nova tarefa
async function handleAddTask() {
    const title = taskInput.value.trim();
    if (!title) return;

    try {
        const savedTask = await api.createTask(title);
        allTasks.push(savedTask);
        taskInput.value = '';
        render();
    } catch (error) {
        alert('Erro ao adicionar tarefa: ' + error.message);
    }
}

// Alternar status da tarefa (concluída/pendente)
window.toggleTask = async function(id, currentStatus) {
    try {
        await api.updateTask(id, !currentStatus);
        const index = allTasks.findIndex(t => t.id === id);
        if (index > -1) {
            allTasks[index].done = !currentStatus;
            render();
        }
    } catch (error) {
        alert('Erro ao atualizar status: ' + error.message);
    }
};

// Excluir tarefa
window.deleteTask = async function(id) {
    if (!confirm('Deseja realmente excluir esta tarefa?')) return;
    try {
        await api.deleteTask(id);
        allTasks = allTasks.filter(t => t.id !== id);
        render();
    } catch (error) {
        alert('Erro ao excluir tarefa: ' + error.message);
    }
};

// Renderização da lista com suporte a filtros
function render() {
    if (!taskList) return;
    taskList.innerHTML = '';
    
    // Aplica o filtro selecionado
    const filtered = allTasks.filter(t => {
        if (currentFilter === 'pending') return !t.done;
        if (currentFilter === 'completed') return t.done;
        return true;
    });

    if (filtered.length === 0) {
        taskList.innerHTML = '<li>Nenhuma tarefa encontrada.</li>';
        return;
    }

    // Cria os elementos da lista
    filtered.forEach(t => {
        const li = document.createElement('li');
        li.style.display = 'flex';
        li.style.justifyContent = 'space-between';
        li.style.marginBottom = '8px';

        const div = document.createElement('div');
        
        const check = document.createElement('input');
        check.type = 'checkbox';
        check.checked = t.done;
        check.style.marginRight = '8px';
        check.addEventListener('change', () => window.toggleTask(t.id, t.done));

        const span = document.createElement('span');
        span.textContent = t.title;
        if (t.done) span.style.textDecoration = 'line-through';

        div.appendChild(check);
        div.appendChild(span);

        const btn = document.createElement('button');
        btn.textContent = 'Excluir';
        btn.className = 'btn-delete';
        btn.addEventListener('click', () => window.deleteTask(t.id));

        li.appendChild(div);
        li.appendChild(btn);
        taskList.appendChild(li);
    });
}

// Configuração dos ouvintes de eventos
filterBtns.forEach(btn => {
    btn.addEventListener('click', (e) => {
        currentFilter = e.target.getAttribute('data-filter');
        render();
    });
});

addBtn?.addEventListener('click', handleAddTask);

taskInput?.addEventListener('keypress', (e) => { 
    if (e.key === 'Enter') handleAddTask(); 
});

// Inicialização automática
if (taskList) {
    window.loadTasks();
}