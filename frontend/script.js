const authDiv = document.getElementById('auth');
const storageDiv = document.getElementById('storage');
const errorDiv = document.getElementById('error');
const messageDiv = document.getElementById('message');

const loginInput = document.getElementById('login');
const passwordInput = document.getElementById('password');
const btnLogin = document.getElementById('btnLogin');
const btnRegister = document.getElementById('btnRegister');

const fileInput = document.getElementById('fileInput');
const uploadBtn = document.getElementById('uploadBtn');
const logoutBtn = document.getElementById('logoutBtn');
const filesTableBody = document.querySelector('#filesTable tbody');

let token = localStorage.getItem('token');
let auth_api = "http://localhost:8082"
let storage_api = "http://localhost:8081"

function showError(msg) {
    errorDiv.textContent = msg;
}
function clearError() {
    errorDiv.textContent = '';
}
function showMessage(msg) {
    messageDiv.textContent = msg;
}
function clearMessage() {
    messageDiv.textContent = '';
}

function formatDate(dateStr) {
    return new Date(dateStr).toLocaleString();
}
function formatSize(bytes) {
    return (bytes / 1024).toFixed(2);
}

function createButton(text, onClick) {
    const btn = document.createElement('button');
    btn.textContent = text;
    btn.onclick = onClick;
    return btn;
}

async function login() {
    clearError();

    const login = loginInput.value.trim();
    const password = passwordInput.value.trim();

    if (!login || !password) {
        showError('Введите логин и пароль');
        return;
    }
    try
    {
        const res = await fetch(auth_api + '/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({login, password})
        });

        if (!res.ok)
            throw new Error('Ошибка авторизации');

        const data = await res.json();

        if (!data.token)
            throw new Error('Токен не получен');

        token = data.token;
        localStorage.setItem('token', token);

        showStorage();
    }
    catch (e)
    {
        showError(e.message);
    }
}

async function register() {
    clearError();

    const login = loginInput.value.trim();
    const password = passwordInput.value.trim();

    if (!login || !password) {
        showError('Введите логин и пароль');
        return;
    }

    try {
        const res = await fetch(auth_api + '/register', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({login, password})
        });

        if (!res.ok) {
            const data = await res.json();
            throw new Error(data.error || 'Ошибка регистрации');
        }

        alert('Регистрация прошла успешно, теперь можно войти');
    } catch (e) {
        showError(e.message);
    }
};

async function fetchFiles() {
    clearMessage();
    try
    {
        const res = await fetch(storage_api + '/GetFiles?token=' + token, {
            headers: {'Authorization': 'Bearer ' + token}
        });

        if (!res.ok)
            throw new Error('Ошибка загрузки списка файлов');

        const files = await res.json();

        renderFiles(files);
    }
    catch (e)
    {
        showMessage(e.message);
    }
}

function renderFiles(files) {
    filesTableBody.innerHTML = '';
    files.forEach(f => {
        const tr = document.createElement('tr');

        tr.appendChild(document.createElement('td')).textContent = f.filename;
        tr.appendChild(document.createElement('td')).textContent = f.mime_type;
        tr.appendChild(document.createElement('td')).textContent = formatSize(f.size);
        tr.appendChild(document.createElement('td')).textContent = formatDate(f.created_at);
        tr.appendChild(document.createElement('td')).textContent = formatDate(f.updated_at);

        const actionsTd = document.createElement('td');
        actionsTd.appendChild(createButton('Скачать', () => downloadFile(f.file_id)));
        actionsTd.appendChild(createButton('Удалить', () => deleteFile(f.file_id)));
        tr.appendChild(actionsTd);

        filesTableBody.appendChild(tr);
    });
}

async function downloadFile(fileID) {
    clearMessage();

    try
    {
        const res = await fetch(storage_api + '/DownloadFile', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({token, file_id: fileID})
        });

        if (!res.ok)
            throw new Error('Ошибка скачивания файла');

        const blob = await res.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');

        let filename = 'file';
        const disposition = res.headers.get('Content-Disposition');
        console.log("0")
        console.log(disposition)
        if (disposition && disposition.includes('filename=')) {
            console.log("1")
            const filenameMatch = disposition.match(/filename="?([^"]+)"?/);
            if (filenameMatch.length > 1) {
                console.log("2")
                filename = filenameMatch[1];
            }
        }

        a.href = url;
        a.download = filename;
        a.click();

        window.URL.revokeObjectURL(url);
    } catch (e)
    {
        showMessage(e.message);
    }
}

async function deleteFile(fileID) {
    clearMessage();

    if (!confirm('Удалить файл?'))
        return;

    try
    {
        const res = await fetch(storage_api + '/DeleteFile', {
            method: 'DELETE',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({token, file_id: fileID})
        });

        if (!res.ok)
            throw new Error('Ошибка удаления файла');

        await fetchFiles();
    }
    catch (e)
    {
        showMessage(e.message);
    }
}

uploadBtn.onclick = async () => {
    clearMessage();
    const file = fileInput.files[0];
    if (!file)
    {
        showMessage('Выберите файл для загрузки');
        return;
    }
    try
    {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('token', token);

        const res = await fetch(storage_api + '/UploadFile', {
            method: 'POST',
            body: formData
        });

        if (!res.ok)
            throw new Error('Ошибка загрузки файла');

        fileInput.value = '';

        await fetchFiles();
    }
    catch (e)
    {
        showMessage(e.message);
    }
};

logoutBtn.onclick = () => {

    localStorage.removeItem('token');
    token = null;
    storageDiv.style.display = 'none';
    authDiv.style.display = 'block';

    clearMessage();
};

function showStorage() {
    authDiv.style.display = 'none';
    storageDiv.style.display = 'block';
    clearError();
    fetchFiles();
}

if (token) {
    showStorage();
}

btnLogin.onclick = login;
btnRegister.onclick = register;
