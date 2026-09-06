const elements = {
  refreshButton: document.getElementById('refreshButton'),
  statusPanel: document.getElementById('statusPanel'),
  workspaceCount: document.getElementById('workspaceCount'),
  environmentCount: document.getElementById('environmentCount'),
  execCount: document.getElementById('execCount'),
  mcpCount: document.getElementById('mcpCount'),
  skillCount: document.getElementById('skillCount'),
  memoryCount: document.getElementById('memoryCount'),
  workspaceBadge: document.getElementById('workspaceBadge'),
  environmentBadge: document.getElementById('environmentBadge'),
  workspaceList: document.getElementById('workspaceList'),
  environmentList: document.getElementById('environmentList'),
  execList: document.getElementById('execList'),
  catalogSummary: document.getElementById('catalogSummary'),
};

function desktopAdapter() {
  const adapter = window.go?.desktop?.Adapter;
  if (!adapter?.GetSnapshot) {
    throw new Error('Wails desktop binding is not ready');
  }
  return adapter;
}

function safeArray(value) {
  return Array.isArray(value) ? value : [];
}

function setStatus(message, kind = 'normal') {
  elements.statusPanel.textContent = message;
  elements.statusPanel.dataset.kind = kind;
}

function setMetric(element, value) {
  element.textContent = String(value);
}

function emptyMessage(container, message) {
  container.replaceChildren();
  container.classList.add('empty');
  container.textContent = message;
}

function renderWorkspaces(workspaces) {
  elements.workspaceBadge.textContent = String(workspaces.length);
  if (workspaces.length === 0) {
    emptyMessage(elements.workspaceList, '暂无 Workspace');
    return;
  }
  elements.workspaceList.replaceChildren();
  elements.workspaceList.classList.remove('empty');
  for (const workspace of workspaces) {
    const item = document.createElement('article');
    item.className = 'list-item';

    const title = document.createElement('strong');
    title.textContent = workspace.name || workspace.workspace_id;
    const id = document.createElement('code');
    id.textContent = workspace.workspace_id || '';
    const path = document.createElement('span');
    path.textContent = workspace.path || '';

    item.append(title, id, path);
    elements.workspaceList.append(item);
  }
}

function renderEnvironments(environments) {
  elements.environmentBadge.textContent = String(environments.length);
  if (environments.length === 0) {
    emptyMessage(elements.environmentList, '暂无 Environment');
    return;
  }
  elements.environmentList.replaceChildren();
  elements.environmentList.classList.remove('empty');
  for (const environment of environments) {
    const item = document.createElement('article');
    item.className = 'list-item';

    const title = document.createElement('strong');
    title.textContent = environment.name || environment.environment_id;
    const id = document.createElement('code');
    id.textContent = environment.environment_id || '';
    const root = document.createElement('span');
    root.textContent = environment.root || '';
    const meta = document.createElement('small');
    const memoryCount = Number(environment.private_memory_count || 0);
    meta.textContent = `Workspace ${environment.workspace_id || '—'} · Private Memory ${memoryCount}`;

    item.append(title, id, root, meta);
    elements.environmentList.append(item);
  }
}

function renderExecutables(executables) {
  if (executables.length === 0) {
    emptyMessage(elements.execList, '暂无允许的 executable');
    return;
  }
  elements.execList.replaceChildren();
  elements.execList.classList.remove('empty');
  for (const executable of executables) {
    const chip = document.createElement('code');
    chip.className = 'chip';
    chip.textContent = executable;
    elements.execList.append(chip);
  }
}

function renderSnapshot(snapshot) {
  const workspaces = safeArray(snapshot.workspaces);
  const environments = safeArray(snapshot.environments);
  const executables = safeArray(snapshot.allowed_executables);
  const mcps = safeArray(snapshot.mcps);
  const skills = safeArray(snapshot.skills);
  const globalMemoryCount = Number(snapshot.global_memory_count || 0);

  setMetric(elements.workspaceCount, workspaces.length);
  setMetric(elements.environmentCount, environments.length);
  setMetric(elements.execCount, executables.length);
  setMetric(elements.mcpCount, mcps.length);
  setMetric(elements.skillCount, skills.length);
  setMetric(elements.memoryCount, globalMemoryCount);

  renderWorkspaces(workspaces);
  renderEnvironments(environments);
  renderExecutables(executables);
  elements.catalogSummary.textContent = `${mcps.length} MCP · ${skills.length} Skill`;
}

async function refreshSnapshot() {
  elements.refreshButton.disabled = true;
  setStatus('正在读取 ADM 状态…', 'loading');
  try {
    const snapshot = await desktopAdapter().GetSnapshot();
    renderSnapshot(snapshot);
    setStatus(`已刷新 · ${new Date().toLocaleTimeString()}`, 'success');
  } catch (error) {
    const message = error?.message || String(error);
    setStatus(`读取失败：${message}`, 'error');
  } finally {
    elements.refreshButton.disabled = false;
  }
}

elements.refreshButton.addEventListener('click', refreshSnapshot);
window.addEventListener('DOMContentLoaded', refreshSnapshot);
