const elements = {
  refreshButton: document.getElementById('refreshButton'),
  statusPanel: document.getElementById('statusPanel'),
  gatewayState: document.getElementById('gatewayState'),
  gatewayBaseURL: document.getElementById('gatewayBaseURL'),
  gatewayHealthURL: document.getElementById('gatewayHealthURL'),
  gatewayURL: document.getElementById('gatewayURL'),
  gatewayAdminURL: document.getElementById('gatewayAdminURL'),
  gatewayProcess: document.getElementById('gatewayProcess'),
  gatewayDetail: document.getElementById('gatewayDetail'),
  gatewayRefreshButton: document.getElementById('gatewayRefreshButton'),
  gatewayStartButton: document.getElementById('gatewayStartButton'),
  gatewayStopButton: document.getElementById('gatewayStopButton'),
  workspaceCount: document.getElementById('workspaceCount'), environmentCount: document.getElementById('environmentCount'),
  execCount: document.getElementById('execCount'), mcpCount: document.getElementById('mcpCount'), skillCount: document.getElementById('skillCount'), memoryCount: document.getElementById('memoryCount'),
  workspaceBadge: document.getElementById('workspaceBadge'), environmentBadge: document.getElementById('environmentBadge'), execBadge: document.getElementById('execBadge'),
  mcpBadge: document.getElementById('mcpBadge'), skillBadge: document.getElementById('skillBadge'),
  managementEnvironment: document.getElementById('managementEnvironment'), managementEnvironmentHint: document.getElementById('managementEnvironmentHint'),
  runtimeRefreshButton: document.getElementById('runtimeRefreshButton'), runtimeHint: document.getElementById('runtimeHint'), verifierList: document.getElementById('verifierList'), processList: document.getElementById('processList'), runList: document.getElementById('runList'), runtimeOutput: document.getElementById('runtimeOutput'),
  mcpTotalCount: document.getElementById('mcpTotalCount'), mcpDefaultCount: document.getElementById('mcpDefaultCount'), mcpEnvironmentCount: document.getElementById('mcpEnvironmentCount'), mcpIssueCount: document.getElementById('mcpIssueCount'),
  mcpForm: document.getElementById('mcpForm'), mcpName: document.getElementById('mcpName'), mcpTransport: document.getElementById('mcpTransport'), mcpEndpointField: document.getElementById('mcpEndpointField'), mcpEndpoint: document.getElementById('mcpEndpoint'),
  mcpExecutableField: document.getElementById('mcpExecutableField'), mcpExecutable: document.getElementById('mcpExecutable'), mcpArgsField: document.getElementById('mcpArgsField'), mcpArgs: document.getElementById('mcpArgs'),
  mcpAuthField: document.getElementById('mcpAuthField'), mcpAuthMode: document.getElementById('mcpAuthMode'), mcpRefsField: document.getElementById('mcpRefsField'), mcpReferencePairs: document.getElementById('mcpReferencePairs'),
  mcpHealthEnabled: document.getElementById('mcpHealthEnabled'), mcpAutoReconnect: document.getElementById('mcpAutoReconnect'), mcpProbeTimeout: document.getElementById('mcpProbeTimeout'), mcpCheckInterval: document.getElementById('mcpCheckInterval'), mcpDefault: document.getElementById('mcpDefault'),
  mcpImportForm: document.getElementById('mcpImportForm'), mcpImportFormat: document.getElementById('mcpImportFormat'), mcpImportConflict: document.getElementById('mcpImportConflict'), mcpImportDefault: document.getElementById('mcpImportDefault'), mcpImportContent: document.getElementById('mcpImportContent'),
  mcpImportApplyButton: document.getElementById('mcpImportApplyButton'), mcpImportPreview: document.getElementById('mcpImportPreview'), mcpList: document.getElementById('mcpList'),
  skillSourceCount: document.getElementById('skillSourceCount'), skillTotalCount: document.getElementById('skillTotalCount'), skillEnvironmentCount: document.getElementById('skillEnvironmentCount'), skillIssueCount: document.getElementById('skillIssueCount'),
  skillSourceForm: document.getElementById('skillSourceForm'), skillSourceRoot: document.getElementById('skillSourceRoot'), skillSupportRoots: document.getElementById('skillSupportRoots'), skillSourceDefault: document.getElementById('skillSourceDefault'), skillSourceList: document.getElementById('skillSourceList'), skillList: document.getElementById('skillList'),
  workspaceForm: document.getElementById('workspaceForm'), workspacePath: document.getElementById('workspacePath'), workspaceName: document.getElementById('workspaceName'), workspaceList: document.getElementById('workspaceList'),
  environmentForm: document.getElementById('environmentForm'), environmentWorkspace: document.getElementById('environmentWorkspace'), environmentName: document.getElementById('environmentName'), environmentRoot: document.getElementById('environmentRoot'), environmentList: document.getElementById('environmentList'),
  environmentDetailBackdrop: document.getElementById('environmentDetailBackdrop'), environmentDetailPanel: document.getElementById('environmentDetailPanel'), environmentDetailTitle: document.getElementById('environmentDetailTitle'), environmentDetail: document.getElementById('environmentDetail'),
  environmentMCPSelections: document.getElementById('environmentMCPSelections'), environmentSkillSelections: document.getElementById('environmentSkillSelections'), closeEnvironmentDetail: document.getElementById('closeEnvironmentDetail'),
  loadEnvironmentMemory: document.getElementById('loadEnvironmentMemory'), environmentMemoryForm: document.getElementById('environmentMemoryForm'), environmentMemoryKey: document.getElementById('environmentMemoryKey'), environmentMemoryValue: document.getElementById('environmentMemoryValue'), environmentMemoryList: document.getElementById('environmentMemoryList'),
  execForm: document.getElementById('execForm'), execExecutable: document.getElementById('execExecutable'), execList: document.getElementById('execList'),
  loadGlobalMemory: document.getElementById('loadGlobalMemory'), globalMemoryForm: document.getElementById('globalMemoryForm'), globalMemoryKey: document.getElementById('globalMemoryKey'), globalMemoryValue: document.getElementById('globalMemoryValue'), globalMemoryList: document.getElementById('globalMemoryList'),
};

let currentSnapshot = null;
let skillSources = [];
let selectedEnvironmentID = '';
let managementEnvironmentID = '';
let managementInspection = null;
let managementCapabilityFacts = new Map();
let skillAvailabilityByID = new Map();
let mcpHealthByKey = new Map();
let runtimeRunsByID = new Map();
let pendingMCPImport = null;
let globalMemoryLoaded = false;
let environmentMemoryLoaded = false;
let statusTimer = null;

function desktopAdapter() {
  const adapter = window.go?.desktop?.Adapter;
  if (!adapter?.GetSnapshot) throw new Error('Wails desktop binding is not ready');
  return adapter;
}
function safeArray(value) { return Array.isArray(value) ? value : []; }
function safeNumber(value, fallback = 0) { const n = Number(value); return Number.isFinite(n) ? n : fallback; }
function setMetric(element, value) { element.textContent = String(value); }
function emptyMessage(container, message) { container.replaceChildren(); container.classList.add('empty'); container.textContent = message; }
function textOrDash(value) { const text = String(value ?? '').trim(); return text || '—'; }
function lineValues(value) { return String(value || '').split(/\r?\n/).map((item) => item.trim()).filter(Boolean); }
function referenceMap(value) {
  const result = {};
  for (const line of lineValues(value)) {
    const index = line.indexOf('=');
    if (index <= 0 || index === line.length - 1) throw new Error(`引用映射必须使用 KEY=REFERENCE：${line}`);
    result[line.slice(0, index).trim()] = line.slice(index + 1).trim();
  }
  return result;
}
function normalizedState(value) { return String(value || 'unknown').toLowerCase().replace(/[^a-z0-9_-]+/g, '-'); }
function setStatus(message, kind = 'normal') {
  if (statusTimer) clearTimeout(statusTimer);
  elements.statusPanel.textContent = message; elements.statusPanel.dataset.kind = kind; elements.statusPanel.hidden = false;
  if (kind !== 'loading') statusTimer = setTimeout(() => { elements.statusPanel.hidden = true; statusTimer = null; }, kind === 'error' ? 7000 : 2400);
}
function stateBadge(label, state = label) {
  const badge = document.createElement('span'); badge.className = 'resource-badge'; badge.dataset.state = normalizedState(state); badge.textContent = label || 'unknown'; return badge;
}
function createActionButton(label, action, id, kind = 'secondary') {
  const button = document.createElement('button'); button.type = 'button'; button.className = `small-button ${kind === 'danger' ? 'danger-button' : 'secondary-button'}`;
  button.textContent = label; button.dataset.action = action; button.dataset.id = id; return button;
}
function checkControl(labelText, checked, action, id, disabled = false) {
  const label = document.createElement('label'); label.className = 'check-field compact-check';
  const input = document.createElement('input'); input.type = 'checkbox'; input.checked = Boolean(checked); input.disabled = disabled; input.dataset.action = action; input.dataset.id = id;
  const text = document.createElement('span'); text.textContent = labelText; label.append(input, text); return label;
}
function currentEnvironment() { return safeArray(currentSnapshot?.environments).find((item) => item.environment_id === managementEnvironmentID) || null; }
function capabilityMap(report) { const map = new Map(); for (const fact of safeArray(report?.facts)) if (fact?.key) map.set(fact.key, fact); return map; }
function mcpHealthKey(environmentID, mcpID) { return `${environmentID || ''}:${mcpID || ''}`; }

const defaultADMBaseURL = 'http://127.0.0.1:43137';
const admProfileStorageKey = 'adm-v2.desktop.base-url';
function currentADMBaseURL() { return elements.gatewayBaseURL.value.trim() || defaultADMBaseURL; }
function saveADMBaseURL(value) { try { window.localStorage.setItem(admProfileStorageKey, value); } catch (_) {} }
function loadADMBaseURL() { try { return window.localStorage.getItem(admProfileStorageKey) || defaultADMBaseURL; } catch (_) { return defaultADMBaseURL; } }
function renderGatewayStatus(status) {
  const state = status?.state || 'unknown'; elements.gatewayState.textContent = state; elements.gatewayState.dataset.state = state;
  const baseURL = status?.base_url || currentADMBaseURL();
  elements.gatewayBaseURL.value = baseURL; elements.gatewayHealthURL.textContent = status?.health_url || `${baseURL.replace(/\/$/, '')}/healthz`;
  elements.gatewayURL.textContent = status?.agent_mcp_url || status?.mcp_url || `${baseURL.replace(/\/$/, '')}/mcp`;
  elements.gatewayAdminURL.textContent = status?.admin_mcp_url || `${baseURL.replace(/\/$/, '')}/admin/mcp`;
  elements.gatewayProcess.textContent = `PID ${status?.pid || '—'} · Version ${status?.version || '—'}`;
  const localEligible = Boolean(status?.local_bootstrap_eligible);
  elements.gatewayDetail.textContent = status?.detail || (state === 'running' ? 'ADM health check 通过；Desktop 管理数据通过 Admin MCP 读取。' : 'ADM 未连接时 Desktop 不读取或修改本地 state。');
  elements.gatewayStartButton.disabled = !localEligible || state === 'running' || state === 'incompatible';
  elements.gatewayStopButton.disabled = !localEligible || state === 'stopped' || state === 'incompatible' || state === 'unknown';
  saveADMBaseURL(baseURL);
}
async function refreshGatewayStatus(showMessage = false) {
  try { const status = await desktopAdapter().ConnectADM({base_url: currentADMBaseURL()}); renderGatewayStatus(status); if (showMessage) setStatus('ADM 连接状态已刷新', 'success'); return status; }
  catch (error) { elements.gatewayState.textContent = 'unknown'; elements.gatewayState.dataset.state = 'unknown'; elements.gatewayDetail.textContent = error?.message || String(error); elements.gatewayStartButton.disabled = true; elements.gatewayStopButton.disabled = true; if (showMessage) setStatus(`ADM 连接检查失败：${error?.message || String(error)}`, 'error'); throw error; }
}
async function runGatewayAction(label, action) {
  elements.gatewayStartButton.disabled = true; elements.gatewayStopButton.disabled = true; setStatus(`${label}…`, 'loading');
  try {
    const status = await action({base_url: currentADMBaseURL()}); renderGatewayStatus(status);
    if (status?.state === 'running') await refreshSnapshot(`${label}完成`); else { clearManagementData(); setStatus(`${label}完成`, 'success'); }
  }
  catch (error) { clearManagementData(); setStatus(`${label}失败：${error?.message || String(error)}`, 'error'); try { await refreshGatewayStatus(false); } catch (_) {} }
}

function renderWorkspaceOptions(workspaces) {
  const current = elements.environmentWorkspace.value; elements.environmentWorkspace.replaceChildren();
  if (!workspaces.length) { const option = document.createElement('option'); option.value = ''; option.textContent = '先添加 Workspace'; elements.environmentWorkspace.append(option); elements.environmentWorkspace.disabled = true; return; }
  elements.environmentWorkspace.disabled = false;
  for (const workspace of workspaces) { const option = document.createElement('option'); option.value = workspace.workspace_id || ''; option.textContent = `${workspace.name || workspace.workspace_id} · ${workspace.path || ''}`; elements.environmentWorkspace.append(option); }
  if (workspaces.some((workspace) => workspace.workspace_id === current)) elements.environmentWorkspace.value = current;
}
function renderManagementEnvironmentOptions(environments) {
  const previous = managementEnvironmentID; elements.managementEnvironment.replaceChildren();
  const none = document.createElement('option'); none.value = ''; none.textContent = '不选择 Environment（只管理全局定义）'; elements.managementEnvironment.append(none);
  for (const environment of environments) { const option = document.createElement('option'); option.value = environment.environment_id || ''; option.textContent = `${environment.name || environment.environment_id} · ${environment.root || ''}`; elements.managementEnvironment.append(option); }
  if (previous && environments.some((item) => item.environment_id === previous)) managementEnvironmentID = previous;
  else if (environments.length === 1) managementEnvironmentID = environments[0].environment_id;
  else managementEnvironmentID = '';
  elements.managementEnvironment.value = managementEnvironmentID;
  updateManagementHint();
}
function updateManagementHint() {
  const environment = currentEnvironment();
  elements.managementEnvironmentHint.textContent = environment
    ? `当前查看 ${environment.name || environment.environment_id}。Environment 开关只修改这个上下文；全局删除仍影响 catalog。`
    : '未选择 Environment：可以管理全局定义和默认值，但不会显示 Environment 选择、MCP 探测或 Skill 可用性。';
}
function runtimeWriterOwner() { return currentEnvironment()?.writer?.owner || ''; }
function runtimeItem(titleText, idText, state, detailText, actions = []) {
  const item = document.createElement('article'); item.className = 'runtime-item';
  const head = document.createElement('div'); head.className = 'runtime-item-head';
  const title = document.createElement('strong'); title.textContent = titleText || idText || 'runtime item';
  head.append(title, stateBadge(state || 'unknown', state || 'unknown'));
  const id = document.createElement('code'); id.textContent = idText || '';
  const detail = document.createElement('small'); detail.textContent = detailText || '';
  item.append(head, id, detail);
  if (actions.length) { const actionRow = document.createElement('div'); actionRow.className = 'runtime-actions'; actionRow.append(...actions); item.append(actionRow); }
  return item;
}
function renderRuntime(verifiers, processes, runs) {
  const environment = currentEnvironment();
  if (!environment) {
    elements.runtimeHint.textContent = '选择 Management Environment 后可查看当前 Gateway owner 的运行状态。';
    emptyMessage(elements.verifierList, '请选择 Environment'); emptyMessage(elements.processList, '请选择 Environment'); emptyMessage(elements.runList, '请选择 Environment');
    elements.runtimeOutput.textContent = '尚无输出'; runtimeRunsByID = new Map(); return;
  }
  const writerOwner = runtimeWriterOwner();
  elements.runtimeHint.textContent = writerOwner ? `当前 Environment Writer: ${writerOwner}。停止/取消/运行操作必须使用这个现有 Writer。` : '当前 Environment 没有 active writer；可以查看状态和日志，但运行/停止/取消操作不可用。';

  if (!verifiers.length) emptyMessage(elements.verifierList, '没有配置 verifier');
  else {
    elements.verifierList.replaceChildren(); elements.verifierList.classList.remove('empty');
    for (const verifier of verifiers) {
      const run = createActionButton('运行', 'run-verifier', verifier.verifier_id); run.disabled = !writerOwner || verifier.enabled === false;
      const args = safeArray(verifier.args).join(' '); const detail = `${verifier.kind || 'custom'} · ${verifier.executable || ''}${args ? ` ${args}` : ''}${verifier.cwd ? ` · cwd ${verifier.cwd}` : ''}`;
      elements.verifierList.append(runtimeItem(verifier.name || verifier.verifier_id, verifier.verifier_id, verifier.enabled === false ? 'disabled' : 'configured', detail, [run]));
    }
  }

  if (!processes.length) emptyMessage(elements.processList, '当前 Gateway owner 没有 process');
  else {
    elements.processList.replaceChildren(); elements.processList.classList.remove('empty');
    for (const process of processes) {
      const actions = [createActionButton('日志', 'process-logs', process.id)];
      if (process.state === 'running') { const stop = createActionButton('停止', 'stop-process', process.id, 'danger'); stop.disabled = !writerOwner; actions.push(stop); }
      const ports = safeArray(process.listening_ports); const detail = `PID ${process.pid || '—'}${ports.length ? ` · ports ${ports.join(', ')}` : ''}${process.exit_code !== undefined && process.exit_code !== null ? ` · exit ${process.exit_code}` : ''}${process.error_kind ? ` · ${process.error_kind}` : ''}`;
      elements.processList.append(runtimeItem(process.id, process.id, process.state, detail, actions));
    }
  }

  runtimeRunsByID = new Map();
  if (!runs.length) emptyMessage(elements.runList, '当前 Gateway owner 没有 generic run');
  else {
    elements.runList.replaceChildren(); elements.runList.classList.remove('empty');
    for (const run of runs) {
      runtimeRunsByID.set(run.id, run); const actions = [createActionButton('输出', 'run-output', run.id)];
      if (run.state === 'running') { const cancel = createActionButton('取消', 'cancel-run', run.id, 'danger'); cancel.disabled = !writerOwner; actions.push(cancel); }
      const detail = `${run.executable || ''}${safeArray(run.args).length ? ` ${safeArray(run.args).join(' ')}` : ''}${run.cwd ? ` · cwd ${run.cwd}` : ''}${run.exit_code !== undefined && run.exit_code !== null ? ` · exit ${run.exit_code}` : ''}${run.error_kind ? ` · ${run.error_kind}` : ''}`;
      elements.runList.append(runtimeItem(run.id, run.id, run.state, detail, actions));
    }
  }
}
async function refreshRuntimeContext(showMessage = false) {
  const environment = currentEnvironment();
  if (!environment) { renderRuntime([], [], []); return; }
  if (showMessage) setStatus('正在刷新 Runtime 状态…', 'loading');
  const [verifiers, processes, runs] = await Promise.all([
    desktopAdapter().ListVerifiers(environment.environment_id), desktopAdapter().ListProcesses(environment.environment_id), desktopAdapter().ListRuns(environment.environment_id),
  ]);
  renderRuntime(safeArray(verifiers), safeArray(processes), safeArray(runs));
  if (showMessage) setStatus('Runtime 状态已刷新', 'success');
}
function showRuntimeOutput(title, stdout = '', stderr = '', meta = '') {
  const sections = [title]; if (meta) sections.push(meta); if (stdout) sections.push(`STDOUT\n${stdout}`); if (stderr) sections.push(`STDERR\n${stderr}`);
  elements.runtimeOutput.textContent = sections.filter(Boolean).join('\n\n') || '尚无输出';
}
async function runRuntimeAction(label, action) {
  setStatus(`${label}…`, 'loading');
  try { const result = await action(); await refreshRuntimeContext(); setStatus(`${label}完成`, 'success'); return result; }
  catch (error) { setStatus(`${label}失败：${error?.message || String(error)}`, 'error'); return null; }
}
function renderWorkspaces(workspaces) {
  elements.workspaceBadge.textContent = String(workspaces.length); renderWorkspaceOptions(workspaces);
  if (!workspaces.length) return emptyMessage(elements.workspaceList, '暂无 Workspace');
  elements.workspaceList.replaceChildren(); elements.workspaceList.classList.remove('empty');
  for (const workspace of workspaces) {
    const item = document.createElement('article'); item.className = 'list-item managed-item'; const content = document.createElement('div'); content.className = 'item-content';
    const title = document.createElement('strong'); title.textContent = workspace.name || workspace.workspace_id; const id = document.createElement('code'); id.textContent = workspace.workspace_id || ''; const path = document.createElement('span'); path.textContent = workspace.path || ''; content.append(title, id, path);
    const actions = document.createElement('div'); actions.className = 'item-actions'; actions.append(createActionButton('改名', 'rename-workspace', workspace.workspace_id), createActionButton('删除记录', 'remove-workspace', workspace.workspace_id, 'danger')); item.append(content, actions); elements.workspaceList.append(item);
  }
}
function renderEnvironments(environments) {
  elements.environmentBadge.textContent = String(environments.length);
  if (!environments.length) return emptyMessage(elements.environmentList, '暂无 Environment');
  elements.environmentList.replaceChildren(); elements.environmentList.classList.remove('empty');
  for (const environment of environments) {
    const item = document.createElement('article'); item.className = 'list-item managed-item'; const content = document.createElement('div'); content.className = 'item-content';
    const title = document.createElement('strong'); title.textContent = environment.name || environment.environment_id; const id = document.createElement('code'); id.textContent = environment.environment_id || ''; const root = document.createElement('span'); root.textContent = environment.root || '';
    const meta = document.createElement('small'); meta.textContent = `Workspace ${environment.workspace_id || '—'} · Private Memory ${safeNumber(environment.private_memory_count)} · ${environment.writer?.owner ? `Writer ${environment.writer.owner}` : 'No writer'}`; content.append(title, id, root, meta);
    const actions = document.createElement('div'); actions.className = 'item-actions'; actions.append(createActionButton('详情', 'inspect-environment', environment.environment_id), createActionButton('改名', 'rename-environment', environment.environment_id), createActionButton('删除记录', 'remove-environment', environment.environment_id, 'danger')); item.append(content, actions); elements.environmentList.append(item);
  }
}
function renderExecutables(executables) {
  elements.execBadge.textContent = String(executables.length); if (!executables.length) return emptyMessage(elements.execList, '暂无允许的 executable');
  elements.execList.replaceChildren(); elements.execList.classList.remove('empty');
  for (const executable of executables) { const row = document.createElement('div'); row.className = 'managed-row'; const code = document.createElement('code'); code.textContent = executable; row.append(code, createActionButton('移除', 'remove-executable', executable, 'danger')); elements.execList.append(row); }
}

function resourceHeader(titleText, idText, badges) {
  const header = document.createElement('div'); header.className = 'resource-title-line'; const title = document.createElement('strong'); title.textContent = titleText; header.append(title, ...badges);
  const id = document.createElement('code'); id.textContent = idText || ''; return {header, id};
}
function mcpConfigured(entry) { return entry.transport === 'stdio' ? Boolean(entry.executable) : Boolean(entry.endpoint); }
function renderMCPManager(mcps) {
  const environment = currentEnvironment(); const selected = new Set(safeArray(environment?.enabled_mcp_ids));
  const issueStates = new Set(['unavailable', 'degraded', 'unconfigured', 'error']);
  let issues = 0;
  elements.mcpBadge.textContent = String(mcps.length); setMetric(elements.mcpTotalCount, mcps.length); setMetric(elements.mcpDefaultCount, mcps.filter((mcp) => mcp.default_include_in_environment).length); setMetric(elements.mcpEnvironmentCount, environment ? mcps.filter((mcp) => selected.has(mcp.id)).length : 0);
  if (!mcps.length) { setMetric(elements.mcpIssueCount, 0); return emptyMessage(elements.mcpList, '暂无 MCP 定义。可以添加 typed MCP，或先预览再导入 JSON/JSONC。'); }
  elements.mcpList.replaceChildren(); elements.mcpList.classList.remove('empty');
  for (const entry of mcps) {
    const fact = managementCapabilityFacts.get(`mcp/${entry.id}`); const enabled = environment ? selected.has(entry.id) : false; const health = managementEnvironmentID ? mcpHealthByKey.get(mcpHealthKey(managementEnvironmentID, entry.id)) : null;
    const configured = mcpConfigured(entry); const effectiveState = !configured ? 'unconfigured' : environment && !enabled ? 'disabled' : environment ? (health?.state || 'not_observed') : (fact?.state || 'configured');
    if (issueStates.has(normalizedState(effectiveState))) issues++;
    const row = document.createElement('article'); row.className = 'resource-row';
    const main = document.createElement('div'); main.className = 'resource-main'; const badges = [stateBadge(entry.transport || 'streamable-http', 'transport'), stateBadge(effectiveState, effectiveState)]; if (fact?.state && fact.state !== effectiveState) badges.push(stateBadge(`config ${fact.state}`, fact.state)); const {header, id} = resourceHeader(entry.name || entry.id, entry.id, badges);
    const detail = document.createElement('div'); detail.className = 'resource-detail'; detail.textContent = entry.transport === 'stdio' ? `Executable: ${textOrDash(entry.executable)}${safeArray(entry.args).length ? ` · ${entry.args.length} args` : ''}` : `Endpoint: ${textOrDash(entry.endpoint)} · Auth: ${entry.auth_mode || 'none'}`;
    const note = document.createElement('small'); note.textContent = health?.message || fact?.message || (fact?.reason_code ? `Reason: ${fact.reason_code}` : environment && enabled && !health ? '尚未显式探测；刷新页面不会自动连接 MCP。' : '');
    main.append(header, id, detail); if (note.textContent) main.append(note);
    const controls = document.createElement('div'); controls.className = 'resource-actions'; controls.append(checkControl('新环境默认', entry.default_include_in_environment, 'default-mcp', entry.id));
    controls.append(checkControl(environment ? '当前环境启用' : '选择环境后启用', enabled, 'environment-mcp', entry.id, !environment));
    const probe = createActionButton(health ? '重新探测' : '探测', 'probe-mcp', entry.id); probe.disabled = !environment || !enabled || !configured; controls.append(probe, createActionButton('删除全局定义', 'remove-mcp', entry.id, 'danger'));
    row.append(main, controls); elements.mcpList.append(row);
  }
  setMetric(elements.mcpIssueCount, issues);
}
function skillAvailabilityState(entry, environment, selected) {
  const availability = skillAvailabilityByID.get(entry.id); if (!environment) return entry.artifact_path && entry.source_root ? 'configured' : 'unconfigured'; if (availability?.state) return availability.state; return selected ? 'not_observed' : 'disabled';
}
function renderSkillSources() {
  setMetric(elements.skillSourceCount, skillSources.length);
  if (!skillSources.length) return emptyMessage(elements.skillSourceList, '暂无 Skill source。添加显式 source root 后由 Core 扫描 SKILL.md。');
  elements.skillSourceList.replaceChildren(); elements.skillSourceList.classList.remove('empty');
  for (const source of skillSources) {
    const discovered = safeArray(currentSnapshot?.skills).filter((skill) => skill.source_id === source.skill_source_id); const status = source.last_refresh_status || 'pending';
    const row = document.createElement('article'); row.className = 'resource-row'; const main = document.createElement('div'); main.className = 'resource-main'; const {header, id} = resourceHeader(source.root || source.skill_source_id, source.skill_source_id, [stateBadge(status, status), stateBadge(`${discovered.length} skills`, 'count')]);
    const detail = document.createElement('div'); detail.className = 'resource-detail'; detail.textContent = `Support roots: ${safeArray(source.support_roots).length} · 新环境默认: ${source.default_include_in_environment ? '是' : '否'}${source.last_refresh_at ? ` · Last refresh ${new Date(source.last_refresh_at).toLocaleString()}` : ''}`;
    const note = document.createElement('small'); note.textContent = source.last_refresh_error || '刷新只更新这个 source；失败不会污染其他 Skill。'; main.append(header, id, detail, note);
    const controls = document.createElement('div'); controls.className = 'resource-actions'; controls.append(createActionButton('刷新 Source', 'refresh-skill-source', source.skill_source_id), createActionButton('删除 Source', 'remove-skill-source', source.skill_source_id, 'danger')); row.append(main, controls); elements.skillSourceList.append(row);
  }
}
function renderSkillManager(skills) {
  const environment = currentEnvironment(); const selected = new Set(safeArray(environment?.enabled_skill_ids)); let issues = 0;
  elements.skillBadge.textContent = String(skills.length); setMetric(elements.skillTotalCount, skills.length); setMetric(elements.skillEnvironmentCount, environment ? skills.filter((skill) => selected.has(skill.id)).length : 0);
  renderSkillSources();
  if (!skills.length) { setMetric(elements.skillIssueCount, 0); return emptyMessage(elements.skillList, '暂无已发现 Skill。添加并刷新 Skill source 后会显示在这里。'); }
  elements.skillList.replaceChildren(); elements.skillList.classList.remove('empty');
  for (const entry of skills) {
    const enabled = environment ? selected.has(entry.id) : false; const availability = skillAvailabilityByID.get(entry.id); const state = skillAvailabilityState(entry, environment, enabled); if (['unavailable', 'degraded', 'unconfigured'].includes(normalizedState(state))) issues++;
    const fact = managementCapabilityFacts.get(`skill/${entry.id}`); const row = document.createElement('article'); row.className = 'resource-row'; const main = document.createElement('div'); main.className = 'resource-main';
    const {header, id} = resourceHeader(entry.name || entry.id, entry.id, [stateBadge(state, state), entry.source_id ? stateBadge('source-managed', 'source') : stateBadge('legacy', 'legacy')]);
    const detail = document.createElement('div'); detail.className = 'resource-detail'; detail.textContent = `Artifact: ${textOrDash(entry.relative_artifact_path || entry.artifact_path)} · Support roots: ${safeArray(entry.support_roots).length}`;
    const note = document.createElement('small'); note.textContent = availability?.reason || fact?.message || (availability?.missing_support_roots?.length ? `Missing support roots: ${availability.missing_support_roots.join(', ')}` : '');
    main.append(header, id, detail); if (note.textContent) main.append(note);
    const controls = document.createElement('div'); controls.className = 'resource-actions'; controls.append(checkControl('新环境默认', entry.default_include_in_environment, 'default-skill', entry.id)); controls.append(checkControl(environment ? '当前环境启用' : '选择环境后启用', enabled, 'environment-skill', entry.id, !environment));
    if (!entry.source_id) controls.append(createActionButton('删除条目', 'remove-skill', entry.id, 'danger'));
    row.append(main, controls); elements.skillList.append(row);
  }
  setMetric(elements.skillIssueCount, issues);
}
function renderSnapshotBase(snapshot) {
  currentSnapshot = snapshot; const workspaces = safeArray(snapshot.workspaces), environments = safeArray(snapshot.environments), executables = safeArray(snapshot.allowed_executables), mcps = safeArray(snapshot.mcps), skills = safeArray(snapshot.skills);
  setMetric(elements.workspaceCount, workspaces.length); setMetric(elements.environmentCount, environments.length); setMetric(elements.execCount, executables.length); setMetric(elements.mcpCount, mcps.length); setMetric(elements.skillCount, skills.length); setMetric(elements.memoryCount, safeNumber(snapshot.global_memory_count));
  renderWorkspaces(workspaces); renderEnvironments(environments); renderExecutables(executables); renderManagementEnvironmentOptions(environments); renderMCPManager(mcps); renderSkillManager(skills);
  if (selectedEnvironmentID && !environments.some((env) => env.environment_id === selectedEnvironmentID)) closeEnvironmentDetail();
}
function clearManagementData(message = 'ADM 未连接。连接 Admin MCP 后加载管理数据。') {
  skillSources = []; managementInspection = null; managementCapabilityFacts = new Map(); skillAvailabilityByID = new Map(); mcpHealthByKey = new Map(); runtimeRunsByID = new Map();
  renderSnapshotBase({workspaces: [], environments: [], allowed_executables: [], mcps: [], skills: [], global_memory_count: 0});
  emptyMessage(elements.globalMemoryList, message); globalMemoryLoaded = false; closeEnvironmentDetail();
  renderRuntime([], [], []);
}
async function refreshManagementContext() {
  managementInspection = null; managementCapabilityFacts = new Map(); skillAvailabilityByID = new Map(); updateManagementHint();
  if (managementEnvironmentID) {
    const [inspection, availability] = await Promise.all([desktopAdapter().InspectEnvironment(managementEnvironmentID), desktopAdapter().ListEnvironmentSkills(managementEnvironmentID)]);
    managementInspection = inspection; managementCapabilityFacts = capabilityMap(inspection.capability_report); for (const item of safeArray(availability?.skills)) if (item?.skill_id) skillAvailabilityByID.set(item.skill_id, item);
  }
  renderMCPManager(safeArray(currentSnapshot?.mcps)); renderSkillManager(safeArray(currentSnapshot?.skills));
  await refreshRuntimeContext();
}

function renderImportPreview(preview) {
  const candidates = safeArray(preview?.candidates); elements.mcpImportPreview.replaceChildren(); elements.mcpImportPreview.classList.remove('empty');
  if (!candidates.length) { emptyMessage(elements.mcpImportPreview, '没有可导入候选项。'); elements.mcpImportApplyButton.disabled = true; return; }
  let hasErrors = false;
  const head = document.createElement('div'); head.className = 'preview-summary'; head.textContent = `识别格式：${preview.format || 'unknown'} · ${candidates.length} 个候选项。预览不会修改 catalog 或 Environment。`; elements.mcpImportPreview.append(head);
  for (const candidate of candidates) {
    const errors = safeArray(candidate.errors), warnings = safeArray(candidate.warnings), refs = safeArray(candidate.reference_requirements); if (errors.length) hasErrors = true;
    const row = document.createElement('div'); row.className = 'preview-row'; const title = document.createElement('strong'); title.textContent = candidate.name || 'unnamed';
    const badges = document.createElement('div'); badges.className = 'inline-badges'; badges.append(stateBadge(candidate.transport || 'unknown', 'transport'), stateBadge(errors.length ? 'error' : warnings.length ? 'warning' : 'ready', errors.length ? 'unavailable' : warnings.length ? 'degraded' : 'available'));
    const detail = document.createElement('small'); detail.textContent = `${candidate.endpoint_configured ? 'Endpoint configured' : candidate.executable ? `Executable ${candidate.executable}` : 'No endpoint/executable'}${refs.length ? ` · ${refs.length} reference requirements` : ''}`;
    row.append(title, badges, detail);
    for (const issue of [...errors, ...warnings]) { const issueLine = document.createElement('small'); issueLine.className = errors.includes(issue) ? 'error-text' : 'warning-text'; issueLine.textContent = `${issue.field_path ? `${issue.field_path}: ` : ''}${issue.message || issue.kind}`; row.append(issueLine); }
    elements.mcpImportPreview.append(row);
  }
  elements.mcpImportApplyButton.disabled = hasErrors; if (hasErrors) pendingMCPImport = null;
}
function renderImportApplyResult(result) {
  const mutations = safeArray(result?.result?.mutations); elements.mcpImportPreview.replaceChildren(); elements.mcpImportPreview.classList.remove('empty'); const summary = document.createElement('div'); summary.className = 'preview-summary'; summary.textContent = `已应用 ${mutations.length} 个 MCP 定义；现有 Environment 选择未被自动修改。`; elements.mcpImportPreview.append(summary);
  for (const mutation of mutations) { const row = document.createElement('div'); row.className = 'preview-row'; const name = document.createElement('strong'); name.textContent = mutation.definition?.name || mutation.definition?.id || 'MCP'; const detail = document.createElement('small'); detail.textContent = `${mutation.action || 'created'} · ${mutation.definition?.transport || ''}`; row.append(name, detail); elements.mcpImportPreview.append(row); }
}

function detailRow(label, value) { const row = document.createElement('div'); row.className = 'detail-row'; const term = document.createElement('strong'); term.textContent = label; const content = document.createElement('span'); content.textContent = value || '—'; row.append(term, content); return row; }
function renderSelectionList(container, entries, selectedIDs, kind, facts, availabilityMap) {
  if (!entries.length) return emptyMessage(container, `暂无 ${kind.toUpperCase()} 项`); container.replaceChildren(); container.classList.remove('empty'); const selected = new Set(safeArray(selectedIDs));
  for (const entry of entries) {
    const row = document.createElement('label'); row.className = 'selection-row rich-selection'; const configured = kind === 'mcp' ? mcpConfigured(entry) : Boolean(entry.artifact_path && entry.source_root); const enabled = selected.has(entry.id);
    const checkbox = document.createElement('input'); checkbox.type = 'checkbox'; checkbox.checked = enabled; checkbox.disabled = !configured; checkbox.dataset.kind = kind; checkbox.dataset.id = entry.id;
    const text = document.createElement('span'); text.className = 'selection-copy'; const name = document.createElement('strong'); name.textContent = entry.name || entry.id; const id = document.createElement('code'); id.textContent = entry.id || '';
    let state = enabled ? 'configured' : 'disabled', reason = ''; if (kind === 'mcp') { const fact = facts.get(`mcp/${entry.id}`); state = !configured ? 'unconfigured' : !enabled ? 'disabled' : fact?.state || 'configured'; reason = fact?.message || fact?.reason_code || ''; } else { const availability = availabilityMap.get(entry.id); state = availability?.state || (!configured ? 'unconfigured' : !enabled ? 'disabled' : 'configured'); reason = availability?.reason || facts.get(`skill/${entry.id}`)?.message || ''; }
    const meta = document.createElement('small'); meta.textContent = reason || (configured ? '配置可用' : '缺少配置'); text.append(name, id, meta); row.append(checkbox, text, stateBadge(state, state)); container.append(row);
  }
}
async function renderEnvironmentDetailFromInspection(inspection) {
  const environment = inspection.environment || {}, workspace = inspection.workspace || {}, facts = capabilityMap(inspection.capability_report); const availabilityMap = new Map();
  try { const availability = await desktopAdapter().ListEnvironmentSkills(environment.environment_id); for (const item of safeArray(availability?.skills)) if (item?.skill_id) availabilityMap.set(item.skill_id, item); } catch (_) {}
  elements.environmentDetailTitle.textContent = environment.name || environment.environment_id || 'Environment';
  const unavailable = safeArray(inspection.capability_report?.facts).filter((fact) => ['unavailable', 'degraded', 'unconfigured'].includes(normalizedState(fact.state))).length;
  elements.environmentDetail.replaceChildren(detailRow('Environment ID', environment.environment_id), detailRow('Root', environment.root), detailRow('Workspace', `${workspace.name || workspace.workspace_id || '—'} · ${workspace.path || ''}`), detailRow('State', environment.state), detailRow('Writer', environment.writer?.owner || '无 active writer'), detailRow('Private Memory', `${safeNumber(environment.private_memory_count)} entries`), detailRow('Capability issues', `${unavailable} unavailable/degraded/unconfigured`), detailRow('Unresolved MCP IDs', safeArray(inspection.unresolved_mcp_ids).join(', ') || '无'), detailRow('Unresolved Skill IDs', safeArray(inspection.unresolved_skill_ids).join(', ') || '无'));
  renderSelectionList(elements.environmentMCPSelections, safeArray(currentSnapshot?.mcps), environment.enabled_mcp_ids, 'mcp', facts, availabilityMap); renderSelectionList(elements.environmentSkillSelections, safeArray(currentSnapshot?.skills), environment.enabled_skill_ids, 'skill', facts, availabilityMap);
  elements.environmentDetailBackdrop.hidden = false; elements.environmentDetailPanel.hidden = false;
}
function closeEnvironmentDetail() { selectedEnvironmentID = ''; environmentMemoryLoaded = false; elements.environmentDetailBackdrop.hidden = true; elements.environmentDetailPanel.hidden = true; elements.environmentDetail.replaceChildren(); emptyMessage(elements.environmentMemoryList, '尚未加载 private Memory'); }
async function refreshSelectedEnvironmentDetail() { if (!selectedEnvironmentID) return; await renderEnvironmentDetailFromInspection(await desktopAdapter().InspectEnvironment(selectedEnvironmentID)); }

function renderMemory(container, entries, scope) {
  if (!entries.length) return emptyMessage(container, '没有 Memory 条目'); container.replaceChildren(); container.classList.remove('empty');
  for (const entry of entries) { const row = document.createElement('div'); row.className = 'memory-row'; const content = document.createElement('div'); content.className = 'memory-value'; const key = document.createElement('code'); key.textContent = entry.key || ''; const value = document.createElement('pre'); value.textContent = entry.value || ''; content.append(key, value); row.append(content, createActionButton('删除', `delete-${scope}-memory`, entry.key, 'danger')); container.append(row); }
}
async function loadGlobalMemory() { setStatus('正在显式读取 Global Memory…', 'loading'); try { renderMemory(elements.globalMemoryList, safeArray(await desktopAdapter().ListGlobalMemory()), 'global'); globalMemoryLoaded = true; setStatus('Global Memory 已加载', 'success'); } catch (error) { setStatus(`Global Memory 读取失败：${error?.message || String(error)}`, 'error'); } }
async function loadEnvironmentMemory() { if (!selectedEnvironmentID) return; setStatus('正在显式读取 Environment-private Memory…', 'loading'); try { renderMemory(elements.environmentMemoryList, safeArray(await desktopAdapter().ListEnvironmentMemory(selectedEnvironmentID)), 'environment'); environmentMemoryLoaded = true; setStatus('Environment-private Memory 已加载', 'success'); } catch (error) { setStatus(`Environment-private Memory 读取失败：${error?.message || String(error)}`, 'error'); } }

async function refreshSnapshot(successMessage = '') {
  elements.refreshButton.disabled = true; setStatus('正在通过 Admin MCP 读取 ADM 状态…', 'loading');
  try {
    const [snapshot, sources] = await Promise.all([desktopAdapter().GetSnapshot(), desktopAdapter().ListSkillSources()]); skillSources = safeArray(sources); renderSnapshotBase(snapshot); await refreshManagementContext(); if (selectedEnvironmentID) await refreshSelectedEnvironmentDetail(); setStatus(successMessage || `已刷新 · ${new Date().toLocaleTimeString()}`, 'success');
  } catch (error) { clearManagementData(); setStatus(`Admin MCP 读取失败：${error?.message || String(error)}`, 'error'); }
  finally { elements.refreshButton.disabled = false; }
}
async function refreshConnectedADM(showConnectionMessage = false) {
  const status = await refreshGatewayStatus(showConnectionMessage);
  if (status?.state === 'running') await refreshSnapshot();
  else { clearManagementData(); if (!showConnectionMessage) setStatus('ADM 未连接；管理数据未加载', 'error'); }
  return status;
}
async function runMutation(label, action, after) { setStatus(`${label}…`, 'loading'); try { await action(); await refreshSnapshot(`${label}完成`); if (after) await after(); } catch (error) { setStatus(`${label}失败：${error?.message || String(error)}`, 'error'); } }

function syncMCPTransportForm() {
  const stdio = elements.mcpTransport.value === 'stdio'; elements.mcpEndpointField.hidden = stdio; elements.mcpExecutableField.hidden = !stdio; elements.mcpArgsField.hidden = !stdio; elements.mcpAuthField.hidden = stdio; elements.mcpEndpoint.required = !stdio; elements.mcpExecutable.required = stdio;
  elements.mcpRefsField.querySelector('span').textContent = stdio ? 'Environment reference mappings（每行 KEY=${ENV_NAME}）' : 'Header reference mappings（每行 Header=${ENV_NAME}）';
}

elements.mcpTransport.addEventListener('change', syncMCPTransportForm);
elements.managementEnvironment.addEventListener('change', async () => { managementEnvironmentID = elements.managementEnvironment.value; updateManagementHint(); setStatus('正在加载 Environment MCP/Skill/Runtime 状态…', 'loading'); try { await refreshManagementContext(); setStatus('Environment 管理上下文已切换', 'success'); } catch (error) { setStatus(`Environment 状态读取失败：${error?.message || String(error)}`, 'error'); } });
elements.runtimeRefreshButton.addEventListener('click', () => refreshRuntimeContext(true).catch((error) => setStatus(`Runtime 刷新失败：${error?.message || String(error)}`, 'error')));
elements.verifierList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action="run-verifier"]'); if (!button || !managementEnvironmentID) return;
  const owner = runtimeWriterOwner(); if (!owner) return setStatus('当前 Environment 没有 active writer，不能运行 verifier。', 'error');
  const result = await runRuntimeAction('运行 verifier', () => desktopAdapter().RunVerifier(managementEnvironmentID, owner, button.dataset.id));
  if (result) showRuntimeOutput(`Verifier ${button.dataset.id} · ${result.status || 'unknown'}`, result.stdout || '', result.stderr || '', `${result.summary || ''}${result.exit_code !== undefined ? ` · exit ${result.exit_code}` : ''}`);
});
elements.processList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button || !managementEnvironmentID) return;
  if (button.dataset.action === 'process-logs') {
    setStatus('正在读取 process 日志…', 'loading');
    try { const logs = await desktopAdapter().GetProcessLogs(managementEnvironmentID, button.dataset.id); showRuntimeOutput(`Process ${button.dataset.id}`, logs.stdout || '', logs.stderr || '', `${logs.stdout_truncated ? 'stdout truncated ' : ''}${logs.stderr_truncated ? 'stderr truncated' : ''}`.trim()); setStatus('Process 日志已读取', 'success'); }
    catch (error) { setStatus(`Process 日志读取失败：${error?.message || String(error)}`, 'error'); }
  }
  if (button.dataset.action === 'stop-process') {
    const owner = runtimeWriterOwner(); if (!owner) return setStatus('当前 Environment 没有 active writer，不能停止 process。', 'error');
    await runRuntimeAction('停止 process', () => desktopAdapter().StopProcess(managementEnvironmentID, owner, button.dataset.id));
  }
});
elements.runList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button || !managementEnvironmentID) return; const run = runtimeRunsByID.get(button.dataset.id);
  if (button.dataset.action === 'run-output') { if (!run) return; showRuntimeOutput(`Run ${run.id} · ${run.state || 'unknown'}`, run.stdout || '', run.stderr || '', `${run.message || ''}${run.exit_code !== undefined && run.exit_code !== null ? ` · exit ${run.exit_code}` : ''}`); }
  if (button.dataset.action === 'cancel-run') {
    const owner = runtimeWriterOwner(); if (!owner) return setStatus('当前 Environment 没有 active writer，不能取消 run。', 'error');
    await runRuntimeAction('取消 run', () => desktopAdapter().CancelRun(managementEnvironmentID, owner, button.dataset.id));
  }
});

elements.mcpForm.addEventListener('submit', async (event) => {
  event.preventDefault();
  try {
    const transport = elements.mcpTransport.value; const refs = referenceMap(elements.mcpReferencePairs.value); const input = {
      name: elements.mcpName.value.trim(), transport, auth_mode: transport === 'stdio' ? 'none' : elements.mcpAuthMode.value,
      endpoint: transport === 'stdio' ? '' : elements.mcpEndpoint.value.trim(), executable: transport === 'stdio' ? elements.mcpExecutable.value.trim() : '', args: transport === 'stdio' ? lineValues(elements.mcpArgs.value) : [],
      header_refs: transport === 'stdio' ? {} : refs, env_refs: transport === 'stdio' ? refs : {}, default_include_in_environment: elements.mcpDefault.checked,
      health_policy: {health_check_enabled: elements.mcpHealthEnabled.checked, check_interval_seconds: Math.max(1, safeNumber(elements.mcpCheckInterval.value, 30)), probe_timeout_seconds: Math.max(1, safeNumber(elements.mcpProbeTimeout.value, 5)), auto_reconnect: elements.mcpAutoReconnect.checked, reconnect_interval_seconds: Math.max(1, safeNumber(elements.mcpCheckInterval.value, 30))},
    };
    if (!input.name) throw new Error('MCP 名称不能为空');
    await runMutation('添加全局 MCP', async () => { await desktopAdapter().AddMCP(input); elements.mcpForm.reset(); elements.mcpTransport.value = 'streamable-http'; elements.mcpHealthEnabled.checked = true; elements.mcpProbeTimeout.value = '5'; elements.mcpCheckInterval.value = '30'; syncMCPTransportForm(); });
  } catch (error) { setStatus(`添加 MCP 失败：${error?.message || String(error)}`, 'error'); }
});

elements.mcpImportForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const content = elements.mcpImportContent.value; if (!content.trim()) return;
  const input = {format: elements.mcpImportFormat.value, json_or_jsonc: content, conflict_policy: elements.mcpImportConflict.value, default_include: elements.mcpImportDefault.checked}; setStatus('正在预览 MCP 导入…', 'loading');
  try { const preview = await desktopAdapter().PreviewMCPImport(input); pendingMCPImport = input; renderImportPreview(preview); if (!elements.mcpImportApplyButton.disabled) setStatus('导入预览已生成；确认后再应用', 'success'); else setStatus('导入预览包含错误，不能应用', 'error'); }
  catch (error) { pendingMCPImport = null; elements.mcpImportApplyButton.disabled = true; emptyMessage(elements.mcpImportPreview, `预览失败：${error?.message || String(error)}`); setStatus(`MCP 导入预览失败：${error?.message || String(error)}`, 'error'); }
});
elements.mcpImportApplyButton.addEventListener('click', async () => {
  if (!pendingMCPImport) return; elements.mcpImportApplyButton.disabled = true; setStatus('正在应用 MCP 导入…', 'loading');
  try { const result = await desktopAdapter().ApplyMCPImport(pendingMCPImport); pendingMCPImport = null; renderImportApplyResult(result); elements.mcpImportContent.value = ''; await refreshSnapshot('MCP 导入完成'); }
  catch (error) { setStatus(`MCP 导入失败：${error?.message || String(error)}`, 'error'); }
});

elements.skillSourceForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const root = elements.skillSourceRoot.value.trim(); if (!root) return; setStatus('正在添加 Skill source…', 'loading');
  try {
    const source = await desktopAdapter().AddSkillSource({root, support_roots: lineValues(elements.skillSupportRoots.value), default_include_in_environment: elements.skillSourceDefault.checked});
    try { await desktopAdapter().RefreshSkillSource(source.skill_source_id); elements.skillSourceForm.reset(); await refreshSnapshot('Skill source 已添加并刷新'); }
    catch (refreshError) { await refreshSnapshot(); setStatus(`Skill source 已保存，但刷新失败：${refreshError?.message || String(refreshError)}`, 'error'); }
  } catch (error) { setStatus(`添加 Skill source 失败：${error?.message || String(error)}`, 'error'); }
});

elements.mcpList.addEventListener('change', async (event) => {
  const input = event.target.closest('input[data-action]'); if (!input) return;
  if (input.dataset.action === 'default-mcp') await runMutation('更新 MCP 新环境默认值', () => desktopAdapter().SetMCPDefault(input.dataset.id, input.checked));
  if (input.dataset.action === 'environment-mcp' && managementEnvironmentID) await runMutation('更新当前 Environment MCP 选择', () => desktopAdapter().SetEnvironmentMCP(managementEnvironmentID, input.dataset.id, input.checked));
});
elements.mcpList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id; const entry = safeArray(currentSnapshot?.mcps).find((mcp) => mcp.id === id);
  if (button.dataset.action === 'probe-mcp' && managementEnvironmentID) { setStatus(`正在探测 ${entry?.name || id}…`, 'loading'); try { const health = await desktopAdapter().ProbeMCPHealth(managementEnvironmentID, id); mcpHealthByKey.set(mcpHealthKey(managementEnvironmentID, id), health); renderMCPManager(safeArray(currentSnapshot?.mcps)); setStatus(`MCP 探测完成：${health.state}`, health.state === 'healthy' ? 'success' : 'error'); } catch (error) { setStatus(`MCP 探测失败：${error?.message || String(error)}`, 'error'); } }
  if (button.dataset.action === 'remove-mcp' && window.confirm(`这是全局删除，不是只从当前 Environment 禁用。已有 Environment 中的 ID 引用不会被静默改写，删除后可能显示 unresolved。继续？\n${entry?.name || id}`)) await runMutation('删除全局 MCP 定义', () => desktopAdapter().RemoveMCP(id));
});

elements.skillSourceList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id; const source = skillSources.find((item) => item.skill_source_id === id);
  if (button.dataset.action === 'refresh-skill-source') await runMutation('刷新 Skill source', () => desktopAdapter().RefreshSkillSource(id));
  if (button.dataset.action === 'remove-skill-source' && window.confirm(`删除全局 Skill source 及它发现的 Skill 条目。Environment 中已有 Skill ID 不会被静默改写，之后可能显示 unresolved。继续？\n${source?.root || id}`)) await runMutation('删除 Skill source', () => desktopAdapter().RemoveSkillSource(id));
});
elements.skillList.addEventListener('change', async (event) => {
  const input = event.target.closest('input[data-action]'); if (!input) return;
  if (input.dataset.action === 'default-skill') await runMutation('更新 Skill 新环境默认值', () => desktopAdapter().SetSkillDefault(input.dataset.id, input.checked));
  if (input.dataset.action === 'environment-skill' && managementEnvironmentID) await runMutation('更新当前 Environment Skill 选择', () => desktopAdapter().SetEnvironmentSkill(managementEnvironmentID, input.dataset.id, input.checked));
});
elements.skillList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action="remove-skill"]'); if (!button) return; if (window.confirm(`删除这个全局 Skill 条目？Source-managed Skill 应通过 Source 管理。\n${button.dataset.id}`)) await runMutation('删除 Skill 条目', () => desktopAdapter().RemoveSkill(button.dataset.id)); });

elements.workspaceForm.addEventListener('submit', async (event) => { event.preventDefault(); const path = elements.workspacePath.value.trim(), name = elements.workspaceName.value.trim(); if (path) await runMutation('添加 Workspace', async () => { await desktopAdapter().AddWorkspace({path, name}); elements.workspaceForm.reset(); }); });
elements.environmentForm.addEventListener('submit', async (event) => { event.preventDefault(); const workspaceID = elements.environmentWorkspace.value, name = elements.environmentName.value.trim(), root = elements.environmentRoot.value.trim(); if (workspaceID && name) await runMutation('创建 Environment', async () => { await desktopAdapter().CreateEnvironment({workspace_id: workspaceID, name, root}); elements.environmentName.value = ''; elements.environmentRoot.value = ''; }); });
elements.execForm.addEventListener('submit', async (event) => { event.preventDefault(); const executable = elements.execExecutable.value.trim(); if (executable) await runMutation('更新 exec allowlist', async () => { await desktopAdapter().AllowExecutable(executable); elements.execForm.reset(); }); });
elements.globalMemoryForm.addEventListener('submit', async (event) => { event.preventDefault(); const key = elements.globalMemoryKey.value.trim(), value = elements.globalMemoryValue.value; if (key) await runMutation('写入 Global Memory', () => desktopAdapter().WriteGlobalMemory(key, value), async () => { elements.globalMemoryForm.reset(); if (globalMemoryLoaded) await loadGlobalMemory(); }); });
elements.environmentMemoryForm.addEventListener('submit', async (event) => { event.preventDefault(); if (!selectedEnvironmentID) return; const key = elements.environmentMemoryKey.value.trim(), value = elements.environmentMemoryValue.value; if (key) await runMutation('写入 Environment-private Memory', () => desktopAdapter().WriteEnvironmentMemory(selectedEnvironmentID, key, value), async () => { elements.environmentMemoryForm.reset(); if (environmentMemoryLoaded) await loadEnvironmentMemory(); }); });

elements.workspaceList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id, workspace = safeArray(currentSnapshot?.workspaces).find((item) => item.workspace_id === id); if (button.dataset.action === 'rename-workspace') { const name = window.prompt('新的 Workspace 名称', workspace?.name || ''); if (name !== null) await runMutation('重命名 Workspace', () => desktopAdapter().RenameWorkspace(id, name)); } if (button.dataset.action === 'remove-workspace' && window.confirm(`只移除 ADM Workspace 记录，不删除目录。继续？\n${workspace?.path || id}`)) await runMutation('移除 Workspace', () => desktopAdapter().RemoveWorkspace(id)); });
elements.environmentList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id, environment = safeArray(currentSnapshot?.environments).find((item) => item.environment_id === id);
  if (button.dataset.action === 'inspect-environment') { setStatus('读取 Environment 详情…', 'loading'); try { selectedEnvironmentID = id; managementEnvironmentID = id; elements.managementEnvironment.value = id; await refreshManagementContext(); environmentMemoryLoaded = false; emptyMessage(elements.environmentMemoryList, '尚未加载 private Memory'); await renderEnvironmentDetailFromInspection(await desktopAdapter().InspectEnvironment(id)); setStatus('Environment 详情已加载', 'success'); } catch (error) { setStatus(`读取 Environment 详情失败：${error?.message || String(error)}`, 'error'); } }
  if (button.dataset.action === 'rename-environment') { const name = window.prompt('新的 Environment 名称', environment?.name || ''); if (name !== null) await runMutation('重命名 Environment', () => desktopAdapter().RenameEnvironment(id, name)); }
  if (button.dataset.action === 'remove-environment' && window.confirm(`只移除 ADM Environment 记录，不删除 root 或项目文件。继续？\n${environment?.root || id}`)) await runMutation('移除 Environment', () => desktopAdapter().RemoveEnvironment(id));
});
elements.execList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action="remove-executable"]'); if (button) await runMutation('移除 executable', () => desktopAdapter().RemoveExecutable(button.dataset.id)); });
async function handleSelectionChange(event, kind) { const checkbox = event.target.closest('input[data-kind]'); if (!checkbox || !selectedEnvironmentID) return; await runMutation(`更新 Environment ${kind.toUpperCase()} 选择`, () => kind === 'mcp' ? desktopAdapter().SetEnvironmentMCP(selectedEnvironmentID, checkbox.dataset.id, checkbox.checked) : desktopAdapter().SetEnvironmentSkill(selectedEnvironmentID, checkbox.dataset.id, checkbox.checked)); }
elements.environmentMCPSelections.addEventListener('change', (event) => handleSelectionChange(event, 'mcp')); elements.environmentSkillSelections.addEventListener('change', (event) => handleSelectionChange(event, 'skill'));
elements.globalMemoryList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action="delete-global-memory"]'); if (button) await runMutation('删除 Global Memory', () => desktopAdapter().DeleteGlobalMemory(button.dataset.id), loadGlobalMemory); });
elements.environmentMemoryList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action="delete-environment-memory"]'); if (button && selectedEnvironmentID) await runMutation('删除 Environment-private Memory', () => desktopAdapter().DeleteEnvironmentMemory(selectedEnvironmentID, button.dataset.id), loadEnvironmentMemory); });

elements.loadGlobalMemory.addEventListener('click', loadGlobalMemory); elements.loadEnvironmentMemory.addEventListener('click', loadEnvironmentMemory); elements.closeEnvironmentDetail.addEventListener('click', closeEnvironmentDetail); elements.environmentDetailBackdrop.addEventListener('click', closeEnvironmentDetail);
window.addEventListener('keydown', (event) => { if (event.key === 'Escape' && selectedEnvironmentID) closeEnvironmentDetail(); });
elements.gatewayRefreshButton.addEventListener('click', () => refreshConnectedADM(true).catch((error) => { clearManagementData(); setStatus(`ADM 连接检查失败：${error?.message || String(error)}`, 'error'); }));
elements.gatewayBaseURL.addEventListener('change', () => { const value = currentADMBaseURL(); saveADMBaseURL(value); refreshConnectedADM(true).catch((error) => { clearManagementData(); setStatus(`ADM 连接检查失败：${error?.message || String(error)}`, 'error'); }); });
elements.gatewayBaseURL.addEventListener('keydown', (event) => { if (event.key === 'Enter') { event.preventDefault(); elements.gatewayBaseURL.blur(); } });
elements.gatewayStartButton.addEventListener('click', () => runGatewayAction('启动本地 ADM', (input) => desktopAdapter().StartLocalADM(input)));
elements.gatewayStopButton.addEventListener('click', () => runGatewayAction('停止本地 ADM', (input) => desktopAdapter().StopLocalADM(input)));
elements.refreshButton.addEventListener('click', () => refreshConnectedADM(false).catch((error) => { clearManagementData(); setStatus(`ADM 连接检查失败：${error?.message || String(error)}`, 'error'); }));
window.addEventListener('DOMContentLoaded', () => { elements.gatewayBaseURL.value = loadADMBaseURL(); syncMCPTransportForm(); clearManagementData(); refreshConnectedADM(false).catch((error) => { clearManagementData(); setStatus(`ADM 连接检查失败：${error?.message || String(error)}`, 'error'); }); });
