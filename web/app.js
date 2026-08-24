async function loadSignals() {
  const response = await fetch('/api/signals');
  const data = await response.json();
  document.querySelector('#signals').innerHTML = data.signals.map((signal) =>
    `<li>${signal.id} · ${signal.block} · ${signal.aspect} · ${signal.inspection}</li>`).join('');
}
document.querySelector('#refresh').addEventListener('click', loadSignals);
loadSignals();
