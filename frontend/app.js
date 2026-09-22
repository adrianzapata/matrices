'use strict';

let token = null;

const $ = (id) => document.getElementById(id);

function setStatus(el, message, kind) {
  el.textContent = message;
  el.className = 'status' + (kind ? ' ' + kind : '');
}

function formatMatrix(matrix) {
  return matrix.map((row) => row.map((v) => Number(v.toFixed(4))).join('\t')).join('\n');
}

$('loginBtn').addEventListener('click', async () => {
  const baseUrl = $('goApiUrl').value.replace(/\/+$/, '');
  const username = $('username').value;
  const password = $('password').value;
  const statusEl = $('tokenStatus');

  setStatus(statusEl, 'Solicitando token...');
  try {
    const res = await fetch(`${baseUrl}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || `HTTP ${res.status}`);
    token = data.token;
    setStatus(statusEl, 'Token obtenido correctamente.', 'ok');
  } catch (err) {
    token = null;
    setStatus(statusEl, `Error: ${err.message}`, 'error');
  }
});

$('processBtn').addEventListener('click', async () => {
  const baseUrl = $('goApiUrl').value.replace(/\/+$/, '');
  const statusEl = $('processStatus');
  const resultsEl = $('results');

  let matrix;
  try {
    matrix = JSON.parse($('matrixInput').value);
  } catch (err) {
    setStatus(statusEl, 'La matriz no es un JSON válido.', 'error');
    return;
  }

  if (!token) {
    setStatus(statusEl, 'Primero obtén un token (paso 1).', 'error');
    return;
  }

  setStatus(statusEl, 'Procesando...');
  resultsEl.hidden = true;

  try {
    const res = await fetch(`${baseUrl}/api/matrix/process`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ matrix }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || `HTTP ${res.status}`);

    $('originalOut').textContent = formatMatrix(data.original);
    $('rotatedOut').textContent = formatMatrix(data.rotated);
    $('qOut').textContent = formatMatrix(data.q);
    $('rOut').textContent = formatMatrix(data.r);

    $('statMax').textContent = Number(data.stats.max.toFixed(4));
    $('statMin').textContent = Number(data.stats.min.toFixed(4));
    $('statAvg').textContent = Number(data.stats.average.toFixed(4));
    $('statSum').textContent = Number(data.stats.sum.toFixed(4));
    $('statDiag').textContent =
      `${data.stats.diagonal.rotated} / ${data.stats.diagonal.q} / ${data.stats.diagonal.r}`;

    resultsEl.hidden = false;
    setStatus(statusEl, 'Listo.', 'ok');
  } catch (err) {
    setStatus(statusEl, `Error: ${err.message}`, 'error');
  }
});
