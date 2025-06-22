let tgUser = null;
let userId = null;
let carriers = [];

// Получаем данные пользователя из Telegram WebApp API
function getTelegramUser() {
  if (window.Telegram && Telegram.WebApp && Telegram.WebApp.initDataUnsafe && Telegram.WebApp.initDataUnsafe.user) {
    return Telegram.WebApp.initDataUnsafe.user;
  }
  return null;
}

document.addEventListener('DOMContentLoaded', async () => {
  tgUser = getTelegramUser();
  if (!tgUser) {
    document.getElementById('user-info').innerText = 'Ошибка: не удалось получить данные пользователя Telegram.';
    return;
  }
  userId = tgUser.id;
  document.getElementById('user-info').innerText = `Пользователь: @${tgUser.username || 'без ника'} (${tgUser.id})`;

  // Загрузка перевозчиков
  await loadCarriers();

  // Регистрируем пользователя на сервере
  await fetch('/api/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ telegram_id: tgUser.id, username: tgUser.username })
  });

  // Загрузка списка треков
  await loadTrackings();

  // Обработка добавления трека
  document.getElementById('add-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const tracking_number = document.getElementById('tracking_number').value.trim();
    const carrier_code = document.getElementById('carrier_code').value.trim();
    if (!tracking_number || !carrier_code) return;
    await fetch('/api/tracking', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: userId, tracking_number, carrier_code })
    });
    document.getElementById('tracking_number').value = '';
    document.getElementById('carrier_code').value = '';
    await loadTrackings();
  });

  // Модальное окно
  document.getElementById('modal-close').onclick = () => {
    document.getElementById('modal').style.display = 'none';
  };
  document.getElementById('modal').onclick = (e) => {
    if (e.target === document.getElementById('modal')) {
      document.getElementById('modal').style.display = 'none';
    }
  };
});

async function loadCarriers() {
  const select = document.getElementById('carrier_code');
  select.innerHTML = '<option value="">Загрузка...</option>';
  try {
    const resp = await fetch('/api/carriers');
    if (!resp.ok) throw new Error();
    carriers = await resp.json();
    select.innerHTML = '<option value="">Выберите перевозчика</option>';
    carriers.forEach(c => {
      const opt = document.createElement('option');
      opt.value = c.code;
      opt.innerText = c.name;
      select.appendChild(opt);
    });
  } catch {
    select.innerHTML = '<option value="">Ошибка загрузки</option>';
  }
}

async function loadTrackings() {
  const list = document.getElementById('trackings-list');
  list.innerHTML = '<li>Загрузка...</li>';
  const resp = await fetch(`/api/tracking?user_id=${userId}`);
  if (!resp.ok) {
    list.innerHTML = '<li>Ошибка загрузки</li>';
    return;
  }
  const data = await resp.json();
  if (!data.length) {
    list.innerHTML = '<li>Нет треков</li>';
    return;
  }
  list.innerHTML = '';
  data.forEach(tr => {
    const li = document.createElement('li');
    li.innerHTML = `<span class="track-link" style="cursor:pointer;">${tr.tracking_number} <span class="status">[${tr.status || '—'}]</span></span>`;
    li.querySelector('.track-link').onclick = () => showTrackingDetails(tr.id);
    const delBtn = document.createElement('button');
    delBtn.innerText = 'Удалить';
    delBtn.onclick = async () => {
      await fetch(`/api/tracking?id=${tr.id}`, { method: 'DELETE' });
      await loadTrackings();
    };
    li.appendChild(delBtn);
    list.appendChild(li);
  });
}

async function showTrackingDetails(id) {
  const resp = await fetch(`/api/tracking/status?id=${id}`);
  if (!resp.ok) {
    showModal('Ошибка загрузки статуса');
    return;
  }
  const tr = await resp.json();
  let html = `<b>Трек-номер:</b> ${tr.tracking_number}<br>`;
  html += `<b>Перевозчик:</b> ${tr.carrier_code}<br>`;
  html += `<b>Статус:</b> ${tr.status || '—'}<br>`;
  html += `<b>Обновлено:</b> ${tr.last_update ? new Date(tr.last_update).toLocaleString() : '—'}<br>`;
  showModal(html);
}

function showModal(html) {
  document.getElementById('modal-body').innerHTML = html;
  document.getElementById('modal').style.display = 'flex';
} 