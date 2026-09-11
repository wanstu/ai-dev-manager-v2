const elements = {
  refreshButton: document.getElementById('refreshButton'),
  launchAtLogin: document.getElementById('launchAtLogin'), desktopShellHint: document.getElementById('desktopShellHint'),
  statusPanel: document.getElementById('statusPanel'),
  dashboardDataState: document.getElementById('dashboardDataState'), dashboardLastSuccess: document.getElementById('dashboardLastSuccess'), dashboardDataDetail: document.getElementById('dashboardDataDetail'),
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
  managementEnvironment: document.getElementById('managementEnvironment'), managementEnvironmentHint: document.getElementById('managementEnvironmentHint'), createEnvironmentButton: document.getElementById('createEnvironmentButton'), editEnvironmentButton: document.getElementById('editEnvironmentButton'),
  runtimeRefreshButton: document.getElementById('runtimeRefreshButton'), runtimeHint: document.getElementById('runtimeHint'), verifierList: document.getElementById('verifierList'), processList: document.getElementById('processList'), runList: document.getElementById('runList'), runtimeOutput: document.getElementById('runtimeOutput'),
  mcpTotalCount: document.getElementById('mcpTotalCount'), mcpDefaultCount: document.getElementById('mcpDefaultCount'), mcpEnvironmentCount: document.getElementById('mcpEnvironmentCount'), mcpIssueCount: document.getElementById('mcpIssueCount'),
  mcpForm: document.getElementById('mcpForm'), mcpName: document.getElementById('mcpName'), mcpTransport: document.getElementById('mcpTransport'), mcpEndpointField: document.getElementById('mcpEndpointField'), mcpEndpoint: document.getElementById('mcpEndpoint'),
  mcpExecutableField: document.getElementById('mcpExecutableField'), mcpExecutable: document.getElementById('mcpExecutable'), mcpArgsField: document.getElementById('mcpArgsField'), mcpArgs: document.getElementById('mcpArgs'),
  mcpAuthField: document.getElementById('mcpAuthField'), mcpAuthMode: document.getElementById('mcpAuthMode'), mcpRefsField: document.getElementById('mcpRefsField'), mcpReferencePairs: document.getElementById('mcpReferencePairs'),
  mcpHealthEnabled: document.getElementById('mcpHealthEnabled'), mcpAutoReconnect: document.getElementById('mcpAutoReconnect'), mcpProbeTimeout: document.getElementById('mcpProbeTimeout'), mcpCheckInterval: document.getElementById('mcpCheckInterval'), mcpReconnectInterval: document.getElementById('mcpReconnectInterval'), mcpDefault: document.getElementById('mcpDefault'),
  mcpEditorFlow: document.getElementById('mcpEditorFlow'), mcpEditorSummary: document.getElementById('mcpEditorSummary'), mcpEditorHint: document.getElementById('mcpEditorHint'), mcpSubmitButton: document.getElementById('mcpSubmitButton'), mcpEditCancelButton: document.getElementById('mcpEditCancelButton'),
  mcpImportForm: document.getElementById('mcpImportForm'), mcpImportFormat: document.getElementById('mcpImportFormat'), mcpImportConflict: document.getElementById('mcpImportConflict'), mcpImportDefault: document.getElementById('mcpImportDefault'), mcpImportContent: document.getElementById('mcpImportContent'),
  mcpImportApplyButton: document.getElementById('mcpImportApplyButton'), mcpImportPreview: document.getElementById('mcpImportPreview'), mcpList: document.getElementById('mcpList'),
  mcpFilter: document.getElementById('mcpFilter'), mcpStateFilter: document.getElementById('mcpStateFilter'), mcpVisibleCount: document.getElementById('mcpVisibleCount'), mcpListTotalCount: document.getElementById('mcpListTotalCount'),
  skillSourceCount: document.getElementById('skillSourceCount'), skillTotalCount: document.getElementById('skillTotalCount'), skillEnvironmentCount: document.getElementById('skillEnvironmentCount'), skillIssueCount: document.getElementById('skillIssueCount'),
  skillSourceForm: document.getElementById('skillSourceForm'), skillSourceRoot: document.getElementById('skillSourceRoot'), skillSupportRoots: document.getElementById('skillSupportRoots'), skillSourceDefault: document.getElementById('skillSourceDefault'), skillSourceList: document.getElementById('skillSourceList'), skillList: document.getElementById('skillList'),
  skillFilter: document.getElementById('skillFilter'), skillStateFilter: document.getElementById('skillStateFilter'), skillVisibleCount: document.getElementById('skillVisibleCount'), skillListTotalCount: document.getElementById('skillListTotalCount'),
  skillProbeAllButton: document.getElementById('skillProbeAllButton'), skillSelectVisibleButton: document.getElementById('skillSelectVisibleButton'), skillClearSelectionButton: document.getElementById('skillClearSelectionButton'), skillSelectedCount: document.getElementById('skillSelectedCount'), skillDeleteSelectedButton: document.getElementById('skillDeleteSelectedButton'), skillClearUnavailableButton: document.getElementById('skillClearUnavailableButton'), skillBulkHint: document.getElementById('skillBulkHint'),
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
let editingMCPID = '';
let globalMemoryLoaded = false;
let environmentMemoryLoaded = false;
let statusTimer = null;
let managementNavigation = null;
let environmentDetailOpener = null;
let connectionGeneration = 0;
let environmentGeneration = 0;
let detailGeneration = 0;
let lastSnapshotSuccessAt = 0;
let skillSourcesState = 'unloaded';
let skillSourcesError = '';
let managementContextError = '';
let selectedSkillIDs = new Set();
let explicitSkillAvailabilityProbe = null;
let skillBulkBusy = false;

function desktopAdapter() {
  const adapter = window.go?.desktop?.Adapter;
  if (!adapter?.GetSnapshot) throw new Error('Wails desktop binding is not ready');
  return trackDesktopAdapter(adapter);
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
function referenceLines(values) { return Object.entries(values || {}).sort(([a], [b]) => a.localeCompare(b)).map(([key, value]) => `${key}=${value}`).join('\n'); }
function humanConfigState(state) {
  return ({available: '可用', configured: '已配置', degraded: '降级', unavailable: '不可用', unconfigured: '未配置', disabled: '未启用', unknown: '未知'})[normalizedState(state)] || String(state || '未知');
}
function humanRuntimeState(state) {
  return ({available: '可用', healthy: '健康', error: '错误', degraded: '降级', disabled: '未启用', not_observed: '尚未探测', configured: '已配置', unavailable: '不可用', unconfigured: '未配置', source_missing: 'Source 缺失', artifact_missing: 'Artifact 缺失', artifact_unreadable: 'Artifact 不可读', support_root_missing: 'Support root 缺失', unresolved: '未解析', no_environment: '未选择环境', unknown: '未知'})[normalizedState(state)] || String(state || '未知');
}
function humanRefreshState(state) { return ({ok: '刷新正常', pending: '待刷新', error: '刷新失败'})[normalizedState(state)] || String(state || '未知'); }
function mcpReferenceVariableNames(entry) {
  const refs = entry?.transport === 'stdio' ? entry?.env_refs : entry?.header_refs; const names = new Set();
  for (const value of Object.values(refs || {})) { const text = String(value || ''); for (const match of text.matchAll(/\$\{([A-Za-z_][A-Za-z0-9_]*)/g)) names.add(match[1]); }
  return [...names].sort();
}
function setStatus(message, kind = 'normal') {
  if (statusTimer) clearTimeout(statusTimer);
  const dialogMessage = activeEditorDialog()?.querySelector('.dialog-message');
  if (dialogMessage) { dialogMessage.textContent = message; dialogMessage.hidden = kind !== 'error'; }
  elements.statusPanel.textContent = message; elements.statusPanel.dataset.kind = kind; elements.statusPanel.hidden = false;
  if (kind !== 'loading') statusTimer = setTimeout(() => { elements.statusPanel.hidden = true; statusTimer = null; }, kind === 'error' ? 7000 : 2400);
}
function initializeManagementNavigation() {
  if (!window.ADMNavigation?.createNavigation) throw new Error('management navigation helper is not ready');
  managementNavigation = window.ADMNavigation.createNavigation({
    beforeNavigate: () => !connectionSwitching && !activeEditorDialog() && elements.environmentDetailPanel.hidden,
  });
}

function renderDashboardState(state, error = '') {
  if (!window.ADMDashboard?.renderDashboard) return;
  window.ADMDashboard.renderDashboard({
    state: elements.dashboardDataState,
    lastSuccess: elements.dashboardLastSuccess,
    detail: elements.dashboardDataDetail,
    metrics: {
      workspaceCount: elements.workspaceCount,
      environmentCount: elements.environmentCount,
      execCount: elements.execCount,
      mcpCount: elements.mcpCount,
      skillCount: elements.skillCount,
      memoryCount: elements.memoryCount,
    },
  }, {state, snapshot: currentSnapshot, lastSuccessAt: lastSnapshotSuccessAt, error});
}
function captureEnvironmentScope() { return {connectionGeneration, environmentGeneration, environmentID: managementEnvironmentID}; }
function environmentScopeIsCurrent(scope) { return window.ADMDashboard?.scopeMatches(scope, {connectionGeneration, environmentGeneration, environmentID: managementEnvironmentID}) ?? false; }
function errorText(error) { return error?.message || String(error); }

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
function currentADMBaseURL() { return elements.gatewayBaseURL.value.trim() || defaultADMBaseURL; }
async function loadDesktopPreferences(showMessage = false) {
  try {
    const preferences = await desktopAdapter().GetDesktopPreferences();
    const supported = Boolean(preferences?.launch_at_login_supported);
    elements.launchAtLogin.disabled = !supported;
    elements.launchAtLogin.checked = supported && Boolean(preferences?.launch_at_login);
    elements.desktopShellHint.textContent = supported ? '关闭窗口后继续驻留托盘' : '当前平台不支持开机启动';
    if (showMessage) setStatus('Desktop 设置已刷新', 'success');
    return preferences;
  } catch (error) {
    elements.launchAtLogin.disabled = true;
    elements.desktopShellHint.textContent = `设置读取失败：${error?.message || String(error)}`;
    if (showMessage) setStatus(`Desktop 设置读取失败：${error?.message || String(error)}`, 'error');
    return null;
  }
}
async function updateLaunchAtLogin() {
  const desired = elements.launchAtLogin.checked;
  elements.launchAtLogin.disabled = true;
  elements.desktopShellHint.textContent = desired ? '正在启用开机启动…' : '正在关闭开机启动…';
  try {
    const preferences = await desktopAdapter().SetLaunchAtLogin(desired);
    window.runtime?.EventsEmit?.('desktop:preferences-changed');
    elements.launchAtLogin.checked = Boolean(preferences?.launch_at_login);
    elements.desktopShellHint.textContent = '关闭窗口后继续驻留托盘';
    setStatus(preferences?.launch_at_login ? '已启用开机启动' : '已关闭开机启动', 'success');
  } catch (error) {
    elements.launchAtLogin.checked = !desired;
    elements.desktopShellHint.textContent = `设置失败：${error?.message || String(error)}`;
    setStatus(`开机启动设置失败：${error?.message || String(error)}`, 'error');
  } finally {
    elements.launchAtLogin.disabled = false;
  }
}
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

}
async function refreshGatewayStatus(showMessage = false) {
  if (connectionProfilesLoaded && !connectionProfiles.active_id) throw new Error('请先选择 ADM 连接');
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
  elements.createEnvironmentButton.disabled = !workspaces.length;
  if (!workspaces.length) { const option = document.createElement('option'); option.value = ''; option.textContent = '先添加 Workspace'; elements.environmentWorkspace.append(option); elements.environmentWorkspace.disabled = true; return; }
  elements.environmentWorkspace.disabled = false;
  for (const workspace of workspaces) { const option = document.createElement('option'); option.value = workspace.workspace_id || ''; option.textContent = `${workspace.name || workspace.workspace_id} · ${workspace.path || ''}`; elements.environmentWorkspace.append(option); }
  if (workspaces.some((workspace) => workspace.workspace_id === current)) elements.environmentWorkspace.value = current;
}
function renderManagementEnvironmentOptions(environments) {
  const previous = managementEnvironmentID; elements.managementEnvironment.replaceChildren(); elements.managementEnvironment.disabled = false;
  const none = document.createElement('option'); none.value = ''; none.textContent = '不选择 Environment（只管理全局定义）'; elements.managementEnvironment.append(none);
  for (const environment of environments) { const option = document.createElement('option'); option.value = environment.environment_id || ''; option.textContent = `${environment.name || environment.environment_id} · ${environment.root || ''}`; elements.managementEnvironment.append(option); }
  if (previous && environments.some((item) => item.environment_id === previous)) managementEnvironmentID = previous;
  else if (environments.length === 1) managementEnvironmentID = environments[0].environment_id;
  else managementEnvironmentID = '';
  if (managementEnvironmentID !== previous) { environmentGeneration++; explicitSkillAvailabilityProbe = null; }
  elements.managementEnvironment.value = managementEnvironmentID;
  updateManagementHint();
}
function updateManagementHint() {
  const environment = currentEnvironment();
  elements.editEnvironmentButton.disabled = !environment;
  const base = environment
    ? `当前查看 ${environment.name || environment.environment_id}。可在这里编辑名称；Workspace / Root 仍由 Core 视为固定绑定。Environment 开关只修改这个上下文。`
    : '未选择 Environment：可以新建 Environment 或管理全局定义和默认值，但不会显示 Environment 选择、MCP 探测或 Skill 可用性。';
  elements.managementEnvironmentHint.textContent = managementContextError ? `${base} 部分状态不可用：${managementContextError}` : base;
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
async function refreshRuntimeContext(showMessage = false, scope = captureEnvironmentScope()) {
  const environment = currentEnvironment();
  if (!environment || !scope.environmentID) { if (environmentScopeIsCurrent(scope)) renderRuntime([], [], []); return {stale: false, errors: []}; }
  if (environment.environment_id !== scope.environmentID || !environmentScopeIsCurrent(scope)) return {stale: true, errors: []};
  if (showMessage) setStatus('正在刷新 Runtime 状态…', 'loading');
  const results = await Promise.allSettled([
    desktopAdapter().ListVerifiers(scope.environmentID), desktopAdapter().ListProcesses(scope.environmentID), desktopAdapter().ListRuns(scope.environmentID),
  ]);
  if (!environmentScopeIsCurrent(scope)) return {stale: true, errors: []};
  const [verifierResult, processResult, runResult] = results;
  renderRuntime(verifierResult.status === 'fulfilled' ? safeArray(verifierResult.value) : [], processResult.status === 'fulfilled' ? safeArray(processResult.value) : [], runResult.status === 'fulfilled' ? safeArray(runResult.value) : []);
  const errors = [];
  if (verifierResult.status === 'rejected') { const message = `Verifiers: ${errorText(verifierResult.reason)}`; errors.push(message); emptyMessage(elements.verifierList, `Verifier 状态不可用：${errorText(verifierResult.reason)}`); }
  if (processResult.status === 'rejected') { const message = `Processes: ${errorText(processResult.reason)}`; errors.push(message); emptyMessage(elements.processList, `Process 状态不可用：${errorText(processResult.reason)}`); }
  if (runResult.status === 'rejected') { const message = `Runs: ${errorText(runResult.reason)}`; errors.push(message); emptyMessage(elements.runList, `Run 状态不可用：${errorText(runResult.reason)}`); runtimeRunsByID = new Map(); }
  if (showMessage) setStatus(errors.length ? `Runtime 部分状态不可用：${errors.join(' · ')}` : 'Runtime 状态已刷新', errors.length ? 'error' : 'success');
  return {stale: false, errors};
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
function mcpConfigurationState(entry, fact) {
  if (!mcpConfigured(entry)) return 'unconfigured';
  const observed = normalizedState(fact?.state || 'configured');
  return observed === 'disabled' ? 'configured' : observed;
}
function mcpRuntimeState(environment, enabled, health) {
  if (!environment) return 'no_environment';
  if (!enabled) return 'disabled';
  return health?.state || 'not_observed';
}
function renderMCPManager(mcps) {
  const scrollTop = elements.mcpList.scrollTop;
  try { renderMCPManagerContents(mcps); } finally { elements.mcpList.scrollTop = scrollTop; }
}
function renderMCPManagerContents(mcps) {
  const environment = currentEnvironment(); const selected = new Set(safeArray(environment?.enabled_mcp_ids));
  const configIssueStates = new Set(['unavailable', 'degraded', 'unconfigured', 'error']);
  const runtimeIssueStates = new Set(['unavailable', 'degraded', 'unconfigured', 'error']);
  const query = String(elements.mcpFilter?.value || '').trim().toLowerCase(); const filter = elements.mcpStateFilter?.value || 'all';
  let issues = 0, visible = 0;
  elements.mcpBadge.textContent = String(mcps.length); setMetric(elements.mcpTotalCount, mcps.length); setMetric(elements.mcpListTotalCount, mcps.length); setMetric(elements.mcpDefaultCount, mcps.filter((mcp) => mcp.default_include_in_environment).length); setMetric(elements.mcpEnvironmentCount, environment ? mcps.filter((mcp) => selected.has(mcp.id)).length : 0);
  if (!mcps.length) { setMetric(elements.mcpIssueCount, 0); setMetric(elements.mcpVisibleCount, 0); return emptyMessage(elements.mcpList, '暂无 MCP 定义。可以添加 typed MCP，或先预览再导入 JSON/JSONC。'); }
  elements.mcpList.replaceChildren(); elements.mcpList.classList.remove('empty');
  for (const entry of mcps) {
    const fact = managementCapabilityFacts.get(`mcp/${entry.id}`); const enabled = environment ? selected.has(entry.id) : false; const health = managementEnvironmentID ? mcpHealthByKey.get(mcpHealthKey(managementEnvironmentID, entry.id)) : null;
    const configState = mcpConfigurationState(entry, fact); const runtimeState = mcpRuntimeState(environment, enabled, health); const configNormalized = normalizedState(configState); const runtimeNormalized = normalizedState(runtimeState);
    const issue = configIssueStates.has(configNormalized) || runtimeIssueStates.has(runtimeNormalized); if (issue) issues++;
    const haystack = [entry.name, entry.id, entry.endpoint, entry.executable, entry.transport].filter(Boolean).join(' ').toLowerCase();
    const filterMatch = filter === 'all' || (filter === 'selected' && Boolean(environment && enabled)) || (filter === 'issues' && issue) || (filter === 'unobserved' && runtimeNormalized === 'not_observed');
    if ((query && !haystack.includes(query)) || !filterMatch) continue;
    visible++;
    const row = document.createElement('article'); row.className = 'resource-row'; if (entry.id === editingMCPID) row.dataset.editing = 'true';
    const main = document.createElement('div'); main.className = 'resource-main';
    const badges = [stateBadge(entry.transport || 'streamable-http', 'transport'), stateBadge(`配置 · ${humanConfigState(configState)}`, configState), stateBadge(`运行 · ${humanRuntimeState(runtimeState)}`, runtimeState)];
    const {header, id} = resourceHeader(entry.name || entry.id, entry.id, badges);
    const refCount = entry.transport === 'stdio' ? Object.keys(entry.env_refs || {}).length : Object.keys(entry.header_refs || {}).length;
    const detail = document.createElement('div'); detail.className = 'resource-detail'; detail.textContent = entry.transport === 'stdio' ? `Executable: ${textOrDash(entry.executable)}${safeArray(entry.args).length ? ` · ${entry.args.length} args` : ''}${refCount ? ` · ${refCount} env refs` : ''}` : `Endpoint: ${textOrDash(entry.endpoint)} · Auth: ${entry.auth_mode || 'none'}${refCount ? ` · ${refCount} header refs` : ''}`;
    main.append(header, id, detail);
    const configMessage = fact?.message || (fact?.reason_code ? `Reason: ${fact.reason_code}` : '');
    if (configMessage) { const note = document.createElement('small'); note.textContent = `配置：${configMessage}`; main.append(note); }
    const referenceNames = mcpReferenceVariableNames(entry);
    if (referenceNames.length && ['unavailable', 'unconfigured'].includes(configNormalized)) { const note = document.createElement('small'); note.className = 'reference-text'; note.textContent = `配置引用：${referenceNames.join(', ')}（至少一个当前未解析）`; main.append(note); }
    const runtimeMessage = health?.message || (environment && enabled && !health ? '尚未显式探测；刷新页面不会自动连接 MCP。' : '');
    if (runtimeMessage) { const note = document.createElement('small'); note.textContent = `运行：${runtimeMessage}`; main.append(note); }
    const controls = document.createElement('div'); controls.className = 'resource-actions'; controls.append(checkControl('新环境默认', entry.default_include_in_environment, 'default-mcp', entry.id));
    controls.append(checkControl(environment ? '当前环境启用' : '选择环境后启用', enabled, 'environment-mcp', entry.id, !environment));
    const edit = createActionButton(entry.id === editingMCPID ? '正在编辑' : '编辑', 'edit-mcp', entry.id); edit.disabled = entry.id === editingMCPID;
    const probe = createActionButton(health ? '重新探测' : '探测', 'probe-mcp', entry.id); probe.disabled = !environment || !enabled || ['unavailable', 'unconfigured', 'error'].includes(configNormalized); if (probe.disabled && environment && enabled) probe.title = '先解决配置问题，再进行运行探测';
    controls.append(edit, probe, createActionButton('删除全局定义', 'remove-mcp', entry.id, 'danger'));
    row.append(main, controls); elements.mcpList.append(row);
  }
  setMetric(elements.mcpIssueCount, issues); setMetric(elements.mcpVisibleCount, visible);
  if (!visible) emptyMessage(elements.mcpList, query || filter !== 'all' ? '没有符合当前筛选条件的 MCP。' : '暂无 MCP 定义。');
}
function skillAvailabilityState(entry, environment, selected) {
  const availability = skillAvailabilityByID.get(entry.id); if (!environment) return entry.artifact_path && entry.source_root ? 'configured' : 'unconfigured'; if (availability?.state) return availability.state; return selected ? 'not_observed' : 'disabled';
}
function skillCatalogFingerprint(skills = safeArray(currentSnapshot?.skills)) {
  return window.ADMSkillBulk?.catalogFingerprint(skills) || skills.map((skill) => skill.id || '').filter(Boolean).sort().join('|');
}
function currentSkillProbeIdentity() {
  return {connectionGeneration, environmentGeneration, environmentID: managementEnvironmentID, catalogFingerprint: skillCatalogFingerprint()};
}
function explicitSkillProbeIsCurrent() {
  return Boolean(window.ADMSkillBulk?.probeMatches(explicitSkillAvailabilityProbe, currentSkillProbeIdentity()));
}
function skillAvailabilitySummary() {
  return window.ADMSkillBulk?.summarizeAvailability(safeArray(currentSnapshot?.skills), [...skillAvailabilityByID.values()]) || {available: 0, disabled: 0, unavailable: 0, unknown: 0, cleanupIDs: []};
}
function pruneSkillSelection() {
  const retained = window.ADMSkillBulk?.retainExistingSelection([...selectedSkillIDs], safeArray(currentSnapshot?.skills)) || [];
  selectedSkillIDs = new Set(retained);
}
function updateSkillBulkControls() {
  pruneSkillSelection();
  const environment = currentEnvironment(); const selectedCount = selectedSkillIDs.size; const probeCurrent = explicitSkillProbeIsCurrent(); const summary = probeCurrent ? skillAvailabilitySummary() : null;
  elements.skillSelectedCount.textContent = String(selectedCount);
  elements.skillProbeAllButton.disabled = skillBulkBusy || !environment;
  elements.skillSelectVisibleButton.disabled = skillBulkBusy || !safeArray(currentSnapshot?.skills).length;
  elements.skillClearSelectionButton.disabled = skillBulkBusy || selectedCount === 0;
  elements.skillDeleteSelectedButton.disabled = skillBulkBusy || selectedCount === 0;
  elements.skillClearUnavailableButton.disabled = skillBulkBusy || !environment || !probeCurrent || !summary?.cleanupIDs?.length;
  elements.skillClearUnavailableButton.textContent = probeCurrent && summary?.cleanupIDs?.length ? `一键清除不可用 (${summary.cleanupIDs.length})` : '一键清除不可用';
  if (!environment) elements.skillBulkHint.textContent = '选择 Management Environment 后可显式批量检查可用性；未启用 Skill 不会被当作不可用。';
  else if (!probeCurrent) elements.skillBulkHint.textContent = `当前 Environment：${environment.name || environment.environment_id}。点击“批量检查可用性”后才允许一键清理；自动加载的状态不会授权删除。`;
  else elements.skillBulkHint.textContent = `最近显式检查：可用 ${summary.available} · 不可用 ${summary.unavailable} · 未启用 ${summary.disabled} · 未知 ${summary.unknown}。清理只删除 ADM catalog metadata，不删除磁盘文件。`;
}
async function probeAllSkillAvailability() {
  const scope = captureEnvironmentScope();
  if (!scope.environmentID || !currentEnvironment()) return setStatus('请先选择 Management Environment，再批量检查 Skill 可用性。', 'error');
  const catalogFingerprint = skillCatalogFingerprint(); skillBulkBusy = true; explicitSkillAvailabilityProbe = null; updateSkillBulkControls(); setStatus('正在批量检查当前 Environment 的 Skill 可用性…', 'loading');
  try {
    const result = await desktopAdapter().ListEnvironmentSkills(scope.environmentID);
    if (!environmentScopeIsCurrent(scope) || catalogFingerprint !== skillCatalogFingerprint()) return false;
    skillAvailabilityByID = new Map(); for (const item of safeArray(result?.skills)) if (item?.skill_id) skillAvailabilityByID.set(item.skill_id, item);
    explicitSkillAvailabilityProbe = {...scope, catalogFingerprint, checkedAt: Date.now()};
    renderSkillManager(safeArray(currentSnapshot?.skills));
    const summary = skillAvailabilitySummary();
    setStatus(`Skill 可用性检查完成：可用 ${summary.available} · 不可用 ${summary.unavailable} · 未启用 ${summary.disabled} · 未知 ${summary.unknown}`, summary.unavailable ? 'error' : 'success');
    return true;
  } catch (error) {
    if (environmentScopeIsCurrent(scope)) { explicitSkillAvailabilityProbe = null; updateSkillBulkControls(); setStatus(`Skill 可用性检查失败：${errorText(error)}`, 'error'); }
    return false;
  } finally { skillBulkBusy = false; renderSkillManager(safeArray(currentSnapshot?.skills)); }
}
async function removeSkillCatalogEntries(ids, label) {
  const entriesByID = new Map(safeArray(currentSnapshot?.skills).map((entry) => [entry.id, entry]));
  const targets = [...new Set(safeArray(ids))].map((id) => entriesByID.get(id)).filter(Boolean);
  if (!targets.length) return;
  const sourceManaged = targets.filter((entry) => entry.source_id).length; const legacy = targets.length - sourceManaged;
  const confirmation = `${label} ${targets.length} 个 Skill？\n\n只删除 ADM catalog metadata，不删除磁盘文件。已有 Environment 中的 Skill ID 引用不会被静默改写，删除后可能显示 unresolved。${sourceManaged ? `\n其中 ${sourceManaged} 个是 source-managed，之后刷新 Source 可能重新发现并恢复。` : ''}\nLegacy: ${legacy} · Source-managed: ${sourceManaged}`;
  if (!window.confirm(confirmation)) return;
  const requestGeneration = connectionGeneration; skillBulkBusy = true; explicitSkillAvailabilityProbe = null; updateSkillBulkControls(); setStatus(`${label}…`, 'loading');
  const failures = []; let removed = 0;
  try {
    for (const entry of targets) {
      if (requestGeneration !== connectionGeneration) { failures.push({id: entry.id, error: 'connection changed before deletion'}); break; }
      try { await desktopAdapter().RemoveSkill(entry.id); removed++; selectedSkillIDs.delete(entry.id); }
      catch (error) { failures.push({id: entry.id, error: errorText(error)}); }
    }
    if (requestGeneration === connectionGeneration) await refreshSnapshot();
    const summary = failures.length ? `${label}完成：已删除 ${removed}/${targets.length}；失败 ${failures.length}。${failures.slice(0, 3).map((item) => `${item.id}: ${item.error}`).join(' · ')}` : `${label}完成：已删除 ${removed} 个 Skill catalog 条目。`;
    setStatus(summary, failures.length ? 'error' : 'success');
  } finally { skillBulkBusy = false; renderSkillManager(safeArray(currentSnapshot?.skills)); }
}
function renderSkillSources() {
  if (skillSourcesState === 'loading') { setMetric(elements.skillSourceCount, '—'); return emptyMessage(elements.skillSourceList, '正在读取 Skill sources…'); }
  if (skillSourcesState === 'error') { setMetric(elements.skillSourceCount, '—'); return emptyMessage(elements.skillSourceList, `Skill source 列表不可用：${skillSourcesError || 'unknown error'}`); }
  if (skillSourcesState !== 'success') { setMetric(elements.skillSourceCount, '—'); return emptyMessage(elements.skillSourceList, 'Skill sources 尚未加载'); }
  setMetric(elements.skillSourceCount, skillSources.length);
  if (!skillSources.length) {
    const legacyCount = safeArray(currentSnapshot?.skills).filter((skill) => !skill.source_id).length;
    return emptyMessage(elements.skillSourceList, legacyCount ? `暂无 Skill source。下方 ${legacyCount} 个 Skill 是历史 legacy 条目；添加 source root 后新发现项会标记为 source-managed。` : '暂无 Skill source。添加显式 source root 后由 Core 扫描 SKILL.md。');
  }
  elements.skillSourceList.replaceChildren(); elements.skillSourceList.classList.remove('empty');
  for (const source of skillSources) {
    const discovered = safeArray(currentSnapshot?.skills).filter((skill) => skill.source_id === source.skill_source_id); const status = source.last_refresh_status || 'pending';
    const row = document.createElement('article'); row.className = 'resource-row'; const main = document.createElement('div'); main.className = 'resource-main'; const {header, id} = resourceHeader(source.root || source.skill_source_id, source.skill_source_id, [stateBadge(humanRefreshState(status), status), stateBadge(`${discovered.length} skills`, 'count')]);
    const detail = document.createElement('div'); detail.className = 'resource-detail'; detail.textContent = `Support roots: ${safeArray(source.support_roots).length} · 新环境默认: ${source.default_include_in_environment ? '是' : '否'}${source.last_refresh_at ? ` · Last refresh ${new Date(source.last_refresh_at).toLocaleString()}` : ''}`;
    const note = document.createElement('small'); note.textContent = source.last_refresh_error || '刷新只更新这个 source；失败不会污染其他 Skill。'; main.append(header, id, detail, note);
    const controls = document.createElement('div'); controls.className = 'resource-actions'; controls.append(createActionButton('刷新 Source', 'refresh-skill-source', source.skill_source_id), createActionButton('删除 Source', 'remove-skill-source', source.skill_source_id, 'danger')); row.append(main, controls); elements.skillSourceList.append(row);
  }
}
function renderSkillManager(skills) {
  pruneSkillSelection();
  const environment = currentEnvironment(); const selected = new Set(safeArray(environment?.enabled_skill_ids));
  const query = String(elements.skillFilter?.value || '').trim().toLowerCase(); const filter = elements.skillStateFilter?.value || 'all';
  let issues = 0, visible = 0;
  elements.skillBadge.textContent = String(skills.length); setMetric(elements.skillTotalCount, skills.length); setMetric(elements.skillListTotalCount, skills.length); setMetric(elements.skillEnvironmentCount, environment ? skills.filter((skill) => selected.has(skill.id)).length : 0);
  renderSkillSources();
  if (!skills.length) { setMetric(elements.skillIssueCount, 0); setMetric(elements.skillVisibleCount, 0); emptyMessage(elements.skillList, '暂无已发现 Skill。添加并刷新 Skill source 后会显示在这里。'); updateSkillBulkControls(); return; }
  elements.skillList.replaceChildren(); elements.skillList.classList.remove('empty');
  for (const entry of skills) {
    const enabled = environment ? selected.has(entry.id) : false; const availability = skillAvailabilityByID.get(entry.id); const state = skillAvailabilityState(entry, environment, enabled); const normalized = normalizedState(state); const issue = normalized === 'degraded' || Boolean(window.ADMSkillBulk?.isCleanupState(normalized)) || ['unavailable', 'unconfigured', 'error'].includes(normalized); if (issue) issues++;
    const haystack = [entry.name, entry.id, entry.relative_artifact_path, entry.artifact_path, entry.source_root].filter(Boolean).join(' ').toLowerCase();
    const filterMatch = filter === 'all' || (filter === 'selected' && Boolean(environment && enabled)) || (filter === 'issues' && issue) || (filter === 'source' && Boolean(entry.source_id)) || (filter === 'legacy' && !entry.source_id);
    if ((query && !haystack.includes(query)) || !filterMatch) continue;
    visible++;
    const fact = managementCapabilityFacts.get(`skill/${entry.id}`); const row = document.createElement('article'); row.className = 'resource-row'; const main = document.createElement('div'); main.className = 'resource-main';
    const {header, id} = resourceHeader(entry.name || entry.id, entry.id, [stateBadge(`可用性 · ${humanRuntimeState(state)}`, state), entry.source_id ? stateBadge('Source 管理', 'source') : stateBadge('Legacy', 'legacy')]);
    const detail = document.createElement('div'); detail.className = 'resource-detail'; detail.textContent = `Artifact: ${textOrDash(entry.relative_artifact_path || entry.artifact_path)} · Support roots: ${safeArray(entry.support_roots).length}`;
    const note = document.createElement('small'); const availabilityNote = availability?.reason || fact?.message || (availability?.missing_support_roots?.length ? `Missing support roots: ${availability.missing_support_roots.join(', ')}` : ''); note.textContent = availabilityNote ? `可用性：${availabilityNote}` : '';
    main.append(header, id, detail); if (note.textContent) main.append(note);
    const controls = document.createElement('div'); controls.className = 'resource-actions';
    const selectControl = checkControl('选择', selectedSkillIDs.has(entry.id), 'select-skill', entry.id, skillBulkBusy); selectControl.classList.add('skill-select-control');
    controls.append(selectControl, checkControl('新环境默认', entry.default_include_in_environment, 'default-skill', entry.id)); controls.append(checkControl(environment ? '当前环境启用' : '选择环境后启用', enabled, 'environment-skill', entry.id, !environment));
    if (!entry.source_id) controls.append(createActionButton('删除条目', 'remove-skill', entry.id, 'danger'));
    row.append(main, controls); elements.skillList.append(row);
  }
  setMetric(elements.skillIssueCount, issues); setMetric(elements.skillVisibleCount, visible);
  if (!visible) emptyMessage(elements.skillList, query || filter !== 'all' ? '没有符合当前筛选条件的 Skill。' : '暂无已发现 Skill。');
  updateSkillBulkControls();
}
function renderSnapshotBase(snapshot) {
  currentSnapshot = snapshot; const workspaces = safeArray(snapshot.workspaces), environments = safeArray(snapshot.environments), executables = safeArray(snapshot.allowed_executables), mcps = safeArray(snapshot.mcps), skills = safeArray(snapshot.skills);
  if (editingMCPID && !mcps.some((mcp) => mcp.id === editingMCPID)) resetMCPEditor(false, false);
  renderWorkspaces(workspaces); renderEnvironments(environments); renderExecutables(executables); renderManagementEnvironmentOptions(environments); renderMCPManager(mcps); renderSkillManager(skills);
  renderDashboardState('success');
  if (selectedEnvironmentID && !environments.some((env) => env.environment_id === selectedEnvironmentID)) closeEnvironmentDetail();
}
function renderManagementUnavailable(message) {
  for (const element of [elements.workspaceBadge, elements.environmentBadge, elements.execBadge, elements.mcpBadge, elements.skillBadge, elements.mcpTotalCount, elements.mcpDefaultCount, elements.mcpEnvironmentCount, elements.mcpIssueCount, elements.mcpVisibleCount, elements.mcpListTotalCount, elements.skillSourceCount, elements.skillTotalCount, elements.skillEnvironmentCount, elements.skillIssueCount, elements.skillVisibleCount, elements.skillListTotalCount]) setMetric(element, '—');
  renderWorkspaceOptions([]);
  elements.managementEnvironment.replaceChildren(new Option('管理数据未加载', '')); elements.managementEnvironment.value = ''; elements.managementEnvironment.disabled = true; elements.editEnvironmentButton.disabled = true;
  elements.managementEnvironmentHint.textContent = message;
  emptyMessage(elements.workspaceList, message); emptyMessage(elements.environmentList, message); emptyMessage(elements.execList, message); emptyMessage(elements.mcpList, message); emptyMessage(elements.skillSourceList, message); emptyMessage(elements.skillList, message); emptyMessage(elements.globalMemoryList, message);
  elements.runtimeHint.textContent = message; emptyMessage(elements.verifierList, message); emptyMessage(elements.processList, message); emptyMessage(elements.runList, message); elements.runtimeOutput.textContent = '尚无输出'; runtimeRunsByID = new Map();
  updateSkillBulkControls();
}
function clearManagementData(message = 'ADM 未连接。连接 Admin MCP 后加载管理数据。') {
  connectionGeneration++; environmentGeneration++; detailGeneration++;
  resetMCPEditor(true, false);
  pendingMCPImport = null; elements.mcpImportApplyButton.disabled = true;
  elements.mcpImportContent.value = ''; emptyMessage(elements.mcpImportPreview, '尚未预览。导入不会修改 Environment 选择。');
  managementEnvironmentID = ''; currentSnapshot = null; lastSnapshotSuccessAt = 0;
  skillSources = []; skillSourcesState = 'unloaded'; skillSourcesError = ''; managementContextError = '';
  selectedSkillIDs = new Set(); explicitSkillAvailabilityProbe = null; skillBulkBusy = false;
  managementInspection = null; managementCapabilityFacts = new Map(); skillAvailabilityByID = new Map(); mcpHealthByKey = new Map(); runtimeRunsByID = new Map();
  renderDashboardState('unloaded', message); renderManagementUnavailable(message);
  globalMemoryLoaded = false; if (!elements.environmentDetailPanel.hidden) closeEnvironmentDetail();
}
async function refreshManagementContext(scope = captureEnvironmentScope()) {
  if (!environmentScopeIsCurrent(scope)) return {stale: true, errors: []};
  managementInspection = null; managementCapabilityFacts = new Map(); skillAvailabilityByID = new Map(); managementContextError = ''; updateManagementHint();
  if (!scope.environmentID) { renderMCPManager(safeArray(currentSnapshot?.mcps)); renderSkillManager(safeArray(currentSnapshot?.skills)); renderRuntime([], [], []); return {stale: false, errors: []}; }
  const [inspectionResult, availabilityResult] = await Promise.allSettled([
    desktopAdapter().InspectEnvironment(scope.environmentID),
    desktopAdapter().ListEnvironmentSkills(scope.environmentID),
  ]);
  if (!environmentScopeIsCurrent(scope)) return {stale: true, errors: []};
  const errors = [];
  if (inspectionResult.status === 'fulfilled') { managementInspection = inspectionResult.value; managementCapabilityFacts = capabilityMap(inspectionResult.value?.capability_report); }
  else errors.push(`Environment inspection: ${errorText(inspectionResult.reason)}`);
  if (availabilityResult.status === 'fulfilled') { for (const item of safeArray(availabilityResult.value?.skills)) if (item?.skill_id) skillAvailabilityByID.set(item.skill_id, item); }
  else errors.push(`Skill availability: ${errorText(availabilityResult.reason)}`);
  managementContextError = errors.join(' · '); updateManagementHint();
  renderMCPManager(safeArray(currentSnapshot?.mcps)); renderSkillManager(safeArray(currentSnapshot?.skills));
  const runtimeResult = await refreshRuntimeContext(false, scope);
  return {stale: Boolean(runtimeResult?.stale), errors: [...errors, ...safeArray(runtimeResult?.errors)]};
}

function renderImportPreview(preview) {
  const candidates = safeArray(preview?.candidates); elements.mcpImportPreview.replaceChildren(); elements.mcpImportPreview.classList.remove('empty');
  if (!candidates.length) { emptyMessage(elements.mcpImportPreview, '没有可导入候选项。'); elements.mcpImportApplyButton.disabled = true; return; }
  let hasErrors = false;
  const head = document.createElement('div'); head.className = 'preview-summary'; head.textContent = `识别格式：${preview.format || 'unknown'} · ${candidates.length} 个候选项。预览不会修改 catalog 或 Environment；源配置中的 literal env/header 值不会被保存。`; elements.mcpImportPreview.append(head);
  for (const candidate of candidates) {
    const errors = safeArray(candidate.errors), warnings = safeArray(candidate.warnings), refs = safeArray(candidate.reference_requirements); if (errors.length) hasErrors = true;
    const row = document.createElement('div'); row.className = 'preview-row'; const title = document.createElement('strong'); title.textContent = candidate.name || 'unnamed';
    const badges = document.createElement('div'); badges.className = 'inline-badges'; badges.append(stateBadge(candidate.transport || 'unknown', 'transport'), stateBadge(errors.length ? 'error' : warnings.length ? 'warning' : 'ready', errors.length ? 'unavailable' : warnings.length ? 'degraded' : 'available'));
    const detail = document.createElement('small'); detail.textContent = `${candidate.endpoint_configured ? 'Endpoint configured' : candidate.executable ? `Executable ${candidate.executable}` : 'No endpoint/executable'}${refs.length ? ` · ${refs.length} 个环境变量引用待配置` : ''}`;
    row.append(title, badges, detail);
    for (const ref of refs) { const refLine = document.createElement('small'); refLine.className = 'reference-text'; refLine.textContent = `${ref.field_path || 'reference'} → ${ref.reference_name}`; row.append(refLine); }
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
async function renderEnvironmentDetailFromInspection(inspection, token = detailGeneration) {
  const environment = inspection.environment || {}, workspace = inspection.workspace || {}, facts = capabilityMap(inspection.capability_report); const availabilityMap = new Map();
  try { const availability = await desktopAdapter().ListEnvironmentSkills(environment.environment_id); for (const item of safeArray(availability?.skills)) if (item?.skill_id) availabilityMap.set(item.skill_id, item); } catch (_) {}
  if (token !== detailGeneration || environment.environment_id !== selectedEnvironmentID) return false;
  elements.environmentDetailTitle.textContent = environment.name || environment.environment_id || 'Environment';
  const unavailable = safeArray(inspection.capability_report?.facts).filter((fact) => ['unavailable', 'degraded', 'unconfigured'].includes(normalizedState(fact.state))).length;
  elements.environmentDetail.replaceChildren(detailRow('Environment ID', environment.environment_id), detailRow('Root', environment.root), detailRow('Workspace', `${workspace.name || workspace.workspace_id || '—'} · ${workspace.path || ''}`), detailRow('State', environment.state), detailRow('Writer', environment.writer?.owner || '无 active writer'), detailRow('Private Memory', `${safeNumber(environment.private_memory_count)} entries`), detailRow('Capability issues', `${unavailable} unavailable/degraded/unconfigured`), detailRow('Unresolved MCP IDs', safeArray(inspection.unresolved_mcp_ids).join(', ') || '无'), detailRow('Unresolved Skill IDs', safeArray(inspection.unresolved_skill_ids).join(', ') || '无'));
  renderSelectionList(elements.environmentMCPSelections, safeArray(currentSnapshot?.mcps), environment.enabled_mcp_ids, 'mcp', facts, availabilityMap); renderSelectionList(elements.environmentSkillSelections, safeArray(currentSnapshot?.skills), environment.enabled_skill_ids, 'skill', facts, availabilityMap);
  elements.environmentDetailBackdrop.hidden = false; elements.environmentDetailPanel.hidden = false; return true;
}
function closeEnvironmentDetail() { const wasOpen = !elements.environmentDetailPanel.hidden; const opener = environmentDetailOpener; detailGeneration++; selectedEnvironmentID = ''; environmentDetailOpener = null; environmentMemoryLoaded = false; elements.environmentDetailBackdrop.hidden = true; elements.environmentDetailPanel.hidden = true; elements.environmentDetail.replaceChildren(); emptyMessage(elements.environmentMemoryList, '尚未加载 private Memory'); if (!wasOpen) return; if (opener?.isConnected && !opener.closest('[hidden]')) opener.focus({preventScroll: true}); else document.querySelector('[data-management-page]:not([hidden]) [data-page-heading]')?.focus({preventScroll: true}); }
async function refreshSelectedEnvironmentDetail() { if (!selectedEnvironmentID) return false; const token = detailGeneration, id = selectedEnvironmentID; const inspection = await desktopAdapter().InspectEnvironment(id); if (token !== detailGeneration || id !== selectedEnvironmentID) return false; return renderEnvironmentDetailFromInspection(inspection, token); }

function renderMemory(container, entries, scope) {
  if (!entries.length) return emptyMessage(container, '没有 Memory 条目'); container.replaceChildren(); container.classList.remove('empty');
  for (const entry of entries) { const row = document.createElement('div'); row.className = 'memory-row'; const content = document.createElement('div'); content.className = 'memory-value'; const key = document.createElement('code'); key.textContent = entry.key || ''; const value = document.createElement('pre'); value.textContent = entry.value || ''; content.append(key, value); const edit = createActionButton('编辑', `edit-${scope}-memory`, entry.key); edit.addEventListener('click', () => { const isGlobal = scope === 'global'; (isGlobal ? elements.globalMemoryKey : elements.environmentMemoryKey).value = entry.key || ''; (isGlobal ? elements.globalMemoryValue : elements.environmentMemoryValue).value = entry.value || ''; openEditorDialog(isGlobal ? 'globalMemoryDialog' : 'environmentMemoryDialog'); }); row.append(content, edit, createActionButton('删除', `delete-${scope}-memory`, entry.key, 'danger')); container.append(row); }
}
async function loadGlobalMemory() { setStatus('正在显式读取 Global Memory…', 'loading'); try { renderMemory(elements.globalMemoryList, safeArray(await desktopAdapter().ListGlobalMemory()), 'global'); globalMemoryLoaded = true; setStatus('Global Memory 已加载', 'success'); } catch (error) { setStatus(`Global Memory 读取失败：${error?.message || String(error)}`, 'error'); } }
async function loadEnvironmentMemory() { if (!selectedEnvironmentID) return; setStatus('正在显式读取 Environment-private Memory…', 'loading'); try { renderMemory(elements.environmentMemoryList, safeArray(await desktopAdapter().ListEnvironmentMemory(selectedEnvironmentID)), 'environment'); environmentMemoryLoaded = true; setStatus('Environment-private Memory 已加载', 'success'); } catch (error) { setStatus(`Environment-private Memory 读取失败：${error?.message || String(error)}`, 'error'); } }

async function refreshSnapshot(successMessage = '') {
  const requestGeneration = connectionGeneration;
  elements.refreshButton.disabled = true; renderDashboardState('loading');
  setStatus('正在通过 Admin MCP 读取 ADM 状态…', 'loading');
  try {
    let snapshot;
    try { snapshot = await desktopAdapter().GetSnapshot(); }
    catch (error) {
      if (requestGeneration !== connectionGeneration) return false;
      const message = `Admin MCP 管理快照读取失败：${errorText(error)}`;
      renderDashboardState(currentSnapshot ? 'stale' : 'error', message);
      if (!currentSnapshot) renderManagementUnavailable(message);
      setStatus(message, 'error'); return false;
    }
    if (requestGeneration !== connectionGeneration) return false;
    lastSnapshotSuccessAt = Date.now(); explicitSkillAvailabilityProbe = null; skillSourcesState = 'loading'; skillSourcesError = ''; renderSnapshotBase(snapshot);

    let sourceError = '';
    try { skillSources = safeArray(await desktopAdapter().ListSkillSources()); skillSourcesState = 'success'; }
    catch (error) { skillSources = []; skillSourcesState = 'error'; skillSourcesError = errorText(error); sourceError = `Skill sources: ${skillSourcesError}`; }
    if (requestGeneration !== connectionGeneration) return false;
    renderSkillManager(safeArray(currentSnapshot?.skills));

    const contextResult = await refreshManagementContext(captureEnvironmentScope());
    if (requestGeneration !== connectionGeneration) return false;
    let detailError = '';
    if (selectedEnvironmentID) {
      try { await refreshSelectedEnvironmentDetail(); }
      catch (error) { detailError = `Environment detail: ${errorText(error)}`; }
    }
    const auxiliaryErrors = [sourceError, ...safeArray(contextResult?.errors), detailError].filter(Boolean);
    if (auxiliaryErrors.length) setStatus(`基础快照已刷新；部分状态不可用：${auxiliaryErrors.join(' · ')}`, 'error');
    else setStatus(successMessage || `已刷新 · ${new Date().toLocaleTimeString()}`, 'success');
    return true;
  } finally { if (requestGeneration === connectionGeneration) elements.refreshButton.disabled = false; }
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
function syncMCPHealthForm() {
  const healthEnabled = elements.mcpHealthEnabled.checked; const reconnectEnabled = elements.mcpAutoReconnect.checked;
  elements.mcpProbeTimeout.disabled = !healthEnabled; elements.mcpCheckInterval.disabled = !healthEnabled; elements.mcpReconnectInterval.disabled = !reconnectEnabled;
}
function clearMCPObservedHealth(mcpID) {
  for (const key of [...mcpHealthByKey.keys()]) if (key.endsWith(`:${mcpID}`)) mcpHealthByKey.delete(key);
}
function resetMCPEditor(close = false, rerender = true) {
  editingMCPID = ''; elements.mcpForm.reset(); elements.mcpTransport.value = 'streamable-http'; elements.mcpHealthEnabled.checked = true; elements.mcpAutoReconnect.checked = false; elements.mcpProbeTimeout.value = '5'; elements.mcpCheckInterval.value = '30'; elements.mcpReconnectInterval.value = '30';
  elements.mcpEditorSummary.textContent = '添加 MCP 定义'; elements.mcpEditorHint.hidden = true; elements.mcpSubmitButton.textContent = '添加全局 MCP'; elements.mcpEditCancelButton.hidden = true; syncMCPTransportForm(); syncMCPHealthForm();
  if (close) closeEditorDialog('mcpEditorFlow'); if (rerender) renderMCPManager(safeArray(currentSnapshot?.mcps));
}
function beginMCPEdit(entry) {
  if (!entry) return;
  editingMCPID = entry.id; elements.mcpName.value = entry.name || ''; elements.mcpTransport.value = entry.transport || 'streamable-http'; elements.mcpEndpoint.value = entry.endpoint || ''; elements.mcpExecutable.value = entry.executable || ''; elements.mcpArgs.value = safeArray(entry.args).join('\n'); elements.mcpAuthMode.value = entry.auth_mode || 'none';
  elements.mcpReferencePairs.value = referenceLines(entry.transport === 'stdio' ? entry.env_refs : entry.header_refs); elements.mcpHealthEnabled.checked = Boolean(entry.health_policy?.health_check_enabled); elements.mcpAutoReconnect.checked = Boolean(entry.health_policy?.auto_reconnect); elements.mcpProbeTimeout.value = String(entry.health_policy?.probe_timeout_seconds || 5); elements.mcpCheckInterval.value = String(entry.health_policy?.check_interval_seconds || 30); elements.mcpReconnectInterval.value = String(entry.health_policy?.reconnect_interval_seconds || 30); elements.mcpDefault.checked = Boolean(entry.default_include_in_environment);
  elements.mcpEditorSummary.textContent = `编辑 MCP · ${entry.name || entry.id}`; elements.mcpEditorHint.hidden = false; elements.mcpSubmitButton.textContent = '保存修改'; elements.mcpEditCancelButton.hidden = false; syncMCPTransportForm(); syncMCPHealthForm(); openEditorDialog('mcpEditorFlow');
}

elements.mcpTransport.addEventListener('change', syncMCPTransportForm);
elements.mcpHealthEnabled.addEventListener('change', syncMCPHealthForm);
elements.mcpAutoReconnect.addEventListener('change', syncMCPHealthForm);
elements.mcpEditCancelButton.addEventListener('click', () => resetMCPEditor(true));
elements.mcpFilter.addEventListener('input', () => renderMCPManager(safeArray(currentSnapshot?.mcps)));
elements.mcpStateFilter.addEventListener('change', () => renderMCPManager(safeArray(currentSnapshot?.mcps)));
elements.skillFilter.addEventListener('input', () => renderSkillManager(safeArray(currentSnapshot?.skills)));
elements.skillStateFilter.addEventListener('change', () => renderSkillManager(safeArray(currentSnapshot?.skills)));
elements.skillProbeAllButton.addEventListener('click', () => probeAllSkillAvailability());
elements.skillSelectVisibleButton.addEventListener('click', () => { for (const input of elements.skillList.querySelectorAll('input[data-action="select-skill"]')) selectedSkillIDs.add(input.dataset.id); renderSkillManager(safeArray(currentSnapshot?.skills)); });
elements.skillClearSelectionButton.addEventListener('click', () => { selectedSkillIDs.clear(); renderSkillManager(safeArray(currentSnapshot?.skills)); });
elements.skillDeleteSelectedButton.addEventListener('click', () => removeSkillCatalogEntries([...selectedSkillIDs], '批量删除'));
elements.skillClearUnavailableButton.addEventListener('click', () => {
  if (!explicitSkillProbeIsCurrent()) return setStatus('当前 Skill 可用性检查结果已过期，请重新批量检查后再清理。', 'error');
  const ids = skillAvailabilitySummary().cleanupIDs;
  if (!ids.length) return setStatus('当前显式检查没有发现可清理的不可用 Skill。', 'success');
  return removeSkillCatalogEntries(ids, '一键清除不可用 Skill');
});
elements.managementEnvironment.addEventListener('change', async () => {
  environmentGeneration++; managementEnvironmentID = elements.managementEnvironment.value; explicitSkillAvailabilityProbe = null; managementContextError = ''; managementInspection = null; managementCapabilityFacts = new Map(); skillAvailabilityByID = new Map(); runtimeRunsByID = new Map(); elements.runtimeOutput.textContent = '尚无输出'; updateManagementHint(); updateSkillBulkControls();
  const scope = captureEnvironmentScope();
  if (scope.environmentID) { emptyMessage(elements.verifierList, '正在读取 verifier…'); emptyMessage(elements.processList, '正在读取 process…'); emptyMessage(elements.runList, '正在读取 run…'); }
  setStatus('正在加载 Environment MCP/Skill/Runtime 状态…', 'loading');
  try { const result = await refreshManagementContext(scope); if (result?.stale) return; setStatus(result?.errors?.length ? `Environment 已切换；部分状态不可用：${result.errors.join(' · ')}` : 'Environment 管理上下文已切换', result?.errors?.length ? 'error' : 'success'); }
  catch (error) { if (environmentScopeIsCurrent(scope)) setStatus(`Environment 状态读取失败：${errorText(error)}`, 'error'); }
});
elements.editEnvironmentButton.addEventListener('click', () => {
  const environment = currentEnvironment();
  if (!environment) return;
  openRenameDialog('Environment', environment.name || '', (name) => desktopAdapter().RenameEnvironment(environment.environment_id, name));
});
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
  const editingID = editingMCPID;
  try {
    const transport = elements.mcpTransport.value; const refs = referenceMap(elements.mcpReferencePairs.value); const input = {
      name: elements.mcpName.value.trim(), transport, auth_mode: transport === 'stdio' ? 'none' : elements.mcpAuthMode.value,
      endpoint: transport === 'stdio' ? '' : elements.mcpEndpoint.value.trim(), executable: transport === 'stdio' ? elements.mcpExecutable.value.trim() : '', args: transport === 'stdio' ? lineValues(elements.mcpArgs.value) : [],
      header_refs: transport === 'stdio' ? {} : refs, env_refs: transport === 'stdio' ? refs : {}, default_include_in_environment: elements.mcpDefault.checked,
      health_policy: {health_check_enabled: elements.mcpHealthEnabled.checked, check_interval_seconds: Math.max(1, safeNumber(elements.mcpCheckInterval.value, 30)), probe_timeout_seconds: Math.max(1, safeNumber(elements.mcpProbeTimeout.value, 5)), auto_reconnect: elements.mcpAutoReconnect.checked, reconnect_interval_seconds: Math.max(1, safeNumber(elements.mcpReconnectInterval.value, 30))},
    };
    if (!input.name) throw new Error('MCP 名称不能为空');
    if (editingID) {
      await runMutation('更新全局 MCP', async () => { await desktopAdapter().UpdateMCP(editingID, input); clearMCPObservedHealth(editingID); resetMCPEditor(true); });
    } else {
      await runMutation('添加全局 MCP', async () => { await desktopAdapter().AddMCP(input); resetMCPEditor(true); });
    }
  } catch (error) { setStatus(`${editingID ? '更新' : '添加'} MCP 失败：${error?.message || String(error)}`, 'error'); }
});

elements.mcpImportForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const content = elements.mcpImportContent.value; if (!content.trim()) return;
  const input = {format: elements.mcpImportFormat.value, json_or_jsonc: content, conflict_policy: elements.mcpImportConflict.value, default_include: elements.mcpImportDefault.checked}; setStatus('正在预览 MCP 导入…', 'loading');
  try { const preview = await desktopAdapter().PreviewMCPImport(input); pendingMCPImport = input; renderImportPreview(preview); if (!elements.mcpImportApplyButton.disabled) setStatus('导入预览已生成；确认后再应用', 'success'); else setStatus('导入预览包含错误，不能应用', 'error'); }
  catch (error) { pendingMCPImport = null; elements.mcpImportApplyButton.disabled = true; emptyMessage(elements.mcpImportPreview, `预览失败：${error?.message || String(error)}`); setStatus(`MCP 导入预览失败：${error?.message || String(error)}`, 'error'); }
});
elements.mcpImportApplyButton.addEventListener('click', async () => {
  if (!pendingMCPImport) return; elements.mcpImportApplyButton.disabled = true; setStatus('正在应用 MCP 导入…', 'loading');
  try { const result = await desktopAdapter().ApplyMCPImport(pendingMCPImport); pendingMCPImport = null; renderImportApplyResult(result); elements.mcpImportContent.value = ''; closeFormDialog(elements.mcpImportForm); await refreshSnapshot('MCP 导入完成'); }
  catch (error) { setStatus(`MCP 导入失败：${error?.message || String(error)}`, 'error'); }
});

elements.skillSourceForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const root = elements.skillSourceRoot.value.trim(); if (!root) return; setStatus('正在添加 Skill source…', 'loading');
  try {
    const source = await desktopAdapter().AddSkillSource({root, support_roots: lineValues(elements.skillSupportRoots.value), default_include_in_environment: elements.skillSourceDefault.checked});
    elements.skillSourceForm.reset(); closeFormDialog(elements.skillSourceForm);
    try { await desktopAdapter().RefreshSkillSource(source.skill_source_id); elements.skillSourceForm.reset(); closeFormDialog(elements.skillSourceForm); await refreshSnapshot('Skill source 已添加并刷新'); }
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
  if (button.dataset.action === 'edit-mcp') { beginMCPEdit(entry); return; }
  if (button.dataset.action === 'probe-mcp' && managementEnvironmentID) { setStatus(`正在探测 ${entry?.name || id}…`, 'loading'); try { const health = await desktopAdapter().ProbeMCPHealth(managementEnvironmentID, id); mcpHealthByKey.set(mcpHealthKey(managementEnvironmentID, id), health); renderMCPManager(safeArray(currentSnapshot?.mcps)); setStatus(`MCP 探测完成：${health.state}`, health.state === 'healthy' ? 'success' : 'error'); } catch (error) { setStatus(`MCP 探测失败：${error?.message || String(error)}`, 'error'); } }
  if (button.dataset.action === 'remove-mcp' && window.confirm(`这是全局删除，不是只从当前 Environment 禁用。已有 Environment 中的 ID 引用不会被静默改写，删除后可能显示 unresolved。继续？\n${entry?.name || id}`)) await runMutation('删除全局 MCP 定义', async () => { await desktopAdapter().RemoveMCP(id); clearMCPObservedHealth(id); if (editingMCPID === id) resetMCPEditor(true, false); });
});

elements.skillSourceList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id; const source = skillSources.find((item) => item.skill_source_id === id);
  if (button.dataset.action === 'refresh-skill-source') await runMutation('刷新 Skill source', () => desktopAdapter().RefreshSkillSource(id));
  if (button.dataset.action === 'remove-skill-source' && window.confirm(`删除全局 Skill source 及它发现的 Skill 条目。Environment 中已有 Skill ID 不会被静默改写，之后可能显示 unresolved。继续？\n${source?.root || id}`)) await runMutation('删除 Skill source', () => desktopAdapter().RemoveSkillSource(id));
});
elements.skillList.addEventListener('change', async (event) => {
  const input = event.target.closest('input[data-action]'); if (!input) return;
  if (input.dataset.action === 'select-skill') { if (input.checked) selectedSkillIDs.add(input.dataset.id); else selectedSkillIDs.delete(input.dataset.id); updateSkillBulkControls(); return; }
  if (input.dataset.action === 'default-skill') await runMutation('更新 Skill 新环境默认值', () => desktopAdapter().SetSkillDefault(input.dataset.id, input.checked));
  if (input.dataset.action === 'environment-skill' && managementEnvironmentID) await runMutation('更新当前 Environment Skill 选择', () => desktopAdapter().SetEnvironmentSkill(managementEnvironmentID, input.dataset.id, input.checked));
});
elements.skillList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action="remove-skill"]'); if (!button) return; if (window.confirm(`删除这个全局 Skill 条目？Source-managed Skill 应通过 Source 管理。\n${button.dataset.id}`)) await runMutation('删除 Skill 条目', () => desktopAdapter().RemoveSkill(button.dataset.id)); });

elements.workspaceForm.addEventListener('submit', async (event) => { event.preventDefault(); const path = elements.workspacePath.value.trim(), name = elements.workspaceName.value.trim(); if (path) await runMutation('添加 Workspace', async () => { await desktopAdapter().AddWorkspace({path, name}); elements.workspaceForm.reset(); closeFormDialog(elements.workspaceForm); }); });
elements.environmentForm.addEventListener('submit', async (event) => { event.preventDefault(); const workspaceID = elements.environmentWorkspace.value, name = elements.environmentName.value.trim(), root = elements.environmentRoot.value.trim(); if (workspaceID && name) await runMutation('创建 Environment', async () => { await desktopAdapter().CreateEnvironment({workspace_id: workspaceID, name, root}); elements.environmentName.value = ''; elements.environmentRoot.value = ''; closeFormDialog(elements.environmentForm); }); });
elements.execForm.addEventListener('submit', async (event) => { event.preventDefault(); const executable = elements.execExecutable.value.trim(); if (executable) await runMutation('更新 exec allowlist', async () => { await desktopAdapter().AllowExecutable(executable); elements.execForm.reset(); closeFormDialog(elements.execForm); }); });
elements.globalMemoryForm.addEventListener('submit', async (event) => { event.preventDefault(); const key = elements.globalMemoryKey.value.trim(), value = elements.globalMemoryValue.value; if (key) await runMutation('写入 Global Memory', () => desktopAdapter().WriteGlobalMemory(key, value), async () => { elements.globalMemoryForm.reset(); closeFormDialog(elements.globalMemoryForm); if (globalMemoryLoaded) await loadGlobalMemory(); }); });
elements.environmentMemoryForm.addEventListener('submit', async (event) => { event.preventDefault(); if (!selectedEnvironmentID) return; const key = elements.environmentMemoryKey.value.trim(), value = elements.environmentMemoryValue.value; if (key) await runMutation('写入 Environment-private Memory', () => desktopAdapter().WriteEnvironmentMemory(selectedEnvironmentID, key, value), async () => { elements.environmentMemoryForm.reset(); closeFormDialog(elements.environmentMemoryForm); if (environmentMemoryLoaded) await loadEnvironmentMemory(); }); });

elements.workspaceList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id, workspace = safeArray(currentSnapshot?.workspaces).find((item) => item.workspace_id === id); if (button.dataset.action === 'rename-workspace') { openRenameDialog('Workspace', workspace?.name || '', (name) => desktopAdapter().RenameWorkspace(id, name)); } if (button.dataset.action === 'remove-workspace' && window.confirm(`只移除 ADM Workspace 记录，不删除目录。继续？\n${workspace?.path || id}`)) await runMutation('移除 Workspace', () => desktopAdapter().RemoveWorkspace(id)); });
elements.environmentList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id, environment = safeArray(currentSnapshot?.environments).find((item) => item.environment_id === id);
  if (button.dataset.action === 'inspect-environment') {
    environmentDetailOpener = button; const token = ++detailGeneration; selectedEnvironmentID = id;
    if (managementEnvironmentID !== id) { environmentGeneration++; explicitSkillAvailabilityProbe = null; }
    managementEnvironmentID = id; elements.managementEnvironment.value = id; managementContextError = ''; updateSkillBulkControls(); setStatus('读取 Environment 详情…', 'loading');
    try {
      const contextResult = await refreshManagementContext(captureEnvironmentScope());
      if (token !== detailGeneration || selectedEnvironmentID !== id) return;
      environmentMemoryLoaded = false; emptyMessage(elements.environmentMemoryList, '尚未加载 private Memory');
      const inspection = await desktopAdapter().InspectEnvironment(id);
      if (token !== detailGeneration || selectedEnvironmentID !== id) return;
      await renderEnvironmentDetailFromInspection(inspection, token);
      if (token !== detailGeneration || selectedEnvironmentID !== id) return;
      setStatus(contextResult?.errors?.length ? `Environment 详情已加载；部分状态不可用：${contextResult.errors.join(' · ')}` : 'Environment 详情已加载', contextResult?.errors?.length ? 'error' : 'success');
    } catch (error) { if (token === detailGeneration && selectedEnvironmentID === id) { environmentDetailOpener = null; setStatus(`读取 Environment 详情失败：${errorText(error)}`, 'error'); } }
  }
  if (button.dataset.action === 'rename-environment') { openRenameDialog('Environment', environment?.name || '', (name) => desktopAdapter().RenameEnvironment(id, name)); }
  if (button.dataset.action === 'remove-environment' && window.confirm(`只移除 ADM Environment 记录，不删除 root 或项目文件。继续？\n${environment?.root || id}`)) await runMutation('移除 Environment', () => desktopAdapter().RemoveEnvironment(id));
});
elements.execList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action="remove-executable"]'); if (button) await runMutation('移除 executable', () => desktopAdapter().RemoveExecutable(button.dataset.id)); });
async function handleSelectionChange(event, kind) { const checkbox = event.target.closest('input[data-kind]'); if (!checkbox || !selectedEnvironmentID) return; await runMutation(`更新 Environment ${kind.toUpperCase()} 选择`, () => kind === 'mcp' ? desktopAdapter().SetEnvironmentMCP(selectedEnvironmentID, checkbox.dataset.id, checkbox.checked) : desktopAdapter().SetEnvironmentSkill(selectedEnvironmentID, checkbox.dataset.id, checkbox.checked)); }
elements.environmentMCPSelections.addEventListener('change', (event) => handleSelectionChange(event, 'mcp')); elements.environmentSkillSelections.addEventListener('change', (event) => handleSelectionChange(event, 'skill'));
elements.globalMemoryList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action="delete-global-memory"]'); if (button) await runMutation('删除 Global Memory', () => desktopAdapter().DeleteGlobalMemory(button.dataset.id), loadGlobalMemory); });
elements.environmentMemoryList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action="delete-environment-memory"]'); if (button && selectedEnvironmentID) await runMutation('删除 Environment-private Memory', () => desktopAdapter().DeleteEnvironmentMemory(selectedEnvironmentID, button.dataset.id), loadEnvironmentMemory); });

elements.loadGlobalMemory.addEventListener('click', loadGlobalMemory); elements.loadEnvironmentMemory.addEventListener('click', loadEnvironmentMemory); elements.closeEnvironmentDetail.addEventListener('click', closeEnvironmentDetail); elements.environmentDetailBackdrop.addEventListener('click', closeEnvironmentDetail);
elements.launchAtLogin.addEventListener('change', updateLaunchAtLogin);
window.addEventListener('focus', () => { loadDesktopPreferences(false); });
window.addEventListener('keydown', (event) => { if (event.key !== 'Escape' || activeEditorDialog()) return; if (selectedEnvironmentID) closeEnvironmentDetail(); else if (editingMCPID) resetMCPEditor(true); });
elements.gatewayRefreshButton.addEventListener('click', () => refreshConnectedADM(true).catch((error) => { clearManagementData(); setStatus(`ADM 连接检查失败：${error?.message || String(error)}`, 'error'); }));
elements.gatewayStartButton.addEventListener('click', () => runGatewayAction('启动本地 ADM', (input) => desktopAdapter().StartLocalADM(input)));
elements.gatewayStopButton.addEventListener('click', () => runGatewayAction('停止本地 ADM', (input) => desktopAdapter().StopLocalADM(input)));
elements.refreshButton.addEventListener('click', () => refreshConnectedADM(false).catch((error) => { clearManagementData(); setStatus(`ADM 连接检查失败：${error?.message || String(error)}`, 'error'); }));
window.addEventListener('DOMContentLoaded', () => { window.runtime?.EventsOn?.('desktop:preferences-changed', () => loadDesktopPreferences(false)); initializeManagementNavigation(); initializeConnectionProfiles(); syncMCPTransportForm(); clearManagementData(); loadDesktopPreferences(false); });
