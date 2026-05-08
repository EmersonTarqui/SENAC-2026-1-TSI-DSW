const taskInput = document.getElementById('task-input');
const addBtn = document.getElementById('add-btn');
const taskList = document.getElementById('task-list');
const filterBtns = document.querySelectorAll('.filter-btn');

let allTasks = [];
let currentFilter = 'all';

// GET
window.loadTasks = async function() {
    try {
        const data = await api.getTasks();
        allTasks = data || [];
        render();
    } catch (error) {
        taskList.innerHTML = `<li>${error.message}</li>`;
    }
};

// POST
async function handleAddTask() {
    const title = taskInput.value.trim();
    if (!title) return;

    const userId = parseInt(localStorage.getItem('logged_user_id')) || 1;

    try {
        const savedTask = await api.createTask(title, userId);
        allTasks.push(savedTask);
        taskInput.value = '';
        render();
    } catch (error) {
        alert(error.message);
    }
}

// PUT
window.toggleTask = async function(id, currentStatus) {
    try {
        await api.updateTask(id, !currentStatus);
        const index = allTasks.findIndex(t => t.id === id);
        if (index > -1) {
            allTasks[index].done = !currentStatus;
            render();
        }
    } catch (error) {
        alert(error.message);
    }
};

// DELETE
window.deleteTask = async function(id) {
    if (!confirm('Deseja excluir esta tarefa?')) return;
    try {
        await api.deleteTask(id);
        allTasks = allTasks.filter(t => t.id !== id);
        render();
    } catch (error) {
        alert(error.message);
    }
};

// Renderização baseada no Filtro (SPA Style)
function render() {
    taskList.innerHTML = '';
    
    const filtered = allTasks.filter(t => {
        if (currentFilter === 'pending') return !t.done;
        if (currentFilter === 'completed') return t.done;
        return true;
    });

    if (filtered.length === 0) {
        taskList.innerHTML = '<li>Nenhuma tarefa encontrada.</li>';
        return;
    }

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
        btn.addEventListener('click', () => window.deleteTask(t.id));

        li.appendChild(div);
        li.appendChild(btn);
        taskList.appendChild(li);
    });
}

// Eventos
filterBtns.forEach(btn => {
    btn.addEventListener('click', (e) => {
        currentFilter = e.target.getAttribute('data-filter');
        render();
    });
});

addBtn.addEventListener('click', handleAddTask);
taskInput.addEventListener('keypress', (e) => { if (e.key === 'Enter') handleAddTask(); });