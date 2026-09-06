const elements = {
  refreshButton: document.getElementById('refreshButton'),
  statusPanel: document.getElementById('statusPanel'),
  gatewayState: document.getElementById('gatewayState'),
  gatewayURL: document.getElementById('gatewayURL'),
  gatewayProcess: document.getElementById('gatewayProcess'),
  gatewayDetail: document.getElementById('gatewayDetail'),
  gatewayRefreshButton: document.getElementById('gatewayRefreshButton'),
  gatewayStartButton: document.getElementById('gatewayStartButton'),
  gatewayStopButton: document.getElementById('gatewayStopButton'),
  workspaceCount: document.getElementById('workspaceCount'),
  environmentCount: document.getElementById('environmentCount'),
  execCount: document.getElementById('execCount'),
  mcpCount: document.getElementById('mcpCount'),
  skillCount: document.getElementById('skillCount'),
  memoryCount: document.getElementById('memoryCount'),
  workspaceBadge: document.getElementById('workspaceBadge'),
  environmentBadge: document.getElementById('environmentBadge'),
  execBadge: document.getElementById('execBadge'),
  mcpBadge: document.getElementById('mcpBadge'),
  skillBadge: document.getElementById('skillBadge'),
  workspaceForm: document.getElementById('workspaceForm'),
  workspacePath: document.getElementById('workspacePath'),
  workspaceName: document.getElementById('workspaceName'),
  workspaceList: document.getElementById('workspaceList'),
  environmentForm: document.getElementById('environmentForm'),
  environmentWorkspace: document.getElementById('environmentWorkspace'),
  environmentName: document.getElementById('environmentName'),
  environmentRoot: document.getElementById('environmentRoot'),
  environmentList: document.getElementById('environmentList'),
  environmentDetailBackdrop: document.getElementById('environmentDetailBackdrop'),
  environmentDetailPanel: document.getElementById('environmentDetailPanel'),
  environmentDetailTitle: document.getElementById('environmentDetailTitle'),
  environmentDetail: document.getElementById('environmentDetail'),
  environmentMCPSelections: document.getElementById('environmentMCPSelections'),
  environmentSkillSelections: document.getElementById('environmentSkillSelections'),
  closeEnvironmentDetail: document.getElementById('closeEnvironmentDetail'),
  loadEnvironmentMemory: document.getElementById('loadEnvironmentMemory'),
  environmentMemoryForm: document.getElementById('environmentMemoryForm'),
  environmentMemoryKey: document.getElementById('environmentMemoryKey'),
  environmentMemoryValue: document.getElementById('environmentMemoryValue'),
  environmentMemoryList: document.getElementById('environmentMemoryList'),
  execForm: document.getElementById('execForm'),
  execExecutable: document.getElementById('execExecutable'),
  execList: document.getElementById('execList'),
  loadGlobalMemory: document.getElementById('loadGlobalMemory'),
  globalMemoryForm: document.getElementById('globalMemoryForm'),
  globalMemoryKey: document.getElementById('globalMemoryKey'),
  globalMemoryValue: document.getElementById('globalMemoryValue'),
  globalMemoryList: document.getElementById('globalMemoryList'),
  mcpForm: document.getElementById('mcpForm'),
  mcpName: document.getElementById('mcpName'),
  mcpEndpoint: document.getElementById('mcpEndpoint'),
  mcpDefault: document.getElementById('mcpDefault'),
  mcpList: document.getElementById('mcpList'),
  skillForm: document.getElementById('skillForm'),
  skillName: document.getElementById('skillName'),
  skillInstructions: document.getElementById('skillInstructions'),
  skillDefault: document.getElementById('skillDefault'),
  skillList: document.getElementById('skillList'),
};

let currentSnapshot = null;
let selectedEnvironmentID = '';
let globalMemoryLoaded = false;
let environmentMemoryLoaded = false;
let statusTimer = null;

function desktopAdapter() {
  const adapter = window.go?.desktop?.Adapter;
  if (!adapter?.GetSnapshot) throw new Error('Wails desktop binding is not ready');
  return adapter;
}

function safeArray(value) { return Array.isArray(value) ? value : []; }
function setStatus(message, kind = 'normal') {
  if (statusTimer) clearTimeout(statusTimer);
  elements.statusPanel.textContent = message;
  elements.statusPanel.dataset.kind = kind;
  elements.statusPanel.hidden = false;
  if (kind !== 'loading') {
    statusTimer = setTimeout(() => { elements.statusPanel.hidden = true; statusTimer = null; }, kind === 'error' ? 6000 : 2200);
  }
}
function setMetric(element, value) { element.textContent = String(value); }
function emptyMessage(container, message) { container.replaceChildren(); container.classList.add('empty'); container.textContent = message; }

function renderGatewayStatus(status) {
  const state = status?.state || 'unknown';
  elements.gatewayState.textContent = state;
  elements.gatewayState.dataset.state = state;
  elements.gatewayURL.textContent = status?.mcp_url || 'http://127.0.0.1:41137/mcp';
  elements.gatewayProcess.textContent = `PID ${status?.pid || '—'} · Version ${status?.version || '—'}`;
  elements.gatewayDetail.textContent = status?.detail || '';
  elements.gatewayStartButton.disabled = state === 'running' || state === 'incompatible';
  elements.gatewayStopButton.disabled = state === 'stopped' || state === 'incompatible' || state === 'unknown';
}

async function refreshGatewayStatus(showMessage = false) {
  try {
    const status = await desktopAdapter().GetGatewayStatus();
    renderGatewayStatus(status);
    if (showMessage) setStatus('Gateway 状态已刷新', 'success');
    return status;
  } catch (error) {
    elements.gatewayState.textContent = 'unknown';
    elements.gatewayState.dataset.state = 'unknown';
    elements.gatewayDetail.textContent = error?.message || String(error);
    elements.gatewayStartButton.disabled = true;
    elements.gatewayStopButton.disabled = true;
    if (showMessage) setStatus(`Gateway 状态读取失败：${error?.message || String(error)}`, 'error');
    throw error;
  }
}

async function runGatewayAction(label, action) {
  elements.gatewayStartButton.disabled = true;
  elements.gatewayStopButton.disabled = true;
  setStatus(`${label}…`, 'loading');
  try {
    const status = await action();
    renderGatewayStatus(status);
    setStatus(`${label}完成`, 'success');
  } catch (error) {
    setStatus(`${label}失败：${error?.message || String(error)}`, 'error');
    try { await refreshGatewayStatus(false); } catch (_) {}
  }
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
  if (workspaces.some((workspace) => workspace.workspace_id === current)) elements.environmentWorkspace.value = current;
}

function renderWorkspaces(workspaces) {
  elements.workspaceBadge.textContent = String(workspaces.length);
  renderWorkspaceOptions(workspaces);
  if (workspaces.length === 0) return emptyMessage(elements.workspaceList, '暂无 Workspace');
  elements.workspaceList.replaceChildren();
  elements.workspaceList.classList.remove('empty');
  for (const workspace of workspaces) {
    const item = document.createElement('article');
    item.className = 'list-item managed-item';
    const content = document.createElement('div');
    content.className = 'item-content';
    const title = document.createElement('strong'); title.textContent = workspace.name || workspace.workspace_id;
    const id = document.createElement('code'); id.textContent = workspace.workspace_id || '';
    const path = document.createElement('span'); path.textContent = workspace.path || '';
    content.append(title, id, path);
    const actions = document.createElement('div'); actions.className = 'item-actions';
    actions.append(createActionButton('改名', 'rename-workspace', workspace.workspace_id), createActionButton('删除', 'remove-workspace', workspace.workspace_id, 'danger'));
    item.append(content, actions);
    elements.workspaceList.append(item);
  }
}

function renderEnvironments(environments) {
  elements.environmentBadge.textContent = String(environments.length);
  if (environments.length === 0) return emptyMessage(elements.environmentList, '暂无 Environment');
  elements.environmentList.replaceChildren();
  elements.environmentList.classList.remove('empty');
  for (const environment of environments) {
    const item = document.createElement('article'); item.className = 'list-item managed-item';
    const content = document.createElement('div'); content.className = 'item-content';
    const title = document.createElement('strong'); title.textContent = environment.name || environment.environment_id;
    const id = document.createElement('code'); id.textContent = environment.environment_id || '';
    const root = document.createElement('span'); root.textContent = environment.root || '';
    const meta = document.createElement('small');
    const writer = environment.writer?.owner ? `Writer ${environment.writer.owner}` : 'No writer';
    meta.textContent = `Workspace ${environment.workspace_id || '—'} · Private Memory ${Number(environment.private_memory_count || 0)} · ${writer}`;
    content.append(title, id, root, meta);
    const actions = document.createElement('div'); actions.className = 'item-actions';
    actions.append(createActionButton('详情', 'inspect-environment', environment.environment_id), createActionButton('改名', 'rename-environment', environment.environment_id), createActionButton('删除', 'remove-environment', environment.environment_id, 'danger'));
    item.append(content, actions);
    elements.environmentList.append(item);
  }
}

function renderExecutables(executables) {
  elements.execBadge.textContent = String(executables.length);
  if (executables.length === 0) return emptyMessage(elements.execList, '暂无允许的 executable');
  elements.execList.replaceChildren(); elements.execList.classList.remove('empty');
  for (const executable of executables) {
    const row = document.createElement('div'); row.className = 'managed-row';
    const code = document.createElement('code'); code.textContent = executable;
    row.append(code, createActionButton('移除', 'remove-executable', executable, 'danger'));
    elements.execList.append(row);
  }
}

function renderCatalog(container, badge, entries, kind) {
  badge.textContent = String(entries.length);
  if (entries.length === 0) return emptyMessage(container, `暂无 ${kind.toUpperCase()} catalog 项`);
  container.replaceChildren(); container.classList.remove('empty');
  for (const entry of entries) {
    const row = document.createElement('div'); row.className = 'managed-row';
    const info = document.createElement('div'); info.className = 'item-content';
    const name = document.createElement('strong'); name.textContent = entry.name || entry.id;
    const id = document.createElement('code'); id.textContent = entry.id || '';
    const detail = document.createElement('small');
    detail.textContent = kind === 'mcp'
      ? (entry.endpoint || '未配置 endpoint')
      : (entry.instructions ? entry.instructions.replace(/\s+/g, ' ').slice(0, 140) : '未配置 instructions');
    info.append(name, id, detail);
    const controls = document.createElement('div'); controls.className = 'item-actions';
    const label = document.createElement('label'); label.className = 'check-field compact-check';
    const checkbox = document.createElement('input'); checkbox.type = 'checkbox'; checkbox.checked = Boolean(entry.default_include_in_environment);
    checkbox.dataset.action = `default-${kind}`; checkbox.dataset.id = entry.id;
    const text = document.createElement('span'); text.textContent = '新 Environment 默认启用';
    label.append(checkbox, text);
    controls.append(label, createActionButton('删除', `remove-${kind}`, entry.id, 'danger'));
    row.append(info, controls); container.append(row);
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
  setMetric(elements.workspaceCount, workspaces.length); setMetric(elements.environmentCount, environments.length);
  setMetric(elements.execCount, executables.length); setMetric(elements.mcpCount, mcps.length); setMetric(elements.skillCount, skills.length); setMetric(elements.memoryCount, globalMemoryCount);
  renderWorkspaces(workspaces); renderEnvironments(environments); renderExecutables(executables);
  renderCatalog(elements.mcpList, elements.mcpBadge, mcps, 'mcp'); renderCatalog(elements.skillList, elements.skillBadge, skills, 'skill');
  if (selectedEnvironmentID && !environments.some((env) => env.environment_id === selectedEnvironmentID)) closeEnvironmentDetail();
}

function detailRow(label, value) {
  const row = document.createElement('div'); row.className = 'detail-row';
  const term = document.createElement('strong'); term.textContent = label;
  const content = document.createElement('span'); content.textContent = value || '—';
  row.append(term, content); return row;
}

function renderSelectionList(container, entries, selectedIDs, kind) {
  if (entries.length === 0) return emptyMessage(container, `暂无 ${kind.toUpperCase()} catalog 项`);
  container.replaceChildren(); container.classList.remove('empty');
  const selected = new Set(safeArray(selectedIDs));
  for (const entry of entries) {
    const label = document.createElement('label'); label.className = 'selection-row';
    const configured = kind === 'mcp' ? Boolean(entry.endpoint) : Boolean(entry.instructions);
    const checkbox = document.createElement('input'); checkbox.type = 'checkbox'; checkbox.checked = selected.has(entry.id); checkbox.disabled = !configured;
    checkbox.dataset.kind = kind; checkbox.dataset.id = entry.id;
    const text = document.createElement('span'); text.textContent = `${entry.name || entry.id}${configured ? '' : ' · 未配置'}`;
    const id = document.createElement('code'); id.textContent = entry.id || '';
    label.append(checkbox, text, id); container.append(label);
  }
}

function renderEnvironmentDetail(inspection) {
  const environment = inspection.environment || {};
  const workspace = inspection.workspace || {};
  elements.environmentDetailTitle.textContent = environment.name || environment.environment_id || 'Environment';
  elements.environmentDetail.replaceChildren(
    detailRow('Environment ID', environment.environment_id), detailRow('Root', environment.root),
    detailRow('Workspace', `${workspace.name || workspace.workspace_id || '—'} · ${workspace.path || ''}`), detailRow('State', environment.state),
    detailRow('Writer', environment.writer?.owner || '无 active writer'), detailRow('Private Memory', `${Number(environment.private_memory_count || 0)} entries`),
    detailRow('Capabilities', safeArray(inspection.capabilities).join(', ') || '无'),
    detailRow('Unresolved MCP IDs', safeArray(inspection.unresolved_mcp_ids).join(', ') || '无'), detailRow('Unresolved Skill IDs', safeArray(inspection.unresolved_skill_ids).join(', ') || '无'),
  );
  renderSelectionList(elements.environmentMCPSelections, safeArray(currentSnapshot?.mcps), environment.enabled_mcp_ids, 'mcp');
  renderSelectionList(elements.environmentSkillSelections, safeArray(currentSnapshot?.skills), environment.enabled_skill_ids, 'skill');
  elements.environmentDetailBackdrop.hidden = false;
  elements.environmentDetailPanel.hidden = false;
}

function closeEnvironmentDetail() {
  selectedEnvironmentID = ''; environmentMemoryLoaded = false;
  elements.environmentDetailBackdrop.hidden = true;
  elements.environmentDetailPanel.hidden = true; elements.environmentDetail.replaceChildren();
  emptyMessage(elements.environmentMemoryList, '尚未加载 private Memory');
}

function renderMemory(container, entries, scope) {
  if (entries.length === 0) return emptyMessage(container, '没有 Memory 条目');
  container.replaceChildren(); container.classList.remove('empty');
  for (const entry of entries) {
    const row = document.createElement('div'); row.className = 'memory-row';
    const content = document.createElement('div'); content.className = 'memory-value';
    const key = document.createElement('code'); key.textContent = entry.key || '';
    const value = document.createElement('pre'); value.textContent = entry.value || '';
    content.append(key, value);
    row.append(content, createActionButton('删除', `delete-${scope}-memory`, entry.key, 'danger'));
    container.append(row);
  }
}

async function loadGlobalMemory() {
  setStatus('正在显式读取 Global Memory…', 'loading');
  try { renderMemory(elements.globalMemoryList, safeArray(await desktopAdapter().ListGlobalMemory()), 'global'); globalMemoryLoaded = true; setStatus('Global Memory 已加载', 'success'); }
  catch (error) { setStatus(`Global Memory 读取失败：${error?.message || String(error)}`, 'error'); }
}

async function loadEnvironmentMemory() {
  if (!selectedEnvironmentID) return;
  setStatus('正在显式读取 Environment-private Memory…', 'loading');
  try { renderMemory(elements.environmentMemoryList, safeArray(await desktopAdapter().ListEnvironmentMemory(selectedEnvironmentID)), 'environment'); environmentMemoryLoaded = true; setStatus('Environment-private Memory 已加载', 'success'); }
  catch (error) { setStatus(`Environment-private Memory 读取失败：${error?.message || String(error)}`, 'error'); }
}

async function refreshSelectedEnvironmentDetail() {
  if (!selectedEnvironmentID) return;
  const inspection = await desktopAdapter().InspectEnvironment(selectedEnvironmentID);
  renderEnvironmentDetail(inspection);
}

async function refreshSnapshot(successMessage = '') {
  elements.refreshButton.disabled = true; setStatus('正在读取 ADM 状态…', 'loading');
  try {
    renderSnapshot(await desktopAdapter().GetSnapshot());
    if (selectedEnvironmentID) await refreshSelectedEnvironmentDetail();
    setStatus(successMessage || `已刷新 · ${new Date().toLocaleTimeString()}`, 'success');
  } catch (error) { setStatus(`读取失败：${error?.message || String(error)}`, 'error'); }
  finally { elements.refreshButton.disabled = false; }
}

async function runMutation(label, action, after) {
  setStatus(`${label}…`, 'loading');
  try { await action(); await refreshSnapshot(`${label}完成`); if (after) await after(); }
  catch (error) { setStatus(`${label}失败：${error?.message || String(error)}`, 'error'); }
}

elements.workspaceForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const path = elements.workspacePath.value.trim(); const name = elements.workspaceName.value.trim(); if (!path) return;
  await runMutation('添加 Workspace', async () => { await desktopAdapter().AddWorkspace({path, name}); elements.workspaceForm.reset(); });
});

elements.environmentForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const workspaceID = elements.environmentWorkspace.value; const name = elements.environmentName.value.trim(); const root = elements.environmentRoot.value.trim(); if (!workspaceID || !name) return;
  await runMutation('创建 Environment', async () => { await desktopAdapter().CreateEnvironment({workspace_id: workspaceID, name, root}); elements.environmentName.value = ''; elements.environmentRoot.value = ''; });
});

elements.execForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const executable = elements.execExecutable.value.trim(); if (!executable) return;
  await runMutation('更新 exec allowlist', async () => { await desktopAdapter().AllowExecutable(executable); elements.execForm.reset(); });
});

elements.mcpForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const name = elements.mcpName.value.trim(); const endpoint = elements.mcpEndpoint.value.trim(); if (!name || !endpoint) return;
  await runMutation('添加 MCP', async () => { await desktopAdapter().AddMCP({name, endpoint, default_include_in_environment: elements.mcpDefault.checked}); elements.mcpForm.reset(); });
});

elements.skillForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const name = elements.skillName.value.trim(); const instructions = elements.skillInstructions.value.trim(); if (!name || !instructions) return;
  await runMutation('添加 Skill', async () => { await desktopAdapter().AddSkill({name, instructions, default_include_in_environment: elements.skillDefault.checked}); elements.skillForm.reset(); });
});

elements.globalMemoryForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const key = elements.globalMemoryKey.value.trim(); const value = elements.globalMemoryValue.value; if (!key) return;
  await runMutation('写入 Global Memory', () => desktopAdapter().WriteGlobalMemory(key, value), async () => { elements.globalMemoryForm.reset(); if (globalMemoryLoaded) await loadGlobalMemory(); });
});

elements.environmentMemoryForm.addEventListener('submit', async (event) => {
  event.preventDefault(); if (!selectedEnvironmentID) return; const key = elements.environmentMemoryKey.value.trim(); const value = elements.environmentMemoryValue.value; if (!key) return;
  await runMutation('写入 Environment-private Memory', () => desktopAdapter().WriteEnvironmentMemory(selectedEnvironmentID, key, value), async () => { elements.environmentMemoryForm.reset(); if (environmentMemoryLoaded) await loadEnvironmentMemory(); });
});

elements.workspaceList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id;
  const workspace = safeArray(currentSnapshot?.workspaces).find((item) => item.workspace_id === id);
  if (button.dataset.action === 'rename-workspace') { const name = window.prompt('新的 Workspace 名称', workspace?.name || ''); if (name !== null) await runMutation('重命名 Workspace', () => desktopAdapter().RenameWorkspace(id, name)); }
  if (button.dataset.action === 'remove-workspace' && window.confirm(`只移除 ADM Workspace 记录，不删除目录。继续？\n${workspace?.path || id}`)) await runMutation('移除 Workspace', () => desktopAdapter().RemoveWorkspace(id));
});

elements.environmentList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id;
  const environment = safeArray(currentSnapshot?.environments).find((item) => item.environment_id === id);
  if (button.dataset.action === 'inspect-environment') {
    setStatus('读取 Environment 详情…', 'loading');
    try {
      selectedEnvironmentID = id;
      environmentMemoryLoaded = false;
      emptyMessage(elements.environmentMemoryList, '尚未加载 private Memory');
      renderEnvironmentDetail(await desktopAdapter().InspectEnvironment(id));
      setStatus('Environment 详情已加载', 'success');
    }
    catch (error) { setStatus(`读取 Environment 详情失败：${error?.message || String(error)}`, 'error'); }
  }
  if (button.dataset.action === 'rename-environment') { const name = window.prompt('新的 Environment 名称', environment?.name || ''); if (name !== null) await runMutation('重命名 Environment', () => desktopAdapter().RenameEnvironment(id, name)); }
  if (button.dataset.action === 'remove-environment' && window.confirm(`只移除 ADM Environment 记录，不删除 root 或项目文件。继续？\n${environment?.root || id}`)) await runMutation('移除 Environment', () => desktopAdapter().RemoveEnvironment(id));
});

elements.execList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action="remove-executable"]'); if (!button) return;
  await runMutation('移除 executable', () => desktopAdapter().RemoveExecutable(button.dataset.id));
});

async function handleCatalogEvent(event, kind) {
  const target = event.target; const adapter = desktopAdapter();
  if (event.type === 'change' && target.matches(`input[data-action="default-${kind}"]`)) {
    await runMutation(`更新 ${kind.toUpperCase()} 默认值`, () => kind === 'mcp' ? adapter.SetMCPDefault(target.dataset.id, target.checked) : adapter.SetSkillDefault(target.dataset.id, target.checked));
    return;
  }
  if (event.type !== 'click') return;
  const button = target.closest(`button[data-action="remove-${kind}"]`); if (!button) return;
  await runMutation(`删除 ${kind.toUpperCase()}`, () => kind === 'mcp' ? adapter.RemoveMCP(button.dataset.id) : adapter.RemoveSkill(button.dataset.id));
}

elements.mcpList.addEventListener('change', (event) => handleCatalogEvent(event, 'mcp'));
elements.skillList.addEventListener('change', (event) => handleCatalogEvent(event, 'skill'));
elements.mcpList.addEventListener('click', (event) => handleCatalogEvent(event, 'mcp'));
elements.skillList.addEventListener('click', (event) => handleCatalogEvent(event, 'skill'));

async function handleSelectionChange(event, kind) {
  const checkbox = event.target.closest('input[data-kind]'); if (!checkbox || !selectedEnvironmentID) return;
  const adapter = desktopAdapter();
  await runMutation(`更新 Environment ${kind.toUpperCase()} 选择`, () => kind === 'mcp' ? adapter.SetEnvironmentMCP(selectedEnvironmentID, checkbox.dataset.id, checkbox.checked) : adapter.SetEnvironmentSkill(selectedEnvironmentID, checkbox.dataset.id, checkbox.checked));
}

elements.environmentMCPSelections.addEventListener('change', (event) => handleSelectionChange(event, 'mcp'));
elements.environmentSkillSelections.addEventListener('change', (event) => handleSelectionChange(event, 'skill'));

elements.globalMemoryList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action="delete-global-memory"]'); if (!button) return;
  await runMutation('删除 Global Memory', () => desktopAdapter().DeleteGlobalMemory(button.dataset.id), loadGlobalMemory);
});

elements.environmentMemoryList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action="delete-environment-memory"]'); if (!button || !selectedEnvironmentID) return;
  await runMutation('删除 Environment-private Memory', () => desktopAdapter().DeleteEnvironmentMemory(selectedEnvironmentID, button.dataset.id), loadEnvironmentMemory);
});

elements.loadGlobalMemory.addEventListener('click', loadGlobalMemory);
elements.loadEnvironmentMemory.addEventListener('click', loadEnvironmentMemory);
elements.closeEnvironmentDetail.addEventListener('click', closeEnvironmentDetail);
elements.environmentDetailBackdrop.addEventListener('click', closeEnvironmentDetail);
window.addEventListener('keydown', (event) => { if (event.key === 'Escape' && selectedEnvironmentID) closeEnvironmentDetail(); });
elements.gatewayRefreshButton.addEventListener('click', () => refreshGatewayStatus(true));
elements.gatewayStartButton.addEventListener('click', () => runGatewayAction('启动 Gateway', () => desktopAdapter().StartGateway()));
elements.gatewayStopButton.addEventListener('click', () => runGatewayAction('停止 Gateway', () => desktopAdapter().StopGateway()));
elements.refreshButton.addEventListener('click', () => { refreshSnapshot(); refreshGatewayStatus(false).catch(() => {}); });
window.addEventListener('DOMContentLoaded', () => { refreshSnapshot(); refreshGatewayStatus(false).catch(() => {}); });
