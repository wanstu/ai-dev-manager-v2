const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const {pathToFileURL} = require('node:url');
const {spawnSync} = require('node:child_process');

const root = path.resolve(__dirname, '..', '..');
const frontend = path.join(root, 'cmd', 'ai-dev-manager-desktop', 'frontend');
const productionHTML = fs.readFileSync(path.join(frontend, 'index.html'), 'utf8');
const browserCandidates = [
  path.join(process.env['ProgramFiles(x86)'] || '', 'Microsoft', 'Edge', 'Application', 'msedge.exe'),
  path.join(process.env.ProgramFiles || '', 'Microsoft', 'Edge', 'Application', 'msedge.exe'),
  path.join(process.env.ProgramFiles || '', 'Google', 'Chrome', 'Application', 'chrome.exe'),
].filter(Boolean);
const browser = browserCandidates.find((candidate) => fs.existsSync(candidate));
if (!browser) throw new Error('No installed Chromium browser found for Desktop UI smoke');

const fakeBridge = String.raw`<script>
(() => {
  const calls = [];
  const confirmations = [];
  const browserErrors = [];
  window.addEventListener('error', (event) => browserErrors.push(event.error?.stack || event.message || String(event.error || 'error')));
  window.addEventListener('unhandledrejection', (event) => browserErrors.push(event.reason?.stack || event.reason?.message || String(event.reason || 'rejection')));
  let activeID = 'profile-a';
  const profiles = [
    {id: 'profile-a', name: 'Profile A', base_url: 'http://127.0.0.1:43137', start_service_on_desktop_launch: true},
    {id: 'profile-b', name: 'Profile B', base_url: 'http://127.0.0.1:43138'},
  ];
  const longRoot = 'C:\\fixtures\\' + 'very-long-segment-'.repeat(18) + 'workspace-a';
  const snapshotA = {
    workspaces: [
      {workspace_id:'ws-a', name:'Workspace A', path:longRoot},
      {workspace_id:'ws-b', name:'Workspace B', path:'C:\\fixtures\\plain-non-git-workspace-b'},
    ],
    environments: [
      {environment_id:'env-a', workspace_id:'ws-a', name:'Environment A', root:longRoot, state:'ready', writer:{owner:'writer-a'}, enabled_mcp_ids:['mcp-a'], enabled_skill_ids:['skill-a','skill-broken'], private_memory_count:1},
      {environment_id:'env-b', workspace_id:'ws-b', name:'Environment B', root:'C:\\fixtures\\plain-non-git-workspace-b', state:'ready', writer:null, enabled_mcp_ids:[], enabled_skill_ids:[], private_memory_count:0},
    ],
    allowed_executables: ['go'],
    mcps: [{id:'mcp-a', name:'MCP A', transport:'streamable-http', endpoint:'http://127.0.0.1:9900/mcp', default_include_in_environment:false, health_policy:{}}],
    skills: [
      {id:'skill-a', name:'Skill A', source_id:'source-a', source_root:'C:\\fixtures\\skills', artifact_path:'C:\\fixtures\\skills\\skill-a\\SKILL.md', relative_artifact_path:'skill-a/SKILL.md', support_roots:[], default_include_in_environment:false},
      {id:'skill-idle', name:'Idle Skill', source_id:'source-a', source_root:'C:\\fixtures\\skills', artifact_path:'C:\\fixtures\\skills\\idle\\SKILL.md', relative_artifact_path:'idle/SKILL.md', support_roots:[], default_include_in_environment:false},
      {id:'skill-broken', name:'Broken Skill', source_id:'source-a', source_root:'C:\\fixtures\\skills', artifact_path:'C:\\fixtures\\skills\\broken\\SKILL.md', relative_artifact_path:'broken/SKILL.md', support_roots:[], default_include_in_environment:false},
      {id:'skill-legacy', name:'Legacy Skill', source_id:'', source_root:'', artifact_path:'', relative_artifact_path:'', support_roots:[], default_include_in_environment:false},
      {id:'skill-legacy-2', name:'Legacy Skill 2', source_id:'', source_root:'', artifact_path:'', relative_artifact_path:'', support_roots:[], default_include_in_environment:false},
    ],
    global_memory_count: 1,
  };
  const snapshotB = {
    workspaces: [{workspace_id:'ws-profile-b', name:'Workspace Profile B', path:'C:\\fixtures\\profile-b'}],
    environments: [{environment_id:'env-profile-b', workspace_id:'ws-profile-b', name:'Environment Profile B', root:'C:\\fixtures\\profile-b', state:'ready', writer:null, enabled_mcp_ids:[], enabled_skill_ids:[], private_memory_count:0}],
    allowed_executables: [], mcps: [], skills: [], global_memory_count: 0,
  };
  const state = {failSkillSources:false, failVerifiers:false, failProcesses:false, hideProcess:false, delayProcessLogs:false, delayEnvironmentAInspection:false};
  const record = (name, args) => calls.push({name, args});
  const adapter = {
    async GetConnectionProfiles(){ record('GetConnectionProfiles',[]); return {profiles, active_id:activeID}; },
    async SelectConnectionProfile(id){ record('SelectConnectionProfile',[id]); activeID=id; return {profiles, active_id:activeID}; },
    async SaveConnectionProfile(profile){ record('SaveConnectionProfile',[profile]); return {profiles, active_id:activeID}; },
    async DeleteConnectionProfile(id){ record('DeleteConnectionProfile',[id]); return {profiles, active_id:activeID}; },
    async DisconnectADM(){ record('DisconnectADM',[]); return true; },
    async ConnectADM(input){ record('ConnectADM',[input]); const base=profiles.find(p=>p.id===activeID)?.base_url || input?.base_url || ''; return {state:'running', base_url:base, health_url:base+'/healthz', agent_mcp_url:base+'/mcp', admin_mcp_url:base+'/admin/mcp', pid:1234, version:'browser-fixture', local_bootstrap_eligible:false}; },
    async StartLocalADM(input){ record('StartLocalADM',[input]); const base=profiles.find(p=>p.id===activeID)?.base_url || input?.base_url || ''; return {state:'running', base_url:base, health_url:base+'/healthz', agent_mcp_url:base+'/mcp', admin_mcp_url:base+'/admin/mcp', pid:4321, version:'browser-fixture', local_bootstrap_eligible:true}; },
    async GetSnapshot(){ record('GetSnapshot',[]); return activeID==='profile-b' ? structuredClone(snapshotB) : structuredClone(snapshotA); },
    async AddWorkspace(input){ record('AddWorkspace',[input]); throw new Error('fixture workspace save failure'); },
    async ListSkillSources(){ record('ListSkillSources',[]); if(state.failSkillSources) throw new Error('skill sources unavailable'); return activeID==='profile-b' ? [] : [{skill_source_id:'source-a', root:'C:\\fixtures\\skills', support_roots:[], last_refresh_status:'ok'}]; },
    async InspectEnvironment(id){ record('InspectEnvironment',[id]); if(state.delayEnvironmentAInspection && id==='env-a') await new Promise(resolve=>setTimeout(resolve,180)); const snapshot=activeID==='profile-b'?snapshotB:snapshotA; const env=snapshot.environments.find(e=>e.environment_id===id); const workspace=snapshot.workspaces.find(w=>w.workspace_id===env?.workspace_id); const facts=id==='env-a'?[{key:'skill/skill-broken', kind:'skill', state:'unavailable', reason_code:'artifact_missing', message:'fixture artifact missing'}]:[]; return {environment:structuredClone(env), workspace:structuredClone(workspace), capability_report:{generated_at:'2026-09-11T15:00:00Z', facts}, unresolved_mcp_ids:[], unresolved_skill_ids:[]}; },
    async ListEnvironmentSkills(id){ record('ListEnvironmentSkills',[id]); const snapshot=activeID==='profile-b'?snapshotB:snapshotA; const env=snapshot.environments.find(e=>e.environment_id===id); const enabled=new Set(env?.enabled_skill_ids || []); return {skills:snapshot.skills.map(skill => ({skill_id:skill.id, source_id:skill.source_id || '', enabled:enabled.has(skill.id), state:!enabled.has(skill.id) ? 'disabled' : skill.id==='skill-broken' ? 'artifact_missing' : skill.id==='skill-legacy' ? 'unconfigured' : 'available', reason:skill.id==='skill-broken' ? 'fixture artifact missing' : skill.id==='skill-legacy' ? 'fixture legacy entry is unconfigured' : !enabled.has(skill.id) ? 'skill is not enabled for this Environment' : ''}))}; },
    async ListSkillAvailability(){ record('ListSkillAvailability',[]); const snapshot=activeID==='profile-b'?snapshotB:snapshotA; return {scope:'catalog', skills:snapshot.skills.map(skill => ({skill_id:skill.id, source_id:skill.source_id || '', enabled:true, state:skill.id==='skill-broken' ? 'artifact_missing' : skill.id.startsWith('skill-legacy') ? 'unconfigured' : 'available', reason:skill.id==='skill-broken' ? 'fixture artifact missing' : skill.id.startsWith('skill-legacy') ? 'fixture legacy entry is unconfigured' : ''}))}; },
    async RemoveSkill(id){ record('RemoveSkill',[id]); const index=snapshotA.skills.findIndex(skill=>skill.id===id); if(index<0) throw new Error('skill not found: '+id); snapshotA.skills.splice(index,1); return null; },
    async ListVerifiers(id){ record('ListVerifiers',[id]); if(state.failVerifiers) throw new Error('verifiers unavailable'); return [{verifier_id:id==='env-b'?'verifier-b':'verifier-a', name:id==='env-b'?'Verifier B':'Verifier A', kind:'test', executable:'go', args:['test','./...'], enabled:true}]; },
    async ListProcesses(id){ record('ListProcesses',[id]); if(state.failProcesses) throw new Error('processes unavailable'); if(state.hideProcess && id==='env-a') return []; return [{id:id==='env-b'?'proc-b':'proc-a', state:'running', pid:id==='env-b'?3333:2222, listening_ports:[8080]}]; },
    async GetProcessLogs(id,processID){ record('GetProcessLogs',[id,processID]); if(state.delayProcessLogs && id==='env-a') await new Promise(resolve=>setTimeout(resolve,180)); return {stdout:id==='env-b'?'PROCESS_B_LOG':'PROCESS_A_LOG', stderr:'', stdout_truncated:id==='env-a', stderr_truncated:false}; },
    async ListRuns(id){ record('ListRuns',[id]); return [{id:id==='env-b'?'run-b':'run-a', state:'succeeded', executable:'go', args:['test'], stdout:id==='env-b'?'RUN_B_OUTPUT':'RUN_A_OUTPUT', stderr:'', exit_code:0, stdout_truncated:false, stderr_truncated:false}]; },
    async ListGlobalMemory(){ record('ListGlobalMemory',[]); return [{key:'sentinel', value:'visible-after-explicit-load'}]; },
    async GetDesktopPreferences(){ record('GetDesktopPreferences',[]); return {launch_at_login_supported:false, launch_at_login:false}; },
    async SetLaunchAtLogin(value){ record('SetLaunchAtLogin',[value]); return {launch_at_login_supported:false, launch_at_login:false}; },
  };
  window.go = {desktop:{Adapter:adapter}};
  window.runtime = {EventsOn(){}, EventsEmit(){}};
  window.confirm = (message) => { confirmations.push(String(message || '')); return true; };
  window.__fakeADM = {calls, confirmations, state, browserErrors, snapshotA};
})();
</script>`;

const runner = String.raw`<script>
window.addEventListener('DOMContentLoaded', async () => {
  const result = {checks:[], failures:[], width:window.innerWidth, height:window.innerHeight, scale:window.devicePixelRatio};
  const sleep = (ms=30) => new Promise(resolve => setTimeout(resolve, ms));
  const waitFor = async (fn, label, timeout=4000) => {
    const start=Date.now();
    while(Date.now()-start<timeout){ if(fn()) return; await sleep(20); }
    throw new Error('timeout: '+label);
  };
  const check = (condition, label) => {
    result.checks.push(label);
    if(!condition) throw new Error(label);
  };
  const visibleRoute = () => document.querySelector('[data-management-page]:not([hidden])')?.dataset.managementPage || '';
  const routeLink = (route) => document.querySelector('.nav-link[data-route-link="'+route+'"]');
  const clickRoute = async (route) => { routeLink(route).click(); await sleep(); check(visibleRoute()===route, 'route '+route+' visible'); };
  try {
    await waitFor(() => document.getElementById('dashboardDataState')?.dataset.state==='success', 'initial dashboard success');
    check(location.hash==='#/overview', 'default hash canonicalizes to overview');
    check(document.getElementById('workspaceCount').textContent==='2', 'successful snapshot renders count');
    check(window.__fakeADM.calls.some(c=>c.name==='StartLocalADM'), 'startup profile can start local ADM Service with Desktop');
    check(document.querySelectorAll('[data-management-page]:not([hidden])').length===1, 'exactly one page initially visible');

    const beforeRoutes = window.__fakeADM.calls.length;
    for(const route of ['workspaces','environments','runtime','mcp','skills','memory','gateway','exec-allowlist','settings','overview']) await clickRoute(route);
    check(window.__fakeADM.calls.length===beforeRoutes, 'routing alone makes no adapter calls');
    check(document.querySelectorAll('.nav-link[aria-current="page"]').length===1, 'one active menu item');

    await clickRoute('memory');
    check(!window.__fakeADM.calls.some(c=>c.name==='ListGlobalMemory'), 'Memory route does not read values');
    document.getElementById('loadGlobalMemory').click();
    await waitFor(() => window.__fakeADM.calls.some(c=>c.name==='ListGlobalMemory'), 'explicit Global Memory read');
    check(document.getElementById('globalMemoryList').textContent.includes('visible-after-explicit-load'), 'explicit Memory load renders value');

    await clickRoute('workspaces');
    check(document.getElementById('workspaceVisibleCount').textContent==='2' && document.getElementById('workspaceListTotalCount').textContent==='2', 'Workspace visible/total counts use loaded snapshot');
    check(document.getElementById('workspaceList').textContent.includes('Environments 1'), 'Workspace rows show joined Environment counts');
    const projectFilterCalls=window.__fakeADM.calls.length;
    const workspaceFilter=document.getElementById('workspaceFilter'); workspaceFilter.value='Workspace B'; workspaceFilter.dispatchEvent(new Event('input',{bubbles:true})); await sleep();
    check(document.getElementById('workspaceVisibleCount').textContent==='1' && document.getElementById('workspaceList').textContent.includes('Workspace B') && !document.getElementById('workspaceList').textContent.includes('Workspace A'), 'Workspace search filters loaded rows locally');
    check(window.__fakeADM.calls.length===projectFilterCalls, 'Workspace filtering makes no adapter calls');
    workspaceFilter.value=''; workspaceFilter.dispatchEvent(new Event('input',{bubbles:true})); await sleep();
    const managementBeforeWorkspaceFilter=document.getElementById('managementEnvironment').value;
    document.querySelector('#workspaceList button[data-action="filter-environments-by-workspace"][data-id="ws-a"]').click(); await sleep();
    check(visibleRoute()==='environments' && document.getElementById('environmentWorkspaceFilter').value==='ws-a', 'Workspace shortcut navigates to filtered Environments');
    check(document.getElementById('environmentVisibleCount').textContent==='1' && document.getElementById('environmentList').textContent.includes('Environment A') && !document.getElementById('environmentList').textContent.includes('Environment B'), 'Workspace shortcut filters Environment rows');
    check(document.getElementById('managementEnvironment').value===managementBeforeWorkspaceFilter, 'Workspace shortcut does not retarget Management Environment');
    check(window.__fakeADM.calls.length===projectFilterCalls, 'Workspace shortcut remains presentation-only');
    document.getElementById('environmentWorkspaceFilter').value=''; document.getElementById('environmentWorkspaceFilter').dispatchEvent(new Event('change',{bubbles:true})); await sleep();
    await clickRoute('workspaces');
    const workspaceOpener = document.querySelector('[data-management-page="workspaces"] [data-dialog-open="workspaceDialog"]');
    workspaceOpener.focus(); workspaceOpener.click(); await sleep();
    const workspaceDialog=document.getElementById('workspaceDialog');
    check(workspaceDialog.open, 'workspace editor opens');
    document.getElementById('workspacePath').value='C:\\draft\\keep-me';
    location.hash='#/mcp'; await sleep(60);
    check(visibleRoute()==='workspaces' && location.hash==='#/workspaces', 'open editor guards route change');
    check(document.getElementById('workspacePath').value==='C:\\draft\\keep-me', 'guarded navigation retains draft');
    document.getElementById('workspaceForm').requestSubmit();
    await waitFor(() => window.__fakeADM.calls.some(c=>c.name==='AddWorkspace'), 'failed Workspace save call');
    await waitFor(() => document.getElementById('statusPanel').dataset.kind==='error', 'failed Workspace save status');
    check(workspaceDialog.open, 'failed Workspace save keeps editor open');
    check(document.getElementById('workspacePath').value==='C:\\draft\\keep-me', 'failed Workspace save retains draft');
    workspaceDialog.querySelector('[data-dialog-close]').click(); await sleep();
    check(document.activeElement===workspaceOpener, 'closing editor restores visible opener focus');

    await clickRoute('environments');
    const detailButton=document.querySelector('#environmentList button[data-action="inspect-environment"][data-id="env-a"]');
    const detailInspectBefore=window.__fakeADM.calls.filter(c=>c.name==='InspectEnvironment' && c.args[0]==='env-a').length;
    const detailMemoryBefore=window.__fakeADM.calls.filter(c=>c.name==='ListEnvironmentMemory').length;
    detailButton.focus(); detailButton.click();
    await waitFor(() => !document.getElementById('environmentDetailPanel').hidden, 'Environment detail open');
    check(window.__fakeADM.calls.filter(c=>c.name==='InspectEnvironment' && c.args[0]==='env-a').length===detailInspectBefore+1, 'Environment detail reuses one scoped inspection instead of duplicate reads');
    check(document.getElementById('environmentDetail').textContent.includes('Identity') && document.getElementById('environmentDetail').textContent.includes('Runtime authority') && document.getElementById('environmentDetail').textContent.includes('Capability issues'), 'Environment detail groups identity authority and capability facts');
    check(document.getElementById('environmentDetail').textContent.includes('artifact_missing'), 'Environment detail exposes capability reason without probing');
    check(window.__fakeADM.calls.filter(c=>c.name==='ListEnvironmentMemory').length===detailMemoryBefore, 'opening Environment detail does not read private Memory values');
    location.hash='#/skills'; await sleep(60);
    check(visibleRoute()==='environments' && location.hash==='#/environments', 'Environment detail guards route change');
    document.getElementById('closeEnvironmentDetail').click(); await sleep();
    check(document.activeElement===detailButton, 'closing Environment detail restores opener focus');
    const currentEnvironmentRow=document.querySelector('#environmentList .managed-item.current-context');
    check(currentEnvironmentRow?.textContent.includes('Environment A') && currentEnvironmentRow?.textContent.includes('当前管理环境'), 'Environment list marks explicit current Management Environment');
    check(currentEnvironmentRow?.textContent.includes('Workspace A') && currentEnvironmentRow?.textContent.includes('MCP 1') && currentEnvironmentRow?.textContent.includes('Skills 2'), 'Environment rows join Workspace identity and selection counts');
    detailButton.click(); await waitFor(() => !document.getElementById('environmentDetailPanel').hidden, 'Environment detail reopen for shortcut');
    document.querySelector('#environmentDetailRoutes button[data-detail-route="skills"]').click(); await sleep();
    check(document.getElementById('environmentDetailPanel').hidden && visibleRoute()==='skills', 'Environment detail shortcut closes modal then navigates');
    check(document.getElementById('managementEnvironment').value==='env-a', 'Environment detail shortcut preserves explicit Management Environment');
    check(window.__fakeADM.calls.filter(c=>c.name==='ListEnvironmentMemory').length===detailMemoryBefore, 'Environment detail shortcut does not implicitly load Memory');

    await clickRoute('skills');
    const noEnvSelect=document.getElementById('managementEnvironment'); noEnvSelect.value=''; noEnvSelect.dispatchEvent(new Event('change',{bubbles:true}));
    await waitFor(() => !document.getElementById('skillProbeAllButton').disabled, 'global Skill probe does not require Environment');
    check(document.getElementById('skillClearUnavailableButton').disabled, 'automatic availability load does not authorize one-click cleanup');
    check(document.getElementById('skillBulkHint').textContent.includes('全局 Skill catalog') && document.getElementById('skillBulkHint').textContent.includes('不依赖当前 Environment'), 'Skill bulk hint explains global catalog scope');
    const probeCallsBefore=window.__fakeADM.calls.filter(c=>c.name==='ListSkillAvailability').length;
    const envSkillCallsBeforeProbe=window.__fakeADM.calls.filter(c=>c.name==='ListEnvironmentSkills').length;
    document.getElementById('skillProbeAllButton').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='ListSkillAvailability').length>probeCallsBefore, 'explicit global Skill availability call');
    check(window.__fakeADM.calls.filter(c=>c.name==='ListEnvironmentSkills').length===envSkillCallsBeforeProbe, 'global Skill bulk check does not call Environment availability');
    await waitFor(() => !document.getElementById('skillClearUnavailableButton').disabled, 'fresh explicit global probe enables cleanup');
    check(document.getElementById('skillClearUnavailableButton').textContent.includes('(3)'), 'cleanup includes all global structural failures');
    check(document.getElementById('skillBulkHint').textContent.includes('可用 2') && document.getElementById('skillBulkHint').textContent.includes('不可用 3'), 'global probe counts healthy idle Skill as available and broken placeholders as unavailable');

    const selectLegacy=document.querySelector('input[data-action="select-skill"][data-id="skill-legacy"]');
    const selectLegacy2=document.querySelector('input[data-action="select-skill"][data-id="skill-legacy-2"]');
    selectLegacy.click(); selectLegacy2.click(); await sleep();
    check(document.getElementById('skillSelectedCount').textContent==='2', 'two Skills can be selected for bulk delete');
    document.getElementById('skillDeleteSelectedButton').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='RemoveSkill' && (c.args[0]==='skill-legacy' || c.args[0]==='skill-legacy-2')).length===2, 'bulk delete removes each selected stable Skill ID');
    await waitFor(() => document.getElementById('skillListTotalCount').textContent==='3', 'bulk delete refreshes authoritative Skill snapshot');
    check(window.__fakeADM.confirmations.some(message=>message.includes('Legacy: 2') && message.includes('catalog metadata')), 'bulk delete confirmation explains metadata-only impact');
    check(document.getElementById('skillClearUnavailableButton').disabled, 'catalog refresh invalidates prior cleanup authorization');

    const secondProbeCalls=window.__fakeADM.calls.filter(c=>c.name==='ListSkillAvailability').length;
    document.getElementById('skillProbeAllButton').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='ListSkillAvailability').length>secondProbeCalls, 'second explicit global Skill availability call');
    await waitFor(() => !document.getElementById('skillClearUnavailableButton').disabled, 'fresh probe re-enables cleanup after catalog change');
    document.getElementById('skillClearUnavailableButton').click();
    await waitFor(() => window.__fakeADM.calls.some(c=>c.name==='RemoveSkill' && c.args[0]==='skill-broken'), 'one-click cleanup removes explicit-probe unavailable Skill');
    await waitFor(() => document.getElementById('skillListTotalCount').textContent==='2', 'one-click cleanup refreshes Skill catalog');
    check(document.getElementById('skillList').textContent.includes('Skill A') && document.getElementById('skillList').textContent.includes('Idle Skill'), 'one-click cleanup preserves available Skills');
    check(window.__fakeADM.confirmations.some(message=>message.includes('Source-managed: 1') && message.includes('刷新 Source')), 'cleanup confirmation warns source-managed Skill may return');

    window.__fakeADM.state.failSkillSources=true;
    document.getElementById('refreshButton').click();
    await waitFor(() => document.getElementById('skillSourceList').textContent.includes('不可用'), 'Skill source local failure');
    check(document.getElementById('workspaceList').textContent.includes('Workspace A'), 'auxiliary failure preserves base snapshot');
    check(document.getElementById('dashboardDataState').dataset.state==='success', 'auxiliary failure does not demote base snapshot');
    window.__fakeADM.state.failSkillSources=false;

    window.__fakeADM.state.failVerifiers=true;
    const envSelect=document.getElementById('managementEnvironment'); envSelect.value='env-a'; envSelect.dispatchEvent(new Event('change',{bubbles:true}));
    await waitFor(() => document.getElementById('verifierList').textContent.includes('不可用'), 'Verifier local failure');
    check(document.getElementById('processList').textContent.includes('proc-a'), 'Process list survives verifier failure');
    check(document.getElementById('runList').textContent.includes('run-a'), 'Run list survives verifier failure');
    window.__fakeADM.state.failVerifiers=false;

    await clickRoute('runtime');
    const runtimeCallsBeforeTabs=window.__fakeADM.calls.length;
    check(!document.querySelector('[data-runtime-subview-panel="verifiers"]').hidden && document.querySelector('[data-runtime-subview-panel="processes"]').hidden, 'Runtime starts on local Verifiers subview');
    document.querySelector('#runtimeSubviewTabs [data-runtime-subview="processes"]').click(); await sleep();
    check(document.querySelector('[data-runtime-subview-panel="verifiers"]').hidden && !document.querySelector('[data-runtime-subview-panel="processes"]').hidden, 'Runtime switches to Processes subview locally');
    check(window.__fakeADM.calls.length===runtimeCallsBeforeTabs, 'Runtime subview switch makes no adapter calls');
    const processLogCalls=window.__fakeADM.calls.filter(c=>c.name==='GetProcessLogs').length;
    document.querySelector('#processList button[data-action="process-logs"][data-id="proc-a"]').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='GetProcessLogs').length>processLogCalls, 'explicit Process logs call');
    await waitFor(() => document.getElementById('runtimeOutput').textContent.includes('PROCESS_A_LOG'), 'Process logs render bounded output');
    check(document.getElementById('runtimeOutputMeta').textContent.includes('current-owner observation') && document.getElementById('runtimeOutputMeta').textContent.includes('processes proc-a') && document.getElementById('runtimeOutputMeta').textContent.includes('stdout truncated'), 'Process output meta identifies resource and truncation');
    const runtimeCallsBeforeRunTab=window.__fakeADM.calls.length;
    document.querySelector('#runtimeSubviewTabs [data-runtime-subview="runs"]').click(); await sleep();
    check(window.__fakeADM.calls.length===runtimeCallsBeforeRunTab && document.getElementById('runtimeOutput').textContent==='尚无输出', 'Runtime tab change clears output without backend reads');
    document.querySelector('#runList button[data-action="run-output"][data-id="run-a"]').click(); await sleep();
    check(document.getElementById('runtimeOutput').textContent.includes('RUN_A_OUTPUT') && document.getElementById('runtimeOutputMeta').textContent.includes('runs run-a'), 'Run output uses existing owner-local list result');
    document.querySelector('#runtimeSubviewTabs [data-runtime-subview="processes"]').click(); await sleep();
    document.querySelector('#processList button[data-action="process-logs"][data-id="proc-a"]').click(); await waitFor(() => document.getElementById('runtimeOutput').textContent.includes('PROCESS_A_LOG'), 'Process output rebound before disappearance');
    window.__fakeADM.state.hideProcess=true; document.getElementById('runtimeRefreshButton').click();
    await waitFor(() => document.getElementById('processList').textContent.includes('没有 process'), 'authoritative Process disappearance refresh');
    check(document.getElementById('runtimeOutput').textContent==='尚无输出' && document.getElementById('runtimeOutputMeta').textContent.includes('资源已不在最新'), 'successful list refresh clears output for disappeared resource');
    window.__fakeADM.state.hideProcess=false; document.getElementById('runtimeRefreshButton').click(); await waitFor(() => document.getElementById('processList').textContent.includes('proc-a'), 'Process restored for stale-log test');
    window.__fakeADM.state.delayProcessLogs=true;
    const delayedLogCalls=window.__fakeADM.calls.filter(c=>c.name==='GetProcessLogs' && c.args[0]==='env-a').length;
    document.querySelector('#processList button[data-action="process-logs"][data-id="proc-a"]').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='GetProcessLogs' && c.args[0]==='env-a').length>delayedLogCalls, 'delayed Environment A process log started');
    envSelect.value='env-b'; envSelect.dispatchEvent(new Event('change',{bubbles:true}));
    await waitFor(() => document.getElementById('processList').textContent.includes('proc-b'), 'Environment B Runtime list wins process-log race');
    await sleep(240);
    check(!document.getElementById('runtimeOutput').textContent.includes('PROCESS_A_LOG') && document.getElementById('runtimeOutputMeta').textContent.includes('Environment 已切换'), 'late Environment A process log cannot populate Environment B output');
    window.__fakeADM.state.delayProcessLogs=false;

    const forbiddenAutoCalls=['ProbeMCPHealth','RefreshSkillSource','RunVerifier','StopProcess','CancelRun','WriteGlobalMemory','WriteEnvironmentMemory'];
    check(!window.__fakeADM.calls.some(c=>forbiddenAutoCalls.includes(c.name)), 'navigation/refresh makes no implicit probe/mutation call');

    await clickRoute('environments');
    window.__fakeADM.state.delayEnvironmentAInspection=true;
    const delayedAButton=document.querySelector('#environmentList button[data-action="inspect-environment"][data-id="env-a"]');
    const fastBButton=document.querySelector('#environmentList button[data-action="inspect-environment"][data-id="env-b"]');
    const delayedACalls=window.__fakeADM.calls.filter(c=>c.name==='InspectEnvironment' && c.args[0]==='env-a').length;
    delayedAButton.click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='InspectEnvironment' && c.args[0]==='env-a').length>delayedACalls, 'delayed Environment A inspection started');
    fastBButton.click();
    await waitFor(() => !document.getElementById('environmentDetailPanel').hidden && document.getElementById('environmentDetailTitle').textContent==='Environment B', 'Environment B detail wins A to B race');
    await sleep(240);
    check(document.getElementById('environmentDetailTitle').textContent==='Environment B' && document.getElementById('environmentDetail').textContent.includes('env-b'), 'late Environment A response cannot overwrite Environment B detail');
    check(document.getElementById('managementEnvironment').value==='env-b', 'A to B detail race leaves explicit Management Environment on B');
    window.__fakeADM.state.delayEnvironmentAInspection=false;
    document.getElementById('closeEnvironmentDetail').click(); await sleep();

    await clickRoute('environments');
    const managementBeforeInvalidFilter=document.getElementById('managementEnvironment').value;
    const environmentWorkspaceFilter=document.getElementById('environmentWorkspaceFilter'); environmentWorkspaceFilter.value='ws-a'; environmentWorkspaceFilter.dispatchEvent(new Event('change',{bubbles:true}));
    window.__fakeADM.snapshotA.workspaces.splice(window.__fakeADM.snapshotA.workspaces.findIndex(workspace=>workspace.workspace_id==='ws-a'),1);
    document.getElementById('refreshButton').click();
    await waitFor(() => !document.getElementById('environmentFilterHint').hidden, 'removed Workspace filter invalid hint');
    check(document.getElementById('environmentWorkspaceFilter').value==='ws-a' && document.getElementById('environmentFilterHint').textContent.includes('ws-a'), 'removed Workspace filter stays visibly invalid');
    check(document.getElementById('managementEnvironment').value===managementBeforeInvalidFilter, 'invalid Workspace filter never retargets Management Environment');

    const profileSelect=document.getElementById('connectionSelect'); profileSelect.value='profile-b'; profileSelect.dispatchEvent(new Event('change',{bubbles:true}));
    await waitFor(() => document.getElementById('workspaceList').textContent.includes('Workspace Profile B'), 'profile B snapshot');
    check(!document.getElementById('workspaceList').textContent.includes('Workspace A'), 'profile switch clears old snapshot');
    check(!document.getElementById('globalMemoryList').textContent.includes('visible-after-explicit-load'), 'profile switch clears loaded Memory values');

    const pageWidth=Math.max(document.documentElement.scrollWidth, document.body.scrollWidth);
    check(pageWidth<=window.innerWidth+2, 'no whole-window horizontal overflow');
    check(document.querySelectorAll('[data-management-page][hidden] :focus').length===0, 'hidden pages do not retain focus');
    result.ok=true;
  } catch(error) {
    result.ok=false; result.failures.push(error?.stack || error?.message || String(error));
    result.diagnostics = {
      dashboardState: document.getElementById('dashboardDataState')?.dataset.state || '',
      dashboardText: document.getElementById('dashboardDataState')?.textContent || '',
      statusText: document.getElementById('statusPanel')?.textContent || '',
      statusKind: document.getElementById('statusPanel')?.dataset.kind || '',
      calls: window.__fakeADM?.calls || [],
      browserErrors: window.__fakeADM?.browserErrors || [],
      navigationType: typeof window.ADMNavigation,
      dashboardType: typeof window.ADMDashboard,
      initializeProfilesType: typeof window.initializeConnectionProfiles,
      appElementsType: typeof window.elements,
      baseURI: document.baseURI,
    };
  } finally {
    const encoded=btoa(unescape(encodeURIComponent(JSON.stringify(result))));
    document.documentElement.setAttribute('data-browser-smoke', encoded);
  }
});
</script>`;

function htmlForFixture() {
  const absoluteAssets = productionHTML.replace(/(src|href)="\.\/([^"]+)"/g, (_match, attribute, relative) => {
    return `${attribute}="${pathToFileURL(path.join(frontend, relative)).href}"`;
  });
  return absoluteAssets
    .replace('<body>', `<body>\n${fakeBridge}`)
    .replace('</body>', `${runner}\n</body>`);
}

function runBrowser(width, height, scale = 1) {
  const temp = fs.mkdtempSync(path.join(os.tmpdir(), 'adm-desktop-browser-'));
  const fixture = path.join(temp, 'index.html');
  fs.writeFileSync(fixture, htmlForFixture());
  const profile = path.join(temp, 'profile');
  const args = [
    '--headless=new', '--disable-gpu', '--no-first-run', '--no-default-browser-check', '--allow-file-access-from-files',
    `--user-data-dir=${profile}`, `--window-size=${width},${height}`, `--force-device-scale-factor=${scale}`,
    '--virtual-time-budget=20000', '--dump-dom', pathToFileURL(fixture).href,
  ];
  const execution = spawnSync(browser, args, {encoding:'utf8', timeout:70000, maxBuffer:20*1024*1024});
  try {
    if (execution.error) throw execution.error;
    if (execution.status !== 0) throw new Error(`browser exit ${execution.status}: ${execution.stderr}`);
    const match = execution.stdout.match(/data-browser-smoke="([^"]+)"/);
    if (!match) throw new Error(`browser smoke result missing; stderr=${execution.stderr.slice(-2000)}`);
    const result = JSON.parse(Buffer.from(match[1], 'base64').toString('utf8'));
    if (!result.ok) throw new Error(`browser smoke failed ${width}x${height}@${scale}: ${result.failures.join('\n')}\nDiagnostics: ${JSON.stringify(result.diagnostics || {})}`);
    return result;
  } finally {
    fs.rmSync(temp, {recursive:true, force:true});
  }
}

const scenarios = [
  [1120, 760, 1],
  [820, 560, 1],
  [1120, 760, 1.25],
];
for (const [width,height,scale] of scenarios) {
  const result = runBrowser(width,height,scale);
  console.log(`PASS browser ${width}x${height} scale=${scale} checks=${result.checks.length}`);
}
