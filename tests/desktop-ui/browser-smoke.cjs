const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const {pathToFileURL} = require('node:url');
const {spawnSync} = require('node:child_process');

const root = path.resolve(__dirname, '..', '..');
const frontend = path.join(root, 'cmd', 'ai-dev-manager-desktop', 'frontend');
const productionHTML = fs.readFileSync(path.join(frontend, 'index.html'), 'utf8');
const browserCandidates = [
  path.join(process.env.ProgramFiles || '', 'Google', 'Chrome', 'Application', 'chrome.exe'),
  path.join(process.env['ProgramFiles(x86)'] || '', 'Microsoft', 'Edge', 'Application', 'msedge.exe'),
  path.join(process.env.ProgramFiles || '', 'Microsoft', 'Edge', 'Application', 'msedge.exe'),
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
    exec_denials: [
      {executable:'blocked_python_fixture', count:5, first_blocked_at:'2026-09-12T03:00:00Z', last_blocked_at:'2026-09-12T04:00:00Z', last_environment_id:'env-a', last_surface:'run_start', last_reason:'executable "blocked_python_fixture" is not allowed'},
      {executable:'blocked_node_fixture', count:2, first_blocked_at:'2026-09-12T02:00:00Z', last_blocked_at:'2026-09-12T03:30:00Z', last_environment_id:'env-b', last_surface:'process_start', last_reason:'executable "blocked_node_fixture" is not allowed'},
    ],
    mcps: [
      {id:'mcp-a', name:'MCP A', transport:'streamable-http', endpoint:'http://127.0.0.1:9900/mcp', default_include_in_environment:false, health_policy:{}},
      {id:'mcp-b', name:'MCP B', transport:'streamable-http', endpoint:'http://127.0.0.1:9901/mcp', default_include_in_environment:false, health_policy:{}},
    ],
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
  const skillSourcesA = [{skill_source_id:'source-a', root:'C:\\fixtures\\skills', support_roots:[], last_refresh_status:'ok'}];
  const environmentMemory = {'env-a':[{key:'private-sentinel', value:'env-a-private-visible'}], 'env-b':[{key:'private-sentinel', value:'env-b-private-visible'}]};
  const state = {failSkillSources:false, failVerifiers:false, failProcesses:false, hideProcess:false, delayProcessLogs:false, delayEnvironmentAInspection:false, delayMCPProbe:false};
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
    async AllowExecutable(executable){ record('AllowExecutable',[executable]); if(!snapshotA.allowed_executables.includes(executable)) snapshotA.allowed_executables.push(executable); snapshotA.exec_denials=snapshotA.exec_denials.filter(item=>item.executable!==executable); return structuredClone(snapshotA.allowed_executables); },
    async RemoveExecutable(executable){ record('RemoveExecutable',[executable]); snapshotA.allowed_executables=snapshotA.allowed_executables.filter(item=>item!==executable); return structuredClone(snapshotA.allowed_executables); },
    async ClearExecDenial(executable){ record('ClearExecDenial',[executable]); snapshotA.exec_denials=snapshotA.exec_denials.filter(item=>item.executable!==executable); return structuredClone(snapshotA.exec_denials); },
    async ClearAllExecDenials(){ record('ClearAllExecDenials',[]); snapshotA.exec_denials=[]; return null; },
    async ListSkillSources(){ record('ListSkillSources',[]); if(state.failSkillSources) throw new Error('skill sources unavailable'); return activeID==='profile-b' ? [] : structuredClone(skillSourcesA); },
    async UpdateSkillSource(id,input){ record('UpdateSkillSource',[id,input]); const source=skillSourcesA.find(item=>item.skill_source_id===id); if(!source) throw new Error('source not found: '+id); source.root=input.root; source.support_roots=input.support_roots || []; source.default_include_in_environment=Boolean(input.default_include_in_environment); source.last_refresh_status='pending'; source.last_refresh_error='source settings changed; refresh required'; return structuredClone(source); },
    async InspectEnvironment(id){ record('InspectEnvironment',[id]); if(state.delayEnvironmentAInspection && id==='env-a') await new Promise(resolve=>setTimeout(resolve,180)); const snapshot=activeID==='profile-b'?snapshotB:snapshotA; const env=snapshot.environments.find(e=>e.environment_id===id); const workspace=snapshot.workspaces.find(w=>w.workspace_id===env?.workspace_id); const facts=id==='env-a'?[{key:'skill/skill-broken', kind:'skill', state:'unavailable', reason_code:'artifact_missing', message:'fixture artifact missing', source:'fixture capability report', observed_at:'2026-09-12T03:00:00Z'}]:[]; return {environment:structuredClone(env), workspace:structuredClone(workspace), capability_report:{generated_at:'2026-09-11T15:00:00Z', facts}, unresolved_mcp_ids:[], unresolved_skill_ids:[]}; },
    async ListEnvironmentSkills(id){ record('ListEnvironmentSkills',[id]); const snapshot=activeID==='profile-b'?snapshotB:snapshotA; const env=snapshot.environments.find(e=>e.environment_id===id); const enabled=new Set(env?.enabled_skill_ids || []); return {skills:snapshot.skills.map(skill => ({skill_id:skill.id, source_id:skill.source_id || '', enabled:enabled.has(skill.id), state:!enabled.has(skill.id) ? 'disabled' : skill.id==='skill-broken' ? 'artifact_missing' : skill.id==='skill-legacy' ? 'unconfigured' : 'available', reason:skill.id==='skill-broken' ? 'fixture artifact missing' : skill.id==='skill-legacy' ? 'fixture legacy entry is unconfigured' : !enabled.has(skill.id) ? 'skill is not enabled for this Environment' : ''}))}; },
    async ListSkillAvailability(){ record('ListSkillAvailability',[]); const snapshot=activeID==='profile-b'?snapshotB:snapshotA; return {scope:'catalog', skills:snapshot.skills.map(skill => ({skill_id:skill.id, source_id:skill.source_id || '', enabled:true, state:skill.id==='skill-broken' ? 'artifact_missing' : skill.id.startsWith('skill-legacy') ? 'unconfigured' : 'available', reason:skill.id==='skill-broken' ? 'fixture artifact missing' : skill.id.startsWith('skill-legacy') ? 'fixture legacy entry is unconfigured' : ''}))}; },
    async SetMCPDefault(id,value){ record('SetMCPDefault',[id,value]); const entry=snapshotA.mcps.find(item=>item.id===id); if(!entry) throw new Error('mcp not found: '+id); entry.default_include_in_environment=Boolean(value); return structuredClone(entry); },
    async SetEnvironmentMCP(environmentID,id,value){ record('SetEnvironmentMCP',[environmentID,id,value]); const env=snapshotA.environments.find(item=>item.environment_id===environmentID); if(!env) throw new Error('env not found: '+environmentID); const set=new Set(env.enabled_mcp_ids || []); if(value) set.add(id); else set.delete(id); env.enabled_mcp_ids=[...set].sort(); return structuredClone(env); },
    async PreviewMCPImport(input){ record('PreviewMCPImport',[input]); return {format: input?.format || 'generic-mcpservers', candidates:[{name:'Imported MCP', transport:'streamable-http', endpoint:'http://127.0.0.1:9902/mcp', reference_requirements:[{field_path:'headers.Authorization', reference_name:'ADM_IMPORTED_TOKEN'}], errors:[], warnings:[]}]}; },
    async ApplyMCPImport(input){ record('ApplyMCPImport',[input]); const definition={id:'mcp-imported', name:'Imported MCP', transport:'streamable-http', endpoint:'http://127.0.0.1:9902/mcp', default_include_in_environment:Boolean(input?.default_include), health_policy:{}}; if(!snapshotA.mcps.some(item=>item.id===definition.id)) snapshotA.mcps.push(definition); return {result:{mutations:[{action:'created', definition}]}}; },
    async ProbeMCPHealth(environmentID,id){ record('ProbeMCPHealth',[environmentID,id]); if(state.delayMCPProbe && environmentID==='env-a') await new Promise(resolve=>setTimeout(resolve,180)); return {state:'healthy', message:'fixture healthy'}; },
    async RemoveMCP(id){ record('RemoveMCP',[id]); const index=snapshotA.mcps.findIndex(item=>item.id===id); if(index>=0) snapshotA.mcps.splice(index,1); return null; },
    async SetSkillDefault(id,value){ record('SetSkillDefault',[id,value]); const entry=snapshotA.skills.find(item=>item.id===id); if(!entry) throw new Error('skill not found: '+id); entry.default_include_in_environment=Boolean(value); return structuredClone(entry); },
    async SetEnvironmentSkill(environmentID,id,value){ record('SetEnvironmentSkill',[environmentID,id,value]); const env=snapshotA.environments.find(item=>item.environment_id===environmentID); if(!env) throw new Error('env not found: '+environmentID); const set=new Set(env.enabled_skill_ids || []); if(value) set.add(id); else set.delete(id); env.enabled_skill_ids=[...set].sort(); return structuredClone(env); },
    async RemoveSkill(id){ record('RemoveSkill',[id]); const index=snapshotA.skills.findIndex(skill=>skill.id===id); if(index<0) throw new Error('skill not found: '+id); snapshotA.skills.splice(index,1); return null; },
    async ListVerifiers(id){ record('ListVerifiers',[id]); if(state.failVerifiers) throw new Error('verifiers unavailable'); return [{verifier_id:id==='env-b'?'verifier-b':'verifier-a', name:id==='env-b'?'Verifier B':'Verifier A', kind:'test', executable:'go', args:['test','./...'], enabled:true}]; },
    async ListProcesses(id){ record('ListProcesses',[id]); if(state.failProcesses) throw new Error('processes unavailable'); if(state.hideProcess && id==='env-a') return []; return [{id:id==='env-b'?'proc-b':'proc-a', state:'running', pid:id==='env-b'?3333:2222, listening_ports:[8080]}]; },
    async GetProcessLogs(id,processID){ record('GetProcessLogs',[id,processID]); if(state.delayProcessLogs && id==='env-a') await new Promise(resolve=>setTimeout(resolve,180)); return {stdout:id==='env-b'?'PROCESS_B_LOG':'PROCESS_A_LOG', stderr:'', stdout_truncated:id==='env-a', stderr_truncated:false}; },
    async ListRuns(id){ record('ListRuns',[id]); return [{id:id==='env-b'?'run-b':'run-a', state:'succeeded', executable:'go', args:['test'], stdout:id==='env-b'?'RUN_B_OUTPUT':'RUN_A_OUTPUT', stderr:'', exit_code:0, stdout_truncated:false, stderr_truncated:false}]; },
    async ListGlobalMemory(){ record('ListGlobalMemory',[]); return [{key:'sentinel', value:'visible-after-explicit-load'}]; },
    async ListEnvironmentMemory(id){ record('ListEnvironmentMemory',[id]); return structuredClone(environmentMemory[id] || []); },
    async WriteEnvironmentMemory(id,key,value){ record('WriteEnvironmentMemory',[id,key,value]); environmentMemory[id] = (environmentMemory[id] || []).filter(entry=>entry.key!==key); environmentMemory[id].push({key,value}); return null; },
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
    for(const route of ['workspaces','environments','runtime','mcp','skills','gateway','diagnostics','exec-allowlist','settings','overview']) await clickRoute(route);
    check(document.getElementById('diagnosticsPageContent').textContent.includes('不会自动 fan-out'), 'Diagnostics route without Environment does not fan out');
    await clickRoute('gateway');
    check(document.getElementById('managementContextPanel').hidden, 'ADM connection route hides Management Context');
    check(document.querySelectorAll('#gatewayState').length===1, 'Gateway status id is unique in the document');
    check(document.body.textContent.includes('Saved connection profile') && document.body.textContent.includes('Runtime endpoints') && document.body.textContent.includes('Local service lifecycle'), 'Gateway page groups profile endpoints and local lifecycle controls');
    check(document.getElementById('gatewayStartButton') && document.getElementById('gatewayStopButton') && document.getElementById('gatewayRefreshButton'), 'Gateway grouped page preserves lifecycle buttons');    await clickRoute('exec-allowlist');
    check(document.getElementById('managementContextPanel').hidden, 'Exec allowlist route hides Management Context');
    check(document.body.textContent.includes('Allowed executables') && document.getElementById('execList') && document.querySelector('[data-dialog-open="execDialog"]'), 'Exec allowlist grouping preserves explicit allow/remove surface');
    check(document.getElementById('execBlockedCount').textContent==='7' && document.getElementById('execBlockedList').textContent.includes('blocked_python_fixture'), 'Exec allowlist renders denied executable observations sorted by count');
    check(document.getElementById('execBlockedList').textContent.indexOf('blocked_python_fixture') < document.getElementById('execBlockedList').textContent.indexOf('blocked_node_fixture'), 'Exec denied executable observations are sorted by count');
    document.querySelector('#execBlockedList button[data-action="allow-blocked-executable"][data-id="blocked_python_fixture"]').click();
    await waitFor(() => window.__fakeADM.calls.some(c=>c.name==='AllowExecutable' && c.args[0]==='blocked_python_fixture'), 'Allow blocked executable action');
    await waitFor(() => document.getElementById('execList').textContent.includes('blocked_python_fixture') && !document.getElementById('execBlockedList').textContent.includes('blocked_python_fixture'), 'Allowing blocked executable adds allowlist and clears observation');
    document.querySelector('#execBlockedList button[data-action="clear-blocked-executable"][data-id="blocked_node_fixture"]').click();
    await waitFor(() => window.__fakeADM.calls.some(c=>c.name==='ClearExecDenial' && c.args[0]==='blocked_node_fixture'), 'Clear one blocked executable observation');
    await waitFor(() => document.getElementById('execBlockedList').textContent.includes('暂无被拦截 executable'), 'Clearing blocked executable removes final observation');
    const afterExecAuthorityActions = window.__fakeADM.calls.length;
    await clickRoute('settings');
    check(document.getElementById('managementContextPanel').hidden, 'Settings route hides Management Context');
    check(document.body.textContent.includes('Desktop shell preferences') && document.getElementById('launchAtLogin'), 'Settings grouping preserves Desktop shell preference');
    check(window.__fakeADM.calls.length===afterExecAuthorityActions, 'settings routing after explicit exec actions makes no adapter calls');
    check(document.querySelectorAll('.nav-link[aria-current="page"]').length===1, 'one active menu item');

    const memoryCallsBefore = window.__fakeADM.calls.filter(c=>c.name==='ListGlobalMemory').length;
    await clickRoute('memory');
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='ListGlobalMemory').length>memoryCallsBefore, 'Memory route lazy-loads Global Memory for the management view');
    check(document.getElementById('globalMemoryList').textContent.includes('visible-after-explicit-load'), 'Memory route lazy-load renders value');

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
    const diagnosticTab=document.querySelector('#environmentDetailSubviewTabs button[data-environment-detail-subview="diagnostics"]');
    const diagnosticsCallCountBefore=window.__fakeADM.calls.length;
    diagnosticTab.click(); await sleep();
    check(!document.getElementById('environmentDiagnostics').hidden && document.getElementById('environmentDetail').hidden, 'Environment Diagnostics subview is local to detail modal');
    check(document.getElementById('environmentDiagnostics').textContent.includes('Existing InspectEnvironment payload only') && document.getElementById('environmentDiagnostics').textContent.includes('fixture capability report') && document.getElementById('environmentDiagnostics').textContent.includes('artifact_missing'), 'Diagnostics renders returned fact source reason and state');
    check(window.__fakeADM.calls.length===diagnosticsCallCountBefore, 'Diagnostics subview does not call probe verifier reconnect or Memory APIs');
    document.querySelector('#environmentDetailSubviewTabs button[data-environment-detail-subview="summary"]').click(); await sleep();
    check(!document.getElementById('environmentDetail').hidden && document.getElementById('environmentDiagnostics').hidden, 'Environment Summary subview restores summary without adapter calls');
    check(window.__fakeADM.calls.length===diagnosticsCallCountBefore, 'Summary/Diagnostics switching stays local');
    check(window.__fakeADM.calls.filter(c=>c.name==='ListEnvironmentMemory').length===detailMemoryBefore, 'opening Environment detail does not read private Memory values');
    location.hash='#/skills'; await sleep(60);
    check(visibleRoute()==='environments' && location.hash==='#/environments', 'Environment detail guards route change');
    document.getElementById('closeEnvironmentDetail').click(); await sleep();
    check(document.activeElement===detailButton, 'closing Environment detail restores opener focus');
    const currentEnvironmentRow=document.querySelector('#environmentList .managed-item.current-context');
    check(currentEnvironmentRow?.textContent.includes('Environment A') && currentEnvironmentRow?.textContent.includes('当前管理环境'), 'Environment list marks explicit current Management Environment');
    check(currentEnvironmentRow?.textContent.includes('Workspace A') && currentEnvironmentRow?.textContent.includes('MCP 1') && currentEnvironmentRow?.textContent.includes('Skills 2'), 'Environment rows join Workspace identity and selection counts');
    const diagnoseButton=document.querySelector('#environmentList button[data-action="diagnose-environment"][data-id="env-a"]');
    const diagnoseBefore=window.__fakeADM.calls.filter(c=>c.name==='InspectEnvironment' && c.args[0]==='env-a').length;
    const diagnoseCallCountBefore=window.__fakeADM.calls.length;
    diagnoseButton.click(); await waitFor(() => visibleRoute()==='diagnostics' && location.hash==='#/diagnostics', 'Environment diagnose action opens standalone Diagnostics route');
    await waitFor(() => document.getElementById('diagnosticsPageContent').textContent.includes('artifact_missing'), 'standalone Diagnostics route renders returned facts');
    check(document.getElementById('environmentDetailPanel').hidden, 'standalone Diagnostics route does not open Environment detail modal');
    check(window.__fakeADM.calls.filter(c=>c.name==='InspectEnvironment' && c.args[0]==='env-a').length===diagnoseBefore+1, 'Diagnose route reuses InspectEnvironment only');
    check(!window.__fakeADM.calls.slice(diagnoseCallCountBefore).some(c=>['ProbeMCPHealth','ListGlobalMemory','ListEnvironmentMemory','RunVerifier'].includes(c.name)), 'standalone Diagnostics route does not probe, run verifier, or read Memory values');
    detailButton.click(); await waitFor(() => !document.getElementById('environmentDetailPanel').hidden, 'Environment detail reopen for shortcut');
    document.querySelector('#environmentDetailRoutes button[data-detail-route="skills"]').click(); await sleep();
    check(document.getElementById('environmentDetailPanel').hidden && visibleRoute()==='skills', 'Environment detail shortcut closes modal then navigates');
    check(document.getElementById('managementEnvironment').value==='env-a', 'Environment detail shortcut preserves explicit Management Environment');
    check(window.__fakeADM.calls.filter(c=>c.name==='ListEnvironmentMemory').length===detailMemoryBefore, 'Environment detail shortcut does not implicitly load Memory');
    await clickRoute('environments');
    detailButton.click(); await waitFor(() => !document.getElementById('environmentDetailPanel').hidden, 'Environment detail reopen for Memory shortcut');
    document.querySelector('#environmentDetailRoutes button[data-detail-route="memory"]').click(); await sleep();
    check(document.getElementById('environmentDetailPanel').hidden && visibleRoute()==='memory', 'Environment detail Memory shortcut opens Memory page');
    check(document.getElementById('managementEnvironment').value==='env-a' && !document.getElementById('loadEnvironmentMemory').disabled, 'Memory page uses current Management Environment scope');
    check(window.__fakeADM.calls.filter(c=>c.name==='ListEnvironmentMemory').length===detailMemoryBefore, 'Memory route does not auto-read Environment-private values');
    document.getElementById('loadEnvironmentMemory').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='ListEnvironmentMemory' && c.args[0]==='env-a').length>detailMemoryBefore, 'explicit Environment-private Memory read');
    check(document.getElementById('environmentMemoryList').textContent.includes('env-a-private-visible'), 'Environment-private Memory renders only after explicit load');
    await clickRoute('environments');
    detailButton.click(); await waitFor(() => !document.getElementById('environmentDetailPanel').hidden, 'Environment detail can reopen after Memory load');
    document.getElementById('closeEnvironmentDetail').click(); await sleep();
    await clickRoute('memory');
    check(document.getElementById('environmentMemoryList').textContent.includes('env-a-private-visible'), 'closing Environment detail does not clear loaded Memory page values');

    await clickRoute('mcp');
    check(!document.getElementById('managementContextPanel').hidden, 'MCP route shows Management Context because Environment selection is relevant');
    const mcpFilter=document.getElementById('mcpFilter');
    const mcpStateFilter=document.getElementById('mcpStateFilter');
    mcpFilter.value='/MCP [AB]/i'; mcpFilter.dispatchEvent(new Event('input',{bubbles:true}));
    mcpStateFilter.value='unselected'; mcpStateFilter.dispatchEvent(new Event('change',{bubbles:true})); await sleep();
    check(document.getElementById('mcpVisibleCount').textContent==='1' && document.getElementById('mcpList').textContent.includes('MCP B') && !document.getElementById('mcpList').textContent.includes('MCP A'), 'MCP regex search and current-env-unselected filter combine locally');
    const setMCPDefaultBefore=window.__fakeADM.calls.filter(c=>c.name==='SetMCPDefault').length;
    document.getElementById('mcpSetVisibleDefaultButton').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='SetMCPDefault').length>setMCPDefaultBefore, 'MCP visible default batch call');
    check(window.__fakeADM.calls.some(c=>c.name==='SetMCPDefault' && c.args[0]==='mcp-b' && c.args[1]===true) && !window.__fakeADM.calls.some(c=>c.name==='SetMCPDefault' && c.args[0]==='mcp-a'), 'MCP default batch uses current filtered result only');
    const setMCPEnvBefore=window.__fakeADM.calls.filter(c=>c.name==='SetEnvironmentMCP').length;
    document.getElementById('mcpEnableVisibleButton').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='SetEnvironmentMCP').length>setMCPEnvBefore, 'MCP visible Environment batch call');
    check(window.__fakeADM.calls.some(c=>c.name==='SetEnvironmentMCP' && c.args[0]==='env-a' && c.args[1]==='mcp-b' && c.args[2]===true), 'MCP Environment batch captures current Environment and filtered ID');
    mcpFilter.value=''; mcpFilter.dispatchEvent(new Event('input',{bubbles:true})); mcpStateFilter.value='all'; mcpStateFilter.dispatchEvent(new Event('change',{bubbles:true})); await sleep();

    const mcpProbeBefore=window.__fakeADM.calls.filter(c=>c.name==='ProbeMCPHealth').length;
    document.querySelector('#mcpList button[data-action="probe-mcp"][data-id="mcp-a"]').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='ProbeMCPHealth').length>mcpProbeBefore, 'explicit MCP probe call');
    await waitFor(() => document.getElementById('mcpList').textContent.includes('fixture healthy'), 'MCP explicit probe renders observation');
    check(window.__fakeADM.calls.some(c=>c.name==='ProbeMCPHealth' && c.args[0]==='env-a' && c.args[1]==='mcp-a'), 'MCP probe captures Environment and definition identity');

    window.__fakeADM.state.delayMCPProbe=true;
    const delayedProbeBefore=window.__fakeADM.calls.filter(c=>c.name==='ProbeMCPHealth' && c.args[0]==='env-a' && c.args[1]==='mcp-a').length;
    document.querySelector('#mcpList button[data-action="probe-mcp"][data-id="mcp-a"]').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='ProbeMCPHealth' && c.args[0]==='env-a' && c.args[1]==='mcp-a').length>delayedProbeBefore, 'delayed MCP probe started');
    const mcpProbeEnvSelect=document.getElementById('managementEnvironment');
    mcpProbeEnvSelect.value='env-b'; mcpProbeEnvSelect.dispatchEvent(new Event('change',{bubbles:true}));
    await waitFor(() => document.getElementById('mcpList').textContent.includes('MCP A'), 'MCP list rerenders after Environment switch');
    await sleep(240);
    check(!document.getElementById('mcpList').textContent.includes('fixture healthy'), 'late MCP probe from old Environment cannot populate current Environment');
    window.__fakeADM.state.delayMCPProbe=false;
    mcpProbeEnvSelect.value='env-a'; mcpProbeEnvSelect.dispatchEvent(new Event('change',{bubbles:true}));
    await waitFor(() => document.getElementById('mcpList').textContent.includes('fixture healthy'), 'returning to original Environment shows its stored probe observation');

    const mcpImportButton=document.querySelector('[data-dialog-open="mcpImportDialog"]');
    mcpImportButton.click(); await waitFor(() => document.getElementById('mcpImportDialog').open, 'MCP import dialog opens');
    const importContent=document.getElementById('mcpImportContent');
    importContent.value='{ "mcpServers": { "imported": { "url": "http://127.0.0.1:9902/mcp", "headers": { "Authorization": "secret-token" } } } }';
    const previewBefore=window.__fakeADM.calls.filter(c=>c.name==='PreviewMCPImport').length;
    document.getElementById('mcpImportForm').requestSubmit();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='PreviewMCPImport').length>previewBefore, 'MCP import preview call');
    await waitFor(() => !document.getElementById('mcpImportApplyButton').disabled, 'valid MCP import preview enables apply');
    check(document.getElementById('mcpImportPreview').textContent.includes('ADM_IMPORTED_TOKEN'), 'MCP import preview keeps reference requirements visible');
    importContent.value=importContent.value.replace('9902','9903'); importContent.dispatchEvent(new Event('input',{bubbles:true})); await sleep();
    check(document.getElementById('mcpImportApplyButton').disabled && document.getElementById('mcpImportPreview').textContent.includes('请重新预览'), 'MCP import input changes invalidate preview before apply');
    document.getElementById('mcpImportForm').requestSubmit();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='PreviewMCPImport').length>previewBefore+1, 'MCP import preview can be regenerated');
    await waitFor(() => !document.getElementById('mcpImportApplyButton').disabled, 'regenerated MCP import preview enables apply');
    const applyBefore=window.__fakeADM.calls.filter(c=>c.name==='ApplyMCPImport').length;
    document.getElementById('mcpImportApplyButton').click(); document.getElementById('mcpImportApplyButton').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='ApplyMCPImport').length>applyBefore, 'MCP import apply call');
    await waitFor(() => document.getElementById('mcpListTotalCount').textContent==='3', 'MCP import apply refreshes authoritative catalog');
    check(window.__fakeADM.calls.filter(c=>c.name==='ApplyMCPImport').length===applyBefore+1, 'MCP import apply is guarded against double submit');

    await clickRoute('skills');
    check(!document.getElementById('skillsPanel').hidden && document.getElementById('skillSourcesPanel').hidden, 'Skill route defaults to Skills subview');
    document.querySelector('[data-skill-subview="sources"]').click(); await sleep();
    check(document.getElementById('skillsPanel').hidden && !document.getElementById('skillSourcesPanel').hidden, 'Skill Sources subview switches locally');
    check(document.getElementById('skillSourceVisibleCount').textContent==='1' && document.getElementById('skillSourceListTotalCount').textContent==='1', 'Skill Sources subview shows visible and total counts');
    const sourceSearch=document.getElementById('skillSourceFilterInput'); sourceSearch.value='/source-a|fixtures/'; sourceSearch.dispatchEvent(new Event('input',{bubbles:true})); await sleep();
    check(document.getElementById('skillSourceVisibleCount').textContent==='1', 'Skill source regex filter is local to source rows');
    const sourceEditBefore=window.__fakeADM.calls.filter(c=>c.name==='UpdateSkillSource').length;
    document.querySelector('#skillSourceList button[data-action="edit-skill-source"]').click();
    await waitFor(() => document.getElementById('skillSourceDialog').open, 'Skill source editor opens from row');
    check(document.getElementById('skillSourceRoot').value.includes('fixtures'), 'Skill source editor pre-fills existing root');
    document.getElementById('skillSupportRoots').value='C:\\fixtures\\support';
    document.getElementById('skillSourceForm').requestSubmit();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='UpdateSkillSource').length>sourceEditBefore, 'Skill source edit uses update API');
    check(!window.__fakeADM.calls.some(c=>c.name==='RefreshSkillSource'), 'Skill source edit does not implicitly refresh source');
    document.querySelector('[data-skill-subview="skills"]').click(); await sleep();
    check(!document.getElementById('skillsPanel').hidden && document.getElementById('skillSourcesPanel').hidden, 'Skill subview returns to installed Skills locally');
    const sourceFilterSelect=document.getElementById('skillSourceFilter'); sourceFilterSelect.value='source-a'; sourceFilterSelect.dispatchEvent(new Event('change',{bubbles:true})); await sleep();
    check(document.getElementById('skillVisibleCount').textContent==='3' && document.getElementById('skillList').textContent.includes('Source：skills · source-a') && !document.getElementById('skillList').textContent.includes('Legacy metadata'), 'Skill Source filter isolates source-managed rows and separate Source fact');
    sourceFilterSelect.value=''; sourceFilterSelect.dispatchEvent(new Event('change',{bubbles:true})); await sleep();

    const skillFilter=document.getElementById('skillFilter');
    const skillStateFilter=document.getElementById('skillStateFilter');
    skillFilter.value='/Idle|Legacy Skill 2/'; skillFilter.dispatchEvent(new Event('input',{bubbles:true}));
    skillStateFilter.value='unselected'; skillStateFilter.dispatchEvent(new Event('change',{bubbles:true})); await sleep();
    check(document.getElementById('skillVisibleCount').textContent==='2' && document.getElementById('skillList').textContent.includes('Idle Skill') && document.getElementById('skillList').textContent.includes('Legacy Skill 2'), 'Skill regex search and current-env-unselected filter combine locally');
    const setSkillDefaultBefore=window.__fakeADM.calls.filter(c=>c.name==='SetSkillDefault').length;
    document.getElementById('skillSetVisibleDefaultButton').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='SetSkillDefault').length>=setSkillDefaultBefore+2, 'Skill visible default batch calls');
    check(window.__fakeADM.calls.some(c=>c.name==='SetSkillDefault' && c.args[0]==='skill-idle' && c.args[1]===true) && window.__fakeADM.calls.some(c=>c.name==='SetSkillDefault' && c.args[0]==='skill-legacy-2' && c.args[1]===true) && !window.__fakeADM.calls.some(c=>c.name==='SetSkillDefault' && c.args[0]==='skill-a'), 'Skill default batch uses current filtered results only');
    const setSkillEnvBefore=window.__fakeADM.calls.filter(c=>c.name==='SetEnvironmentSkill').length;
    document.getElementById('skillEnableVisibleButton').click();
    await waitFor(() => window.__fakeADM.calls.filter(c=>c.name==='SetEnvironmentSkill').length>=setSkillEnvBefore+2, 'Skill visible Environment batch calls');
    check(window.__fakeADM.calls.some(c=>c.name==='SetEnvironmentSkill' && c.args[0]==='env-a' && c.args[1]==='skill-idle' && c.args[2]===true) && window.__fakeADM.calls.some(c=>c.name==='SetEnvironmentSkill' && c.args[0]==='env-a' && c.args[1]==='skill-legacy-2' && c.args[2]===true), 'Skill Environment batch captures current Environment and filtered IDs');
    skillFilter.value=''; skillFilter.dispatchEvent(new Event('input',{bubbles:true})); skillStateFilter.value='all'; skillStateFilter.dispatchEvent(new Event('change',{bubbles:true})); await sleep();

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

    check(window.__fakeADM.calls.filter(c=>c.name==='ProbeMCPHealth').length===2, 'only the two explicit MCP probes were executed');
    const forbiddenAutoCalls=['RefreshSkillSource','RunVerifier','StopProcess','CancelRun','WriteGlobalMemory','WriteEnvironmentMemory'];
    check(!window.__fakeADM.calls.some(c=>forbiddenAutoCalls.includes(c.name)), 'navigation/refresh makes no implicit mutation call');

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
    '--virtual-time-budget=30000', '--dump-dom', pathToFileURL(fixture).href,
  ];
  const execution = spawnSync(browser, args, {encoding:'utf8', timeout:100000, maxBuffer:20*1024*1024});
  try {
    if (execution.error) throw execution.error;
    if (execution.status !== 0) throw new Error(`browser exit ${execution.status}: ${execution.stderr}`);
    const match = execution.stdout.match(/data-browser-smoke="([^"]+)"/);
    if (!match) throw new Error(`browser smoke result missing; status=${execution.status} signal=${execution.signal} error=${execution.error?.message || ''} stderr=${execution.stderr.slice(-2000)}\nstdoutLen=${execution.stdout.length} stdoutTail=${execution.stdout.slice(-4000)}`);
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
