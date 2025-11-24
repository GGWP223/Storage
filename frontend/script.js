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

let access = localStorage.getItem('access');
let refresh = localStorage.getItem('refresh');
let auth_api = "http://localhost:8082/auth";
let storage_api = "http://localhost:8085/storage";

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
    try {
        const res = await fetch(auth_api + '/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({login, password})
        });

        if (!res.ok)
            throw new Error('Ошибка авторизации');

        const data = await res.json();

        if (!data.access)
            throw new Error('access Токен не получен');

        if (!data.refresh)
            throw new Error('refresh Токен не получен');

        access = data.access;
        refresh = data.refresh;
        localStorage.setItem('access', access);
        localStorage.setItem('refresh', refresh);

        showStorage();
    } catch (e) {
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
}

async function fetchFiles() {
    clearMessage();
    try {
        const res = await fetch(storage_api + '/getAllFiles',{
            method: 'POST',
            headers: {'Authorization': 'Bearer ' + access},
            body: JSON.stringify({ "token": access })
        });

        if (!res.ok)
            throw new Error('Ошибка загрузки списка файлов');

        const response = await res.json();

        renderFiles(response.meta);
    } catch (e) {
        showMessage(e.message);
    }
}

function renderFiles(files) {
    filesTableBody.innerHTML = '';
    files.forEach(f => {
        const tr = document.createElement('tr');

        tr.appendChild(document.createElement('td')).textContent = f.fileName;
        tr.appendChild(document.createElement('td')).textContent = f.mimeType;
        tr.appendChild(document.createElement('td')).textContent = formatSize(f.size);
        tr.appendChild(document.createElement('td')).textContent = formatDate(f.createdAt);
        tr.appendChild(document.createElement('td')).textContent = formatDate(f.updatedAt);

        const actionsTd = document.createElement('td');
        actionsTd.appendChild(createButton('Скачать', () => downloadFile(f.fileID)));
        actionsTd.appendChild(createButton('Удалить', () => deleteFile(f.fileID)));
        tr.appendChild(actionsTd);

        filesTableBody.appendChild(tr);
    });
}

async function downloadFile(fileID) {
    clearMessage();

    try {
        const res = await fetch(storage_api + '/getFile', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ token: access, fileID })
        });

        if (!res.ok)
            throw new Error('Ошибка скачивания файла');

        const json = await res.json();

        if (!json.meta || !json.data)
            throw new Error('Неверный формат ответа');

        const base64data = json.data;
        const filename = json.meta.fileName || 'file';

        const byteCharacters = atob(base64data);
        const byteNumbers = new Array(byteCharacters.length);

        for (let i = 0; i < byteCharacters.length; i++) {
            byteNumbers[i] = byteCharacters.charCodeAt(i);
        }

        const byteArray = new Uint8Array(byteNumbers);
        const blob = new Blob([byteArray], { type: json.meta.mimeType });

        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');

        a.href = url;
        a.download = filename;
        a.click();

        window.URL.revokeObjectURL(url);
    } catch (e) {
        showMessage(e.message);
    }
}


async function deleteFile(fileID) {
    clearMessage();

    if (!confirm('Удалить файл?'))
        return;

    try {
        const res = await fetch(storage_api + '/deleteFile', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({token: access, fileID: fileID})
        });

        if (!res.ok)
            throw new Error('Ошибка удаления файла');

        await fetchFiles();
    } catch (e) {
        showMessage(e.message);
    }
}

uploadBtn.onclick = async () => {
    clearMessage();
    const file = fileInput.files[0];
    if (!file) {
        showMessage('Выберите файл для загрузки');
        return;
    }

    try {
        const arrayBuffer = await file.arrayBuffer();
        const uint8Array = new Uint8Array(arrayBuffer);
        const base64Data = uint8ToBase64(uint8Array);

        const json = JSON.stringify({
            token: access,
            fileName: file.name,
            mimeType: file.type,
            data: base64Data
        });

        const res = await fetch(storage_api + '/uploadFile', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: json
        });

        if (!res.ok)
            throw new Error('Ошибка загрузки файла');

        fileInput.value = '';

        await fetchFiles();
    } catch (e) {
        showMessage(e.message);
    }
};


function uint8ToBase64(u8Arr) {
    const CHUNK_SIZE = 0x8000; // 32768
    let index = 0;
    const length = u8Arr.length;
    let result = '';
    let slice;
    while (index < length) {
        slice = u8Arr.subarray(index, Math.min(index + CHUNK_SIZE, length));
        result += String.fromCharCode.apply(null, slice);
        index += CHUNK_SIZE;
    }
    return btoa(result);
}

logoutBtn.onclick = () => {
    localStorage.removeItem('access');
    localStorage.removeItem('refresh');
    access = null;
    refresh = null;
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

if (access) {
    showStorage();
}

btnLogin.onclick = login;
btnRegister.onclick = register;
