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
  workspaceForm: document.getElementById('workspaceForm'),
  workspacePath: document.getElementById('workspacePath'),
  workspaceName: document.getElementById('workspaceName'),
  workspaceList: document.getElementById('workspaceList'),
  environmentForm: document.getElementById('environmentForm'),
  environmentWorkspace: document.getElementById('environmentWorkspace'),
  environmentName: document.getElementById('environmentName'),
  environmentRoot: document.getElementById('environmentRoot'),
  environmentList: document.getElementById('environmentList'),
  environmentDetailPanel: document.getElementById('environmentDetailPanel'),
  environmentDetailTitle: document.getElementById('environmentDetailTitle'),
  environmentDetail: document.getElementById('environmentDetail'),
  closeEnvironmentDetail: document.getElementById('closeEnvironmentDetail'),
  execList: document.getElementById('execList'),
  catalogSummary: document.getElementById('catalogSummary'),
};

let currentSnapshot = null;
let selectedEnvironmentID = '';

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

function createActionButton(label, action, id, kind = 'secondary') {
  const button = document.createElement('button');
  button.type = 'button';
  button.className = `small-button ${kind === 'danger' ? 'danger-button' : 'secondary-button'}`;
  button.textContent = label;
  button.dataset.action = action;
  button.dataset.id = id;
  return button;
}

function renderWorkspaceOptions(workspaces) {
  const current = elements.environmentWorkspace.value;
  elements.environmentWorkspace.replaceChildren();
  if (workspaces.length === 0) {
    const option = document.createElement('option');
    option.value = '';
    option.textContent = '先添加 Workspace';
    elements.environmentWorkspace.append(option);
    elements.environmentWorkspace.disabled = true;
    return;
  }
  elements.environmentWorkspace.disabled = false;
  for (const workspace of workspaces) {
    const option = document.createElement('option');
    option.value = workspace.workspace_id || '';
    option.textContent = `${workspace.name || workspace.workspace_id} · ${workspace.path || ''}`;
    elements.environmentWorkspace.append(option);
  }
  if (workspaces.some((workspace) => workspace.workspace_id === current)) {
    elements.environmentWorkspace.value = current;
  }
}

function renderWorkspaces(workspaces) {
  elements.workspaceBadge.textContent = String(workspaces.length);
  renderWorkspaceOptions(workspaces);
  if (workspaces.length === 0) {
    emptyMessage(elements.workspaceList, '暂无 Workspace');
    return;
  }

  elements.workspaceList.replaceChildren();
  elements.workspaceList.classList.remove('empty');
  for (const workspace of workspaces) {
    const item = document.createElement('article');
    item.className = 'list-item managed-item';

    const content = document.createElement('div');
    content.className = 'item-content';
    const title = document.createElement('strong');
    title.textContent = workspace.name || workspace.workspace_id;
    const id = document.createElement('code');
    id.textContent = workspace.workspace_id || '';
    const path = document.createElement('span');
    path.textContent = workspace.path || '';
    content.append(title, id, path);

    const actions = document.createElement('div');
    actions.className = 'item-actions';
    actions.append(
      createActionButton('改名', 'rename-workspace', workspace.workspace_id),
      createActionButton('删除', 'remove-workspace', workspace.workspace_id, 'danger'),
    );

    item.append(content, actions);
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
    item.className = 'list-item managed-item';

    const content = document.createElement('div');
    content.className = 'item-content';
    const title = document.createElement('strong');
    title.textContent = environment.name || environment.environment_id;
    const id = document.createElement('code');
    id.textContent = environment.environment_id || '';
    const root = document.createElement('span');
    root.textContent = environment.root || '';
    const meta = document.createElement('small');
    const memoryCount = Number(environment.private_memory_count || 0);
    const writer = environment.writer?.owner ? `Writer ${environment.writer.owner}` : 'No writer';
    meta.textContent = `Workspace ${environment.workspace_id || '—'} · Private Memory ${memoryCount} · ${writer}`;
    content.append(title, id, root, meta);

    const actions = document.createElement('div');
    actions.className = 'item-actions';
    actions.append(
      createActionButton('详情', 'inspect-environment', environment.environment_id),
      createActionButton('改名', 'rename-environment', environment.environment_id),
      createActionButton('删除', 'remove-environment', environment.environment_id, 'danger'),
    );

    item.append(content, actions);
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
  currentSnapshot = snapshot;
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

  if (selectedEnvironmentID && !environments.some((env) => env.environment_id === selectedEnvironmentID)) {
    closeEnvironmentDetail();
  }
}

function detailRow(label, value) {
  const row = document.createElement('div');
  row.className = 'detail-row';
  const term = document.createElement('strong');
  term.textContent = label;
  const content = document.createElement('span');
  content.textContent = value || '—';
  row.append(term, content);
  return row;
}

function renderEnvironmentDetail(inspection) {
  const environment = inspection.environment || {};
  const workspace = inspection.workspace || {};
  const capabilities = safeArray(inspection.capabilities);
  const enabledMCPs = safeArray(inspection.enabled_mcps);
  const enabledSkills = safeArray(inspection.enabled_skills);
  const unresolvedMCPs = safeArray(inspection.unresolved_mcp_ids);
  const unresolvedSkills = safeArray(inspection.unresolved_skill_ids);

  elements.environmentDetailTitle.textContent = environment.name || environment.environment_id || 'Environment';
  elements.environmentDetail.replaceChildren(
    detailRow('Environment ID', environment.environment_id),
    detailRow('Root', environment.root),
    detailRow('Workspace', `${workspace.name || workspace.workspace_id || '—'} · ${workspace.path || ''}`),
    detailRow('State', environment.state),
    detailRow('Writer', environment.writer?.owner || '无 active writer'),
    detailRow('Private Memory', `${Number(environment.private_memory_count || 0)} entries`),
    detailRow('Capabilities', capabilities.join(', ') || '无'),
    detailRow('MCP', enabledMCPs.map((item) => item.name || item.id).join(', ') || '无'),
    detailRow('Skill', enabledSkills.map((item) => item.name || item.id).join(', ') || '无'),
    detailRow('Unresolved MCP IDs', unresolvedMCPs.join(', ') || '无'),
    detailRow('Unresolved Skill IDs', unresolvedSkills.join(', ') || '无'),
  );
  elements.environmentDetailPanel.hidden = false;
}

function closeEnvironmentDetail() {
  selectedEnvironmentID = '';
  elements.environmentDetailPanel.hidden = true;
  elements.environmentDetail.replaceChildren();
}

async function refreshSnapshot(successMessage = '') {
  elements.refreshButton.disabled = true;
  setStatus('正在读取 ADM 状态…', 'loading');
  try {
    const snapshot = await desktopAdapter().GetSnapshot();
    renderSnapshot(snapshot);
    setStatus(successMessage || `已刷新 · ${new Date().toLocaleTimeString()}`, 'success');
  } catch (error) {
    const message = error?.message || String(error);
    setStatus(`读取失败：${message}`, 'error');
  } finally {
    elements.refreshButton.disabled = false;
  }
}

async function runMutation(label, action) {
  setStatus(`${label}…`, 'loading');
  try {
    await action();
    await refreshSnapshot(`${label}完成`);
  } catch (error) {
    setStatus(`${label}失败：${error?.message || String(error)}`, 'error');
  }
}

elements.workspaceForm.addEventListener('submit', async (event) => {
  event.preventDefault();
  const path = elements.workspacePath.value.trim();
  const name = elements.workspaceName.value.trim();
  if (!path) return;
  await runMutation('添加 Workspace', async () => {
    await desktopAdapter().AddWorkspace({path, name});
    elements.workspaceForm.reset();
  });
});

elements.environmentForm.addEventListener('submit', async (event) => {
  event.preventDefault();
  const workspaceID = elements.environmentWorkspace.value;
  const name = elements.environmentName.value.trim();
  const root = elements.environmentRoot.value.trim();
  if (!workspaceID || !name) return;
  await runMutation('创建 Environment', async () => {
    await desktopAdapter().CreateEnvironment({workspace_id: workspaceID, name, root});
    elements.environmentName.value = '';
    elements.environmentRoot.value = '';
  });
});

elements.workspaceList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]');
  if (!button) return;
  const id = button.dataset.id;
  const workspace = safeArray(currentSnapshot?.workspaces).find((item) => item.workspace_id === id);
  if (button.dataset.action === 'rename-workspace') {
    const name = window.prompt('新的 Workspace 名称', workspace?.name || '');
    if (name === null) return;
    await runMutation('重命名 Workspace', () => desktopAdapter().RenameWorkspace(id, name));
  }
  if (button.dataset.action === 'remove-workspace') {
    if (!window.confirm(`只移除 ADM Workspace 记录，不删除目录。继续？\n${workspace?.path || id}`)) return;
    await runMutation('移除 Workspace', () => desktopAdapter().RemoveWorkspace(id));
  }
});

elements.environmentList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]');
  if (!button) return;
  const id = button.dataset.id;
  const environment = safeArray(currentSnapshot?.environments).find((item) => item.environment_id === id);

  if (button.dataset.action === 'inspect-environment') {
    setStatus('读取 Environment 详情…', 'loading');
    try {
      const inspection = await desktopAdapter().InspectEnvironment(id);
      selectedEnvironmentID = id;
      renderEnvironmentDetail(inspection);
      setStatus('Environment 详情已加载', 'success');
    } catch (error) {
      setStatus(`读取 Environment 详情失败：${error?.message || String(error)}`, 'error');
    }
  }

  if (button.dataset.action === 'rename-environment') {
    const name = window.prompt('新的 Environment 名称', environment?.name || '');
    if (name === null) return;
    await runMutation('重命名 Environment', () => desktopAdapter().RenameEnvironment(id, name));
  }

  if (button.dataset.action === 'remove-environment') {
    if (!window.confirm(`只移除 ADM Environment 记录，不删除 root 或项目文件。继续？\n${environment?.root || id}`)) return;
    await runMutation('移除 Environment', () => desktopAdapter().RemoveEnvironment(id));
  }
});

elements.closeEnvironmentDetail.addEventListener('click', closeEnvironmentDetail);
elements.refreshButton.addEventListener('click', () => refreshSnapshot());
window.addEventListener('DOMContentLoaded', () => refreshSnapshot());
