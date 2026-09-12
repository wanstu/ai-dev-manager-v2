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
  runtimeRefreshButton: document.getElementById('runtimeRefreshButton'), runtimeHint: document.getElementById('runtimeHint'), runtimeSubviewTabs: document.getElementById('runtimeSubviewTabs'), runtimeVerifierCount: document.getElementById('runtimeVerifierCount'), runtimeProcessCount: document.getElementById('runtimeProcessCount'), runtimeRunCount: document.getElementById('runtimeRunCount'), verifierList: document.getElementById('verifierList'), processList: document.getElementById('processList'), runList: document.getElementById('runList'), runtimeOutputMeta: document.getElementById('runtimeOutputMeta'), runtimeOutput: document.getElementById('runtimeOutput'),
  mcpTotalCount: document.getElementById('mcpTotalCount'), mcpDefaultCount: document.getElementById('mcpDefaultCount'), mcpEnvironmentCount: document.getElementById('mcpEnvironmentCount'), mcpIssueCount: document.getElementById('mcpIssueCount'),
  mcpForm: document.getElementById('mcpForm'), mcpName: document.getElementById('mcpName'), mcpTransport: document.getElementById('mcpTransport'), mcpEndpointField: document.getElementById('mcpEndpointField'), mcpEndpoint: document.getElementById('mcpEndpoint'),
  mcpExecutableField: document.getElementById('mcpExecutableField'), mcpExecutable: document.getElementById('mcpExecutable'), mcpArgsField: document.getElementById('mcpArgsField'), mcpArgs: document.getElementById('mcpArgs'),
  mcpAuthField: document.getElementById('mcpAuthField'), mcpAuthMode: document.getElementById('mcpAuthMode'), mcpRefsField: document.getElementById('mcpRefsField'), mcpReferencePairs: document.getElementById('mcpReferencePairs'),
  mcpHealthEnabled: document.getElementById('mcpHealthEnabled'), mcpAutoReconnect: document.getElementById('mcpAutoReconnect'), mcpProbeTimeout: document.getElementById('mcpProbeTimeout'), mcpCheckInterval: document.getElementById('mcpCheckInterval'), mcpReconnectInterval: document.getElementById('mcpReconnectInterval'), mcpDefault: document.getElementById('mcpDefault'),
  mcpEditorFlow: document.getElementById('mcpEditorFlow'), mcpEditorSummary: document.getElementById('mcpEditorSummary'), mcpEditorHint: document.getElementById('mcpEditorHint'), mcpSubmitButton: document.getElementById('mcpSubmitButton'), mcpEditCancelButton: document.getElementById('mcpEditCancelButton'),
  mcpImportForm: document.getElementById('mcpImportForm'), mcpImportFormat: document.getElementById('mcpImportFormat'), mcpImportConflict: document.getElementById('mcpImportConflict'), mcpImportDefault: document.getElementById('mcpImportDefault'), mcpImportContent: document.getElementById('mcpImportContent'),
  mcpImportApplyButton: document.getElementById('mcpImportApplyButton'), mcpImportPreview: document.getElementById('mcpImportPreview'), mcpList: document.getElementById('mcpList'),
  mcpFilter: document.getElementById('mcpFilter'), mcpStateFilter: document.getElementById('mcpStateFilter'), mcpVisibleCount: document.getElementById('mcpVisibleCount'), mcpListTotalCount: document.getElementById('mcpListTotalCount'),
  mcpSetVisibleDefaultButton: document.getElementById('mcpSetVisibleDefaultButton'), mcpUnsetVisibleDefaultButton: document.getElementById('mcpUnsetVisibleDefaultButton'), mcpEnableVisibleButton: document.getElementById('mcpEnableVisibleButton'), mcpDisableVisibleButton: document.getElementById('mcpDisableVisibleButton'), mcpBulkHint: document.getElementById('mcpBulkHint'),
  skillSourceCount: document.getElementById('skillSourceCount'), skillTotalCount: document.getElementById('skillTotalCount'), skillEnvironmentCount: document.getElementById('skillEnvironmentCount'), skillIssueCount: document.getElementById('skillIssueCount'),
  skillSourceForm: document.getElementById('skillSourceForm'), skillSourceID: document.getElementById('skillSourceID'), skillSourceDialogTitle: document.getElementById('skillSourceDialogTitle'), skillSourceRoot: document.getElementById('skillSourceRoot'), skillSupportRoots: document.getElementById('skillSupportRoots'), skillSourceDefault: document.getElementById('skillSourceDefault'), skillSourceSubmitButton: document.getElementById('skillSourceSubmitButton'), skillSubviewTabs: document.getElementById('skillSubviewTabs'), skillSubviewSkillCount: document.getElementById('skillSubviewSkillCount'), skillSubviewSourceCount: document.getElementById('skillSubviewSourceCount'), skillsPanel: document.getElementById('skillsPanel'), skillSourcesPanel: document.getElementById('skillSourcesPanel'), skillSourceList: document.getElementById('skillSourceList'), skillList: document.getElementById('skillList'),
  skillFilter: document.getElementById('skillFilter'), skillSourceFilter: document.getElementById('skillSourceFilter'), skillStateFilter: document.getElementById('skillStateFilter'), skillSourceFilterInput: document.getElementById('skillSourceFilterInput'), skillSourceVisibleCount: document.getElementById('skillSourceVisibleCount'), skillSourceListTotalCount: document.getElementById('skillSourceListTotalCount'), skillVisibleCount: document.getElementById('skillVisibleCount'), skillListTotalCount: document.getElementById('skillListTotalCount'),
  skillProbeAllButton: document.getElementById('skillProbeAllButton'), skillSelectVisibleButton: document.getElementById('skillSelectVisibleButton'), skillClearSelectionButton: document.getElementById('skillClearSelectionButton'), skillSetVisibleDefaultButton: document.getElementById('skillSetVisibleDefaultButton'), skillUnsetVisibleDefaultButton: document.getElementById('skillUnsetVisibleDefaultButton'), skillEnableVisibleButton: document.getElementById('skillEnableVisibleButton'), skillDisableVisibleButton: document.getElementById('skillDisableVisibleButton'), skillSelectedCount: document.getElementById('skillSelectedCount'), skillDeleteSelectedButton: document.getElementById('skillDeleteSelectedButton'), skillClearUnavailableButton: document.getElementById('skillClearUnavailableButton'), skillBulkHint: document.getElementById('skillBulkHint'),
  workspaceForm: document.getElementById('workspaceForm'), workspacePath: document.getElementById('workspacePath'), workspaceName: document.getElementById('workspaceName'), workspaceList: document.getElementById('workspaceList'), workspaceFilter: document.getElementById('workspaceFilter'), workspaceVisibleCount: document.getElementById('workspaceVisibleCount'), workspaceListTotalCount: document.getElementById('workspaceListTotalCount'),
  environmentForm: document.getElementById('environmentForm'), environmentWorkspace: document.getElementById('environmentWorkspace'), environmentName: document.getElementById('environmentName'), environmentRoot: document.getElementById('environmentRoot'), environmentList: document.getElementById('environmentList'), environmentFilter: document.getElementById('environmentFilter'), environmentWorkspaceFilter: document.getElementById('environmentWorkspaceFilter'), environmentVisibleCount: document.getElementById('environmentVisibleCount'), environmentListTotalCount: document.getElementById('environmentListTotalCount'), environmentFilterHint: document.getElementById('environmentFilterHint'),
  environmentDetailBackdrop: document.getElementById('environmentDetailBackdrop'), environmentDetailPanel: document.getElementById('environmentDetailPanel'), environmentDetailTitle: document.getElementById('environmentDetailTitle'), environmentDetailSubviewTabs: document.getElementById('environmentDetailSubviewTabs'), environmentDetail: document.getElementById('environmentDetail'), environmentDiagnostics: document.getElementById('environmentDiagnostics'), environmentDetailRoutes: document.getElementById('environmentDetailRoutes'),
  environmentMCPSelections: document.getElementById('environmentMCPSelections'), environmentSkillSelections: document.getElementById('environmentSkillSelections'), closeEnvironmentDetail: document.getElementById('closeEnvironmentDetail'),
  loadEnvironmentMemory: document.getElementById('loadEnvironmentMemory'), writeEnvironmentMemoryButton: document.getElementById('writeEnvironmentMemoryButton'), environmentMemoryScopeHint: document.getElementById('environmentMemoryScopeHint'), environmentMemoryForm: document.getElementById('environmentMemoryForm'), environmentMemoryKey: document.getElementById('environmentMemoryKey'), environmentMemoryValue: document.getElementById('environmentMemoryValue'), environmentMemoryList: document.getElementById('environmentMemoryList'),
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
let environmentSkillAvailabilityByID = new Map();
let mcpHealthByKey = new Map();
let runtimeRunsByID = new Map();
let runtimeSubview = 'verifiers';
let environmentDetailSubview = 'summary';
let runtimeSubviewGeneration = 0;
let runtimeLists = {verifiers: [], processes: [], runs: []};
let runtimeListStates = {verifiers: 'unloaded', processes: 'unloaded', runs: 'unloaded'};
let runtimeListErrors = {verifiers: '', processes: '', runs: ''};
let runtimeOutputBinding = null;
let runtimePendingActionKey = '';
let pendingMCPImport = null;
let mcpImportBusy = false;
let editingMCPID = '';
let editingSkillSourceID = '';
let skillSubview = 'skills';
let mcpBulkBusy = false;
let globalMemoryLoaded = false;
let globalMemoryLoading = false;
let environmentMemoryLoaded = false;
let environmentMemoryLoading = false;
let loadedEnvironmentMemoryID = '';
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
let managementSkillAvailabilityError = '';
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
function searchMatcher(raw) {
  const text = String(raw || '').trim();
  if (!text) return () => true;
  const slash = text.match(/^\/(.*)\/([dgimsuvy]*)$/);
  if (slash) {
    try {
      const flags = [...new Set(slash[2].replace(/[gy]/g, '').split(''))].join('');
      const regex = new RegExp(slash[1], flags || 'i');
      return (value) => regex.test(String(value || ''));
    } catch (_) {}
  }
  const lower = text.toLowerCase();
  return (value) => String(value || '').toLowerCase().includes(lower);
}
function visibleResourceIDs(container) { return [...container.querySelectorAll('.resource-row[data-id]')].map((row) => row.dataset.id).filter(Boolean); }
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
    afterNavigate: (route) => { if (route === 'memory') { renderEnvironmentMemoryScope(); maybeAutoLoadGlobalMemory(); } },
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
  if (managementEnvironmentID !== previous) environmentGeneration++;
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
function runtimeContainer(kind) { return ({verifiers: elements.verifierList, processes: elements.processList, runs: elements.runList})[kind]; }
function runtimeCountElement(kind) { return ({verifiers: elements.runtimeVerifierCount, processes: elements.runtimeProcessCount, runs: elements.runtimeRunCount})[kind]; }
function resetRuntimeCollections(state = 'unloaded', error = '') {
  runtimeLists = {verifiers: [], processes: [], runs: []};
  runtimeListStates = {verifiers: state, processes: state, runs: state};
  runtimeListErrors = {verifiers: error, processes: error, runs: error};
  runtimeRunsByID = new Map();
}
function syncRuntimeSubviewUI() {
  runtimeSubview = window.ADMRuntimeView?.normalizeSubview(runtimeSubview) || runtimeSubview || 'verifiers';
  for (const button of elements.runtimeSubviewTabs.querySelectorAll('[data-runtime-subview]')) {
    const active = button.dataset.runtimeSubview === runtimeSubview; button.setAttribute('aria-selected', String(active)); button.tabIndex = active ? 0 : -1;
  }
  for (const panel of document.querySelectorAll('[data-runtime-subview-panel]')) panel.hidden = panel.dataset.runtimeSubviewPanel !== runtimeSubview;
}
function clearRuntimeOutput(reason = '尚未选择资源') {
  runtimeOutputBinding = null;
  elements.runtimeOutputMeta.textContent = `有界 current-owner observation · ${reason}`;
  elements.runtimeOutput.textContent = '尚无输出';
}
function runtimeOutputBindingIsCurrent(binding = runtimeOutputBinding) {
  if (!binding) return false;
  const current = {...captureEnvironmentScope(), kind: runtimeSubview, id: binding.id, viewGeneration: runtimeSubviewGeneration};
  return window.ADMRuntimeView?.outputBindingMatches(binding, current) ?? false;
}
function validateRuntimeOutputBinding() {
  if (!runtimeOutputBinding) return;
  if (!runtimeOutputBindingIsCurrent(runtimeOutputBinding)) { clearRuntimeOutput('上下文或 Runtime 视图已变化'); return; }
  const known = window.ADMRuntimeView?.resourceStillKnown(runtimeOutputBinding, runtimeLists, runtimeListStates) ?? true;
  if (!known) clearRuntimeOutput('资源已不在最新 owner-local 列表');
}
function showRuntimeOutput(kind, id, title, stdout = '', stderr = '', meta = '', truncation = '') {
  if (runtimeSubview !== kind) return false;
  const scope = captureEnvironmentScope();
  runtimeOutputBinding = window.ADMRuntimeView?.makeOutputBinding(scope, kind, id, runtimeSubviewGeneration) || {...scope, kind, id, viewGeneration: runtimeSubviewGeneration};
  const environment = currentEnvironment(); const observedAt = new Date().toLocaleTimeString();
  elements.runtimeOutputMeta.textContent = `有界 current-owner observation · ${environment?.name || scope.environmentID || 'Environment'} · ${kind} ${id} · observed ${observedAt}${truncation ? ` · ${truncation}` : ''}`;
  const sections = [title]; if (meta) sections.push(meta); if (stdout) sections.push(`STDOUT\n${stdout}`); if (stderr) sections.push(`STDERR\n${stderr}`);
  elements.runtimeOutput.textContent = sections.filter(Boolean).join('\n\n') || '尚无输出';
  return true;
}
function renderRuntimeListState(kind, emptyText, renderItems) {
  const container = runtimeContainer(kind), count = runtimeCountElement(kind), state = runtimeListStates[kind], items = safeArray(runtimeLists[kind]);
  setMetric(count, state === 'success' ? items.length : state === 'loading' ? '…' : '—');
  if (!currentEnvironment()) return emptyMessage(container, '请选择 Environment');
  if (state === 'loading') return emptyMessage(container, `正在读取 ${kind}…`);
  if (state === 'error') return emptyMessage(container, `${kind} 状态不可用：${runtimeListErrors[kind] || 'unknown error'}`);
  if (state !== 'success') return emptyMessage(container, `${kind} 尚未加载`);
  if (!items.length) return emptyMessage(container, emptyText);
  container.replaceChildren(); container.classList.remove('empty'); renderItems(items, container);
}
function renderRuntime() {
  syncRuntimeSubviewUI();
  const environment = currentEnvironment(), writerOwner = runtimeWriterOwner(), pending = Boolean(runtimePendingActionKey);
  elements.runtimeHint.textContent = !environment
    ? '选择 Management Environment 后可查看当前 Gateway owner 的运行状态。'
    : writerOwner
      ? `当前 Environment Writer: ${writerOwner}。Writer 只是观测到的现有 authority；停止/取消/运行不会自动抢 lease。`
      : '当前 Environment 没有 active writer；可以查看 owner-local 状态与日志，但运行/停止/取消操作不可用。';

  renderRuntimeListState('verifiers', '没有配置 verifier', (verifiers, container) => {
    for (const verifier of verifiers) {
      const run = createActionButton('运行', 'run-verifier', verifier.verifier_id); run.disabled = pending || !writerOwner || verifier.enabled === false;
      const args = safeArray(verifier.args).join(' '); const detail = `${verifier.kind || 'custom'} · ${verifier.executable || ''}${args ? ` ${args}` : ''}${verifier.cwd ? ` · cwd ${verifier.cwd}` : ''}`;
      container.append(runtimeItem(verifier.name || verifier.verifier_id, verifier.verifier_id, verifier.enabled === false ? 'disabled' : 'configured', detail, [run]));
    }
  });
  renderRuntimeListState('processes', '当前 Gateway owner 没有 process', (processes, container) => {
    for (const process of processes) {
      const logs = createActionButton('日志', 'process-logs', process.id); logs.disabled = pending; const actions = [logs];
      if (process.state === 'running') { const stop = createActionButton('停止', 'stop-process', process.id, 'danger'); stop.disabled = pending || !writerOwner; actions.push(stop); }
      const ports = safeArray(process.listening_ports); const detail = `PID ${process.pid || '—'}${ports.length ? ` · ports ${ports.join(', ')}` : ''}${process.exit_code !== undefined && process.exit_code !== null ? ` · exit ${process.exit_code}` : ''}${process.error_kind ? ` · ${process.error_kind}` : ''}`;
      container.append(runtimeItem(process.id, process.id, process.state, detail, actions));
    }
  });
  runtimeRunsByID = new Map(safeArray(runtimeLists.runs).map((run) => [run.id, run]));
  renderRuntimeListState('runs', '当前 Gateway owner 没有 generic run', (runs, container) => {
    for (const run of runs) {
      const output = createActionButton('输出', 'run-output', run.id); output.disabled = pending; const actions = [output];
      if (run.state === 'running') { const cancel = createActionButton('取消', 'cancel-run', run.id, 'danger'); cancel.disabled = pending || !writerOwner; actions.push(cancel); }
      const detail = `${run.executable || ''}${safeArray(run.args).length ? ` ${safeArray(run.args).join(' ')}` : ''}${run.cwd ? ` · cwd ${run.cwd}` : ''}${run.exit_code !== undefined && run.exit_code !== null ? ` · exit ${run.exit_code}` : ''}${run.error_kind ? ` · ${run.error_kind}` : ''}`;
      container.append(runtimeItem(run.id, run.id, run.state, detail, actions));
    }
  });
  validateRuntimeOutputBinding();
}
async function refreshRuntimeContext(showMessage = false, scope = captureEnvironmentScope()) {
  const environment = currentEnvironment();
  if (!environment || !scope.environmentID) {
    if (environmentScopeIsCurrent(scope)) { resetRuntimeCollections('unloaded'); renderRuntime(); }
    return {stale: false, errors: []};
  }
  if (environment.environment_id !== scope.environmentID || !environmentScopeIsCurrent(scope)) return {stale: true, errors: []};
  runtimeListStates = {verifiers: 'loading', processes: 'loading', runs: 'loading'}; runtimeListErrors = {verifiers: '', processes: '', runs: ''}; renderRuntime();
  if (showMessage) setStatus('正在刷新 Runtime 状态…', 'loading');
  const results = await Promise.allSettled([
    desktopAdapter().ListVerifiers(scope.environmentID), desktopAdapter().ListProcesses(scope.environmentID), desktopAdapter().ListRuns(scope.environmentID),
  ]);
  if (!environmentScopeIsCurrent(scope)) return {stale: true, errors: []};
  const descriptors = [
    ['verifiers', 'Verifiers', results[0]], ['processes', 'Processes', results[1]], ['runs', 'Runs', results[2]],
  ];
  const errors = [];
  for (const [kind, label, result] of descriptors) {
    if (result.status === 'fulfilled') { runtimeLists[kind] = safeArray(result.value); runtimeListStates[kind] = 'success'; runtimeListErrors[kind] = ''; }
    else { runtimeListStates[kind] = 'error'; runtimeListErrors[kind] = errorText(result.reason); errors.push(`${label}: ${runtimeListErrors[kind]}`); }
  }
  renderRuntime();
  if (showMessage) setStatus(errors.length ? `Runtime 部分状态不可用：${errors.join(' · ')}` : 'Runtime 状态已刷新', errors.length ? 'error' : 'success');
  return {stale: false, errors};
}
async function runRuntimeAction(label, actionKey, scope, action) {
  if (runtimePendingActionKey) return null;
  runtimePendingActionKey = actionKey; renderRuntime(); setStatus(`${label}…`, 'loading');
  try {
    let result;
    try { result = await action(); }
    catch (error) { if (environmentScopeIsCurrent(scope)) setStatus(`${label}失败：${error?.message || String(error)}`, 'error'); return null; }
    const refreshErrors = [];
    if (environmentScopeIsCurrent(scope)) {
      try { const refreshResult = await refreshRuntimeContext(false, scope); refreshErrors.push(...safeArray(refreshResult?.errors)); }
      catch (error) { refreshErrors.push(errorText(error)); }
    }
    if (environmentScopeIsCurrent(scope)) setStatus(refreshErrors.length ? `${label}完成；后续状态刷新部分失败：${refreshErrors.join(' · ')}` : `${label}完成`, refreshErrors.length ? 'error' : 'success');
    return result;
  } finally { runtimePendingActionKey = ''; renderRuntime(); }
}
function renderEnvironmentWorkspaceFilter(workspaces) {
  const current = elements.environmentWorkspaceFilter.value;
  elements.environmentWorkspaceFilter.replaceChildren();
  const all = document.createElement('option'); all.value = ''; all.textContent = '全部 Workspace'; elements.environmentWorkspaceFilter.append(all);
  for (const workspace of workspaces) {
    const option = document.createElement('option'); option.value = workspace.workspace_id || ''; option.textContent = workspace.name || workspace.workspace_id; elements.environmentWorkspaceFilter.append(option);
  }
  const exists = !current || workspaces.some((workspace) => workspace.workspace_id === current);
  if (!exists) { const invalid = document.createElement('option'); invalid.value = current; invalid.textContent = `已移除 · ${current}`; elements.environmentWorkspaceFilter.append(invalid); }
  elements.environmentWorkspaceFilter.value = current;
  elements.environmentFilterHint.hidden = exists;
  elements.environmentFilterHint.textContent = exists ? '' : `Workspace 筛选 ${current} 已失效。请选择“全部 Workspace”或其他 Workspace 重置筛选。`;
}
function renderWorkspaces(workspaces) {
  const model = window.ADMProjectPages?.workspaceListModel(workspaces, safeArray(currentSnapshot?.environments), elements.workspaceFilter.value) || {items: workspaces, total: workspaces.length, visible: workspaces.length, environmentCounts: {}};
  elements.workspaceBadge.textContent = String(model.total); renderWorkspaceOptions(workspaces); renderEnvironmentWorkspaceFilter(workspaces); setMetric(elements.workspaceVisibleCount, model.visible); setMetric(elements.workspaceListTotalCount, model.total);
  if (!model.total) return emptyMessage(elements.workspaceList, '暂无已登记 Workspace。添加现有目录即可；Git 不是前置条件。');
  if (!model.visible) return emptyMessage(elements.workspaceList, '没有符合当前筛选条件的 Workspace。');
  elements.workspaceList.replaceChildren(); elements.workspaceList.classList.remove('empty');
  for (const workspace of model.items) {
    const item = document.createElement('article'); item.className = 'list-item managed-item'; const content = document.createElement('div'); content.className = 'item-content';
    const title = document.createElement('strong'); title.textContent = workspace.name || workspace.workspace_id; const id = document.createElement('code'); id.textContent = workspace.workspace_id || ''; const path = document.createElement('span'); path.className = 'project-path'; path.textContent = workspace.path || '';
    const meta = document.createElement('small'); meta.textContent = `Environments ${safeNumber(model.environmentCounts[workspace.workspace_id])}`; content.append(title, id, path, meta);
    const actions = document.createElement('div'); actions.className = 'item-actions'; actions.append(createActionButton('查看 Environments', 'filter-environments-by-workspace', workspace.workspace_id), createActionButton('改名', 'rename-workspace', workspace.workspace_id), createActionButton('删除记录', 'remove-workspace', workspace.workspace_id, 'danger')); item.append(content, actions); elements.workspaceList.append(item);
  }
}
function renderEnvironments(environments) {
  const workspaces = safeArray(currentSnapshot?.workspaces); const workspaceID = elements.environmentWorkspaceFilter.value;
  const model = window.ADMProjectPages?.environmentListModel(environments, workspaces, {query: elements.environmentFilter.value, workspaceID}) || {items: environments, total: environments.length, visible: environments.length, workspaceNames: {}, workspaceFilterValid: true};
  elements.environmentBadge.textContent = String(model.total); setMetric(elements.environmentVisibleCount, model.visible); setMetric(elements.environmentListTotalCount, model.total);
  elements.environmentFilterHint.hidden = model.workspaceFilterValid; if (model.workspaceFilterValid) elements.environmentFilterHint.textContent = '';
  if (!model.total) return emptyMessage(elements.environmentList, '暂无 Environment。可以从任意已登记 Workspace 创建，Git 不是前置条件。');
  if (!model.visible) return emptyMessage(elements.environmentList, model.workspaceFilterValid ? '没有符合当前筛选条件的 Environment。' : '当前 Workspace 筛选已失效；重置筛选后可查看现有 Environment。');
  elements.environmentList.replaceChildren(); elements.environmentList.classList.remove('empty');
  for (const environment of model.items) {
    const currentContext = environment.environment_id === managementEnvironmentID; const item = document.createElement('article'); item.className = 'list-item managed-item'; item.dataset.environmentId = environment.environment_id || ''; if (currentContext) item.classList.add('current-context');
    const content = document.createElement('div'); content.className = 'item-content'; const titleLine = document.createElement('div'); titleLine.className = 'item-title-line';
    const title = document.createElement('strong'); title.textContent = environment.name || environment.environment_id; const contextMarker = stateBadge('当前管理环境', 'selected'); contextMarker.classList.add('context-marker'); contextMarker.hidden = !currentContext; titleLine.append(title, contextMarker);
    const id = document.createElement('code'); id.textContent = environment.environment_id || ''; const root = document.createElement('span'); root.className = 'project-path'; root.textContent = environment.root || '';
    const workspaceName = model.workspaceNames[environment.workspace_id] || environment.workspace_id || '—';
    const meta = document.createElement('small'); meta.textContent = `Workspace ${workspaceName}${environment.workspace_id ? ` (${environment.workspace_id})` : ''} · ${environment.state || 'unknown'} · MCP ${safeArray(environment.enabled_mcp_ids).length} · Skills ${safeArray(environment.enabled_skill_ids).length} · Private Memory ${safeNumber(environment.private_memory_count)} · ${environment.writer?.owner ? `Writer ${environment.writer.owner}` : 'No writer'}`; content.append(titleLine, id, root, meta);
    const actions = document.createElement('div'); actions.className = 'item-actions'; actions.append(createActionButton('详情', 'inspect-environment', environment.environment_id), createActionButton('诊断', 'diagnose-environment', environment.environment_id), createActionButton('改名', 'rename-environment', environment.environment_id), createActionButton('删除记录', 'remove-environment', environment.environment_id, 'danger')); item.append(content, actions); elements.environmentList.append(item);
  }
}
function updateEnvironmentContextMarkers() {
  for (const item of elements.environmentList.querySelectorAll('[data-environment-id]')) {
    const current = item.dataset.environmentId === managementEnvironmentID;
    item.classList.toggle('current-context', current);
    const marker = item.querySelector('.context-marker'); if (marker) marker.hidden = !current;
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
  const queryText = String(elements.mcpFilter?.value || '').trim(); const queryMatches = searchMatcher(queryText); const filter = elements.mcpStateFilter?.value || 'all';
  let issues = 0, visible = 0;
  elements.mcpBadge.textContent = String(mcps.length); setMetric(elements.mcpTotalCount, mcps.length); setMetric(elements.mcpListTotalCount, mcps.length); setMetric(elements.mcpDefaultCount, mcps.filter((mcp) => mcp.default_include_in_environment).length); setMetric(elements.mcpEnvironmentCount, environment ? mcps.filter((mcp) => selected.has(mcp.id)).length : 0);
  if (!mcps.length) { setMetric(elements.mcpIssueCount, 0); setMetric(elements.mcpVisibleCount, 0); emptyMessage(elements.mcpList, '暂无 MCP 定义。可以添加 typed MCP，或先预览再导入 JSON/JSONC。'); updateMCPBulkControls(); return; }
  elements.mcpList.replaceChildren(); elements.mcpList.classList.remove('empty');
  for (const entry of mcps) {
    const fact = managementCapabilityFacts.get(`mcp/${entry.id}`); const enabled = environment ? selected.has(entry.id) : false; const health = managementEnvironmentID ? mcpHealthByKey.get(mcpHealthKey(managementEnvironmentID, entry.id)) : null;
    const configState = mcpConfigurationState(entry, fact); const runtimeState = mcpRuntimeState(environment, enabled, health); const configNormalized = normalizedState(configState); const runtimeNormalized = normalizedState(runtimeState);
    const issue = configIssueStates.has(configNormalized) || runtimeIssueStates.has(runtimeNormalized); if (issue) issues++;
    const haystack = [entry.name, entry.id, entry.endpoint, entry.executable, entry.transport].filter(Boolean).join(' ');
    const filterMatch = filter === 'all' || (filter === 'selected' && Boolean(environment && enabled)) || (filter === 'unselected' && Boolean(environment && !enabled)) || (filter === 'issues' && issue) || (filter === 'unobserved' && runtimeNormalized === 'not_observed');
    if (!queryMatches(haystack) || !filterMatch) continue;
    visible++;
    const row = document.createElement('article'); row.className = 'resource-row'; row.dataset.id = entry.id || ''; if (entry.id === editingMCPID) row.dataset.editing = 'true';
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
  if (!visible) emptyMessage(elements.mcpList, queryText || filter !== 'all' ? '没有符合当前筛选条件的 MCP。' : '暂无 MCP 定义。');
  updateMCPBulkControls();
}
function skillAvailabilityState(entry, environment, selected) {
  const availability = skillAvailabilityByID.get(entry.id); if (availability?.state) return availability.state; if (!environment) return entry.artifact_path && entry.source_root ? 'configured' : 'unconfigured'; return selected ? 'not_observed' : 'disabled';
}
function skillCatalogFingerprint(skills = safeArray(currentSnapshot?.skills)) {
  return window.ADMSkillBulk?.catalogFingerprint(skills) || skills.map((skill) => skill.id || '').filter(Boolean).sort().join('|');
}
function currentSkillProbeIdentity() {
  return {connectionGeneration, scope: 'catalog', catalogFingerprint: skillCatalogFingerprint()};
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
function updateMCPBulkControls() {
  const visibleCount = visibleResourceIDs(elements.mcpList).length;
  const environment = currentEnvironment();
  const disabled = mcpBulkBusy || visibleCount === 0;
  elements.mcpSetVisibleDefaultButton.disabled = disabled;
  elements.mcpUnsetVisibleDefaultButton.disabled = disabled;
  elements.mcpEnableVisibleButton.disabled = disabled || !environment;
  elements.mcpDisableVisibleButton.disabled = disabled || !environment;
  elements.mcpBulkHint.textContent = environment
    ? `批量操作只作用于当前筛选结果（${visibleCount} 项）；当前 Environment：${environment.name || environment.environment_id}。`
    : `批量默认值只作用于当前筛选结果（${visibleCount} 项）；当前 Environment 启用/取消需要先选择 Environment。`;
}
function updateSkillBulkControls() {
  pruneSkillSelection();
  const environment = currentEnvironment(); const selectedCount = selectedSkillIDs.size; const probeCurrent = explicitSkillProbeIsCurrent(); const summary = probeCurrent ? skillAvailabilitySummary() : null; const hasSkills = Boolean(safeArray(currentSnapshot?.skills).length); const visibleCount = visibleResourceIDs(elements.skillList).length;
  elements.skillSelectedCount.textContent = String(selectedCount);
  elements.skillProbeAllButton.disabled = skillBulkBusy || !hasSkills;
  elements.skillSelectVisibleButton.disabled = skillBulkBusy || !safeArray(currentSnapshot?.skills).length;
  elements.skillClearSelectionButton.disabled = skillBulkBusy || selectedCount === 0;
  elements.skillDeleteSelectedButton.disabled = skillBulkBusy || selectedCount === 0;
  elements.skillSetVisibleDefaultButton.disabled = skillBulkBusy || visibleCount === 0;
  elements.skillUnsetVisibleDefaultButton.disabled = skillBulkBusy || visibleCount === 0;
  elements.skillEnableVisibleButton.disabled = skillBulkBusy || visibleCount === 0 || !environment;
  elements.skillDisableVisibleButton.disabled = skillBulkBusy || visibleCount === 0 || !environment;
  elements.skillClearUnavailableButton.disabled = skillBulkBusy || !probeCurrent || !summary?.cleanupIDs?.length;
  elements.skillClearUnavailableButton.textContent = probeCurrent && summary?.cleanupIDs?.length ? `一键清除不可用 (${summary.cleanupIDs.length})` : '一键清除不可用';
  if (!probeCurrent) elements.skillBulkHint.textContent = '批量检查会检查全局 Skill catalog 的 source root、artifact 与 support roots；不依赖当前 Environment。';
  else elements.skillBulkHint.textContent = `最近全局检查：可用 ${summary.available} · 不可用 ${summary.unavailable} · 未知 ${summary.unknown}。清理只删除 ADM catalog metadata，不删除磁盘文件。`;
}
async function probeAllSkillAvailability() {
  const identity = currentSkillProbeIdentity();
  if (!safeArray(currentSnapshot?.skills).length) return setStatus('当前 catalog 没有 Skill 可检查。', 'success');
  skillBulkBusy = true; explicitSkillAvailabilityProbe = null; updateSkillBulkControls(); setStatus('正在批量检查全局 Skill catalog 可用性…', 'loading');
  try {
    const result = await desktopAdapter().ListSkillAvailability();
    if (identity.connectionGeneration !== connectionGeneration || identity.catalogFingerprint !== skillCatalogFingerprint()) return false;
    skillAvailabilityByID = new Map(); for (const item of safeArray(result?.skills)) if (item?.skill_id) skillAvailabilityByID.set(item.skill_id, item);
    explicitSkillAvailabilityProbe = {...identity, checkedAt: Date.now()};
    renderSkillManager(safeArray(currentSnapshot?.skills));
    const summary = skillAvailabilitySummary();
    setStatus(`Skill 全局可用性检查完成：可用 ${summary.available} · 不可用 ${summary.unavailable} · 未知 ${summary.unknown}`, summary.unavailable ? 'error' : 'success');
    return true;
  } catch (error) {
    if (identity.connectionGeneration === connectionGeneration) { explicitSkillAvailabilityProbe = null; updateSkillBulkControls(); setStatus(`Skill 全局可用性检查失败：${errorText(error)}`, 'error'); }
    return false;
  } finally { skillBulkBusy = false; renderSkillManager(safeArray(currentSnapshot?.skills)); }
}
async function runVisibleBatch(kind, label, ids, mutate, rerender) {
  const uniqueIDs = [...new Set(safeArray(ids))].filter(Boolean);
  if (!uniqueIDs.length) return setStatus('当前筛选结果为空，没有可批量操作的项目。', 'success');
  if (!window.confirm(label + '：当前筛选结果中的 ' + uniqueIDs.length + ' 项？')) return;
  const requestGeneration = connectionGeneration;
  if (kind === 'mcp') mcpBulkBusy = true;
  else skillBulkBusy = true;
  rerender();
  setStatus(label + '…', 'loading');
  const failures = []; let changed = 0;
  try {
    for (const id of uniqueIDs) {
      if (requestGeneration !== connectionGeneration) { failures.push({id, error: 'connection changed before batch completed'}); break; }
      try { await mutate(id); changed++; }
      catch (error) { failures.push({id, error: errorText(error)}); }
    }
    if (requestGeneration === connectionGeneration) await refreshSnapshot();
    const message = failures.length
      ? label + '完成：已处理 ' + changed + '/' + uniqueIDs.length + '；失败 ' + failures.length + '。' + failures.slice(0, 3).map((item) => item.id + ': ' + item.error).join(' · ')
      : label + '完成：已处理 ' + changed + ' 项。';
    setStatus(message, failures.length ? 'error' : 'success');
  } finally {
    if (kind === 'mcp') { mcpBulkBusy = false; renderMCPManager(safeArray(currentSnapshot?.mcps)); }
    else { skillBulkBusy = false; renderSkillManager(safeArray(currentSnapshot?.skills)); }
  }
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
function resetSkillSourceEditor(close = false, rerender = true) {
  editingSkillSourceID = '';
  elements.skillSourceID.value = '';
  elements.skillSourceForm.reset();
  elements.skillSourceDialogTitle.textContent = '添加 Skill source';
  elements.skillSourceSubmitButton.textContent = '添加并刷新 Source';
  if (close) closeEditorDialog('skillSourceDialog');
  if (rerender) renderSkillSources();
}
function beginSkillSourceEdit(source) {
  if (!source) return;
  editingSkillSourceID = source.skill_source_id || '';
  elements.skillSourceID.value = editingSkillSourceID;
  elements.skillSourceRoot.value = source.root || '';
  elements.skillSupportRoots.value = safeArray(source.support_roots).join('\n');
  elements.skillSourceDefault.checked = Boolean(source.default_include_in_environment);
  elements.skillSourceDialogTitle.textContent = '编辑 Skill source';
  elements.skillSourceSubmitButton.textContent = '保存 Source';
  openEditorDialog('skillSourceDialog');
}
function setSkillSubview(next) {
  skillSubview = next === 'sources' ? 'sources' : 'skills';
  syncSkillSubviewUI();
}
function syncSkillSubviewUI() {
  for (const button of elements.skillSubviewTabs?.querySelectorAll('[data-skill-subview]') || []) button.setAttribute('aria-selected', button.dataset.skillSubview === skillSubview ? 'true' : 'false');
  for (const panel of [elements.skillsPanel, elements.skillSourcesPanel]) if (panel) panel.hidden = panel.dataset.skillSubviewPanel !== skillSubview;
}
function skillSourceDisplayName(sourceID) {
  if (!sourceID) return 'Legacy';
  const source = skillSources.find((item) => item.skill_source_id === sourceID);
  if (!source) return sourceID;
  const root = source.root || source.skill_source_id || sourceID;
  return root.split(/[\\/]/).filter(Boolean).pop() || root;
}
function renderSkillSourceFilterOptions(skills) {
  if (!elements.skillSourceFilter) return;
  const current = elements.skillSourceFilter.value;
  const seen = new Set();
  const options = [new Option('全部 Source', ''), new Option('Legacy', '__legacy__')];
  for (const skill of skills) {
    const sourceID = skill.source_id || '';
    if (!sourceID || seen.has(sourceID)) continue;
    seen.add(sourceID);
    options.push(new Option(skillSourceDisplayName(sourceID), sourceID));
  }
  elements.skillSourceFilter.replaceChildren(...options);
  elements.skillSourceFilter.value = current && (current === '__legacy__' || seen.has(current)) ? current : '';
}
function appendSkillFact(parent, label, value) {
  const item = document.createElement('small');
  item.className = 'resource-fact';
  item.textContent = label + '：' + textOrDash(value);
  parent.append(item);
}
function renderSkillSources() {
  const queryText = String(elements.skillSourceFilterInput?.value || '').trim();
  const queryMatches = searchMatcher(queryText);
  if (skillSourcesState === 'loading') {
    setMetric(elements.skillSourceCount, '—'); setMetric(elements.skillSubviewSourceCount, '—'); setMetric(elements.skillSourceListTotalCount, '—'); setMetric(elements.skillSourceVisibleCount, 0);
    return emptyMessage(elements.skillSourceList, '正在读取 Skill sources…');
  }
  if (skillSourcesState === 'error') {
    setMetric(elements.skillSourceCount, '—'); setMetric(elements.skillSubviewSourceCount, '—'); setMetric(elements.skillSourceListTotalCount, '—'); setMetric(elements.skillSourceVisibleCount, 0);
    return emptyMessage(elements.skillSourceList, 'Skill source 列表不可用：' + (skillSourcesError || 'unknown error'));
  }
  if (skillSourcesState !== 'success') {
    setMetric(elements.skillSourceCount, '—'); setMetric(elements.skillSubviewSourceCount, '—'); setMetric(elements.skillSourceListTotalCount, '—'); setMetric(elements.skillSourceVisibleCount, 0);
    return emptyMessage(elements.skillSourceList, 'Skill sources 尚未加载');
  }
  setMetric(elements.skillSourceCount, skillSources.length); setMetric(elements.skillSubviewSourceCount, skillSources.length); setMetric(elements.skillSourceListTotalCount, skillSources.length);
  if (!skillSources.length) {
    setMetric(elements.skillSourceVisibleCount, 0);
    const legacyCount = safeArray(currentSnapshot?.skills).filter((skill) => !skill.source_id).length;
    return emptyMessage(elements.skillSourceList, legacyCount ? '暂无 Skill source。下方 ' + legacyCount + ' 个 Skill 是历史 legacy 条目；添加 source root 后新发现项会标记为 source-managed。' : '暂无 Skill source。添加显式 source root 后由 Core 扫描 SKILL.md。');
  }
  elements.skillSourceList.replaceChildren(); elements.skillSourceList.classList.remove('empty');
  let visible = 0;
  for (const source of skillSources) {
    const discovered = safeArray(currentSnapshot?.skills).filter((skill) => skill.source_id === source.skill_source_id);
    const haystack = [source.skill_source_id, source.root, ...safeArray(source.support_roots)].filter(Boolean).join(' ');
    if (!queryMatches(haystack)) continue;
    visible++;
    const status = source.last_refresh_status || 'pending';
    const row = document.createElement('article'); row.className = 'resource-row'; row.dataset.id = source.skill_source_id || '';
    const main = document.createElement('div'); main.className = 'resource-main';
    const {header, id} = resourceHeader(source.root || source.skill_source_id, source.skill_source_id, [stateBadge(humanRefreshState(status), status), stateBadge(discovered.length + ' skills', 'count')]);
    const facts = document.createElement('div'); facts.className = 'resource-facts';
    appendSkillFact(facts, 'Root', source.root);
    appendSkillFact(facts, 'Support roots', safeArray(source.support_roots).length);
    appendSkillFact(facts, '新环境默认', source.default_include_in_environment ? '是' : '否');
    if (source.last_refresh_at) appendSkillFact(facts, 'Last refresh', new Date(source.last_refresh_at).toLocaleString());
    const note = document.createElement('small'); note.textContent = source.last_refresh_error || '刷新只更新这个 source；失败不会污染其他 Skill。';
    main.append(header, id, facts, note);
    const controls = document.createElement('div'); controls.className = 'resource-actions';
    controls.append(createActionButton('编辑 Source', 'edit-skill-source', source.skill_source_id), createActionButton('刷新 Source', 'refresh-skill-source', source.skill_source_id), createActionButton('删除 Source', 'remove-skill-source', source.skill_source_id, 'danger'));
    row.append(main, controls); elements.skillSourceList.append(row);
  }
  setMetric(elements.skillSourceVisibleCount, visible);
  if (!visible) emptyMessage(elements.skillSourceList, queryText ? '没有符合当前筛选条件的 Skill source。' : '暂无 Skill source。');
}
function renderSkillManager(skills) {
  pruneSkillSelection();
  syncSkillSubviewUI();
  renderSkillSourceFilterOptions(skills);
  const environment = currentEnvironment(); const selected = new Set(safeArray(environment?.enabled_skill_ids));
  const queryText = String(elements.skillFilter?.value || '').trim(); const queryMatches = searchMatcher(queryText); const filter = elements.skillStateFilter?.value || 'all'; const sourceFilter = elements.skillSourceFilter?.value || '';
  let issues = 0, visible = 0;
  elements.skillBadge.textContent = String(skills.length); setMetric(elements.skillTotalCount, skills.length); setMetric(elements.skillSubviewSkillCount, skills.length); setMetric(elements.skillListTotalCount, skills.length); setMetric(elements.skillEnvironmentCount, environment ? skills.filter((skill) => selected.has(skill.id)).length : 0);
  renderSkillSources();
  if (!skills.length) { setMetric(elements.skillIssueCount, 0); setMetric(elements.skillVisibleCount, 0); emptyMessage(elements.skillList, '暂无已发现 Skill。添加并刷新 Skill source 后会显示在这里。'); updateSkillBulkControls(); return; }
  elements.skillList.replaceChildren(); elements.skillList.classList.remove('empty');
  for (const entry of skills) {
    const enabled = environment ? selected.has(entry.id) : false; const availability = skillAvailabilityByID.get(entry.id); const state = skillAvailabilityState(entry, environment, enabled); const normalized = normalizedState(state); const issue = normalized === 'degraded' || Boolean(window.ADMSkillBulk?.isCleanupState(normalized)) || ['unavailable', 'unconfigured', 'error'].includes(normalized); if (issue) issues++;
    const haystack = [entry.name, entry.id, entry.relative_artifact_path, entry.artifact_path, entry.source_root, entry.source_id].filter(Boolean).join(' ');
    const sourceMatch = !sourceFilter || (sourceFilter === '__legacy__' ? !entry.source_id : entry.source_id === sourceFilter);
    const filterMatch = filter === 'all' || (filter === 'selected' && Boolean(environment && enabled)) || (filter === 'unselected' && Boolean(environment && !enabled)) || (filter === 'issues' && issue) || (filter === 'source' && Boolean(entry.source_id)) || (filter === 'legacy' && !entry.source_id);
    if (!queryMatches(haystack) || !sourceMatch || !filterMatch) continue;
    visible++;
    const fact = managementCapabilityFacts.get('skill/' + entry.id); const row = document.createElement('article'); row.className = 'resource-row'; row.dataset.id = entry.id || ''; const main = document.createElement('div'); main.className = 'resource-main';
    const {header, id} = resourceHeader(entry.name || entry.id, entry.id, [stateBadge('可用性 · ' + humanRuntimeState(state), state), entry.source_id ? stateBadge('Source 管理', 'source') : stateBadge('Legacy', 'legacy')]);
    const facts = document.createElement('div'); facts.className = 'resource-facts';
    appendSkillFact(facts, 'Source', entry.source_id ? skillSourceDisplayName(entry.source_id) + ' · ' + entry.source_id : 'Legacy metadata');
    appendSkillFact(facts, 'Artifact', entry.relative_artifact_path || entry.artifact_path);
    appendSkillFact(facts, 'Support roots', safeArray(entry.support_roots).length);
    appendSkillFact(facts, '新环境默认', entry.default_include_in_environment ? '是' : '否');
    appendSkillFact(facts, '当前 Environment', environment ? (enabled ? '已启用' : '未启用') : '未选择');
    appendSkillFact(facts, '全局可用性', humanRuntimeState(state));
    const note = document.createElement('small'); const supportMissing = safeArray(availability?.missing_support_roots);
    const availabilityNote = availability?.reason || fact?.message || (supportMissing.length ? 'Missing support roots: ' + supportMissing.join(', ') : ''); note.textContent = availabilityNote ? '可用性：' + availabilityNote : '';
    main.append(header, id, facts); if (note.textContent) main.append(note);
    const controls = document.createElement('div'); controls.className = 'resource-actions';
    const selectControl = checkControl('选择', selectedSkillIDs.has(entry.id), 'select-skill', entry.id, skillBulkBusy); selectControl.classList.add('skill-select-control');
    controls.append(selectControl, checkControl('新环境默认', entry.default_include_in_environment, 'default-skill', entry.id)); controls.append(checkControl(environment ? '当前环境启用' : '选择环境后启用', enabled, 'environment-skill', entry.id, !environment));
    if (!entry.source_id) controls.append(createActionButton('删除条目', 'remove-skill', entry.id, 'danger'));
    row.append(main, controls); elements.skillList.append(row);
  }
  setMetric(elements.skillIssueCount, issues); setMetric(elements.skillVisibleCount, visible);
  if (!visible) emptyMessage(elements.skillList, queryText || filter !== 'all' || sourceFilter ? '没有符合当前筛选条件的 Skill。' : '暂无已发现 Skill。');
  updateSkillBulkControls();
  updateMCPBulkControls();
}
function renderSnapshotBase(snapshot) {
  currentSnapshot = snapshot; const workspaces = safeArray(snapshot.workspaces), environments = safeArray(snapshot.environments), executables = safeArray(snapshot.allowed_executables), mcps = safeArray(snapshot.mcps), skills = safeArray(snapshot.skills);
  if (editingMCPID && !mcps.some((mcp) => mcp.id === editingMCPID)) resetMCPEditor(false, false);
  renderWorkspaces(workspaces); renderManagementEnvironmentOptions(environments); renderEnvironments(environments); renderExecutables(executables); renderMCPManager(mcps); renderSkillManager(skills);
  if (!globalMemoryLoaded) emptyMessage(elements.globalMemoryList, '尚未加载 Global Memory');
  renderDashboardState('success');
  if (selectedEnvironmentID && !environments.some((env) => env.environment_id === selectedEnvironmentID)) closeEnvironmentDetail();
}
function renderManagementUnavailable(message) {
  for (const element of [elements.workspaceBadge, elements.workspaceVisibleCount, elements.workspaceListTotalCount, elements.environmentBadge, elements.environmentVisibleCount, elements.environmentListTotalCount, elements.execBadge, elements.mcpBadge, elements.skillBadge, elements.mcpTotalCount, elements.mcpDefaultCount, elements.mcpEnvironmentCount, elements.mcpIssueCount, elements.mcpVisibleCount, elements.mcpListTotalCount, elements.skillSourceCount, elements.skillTotalCount, elements.skillEnvironmentCount, elements.skillIssueCount, elements.skillVisibleCount, elements.skillListTotalCount]) setMetric(element, '—');
  renderWorkspaceOptions([]); renderEnvironmentWorkspaceFilter([]);
  elements.managementEnvironment.replaceChildren(new Option('管理数据未加载', '')); elements.managementEnvironment.value = ''; elements.managementEnvironment.disabled = true; elements.editEnvironmentButton.disabled = true;
  elements.managementEnvironmentHint.textContent = message;
  emptyMessage(elements.workspaceList, message); emptyMessage(elements.environmentList, message); emptyMessage(elements.execList, message); emptyMessage(elements.mcpList, message); emptyMessage(elements.skillSourceList, message); emptyMessage(elements.skillList, message); emptyMessage(elements.globalMemoryList, message); emptyMessage(elements.environmentMemoryList, message); renderEnvironmentMemoryScope();
  elements.runtimeHint.textContent = message; resetRuntimeCollections('error', message); emptyMessage(elements.verifierList, message); emptyMessage(elements.processList, message); emptyMessage(elements.runList, message); clearRuntimeOutput(message);
  updateSkillBulkControls();
  updateMCPBulkControls();
}
function clearManagementData(message = 'ADM 未连接。连接 Admin MCP 后加载管理数据。') {
  connectionGeneration++; environmentGeneration++; detailGeneration++;
  resetMCPEditor(true, false);
  pendingMCPImport = null; elements.mcpImportApplyButton.disabled = true;
  elements.mcpImportContent.value = ''; resetMCPImportPreview();
  managementEnvironmentID = ''; currentSnapshot = null; lastSnapshotSuccessAt = 0;
  elements.workspaceFilter.value = ''; elements.environmentFilter.value = ''; elements.environmentWorkspaceFilter.value = '';
  skillSources = []; skillSourcesState = 'unloaded'; skillSourcesError = ''; managementContextError = ''; managementSkillAvailabilityError = '';
  selectedSkillIDs = new Set(); explicitSkillAvailabilityProbe = null; skillBulkBusy = false; mcpBulkBusy = false; editingSkillSourceID = ''; skillSubview = 'skills'; syncSkillSubviewUI();
  managementInspection = null; managementCapabilityFacts = new Map(); skillAvailabilityByID = new Map(); environmentSkillAvailabilityByID = new Map(); mcpHealthByKey = new Map();
  runtimeSubview = 'verifiers'; runtimeSubviewGeneration++; runtimePendingActionKey = ''; resetRuntimeCollections('unloaded'); clearRuntimeOutput('管理上下文已清除'); syncRuntimeSubviewUI();
  renderDashboardState('unloaded', message); renderManagementUnavailable(message);
  globalMemoryLoaded = false; globalMemoryLoading = false; resetEnvironmentMemoryScope(message); if (!elements.environmentDetailPanel.hidden) closeEnvironmentDetail();
}
async function refreshManagementContext(scope = captureEnvironmentScope()) {
  if (!environmentScopeIsCurrent(scope)) return {stale: true, errors: []};
  managementInspection = null; managementCapabilityFacts = new Map(); environmentSkillAvailabilityByID = new Map(); managementContextError = ''; managementSkillAvailabilityError = ''; updateManagementHint();
  if (!scope.environmentID) { renderMCPManager(safeArray(currentSnapshot?.mcps)); renderSkillManager(safeArray(currentSnapshot?.skills)); resetEnvironmentMemoryScope(); resetRuntimeCollections('unloaded'); renderRuntime(); return {stale: false, errors: [], inspection: null}; }
  const [inspectionResult, availabilityResult] = await Promise.allSettled([
    desktopAdapter().InspectEnvironment(scope.environmentID),
    desktopAdapter().ListEnvironmentSkills(scope.environmentID),
  ]);
  if (!environmentScopeIsCurrent(scope)) return {stale: true, errors: []};
  const errors = [];
  if (inspectionResult.status === 'fulfilled') { managementInspection = inspectionResult.value; managementCapabilityFacts = capabilityMap(inspectionResult.value?.capability_report); }
  else errors.push(`Environment inspection: ${errorText(inspectionResult.reason)}`);
  if (availabilityResult.status === 'fulfilled') { for (const item of safeArray(availabilityResult.value?.skills)) if (item?.skill_id) environmentSkillAvailabilityByID.set(item.skill_id, item); }
  else { managementSkillAvailabilityError = errorText(availabilityResult.reason); errors.push(`Skill availability: ${managementSkillAvailabilityError}`); }
  managementContextError = errors.join(' · '); updateManagementHint();
  renderMCPManager(safeArray(currentSnapshot?.mcps)); renderSkillManager(safeArray(currentSnapshot?.skills));
  const runtimeResult = await refreshRuntimeContext(false, scope);
  return {stale: Boolean(runtimeResult?.stale), errors: [...errors, ...safeArray(runtimeResult?.errors)], inspection: managementInspection};
}

function currentMCPImportInput() {
  return {
    format: elements.mcpImportFormat.value,
    json_or_jsonc: elements.mcpImportContent.value,
    conflict_policy: elements.mcpImportConflict.value,
    default_include: elements.mcpImportDefault.checked,
  };
}
function mcpImportFingerprint(input) {
  return JSON.stringify({
    format: String(input?.format || ''),
    json_or_jsonc: String(input?.json_or_jsonc || ''),
    conflict_policy: String(input?.conflict_policy || ''),
    default_include: Boolean(input?.default_include),
  });
}
function resetMCPImportPreview(message = '尚未预览。导入不会修改 Environment 选择。') {
  pendingMCPImport = null;
  mcpImportBusy = false;
  elements.mcpImportApplyButton.disabled = true;
  emptyMessage(elements.mcpImportPreview, message);
}
function invalidateMCPImportPreview(message = '导入内容或选项已变化，请重新预览。') {
  if (!pendingMCPImport) return;
  resetMCPImportPreview(message);
}
function currentPendingMCPImport() {
  const pending = pendingMCPImport;
  if (!pending) return null;
  if (pending.connectionGeneration !== connectionGeneration || pending.fingerprint !== mcpImportFingerprint(currentMCPImportInput())) {
    resetMCPImportPreview('导入预览已过期，请重新预览。');
    return null;
  }
  return pending;
}
function renderImportPreview(preview) {
  const candidates = safeArray(preview?.candidates); elements.mcpImportPreview.replaceChildren(); elements.mcpImportPreview.classList.remove('empty');
  if (!candidates.length) { resetMCPImportPreview('没有可导入候选项。'); return; }
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
function detailGroup(title, rows, note = '') {
  const section = document.createElement('section'); section.className = 'detail-group'; const heading = document.createElement('div'); heading.className = 'detail-group-heading'; const label = document.createElement('strong'); label.textContent = title; heading.append(label);
  if (note) { const description = document.createElement('span'); description.textContent = note; heading.append(description); }
  const grid = document.createElement('div'); grid.className = 'detail-group-grid'; grid.append(...rows); section.append(heading, grid); return section;
}
function renderSelectionList(container, entries, selectedIDs, kind, facts, availabilityMap, availabilityError = '') {
  if (!entries.length) return emptyMessage(container, `暂无 ${kind.toUpperCase()} 项`); container.replaceChildren(); container.classList.remove('empty'); const selected = new Set(safeArray(selectedIDs));
  for (const entry of entries) {
    const row = document.createElement('label'); row.className = 'selection-row rich-selection'; const configured = kind === 'mcp' ? mcpConfigured(entry) : Boolean(entry.artifact_path && entry.source_root); const enabled = selected.has(entry.id);
    const checkbox = document.createElement('input'); checkbox.type = 'checkbox'; checkbox.checked = enabled; checkbox.disabled = !configured; checkbox.dataset.kind = kind; checkbox.dataset.id = entry.id;
    const text = document.createElement('span'); text.className = 'selection-copy'; const name = document.createElement('strong'); name.textContent = entry.name || entry.id; const id = document.createElement('code'); id.textContent = entry.id || '';
    let state = enabled ? 'configured' : 'disabled', reason = '';
    if (kind === 'mcp') {
      const fact = facts.get(`mcp/${entry.id}`); state = !configured ? 'unconfigured' : !enabled ? 'disabled' : fact?.state || 'configured'; reason = fact?.message || fact?.reason_code || '';
    } else {
      const availability = availabilityMap.get(entry.id);
      if (availability?.state) { state = availability.state; reason = availability.reason || ''; }
      else if (!configured) state = 'unconfigured';
      else if (!enabled) state = 'disabled';
      else if (availabilityError) { state = 'unknown'; reason = `可用性读取失败：${availabilityError}`; }
      else state = facts.get(`skill/${entry.id}`)?.state || 'configured';
      if (!reason) reason = facts.get(`skill/${entry.id}`)?.message || '';
    }
    const meta = document.createElement('small'); meta.textContent = reason || (configured ? '配置可用' : '缺少配置'); text.append(name, id, meta); row.append(checkbox, text, stateBadge(state, state)); container.append(row);
  }
}
function syncEnvironmentDetailSubviewUI() {
  environmentDetailSubview = environmentDetailSubview === 'diagnostics' ? 'diagnostics' : 'summary';
  for (const button of elements.environmentDetailSubviewTabs.querySelectorAll('[data-environment-detail-subview]')) {
    const active = button.dataset.environmentDetailSubview === environmentDetailSubview;
    button.setAttribute('aria-selected', String(active));
    button.tabIndex = active ? 0 : -1;
  }
  elements.environmentDetail.hidden = environmentDetailSubview !== 'summary';
  elements.environmentDiagnostics.hidden = environmentDetailSubview !== 'diagnostics';
}
function diagnosticFactText(fact) {
  const parts = [humanRuntimeState(fact.state || 'unknown')];
  if (fact.reason_code) parts.push('reason ' + fact.reason_code);
  if (fact.message) parts.push(fact.message);
  if (fact.source) parts.push('source ' + fact.source);
  if (fact.generated_at) parts.push('generated ' + new Date(fact.generated_at).toLocaleString());
  if (fact.observed_at) parts.push('observed ' + new Date(fact.observed_at).toLocaleString());
  return parts.join(' · ');
}
function renderEnvironmentDiagnostics(inspection) {
  const environment = inspection?.environment || {}, report = inspection?.capability_report || {};
  const facts = safeArray(report.facts);
  const unresolvedMCP = safeArray(inspection?.unresolved_mcp_ids), unresolvedSkill = safeArray(inspection?.unresolved_skill_ids);
  const sourceRows = [
    detailRow('Source', 'Existing InspectEnvironment payload only'),
    detailRow('Generated at', report.generated_at ? new Date(report.generated_at).toLocaleString() : 'Unknown'),
    detailRow('Environment ID', environment.environment_id || '—'),
    detailRow('Capability facts', String(facts.length)),
  ];
  const unresolvedRows = [
    detailRow('MCP IDs', unresolvedMCP.join(', ') || 'None'),
    detailRow('Skill IDs', unresolvedSkill.join(', ') || 'None'),
  ];
  const grouped = new Map();
  for (const fact of facts) {
    const group = fact.kind || String(fact.key || 'capability').split('/')[0] || 'capability';
    if (!grouped.has(group)) grouped.set(group, []);
    grouped.get(group).push(fact);
  }
  const groups = [
    detailGroup('Diagnostic source', sourceRows, 'Diagnostics here never probes MCPs, runs verifiers, reconnects, or reads Memory values.'),
    detailGroup('Unresolved references', unresolvedRows, 'Catalog removals surface as unresolved IDs; the view does not rewrite Environment selections.'),
  ];
  if (!facts.length) {
    groups.push(detailGroup('Capability facts', [detailRow('Returned facts', 'None')], 'Unknown means no fact was returned; it is not an all-green live health claim.'));
  } else {
    for (const [group, items] of grouped) {
      const container = document.createElement('section'); container.className = 'detail-group';
      const heading = document.createElement('h3'); heading.textContent = group + ' facts';
      const list = document.createElement('div'); list.className = 'diagnostic-fact-list';
      for (const fact of items) {
        const row = document.createElement('article'); row.className = 'diagnostic-fact';
        const title = document.createElement('strong'); title.textContent = fact.key || fact.kind || 'capability';
        const badge = stateBadge(humanRuntimeState(fact.state || 'unknown'), fact.state || 'unknown');
        const detail = document.createElement('small'); detail.textContent = diagnosticFactText(fact);
        row.append(title, badge, detail); list.append(row);
      }
      const note = document.createElement('small'); note.textContent = 'Static/observed labels come only from returned fact fields; filtering this view performs no new check.';
      container.append(heading, list, note); groups.push(container);
    }
  }
  elements.environmentDiagnostics.replaceChildren(...groups);
}
function renderEnvironmentDetailFromInspection(inspection, token = detailGeneration) {
  const environment = inspection?.environment || {}, workspace = inspection?.workspace || {}, report = inspection?.capability_report || {}, facts = capabilityMap(report);
  if (token !== detailGeneration || environment.environment_id !== selectedEnvironmentID) return false;
  elements.environmentDetailTitle.textContent = environment.name || environment.environment_id || 'Environment';
  const issueFacts = safeArray(report.facts).filter((fact) => ['unavailable', 'degraded', 'unconfigured'].includes(normalizedState(fact.state)));
  const capabilityRows = issueFacts.length
    ? issueFacts.slice(0, 8).map((fact) => detailRow(fact.key || fact.kind || 'capability', humanRuntimeState(fact.state) + (fact.reason_code ? ' · ' + fact.reason_code : '') + (fact.message ? ' · ' + fact.message : '')))
    : [detailRow('Status', '没有 unavailable / degraded / unconfigured capability fact')];
  if (managementSkillAvailabilityError) capabilityRows.unshift(detailRow('Skill availability', '读取失败 · ' + managementSkillAvailabilityError));
  if (report.generated_at) capabilityRows.push(detailRow('Generated at', new Date(report.generated_at).toLocaleString()));
  const identityRows = [
    detailRow('Environment ID', environment.environment_id),
    detailRow('Root', environment.root),
    detailRow('Workspace', (workspace.name || workspace.workspace_id || '—') + (workspace.workspace_id ? ' (' + workspace.workspace_id + ')' : '') + ' · ' + (workspace.path || '')),
    detailRow('State', environment.state),
  ];
  const authorityRows = [
    detailRow('Writer observation', environment.writer?.owner || 'No active writer observed'),
    detailRow('Private Memory count', safeNumber(environment.private_memory_count) + ' entries'),
  ];
  const unresolvedRows = [
    detailRow('MCP IDs', safeArray(inspection?.unresolved_mcp_ids).join(', ') || 'None'),
    detailRow('Skill IDs', safeArray(inspection?.unresolved_skill_ids).join(', ') || 'None'),
  ];
  elements.environmentDetail.replaceChildren(
    detailGroup('Identity', identityRows, 'Stable identity/root/Workspace facts from the current inspection.'),
    detailGroup('Runtime authority', authorityRows, 'Writer is an observation only; this page never acquires or force-releases a lease.'),
    detailGroup('Capability issues', capabilityRows, issueFacts.length + ' issue fact(s); optional failures stay local.'),
    detailGroup('Unresolved references', unresolvedRows, 'Catalog removals never silently rewrite Environment IDs.')
  );
  renderEnvironmentDiagnostics(inspection);
  syncEnvironmentDetailSubviewUI();
  renderSelectionList(elements.environmentMCPSelections, safeArray(currentSnapshot?.mcps), environment.enabled_mcp_ids, 'mcp', facts, skillAvailabilityByID);
  renderSelectionList(elements.environmentSkillSelections, safeArray(currentSnapshot?.skills), environment.enabled_skill_ids, 'skill', facts, environmentSkillAvailabilityByID, managementSkillAvailabilityError);
  elements.environmentDetailBackdrop.hidden = false; elements.environmentDetailPanel.hidden = false; return true;
}
function closeEnvironmentDetail() {
  const wasOpen = !elements.environmentDetailPanel.hidden; const opener = environmentDetailOpener;
  detailGeneration++; selectedEnvironmentID = ''; environmentDetailOpener = null; environmentDetailSubview = 'summary'; syncEnvironmentDetailSubviewUI(); elements.environmentDetailBackdrop.hidden = true; elements.environmentDetailPanel.hidden = true; elements.environmentDetail.replaceChildren(); elements.environmentDiagnostics.replaceChildren();
  if (!wasOpen) return;
  updateEnvironmentContextMarkers();
  if (opener?.isConnected && !opener.closest('[hidden]')) opener.focus({preventScroll: true}); else document.querySelector('[data-management-page]:not([hidden]) [data-page-heading]')?.focus({preventScroll: true});
}
function refreshSelectedEnvironmentDetail() { if (!selectedEnvironmentID) return false; const token = detailGeneration, id = selectedEnvironmentID, inspection = managementInspection; if (token !== detailGeneration || id !== selectedEnvironmentID || inspection?.environment?.environment_id !== id) return false; return renderEnvironmentDetailFromInspection(inspection, token); }

function renderMemory(container, entries, scope) {
  if (!entries.length) return emptyMessage(container, '没有 Memory 条目'); container.replaceChildren(); container.classList.remove('empty');
  for (const entry of entries) { const row = document.createElement('div'); row.className = 'memory-row'; const content = document.createElement('div'); content.className = 'memory-value'; const key = document.createElement('code'); key.textContent = entry.key || ''; const value = document.createElement('pre'); value.textContent = entry.value || ''; content.append(key, value); const edit = createActionButton('编辑', 'edit-' + scope + '-memory', entry.key); edit.addEventListener('click', () => { const isGlobal = scope === 'global'; (isGlobal ? elements.globalMemoryKey : elements.environmentMemoryKey).value = entry.key || ''; (isGlobal ? elements.globalMemoryValue : elements.environmentMemoryValue).value = entry.value || ''; openEditorDialog(isGlobal ? 'globalMemoryDialog' : 'environmentMemoryDialog'); }); row.append(content, edit, createActionButton('删除', 'delete-' + scope + '-memory', entry.key, 'danger')); container.append(row); }
}
function currentManagementRoute() { return managementNavigation?.current?.() || window.ADMNavigation?.normalizeRoute?.(window.location.hash) || 'overview'; }
function currentGlobalMemoryScope() { return {connectionGeneration}; }
function globalMemoryScopeIsCurrent(scope) { return Boolean(scope && scope.connectionGeneration === connectionGeneration); }
function currentEnvironmentMemoryScope() { return {connectionGeneration, environmentGeneration, environmentID: managementEnvironmentID}; }
function environmentMemoryScopeIsCurrent(scope) { return Boolean(scope && scope.connectionGeneration === connectionGeneration && scope.environmentGeneration === environmentGeneration && scope.environmentID === managementEnvironmentID); }
function resetEnvironmentMemoryScope(message = '') {
  environmentMemoryLoaded = false;
  environmentMemoryLoading = false;
  loadedEnvironmentMemoryID = '';
  emptyMessage(elements.environmentMemoryList, message || (managementEnvironmentID ? '尚未加载当前 Environment Memory' : '请选择 Management Environment 后加载 private Memory'));
  renderEnvironmentMemoryScope();
}
function renderEnvironmentMemoryScope() {
  const environment = currentEnvironment();
  const enabled = Boolean(environment && currentSnapshot && !environmentMemoryLoading);
  elements.loadEnvironmentMemory.disabled = !enabled;
  elements.writeEnvironmentMemoryButton.disabled = !Boolean(environment && currentSnapshot);
  elements.environmentMemoryScopeHint.textContent = environment
    ? '当前 Environment：' + (environment.name || environment.environment_id) + ' · ' + environment.environment_id + '。private Memory 必须显式加载。'
    : '请选择 Management Environment 后显式加载。';
  if (!environment) emptyMessage(elements.environmentMemoryList, '请选择 Management Environment 后加载 private Memory');
}
function maybeAutoLoadGlobalMemory() {
  if (currentManagementRoute() !== 'memory' || globalMemoryLoaded || globalMemoryLoading || !currentSnapshot) return false;
  loadGlobalMemory({auto: true});
  return true;
}
async function loadGlobalMemory(options = {}) {
  if (globalMemoryLoading) return false;
  const scope = options.scope || currentGlobalMemoryScope();
  globalMemoryLoading = true;
  setStatus(options.auto ? '正在加载 Global Memory…' : '正在显式读取 Global Memory…', 'loading');
  try {
    const entries = safeArray(await desktopAdapter().ListGlobalMemory());
    if (!globalMemoryScopeIsCurrent(scope)) return false;
    renderMemory(elements.globalMemoryList, entries, 'global');
    globalMemoryLoaded = true;
    setStatus('Global Memory 已加载', 'success');
    return true;
  } catch (error) {
    if (globalMemoryScopeIsCurrent(scope)) setStatus('Global Memory 读取失败：' + (error?.message || String(error)), 'error');
    return false;
  } finally {
    if (globalMemoryScopeIsCurrent(scope)) globalMemoryLoading = false;
  }
}
async function loadEnvironmentMemory(options = {}) {
  if (environmentMemoryLoading) return false;
  const scope = options.scope || currentEnvironmentMemoryScope();
  if (!scope.environmentID) { renderEnvironmentMemoryScope(); return false; }
  environmentMemoryLoading = true;
  renderEnvironmentMemoryScope();
  setStatus('正在显式读取 Environment-private Memory…', 'loading');
  try {
    const entries = safeArray(await desktopAdapter().ListEnvironmentMemory(scope.environmentID));
    if (!environmentMemoryScopeIsCurrent(scope)) return false;
    renderMemory(elements.environmentMemoryList, entries, 'environment');
    environmentMemoryLoaded = true;
    loadedEnvironmentMemoryID = scope.environmentID;
    setStatus('Environment-private Memory 已加载', 'success');
    return true;
  } catch (error) {
    if (environmentMemoryScopeIsCurrent(scope)) setStatus('Environment-private Memory 读取失败：' + (error?.message || String(error)), 'error');
    return false;
  } finally {
    if (environmentMemoryScopeIsCurrent(scope)) { environmentMemoryLoading = false; renderEnvironmentMemoryScope(); }
  }
}

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
      resetEnvironmentMemoryScope(message);
      if (!currentSnapshot) renderManagementUnavailable(message);
      setStatus(message, 'error'); return false;
    }
    if (requestGeneration !== connectionGeneration) return false;
    lastSnapshotSuccessAt = Date.now(); explicitSkillAvailabilityProbe = null; skillSourcesState = 'loading'; skillSourcesError = ''; renderSnapshotBase(snapshot); renderEnvironmentMemoryScope();

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
elements.workspaceFilter.addEventListener('input', () => renderWorkspaces(safeArray(currentSnapshot?.workspaces)));
elements.environmentFilter.addEventListener('input', () => renderEnvironments(safeArray(currentSnapshot?.environments)));
elements.environmentWorkspaceFilter.addEventListener('change', () => renderEnvironments(safeArray(currentSnapshot?.environments)));
elements.mcpFilter.addEventListener('input', () => renderMCPManager(safeArray(currentSnapshot?.mcps)));
elements.mcpStateFilter.addEventListener('change', () => renderMCPManager(safeArray(currentSnapshot?.mcps)));
for (const element of [elements.mcpImportFormat, elements.mcpImportConflict, elements.mcpImportDefault]) element.addEventListener('change', () => invalidateMCPImportPreview());
elements.mcpImportContent.addEventListener('input', () => invalidateMCPImportPreview());
elements.mcpSetVisibleDefaultButton.addEventListener('click', () => runVisibleBatch('mcp', '批量设置 MCP 新环境默认值', visibleResourceIDs(elements.mcpList), (id) => desktopAdapter().SetMCPDefault(id, true), () => renderMCPManager(safeArray(currentSnapshot?.mcps))));
elements.mcpUnsetVisibleDefaultButton.addEventListener('click', () => runVisibleBatch('mcp', '批量取消 MCP 新环境默认值', visibleResourceIDs(elements.mcpList), (id) => desktopAdapter().SetMCPDefault(id, false), () => renderMCPManager(safeArray(currentSnapshot?.mcps))));
elements.mcpEnableVisibleButton.addEventListener('click', () => { const environmentID = managementEnvironmentID; if (!environmentID) return setStatus('请先选择 Management Environment。', 'error'); return runVisibleBatch('mcp', '批量启用当前 Environment MCP', visibleResourceIDs(elements.mcpList), (id) => desktopAdapter().SetEnvironmentMCP(environmentID, id, true), () => renderMCPManager(safeArray(currentSnapshot?.mcps))); });
elements.mcpDisableVisibleButton.addEventListener('click', () => { const environmentID = managementEnvironmentID; if (!environmentID) return setStatus('请先选择 Management Environment。', 'error'); return runVisibleBatch('mcp', '批量取消当前 Environment MCP', visibleResourceIDs(elements.mcpList), (id) => desktopAdapter().SetEnvironmentMCP(environmentID, id, false), () => renderMCPManager(safeArray(currentSnapshot?.mcps))); });
elements.skillSubviewTabs.addEventListener('click', (event) => { const button = event.target.closest('[data-skill-subview]'); if (!button) return; setSkillSubview(button.dataset.skillSubview); });
elements.skillFilter.addEventListener('input', () => renderSkillManager(safeArray(currentSnapshot?.skills)));
elements.skillSourceFilter.addEventListener('change', () => renderSkillManager(safeArray(currentSnapshot?.skills)));
elements.skillSourceFilterInput.addEventListener('input', () => renderSkillSources());
elements.skillStateFilter.addEventListener('change', () => renderSkillManager(safeArray(currentSnapshot?.skills)));
elements.skillProbeAllButton.addEventListener('click', () => probeAllSkillAvailability());
elements.skillSelectVisibleButton.addEventListener('click', () => { for (const input of elements.skillList.querySelectorAll('input[data-action="select-skill"]')) selectedSkillIDs.add(input.dataset.id); renderSkillManager(safeArray(currentSnapshot?.skills)); });
elements.skillClearSelectionButton.addEventListener('click', () => { selectedSkillIDs.clear(); renderSkillManager(safeArray(currentSnapshot?.skills)); });
elements.skillSetVisibleDefaultButton.addEventListener('click', () => runVisibleBatch('skill', '批量设置 Skill 新环境默认值', visibleResourceIDs(elements.skillList), (id) => desktopAdapter().SetSkillDefault(id, true), () => renderSkillManager(safeArray(currentSnapshot?.skills))));
elements.skillUnsetVisibleDefaultButton.addEventListener('click', () => runVisibleBatch('skill', '批量取消 Skill 新环境默认值', visibleResourceIDs(elements.skillList), (id) => desktopAdapter().SetSkillDefault(id, false), () => renderSkillManager(safeArray(currentSnapshot?.skills))));
elements.skillEnableVisibleButton.addEventListener('click', () => { const environmentID = managementEnvironmentID; if (!environmentID) return setStatus('请先选择 Management Environment。', 'error'); return runVisibleBatch('skill', '批量启用当前 Environment Skill', visibleResourceIDs(elements.skillList), (id) => desktopAdapter().SetEnvironmentSkill(environmentID, id, true), () => renderSkillManager(safeArray(currentSnapshot?.skills))); });
elements.skillDisableVisibleButton.addEventListener('click', () => { const environmentID = managementEnvironmentID; if (!environmentID) return setStatus('请先选择 Management Environment。', 'error'); return runVisibleBatch('skill', '批量取消当前 Environment Skill', visibleResourceIDs(elements.skillList), (id) => desktopAdapter().SetEnvironmentSkill(environmentID, id, false), () => renderSkillManager(safeArray(currentSnapshot?.skills))); });
elements.skillDeleteSelectedButton.addEventListener('click', () => removeSkillCatalogEntries([...selectedSkillIDs], '批量删除'));
elements.skillClearUnavailableButton.addEventListener('click', () => {
  if (!explicitSkillProbeIsCurrent()) return setStatus('当前 Skill 可用性检查结果已过期，请重新批量检查后再清理。', 'error');
  const ids = skillAvailabilitySummary().cleanupIDs;
  if (!ids.length) return setStatus('当前显式检查没有发现可清理的不可用 Skill。', 'success');
  return removeSkillCatalogEntries(ids, '一键清除不可用 Skill');
});
elements.managementEnvironment.addEventListener('change', async () => {
  const nextEnvironmentID = elements.managementEnvironment.value; if (!elements.environmentDetailPanel.hidden) closeEnvironmentDetail();
  environmentGeneration++; managementEnvironmentID = nextEnvironmentID; managementContextError = ''; managementSkillAvailabilityError = ''; managementInspection = null; managementCapabilityFacts = new Map(); environmentSkillAvailabilityByID = new Map(); runtimePendingActionKey = ''; resetEnvironmentMemoryScope(nextEnvironmentID ? '尚未加载当前 Environment Memory' : '请选择 Management Environment 后加载 private Memory'); resetRuntimeCollections(nextEnvironmentID ? 'loading' : 'unloaded'); clearRuntimeOutput('Environment 已切换'); updateManagementHint(); updateSkillBulkControls(); updateEnvironmentContextMarkers(); renderRuntime();
  const scope = captureEnvironmentScope();
  setStatus('正在加载 Environment MCP/Skill/Runtime 状态…', 'loading');
  try { const result = await refreshManagementContext(scope); if (result?.stale) return; setStatus(result?.errors?.length ? `Environment 已切换；部分状态不可用：${result.errors.join(' · ')}` : 'Environment 管理上下文已切换', result?.errors?.length ? 'error' : 'success'); }
  catch (error) { if (environmentScopeIsCurrent(scope)) setStatus(`Environment 状态读取失败：${errorText(error)}`, 'error'); }
});
elements.editEnvironmentButton.addEventListener('click', () => {
  const environment = currentEnvironment();
  if (!environment) return;
  openRenameDialog('Environment', environment.name || '', (name) => desktopAdapter().RenameEnvironment(environment.environment_id, name));
});
elements.runtimeSubviewTabs.addEventListener('click', (event) => {
  const button = event.target.closest('button[data-runtime-subview]'); if (!button) return;
  const next = window.ADMRuntimeView?.normalizeSubview(button.dataset.runtimeSubview) || button.dataset.runtimeSubview;
  if (next === runtimeSubview) return;
  runtimeSubview = next; runtimeSubviewGeneration++; clearRuntimeOutput('已切换 Runtime 视图'); renderRuntime();
});
elements.runtimeRefreshButton.addEventListener('click', () => refreshRuntimeContext(true).catch((error) => setStatus(`Runtime 刷新失败：${error?.message || String(error)}`, 'error')));
elements.verifierList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action="run-verifier"]'); if (!button) return;
  const scope = captureEnvironmentScope(), environmentID = scope.environmentID, verifierID = button.dataset.id, owner = runtimeWriterOwner(), viewGeneration = runtimeSubviewGeneration;
  if (!environmentID) return; if (!owner) return setStatus('当前 Environment 没有 active writer，不能运行 verifier。', 'error');
  const result = await runRuntimeAction('运行 verifier', `verifier:${verifierID}`, scope, () => desktopAdapter().RunVerifier(environmentID, owner, verifierID));
  if (result && environmentScopeIsCurrent(scope) && runtimeSubview === 'verifiers' && runtimeSubviewGeneration === viewGeneration) {
    showRuntimeOutput('verifiers', verifierID, `Verifier ${verifierID} · ${result.status || 'unknown'}`, result.stdout || '', result.stderr || '', `${result.summary || ''}${result.exit_code !== undefined ? ` · exit ${result.exit_code}` : ''}`);
  }
});
elements.processList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return;
  const scope = captureEnvironmentScope(), environmentID = scope.environmentID, processID = button.dataset.id, viewGeneration = runtimeSubviewGeneration;
  if (!environmentID) return;
  if (button.dataset.action === 'process-logs') {
    if (runtimePendingActionKey) return; const actionKey = `process-logs:${processID}`; runtimePendingActionKey = actionKey; renderRuntime(); setStatus('正在读取 process 日志…', 'loading');
    try {
      const logs = await desktopAdapter().GetProcessLogs(environmentID, processID);
      if (!environmentScopeIsCurrent(scope) || runtimeSubview !== 'processes' || runtimeSubviewGeneration !== viewGeneration) return;
      const truncation = window.ADMRuntimeView?.truncationSummary(logs) || '';
      showRuntimeOutput('processes', processID, `Process ${processID}`, logs.stdout || '', logs.stderr || '', '', truncation); setStatus('Process 日志已读取', 'success');
    } catch (error) { if (environmentScopeIsCurrent(scope) && runtimeSubview === 'processes' && runtimeSubviewGeneration === viewGeneration) setStatus(`Process 日志读取失败：${error?.message || String(error)}`, 'error'); }
    finally { if (runtimePendingActionKey === actionKey) { runtimePendingActionKey = ''; renderRuntime(); } }
    return;
  }
  if (button.dataset.action === 'stop-process') {
    const owner = runtimeWriterOwner(); if (!owner) return setStatus('当前 Environment 没有 active writer，不能停止 process。', 'error');
    await runRuntimeAction('停止 process', `stop-process:${processID}`, scope, () => desktopAdapter().StopProcess(environmentID, owner, processID));
  }
});
elements.runList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return;
  const scope = captureEnvironmentScope(), environmentID = scope.environmentID, runID = button.dataset.id, run = runtimeRunsByID.get(runID);
  if (!environmentID) return;
  if (button.dataset.action === 'run-output') {
    if (!run) return; const truncation = window.ADMRuntimeView?.truncationSummary(run) || '';
    showRuntimeOutput('runs', run.id, `Run ${run.id} · ${run.state || 'unknown'}`, run.stdout || '', run.stderr || '', `${run.message || ''}${run.exit_code !== undefined && run.exit_code !== null ? ` · exit ${run.exit_code}` : ''}`, truncation); return;
  }
  if (button.dataset.action === 'cancel-run') {
    const owner = runtimeWriterOwner(); if (!owner) return setStatus('当前 Environment 没有 active writer，不能取消 run。', 'error');
    await runRuntimeAction('取消 run', `cancel-run:${runID}`, scope, () => desktopAdapter().CancelRun(environmentID, owner, runID));
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
  event.preventDefault();
  if (mcpImportBusy) return;
  const input = currentMCPImportInput();
  if (!input.json_or_jsonc.trim()) return;
  const fingerprint = mcpImportFingerprint(input);
  const requestGeneration = connectionGeneration;
  pendingMCPImport = null;
  elements.mcpImportApplyButton.disabled = true;
  mcpImportBusy = true;
  setStatus('正在预览 MCP 导入…', 'loading');
  try {
    const preview = await desktopAdapter().PreviewMCPImport(input);
    if (requestGeneration !== connectionGeneration || fingerprint !== mcpImportFingerprint(currentMCPImportInput())) {
      resetMCPImportPreview('导入预览已过期，请重新预览。');
      return;
    }
    pendingMCPImport = {input, fingerprint, connectionGeneration: requestGeneration};
    renderImportPreview(preview);
    if (!elements.mcpImportApplyButton.disabled) setStatus('导入预览已生成；确认后再应用', 'success');
    else setStatus('导入预览包含错误，不能应用', 'error');
  }
  catch (error) { resetMCPImportPreview('预览失败：' + errorText(error)); setStatus('MCP 导入预览失败：' + errorText(error), 'error'); }
  finally { mcpImportBusy = false; }
});
elements.mcpImportApplyButton.addEventListener('click', async () => {
  const pending = currentPendingMCPImport();
  if (!pending || mcpImportBusy) return;
  mcpImportBusy = true;
  elements.mcpImportApplyButton.disabled = true;
  pendingMCPImport = null;
  setStatus('正在应用 MCP 导入…', 'loading');
  try {
    const result = await desktopAdapter().ApplyMCPImport(pending.input);
    renderImportApplyResult(result);
    elements.mcpImportContent.value = '';
    closeFormDialog(elements.mcpImportForm);
    await refreshSnapshot('MCP 导入完成');
  }
  catch (error) { setStatus('MCP 导入失败：' + errorText(error), 'error'); resetMCPImportPreview('导入应用失败；输入已保留，请重新预览后再应用。'); }
  finally { mcpImportBusy = false; }
});
elements.skillSourceForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const root = elements.skillSourceRoot.value.trim(); if (!root) return;
  const supportRoots = lineValues(elements.skillSupportRoots.value), defaultInclude = elements.skillSourceDefault.checked, editingID = editingSkillSourceID;
  if (editingID) {
    await runMutation('更新 Skill source', async () => { await desktopAdapter().UpdateSkillSource(editingID, {root, support_roots: supportRoots, default_include_in_environment: defaultInclude}); resetSkillSourceEditor(true, false); });
    return;
  }
  setStatus('正在添加 Skill source…', 'loading');
  try {
    const source = await desktopAdapter().AddSkillSource({root, support_roots: supportRoots, default_include_in_environment: defaultInclude});
    elements.skillSourceForm.reset(); closeFormDialog(elements.skillSourceForm);
    try { await desktopAdapter().RefreshSkillSource(source.skill_source_id); elements.skillSourceForm.reset(); closeFormDialog(elements.skillSourceForm); await refreshSnapshot('Skill source 已添加并刷新'); }
    catch (refreshError) { await refreshSnapshot(); setStatus('Skill source 已保存，但刷新失败：' + errorText(refreshError), 'error'); }
  } catch (error) { setStatus('添加 Skill source 失败：' + errorText(error), 'error'); }
});

elements.mcpList.addEventListener('change', async (event) => {
  const input = event.target.closest('input[data-action]'); if (!input) return;
  if (input.dataset.action === 'default-mcp') {
    const id = input.dataset.id, checked = input.checked;
    await runMutation('更新 MCP 新环境默认值', () => desktopAdapter().SetMCPDefault(id, checked));
  }
  if (input.dataset.action === 'environment-mcp' && managementEnvironmentID) {
    const environmentID = managementEnvironmentID, id = input.dataset.id, checked = input.checked;
    await runMutation('更新当前 Environment MCP 选择', () => desktopAdapter().SetEnvironmentMCP(environmentID, id, checked));
  }
});
elements.mcpList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id; const entry = safeArray(currentSnapshot?.mcps).find((mcp) => mcp.id === id);
  if (button.dataset.action === 'edit-mcp') { beginMCPEdit(entry); return; }
  if (button.dataset.action === 'probe-mcp' && managementEnvironmentID) {
    const environmentID = managementEnvironmentID;
    const requestGeneration = connectionGeneration;
    setStatus('正在探测 ' + (entry?.name || id) + '…', 'loading');
    try {
      const health = await desktopAdapter().ProbeMCPHealth(environmentID, id);
      if (requestGeneration !== connectionGeneration || environmentID !== managementEnvironmentID) return;
      mcpHealthByKey.set(mcpHealthKey(environmentID, id), health);
      renderMCPManager(safeArray(currentSnapshot?.mcps));
      setStatus('MCP 探测完成：' + health.state, health.state === 'healthy' ? 'success' : 'error');
    } catch (error) {
      if (requestGeneration === connectionGeneration && environmentID === managementEnvironmentID) setStatus('MCP 探测失败：' + errorText(error), 'error');
    }
  }
  if (button.dataset.action === 'remove-mcp' && window.confirm(`这是全局删除，不是只从当前 Environment 禁用。已有 Environment 中的 ID 引用不会被静默改写，删除后可能显示 unresolved。继续？\n${entry?.name || id}`)) await runMutation('删除全局 MCP 定义', async () => { await desktopAdapter().RemoveMCP(id); clearMCPObservedHealth(id); if (editingMCPID === id) resetMCPEditor(true, false); });
});

elements.skillSourceList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id; const source = skillSources.find((item) => item.skill_source_id === id);
  if (button.dataset.action === 'edit-skill-source') { beginSkillSourceEdit(source); return; }
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
elements.globalMemoryForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const key = elements.globalMemoryKey.value.trim(), value = elements.globalMemoryValue.value; if (!key) return;
  const scope = currentGlobalMemoryScope();
  await runMutation('写入 Global Memory', () => desktopAdapter().WriteGlobalMemory(key, value), async () => {
    elements.globalMemoryForm.reset(); closeFormDialog(elements.globalMemoryForm);
    if (globalMemoryLoaded && globalMemoryScopeIsCurrent(scope)) await loadGlobalMemory({scope});
  });
});
elements.environmentMemoryForm.addEventListener('submit', async (event) => {
  event.preventDefault(); const scope = currentEnvironmentMemoryScope(); if (!scope.environmentID) return setStatus('请先选择 Management Environment。', 'error');
  const key = elements.environmentMemoryKey.value.trim(), value = elements.environmentMemoryValue.value; if (!key) return;
  await runMutation('写入 Environment-private Memory', () => desktopAdapter().WriteEnvironmentMemory(scope.environmentID, key, value), async () => {
    elements.environmentMemoryForm.reset(); closeFormDialog(elements.environmentMemoryForm);
    if (environmentMemoryLoaded && environmentMemoryScopeIsCurrent(scope)) await loadEnvironmentMemory({scope});
  });
});

elements.workspaceList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return;
  const id = button.dataset.id, workspace = safeArray(currentSnapshot?.workspaces).find((item) => item.workspace_id === id);
  if (button.dataset.action === 'filter-environments-by-workspace') {
    if (!workspace) return;
    elements.environmentWorkspaceFilter.value = id;
    renderEnvironments(safeArray(currentSnapshot?.environments));
    managementNavigation?.navigate('environments', {focus: true});
    return;
  }
  if (button.dataset.action === 'rename-workspace') { openRenameDialog('Workspace', workspace?.name || '', (name) => desktopAdapter().RenameWorkspace(id, name)); return; }
  if (button.dataset.action === 'remove-workspace' && window.confirm(`只移除 ADM Workspace 记录，不删除目录。继续？\n${workspace?.path || id}`)) await runMutation('移除 Workspace', () => desktopAdapter().RemoveWorkspace(id));
});
elements.environmentList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action]'); if (!button) return; const id = button.dataset.id, environment = safeArray(currentSnapshot?.environments).find((item) => item.environment_id === id);
  if (button.dataset.action === 'inspect-environment' || button.dataset.action === 'diagnose-environment') {
    environmentDetailSubview = button.dataset.action === 'diagnose-environment' ? 'diagnostics' : 'summary';
    environmentDetailOpener = button; const token = ++detailGeneration; selectedEnvironmentID = id;
    if (managementEnvironmentID !== id) environmentGeneration++;
    managementEnvironmentID = id; elements.managementEnvironment.value = id; managementContextError = ''; managementSkillAvailabilityError = ''; updateSkillBulkControls(); updateEnvironmentContextMarkers(); setStatus('读取 Environment 详情…', 'loading');
    try {
      const contextResult = await refreshManagementContext(captureEnvironmentScope());
      if (token !== detailGeneration || selectedEnvironmentID !== id || contextResult?.stale) return;
      if (!contextResult?.inspection) throw new Error(managementContextError || 'Environment inspection 不可用');
      renderEnvironmentDetailFromInspection(contextResult.inspection, token);
      if (token !== detailGeneration || selectedEnvironmentID !== id) return;
      setStatus(contextResult?.errors?.length ? `Environment 详情已加载；部分状态不可用：${contextResult.errors.join(' · ')}` : 'Environment 详情已加载', contextResult?.errors?.length ? 'error' : 'success');
    } catch (error) { if (token === detailGeneration && selectedEnvironmentID === id) { environmentDetailOpener = null; setStatus(`读取 Environment 详情失败：${errorText(error)}`, 'error'); } }
    return;
  }
  if (button.dataset.action === 'rename-environment') { openRenameDialog('Environment', environment?.name || '', (name) => desktopAdapter().RenameEnvironment(id, name)); }
  if (button.dataset.action === 'remove-environment' && window.confirm(`只移除 ADM Environment 记录，不删除 root 或项目文件。继续？\n${environment?.root || id}`)) await runMutation('移除 Environment', () => desktopAdapter().RemoveEnvironment(id));
});
elements.execList.addEventListener('click', async (event) => { const button = event.target.closest('button[data-action="remove-executable"]'); if (button) await runMutation('移除 executable', () => desktopAdapter().RemoveExecutable(button.dataset.id)); });
async function handleSelectionChange(event, kind) {
  const checkbox = event.target.closest('input[data-kind]'); const environmentID = selectedEnvironmentID; if (!checkbox || !environmentID) return;
  const itemID = checkbox.dataset.id, enabled = checkbox.checked;
  await runMutation(`更新 Environment ${kind.toUpperCase()} 选择`, () => kind === 'mcp' ? desktopAdapter().SetEnvironmentMCP(environmentID, itemID, enabled) : desktopAdapter().SetEnvironmentSkill(environmentID, itemID, enabled));
}
elements.environmentMCPSelections.addEventListener('change', (event) => handleSelectionChange(event, 'mcp')); elements.environmentSkillSelections.addEventListener('change', (event) => handleSelectionChange(event, 'skill'));
elements.globalMemoryList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action="delete-global-memory"]'); if (!button) return;
  const scope = currentGlobalMemoryScope();
  await runMutation('删除 Global Memory', () => desktopAdapter().DeleteGlobalMemory(button.dataset.id), async () => { if (globalMemoryScopeIsCurrent(scope)) await loadGlobalMemory({scope}); });
});
elements.environmentMemoryList.addEventListener('click', async (event) => {
  const button = event.target.closest('button[data-action="delete-environment-memory"]'); if (!button) return;
  const scope = currentEnvironmentMemoryScope(); if (!scope.environmentID) return setStatus('请先选择 Management Environment。', 'error');
  await runMutation('删除 Environment-private Memory', () => desktopAdapter().DeleteEnvironmentMemory(scope.environmentID, button.dataset.id), async () => {
    if (environmentMemoryScopeIsCurrent(scope)) await loadEnvironmentMemory({scope});
  });
});
elements.environmentDetailSubviewTabs.addEventListener('click', (event) => {
  const button = event.target.closest('button[data-environment-detail-subview]'); if (!button) return;
  environmentDetailSubview = button.dataset.environmentDetailSubview || 'summary'; syncEnvironmentDetailSubviewUI();
});
elements.environmentDetailRoutes.addEventListener('click', (event) => {
  const button = event.target.closest('button[data-detail-route]'); if (!button) return;
  const route = button.dataset.detailRoute; closeEnvironmentDetail(); managementNavigation?.navigate(route, {focus: true});
});

elements.loadGlobalMemory.addEventListener('click', loadGlobalMemory); elements.loadEnvironmentMemory.addEventListener('click', loadEnvironmentMemory); elements.closeEnvironmentDetail.addEventListener('click', closeEnvironmentDetail); elements.environmentDetailBackdrop.addEventListener('click', closeEnvironmentDetail);
elements.launchAtLogin.addEventListener('change', updateLaunchAtLogin);
window.addEventListener('focus', () => { loadDesktopPreferences(false); });
window.addEventListener('keydown', (event) => { if (event.key !== 'Escape' || activeEditorDialog()) return; if (selectedEnvironmentID) closeEnvironmentDetail(); else if (editingMCPID) resetMCPEditor(true); });
elements.gatewayRefreshButton.addEventListener('click', () => refreshConnectedADM(true).catch((error) => { clearManagementData(); setStatus(`ADM 连接检查失败：${error?.message || String(error)}`, 'error'); }));
elements.gatewayStartButton.addEventListener('click', () => runGatewayAction('启动本地 ADM', (input) => desktopAdapter().StartLocalADM(input)));
elements.gatewayStopButton.addEventListener('click', () => runGatewayAction('停止本地 ADM', (input) => desktopAdapter().StopLocalADM(input)));
elements.refreshButton.addEventListener('click', () => refreshConnectedADM(false).catch((error) => { clearManagementData(); setStatus(`ADM 连接检查失败：${error?.message || String(error)}`, 'error'); }));
window.addEventListener('DOMContentLoaded', () => { window.runtime?.EventsOn?.('desktop:preferences-changed', () => loadDesktopPreferences(false)); initializeManagementNavigation(); initializeConnectionProfiles(); syncMCPTransportForm(); clearManagementData(); loadDesktopPreferences(false); });
