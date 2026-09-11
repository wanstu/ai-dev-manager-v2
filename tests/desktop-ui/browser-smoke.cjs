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
  const browserErrors = [];
  window.addEventListener('error', (event) => browserErrors.push(event.error?.stack || event.message || String(event.error || 'error')));
  window.addEventListener('unhandledrejection', (event) => browserErrors.push(event.reason?.stack || event.reason?.message || String(event.reason || 'rejection')));
  let activeID = 'profile-a';
  const profiles = [
    {id: 'profile-a', name: 'Profile A', base_url: 'http://127.0.0.1:43137'},
    {id: 'profile-b', name: 'Profile B', base_url: 'http://127.0.0.1:43138'},
  ];
  const longRoot = 'C:\\fixtures\\' + 'very-long-segment-'.repeat(18) + 'workspace-a';
  const snapshotA = {
    workspaces: [{workspace_id:'ws-a', name:'Workspace A', path:longRoot}],
    environments: [
      {environment_id:'env-a', workspace_id:'ws-a', name:'Environment A', root:longRoot, state:'ready', writer:{owner:'writer-a'}, enabled_mcp_ids:['mcp-a'], enabled_skill_ids:['skill-a'], private_memory_count:1},
      {environment_id:'env-b', workspace_id:'ws-a', name:'Environment B', root:longRoot+'\\b', state:'ready', writer:null, enabled_mcp_ids:[], enabled_skill_ids:[], private_memory_count:0},
    ],
    allowed_executables: ['go'],
    mcps: [{id:'mcp-a', name:'MCP A', transport:'streamable-http', endpoint:'http://127.0.0.1:9900/mcp', default_include_in_environment:false, health_policy:{}}],
    skills: [{id:'skill-a', name:'Skill A', source_id:'source-a', source_root:'C:\\fixtures\\skills', artifact_path:'C:\\fixtures\\skills\\skill-a\\SKILL.md', relative_artifact_path:'skill-a/SKILL.md', support_roots:[], default_include_in_environment:false}],
    global_memory_count: 1,
  };
  const snapshotB = {
    workspaces: [{workspace_id:'ws-profile-b', name:'Workspace Profile B', path:'C:\\fixtures\\profile-b'}],
    environments: [{environment_id:'env-profile-b', workspace_id:'ws-profile-b', name:'Environment Profile B', root:'C:\\fixtures\\profile-b', state:'ready', writer:null, enabled_mcp_ids:[], enabled_skill_ids:[], private_memory_count:0}],
    allowed_executables: [], mcps: [], skills: [], global_memory_count: 0,
  };
  const state = {failSkillSources:false, failVerifiers:false};
  const record = (name, args) => calls.push({name, args});
  const adapter = {
    async GetConnectionProfiles(){ record('GetConnectionProfiles',[]); return {profiles, active_id:activeID}; },
    async SelectConnectionProfile(id){ record('SelectConnectionProfile',[id]); activeID=id; return {profiles, active_id:activeID}; },
    async SaveConnectionProfile(profile){ record('SaveConnectionProfile',[profile]); return {profiles, active_id:activeID}; },
    async DeleteConnectionProfile(id){ record('DeleteConnectionProfile',[id]); return {profiles, active_id:activeID}; },
    async DisconnectADM(){ record('DisconnectADM',[]); return true; },
    async ConnectADM(input){ record('ConnectADM',[input]); const base=profiles.find(p=>p.id===activeID)?.base_url || input?.base_url || ''; return {state:'running', base_url:base, health_url:base+'/healthz', agent_mcp_url:base+'/mcp', admin_mcp_url:base+'/admin/mcp', pid:1234, version:'browser-fixture', local_bootstrap_eligible:false}; },
    async GetSnapshot(){ record('GetSnapshot',[]); return activeID==='profile-b' ? structuredClone(snapshotB) : structuredClone(snapshotA); },
    async AddWorkspace(input){ record('AddWorkspace',[input]); throw new Error('fixture workspace save failure'); },
    async ListSkillSources(){ record('ListSkillSources',[]); if(state.failSkillSources) throw new Error('skill sources unavailable'); return activeID==='profile-b' ? [] : [{skill_source_id:'source-a', root:'C:\\fixtures\\skills', support_roots:[], last_refresh_status:'ok'}]; },
    async InspectEnvironment(id){ record('InspectEnvironment',[id]); const snapshot=activeID==='profile-b'?snapshotB:snapshotA; const env=snapshot.environments.find(e=>e.environment_id===id); const workspace=snapshot.workspaces.find(w=>w.workspace_id===env?.workspace_id); return {environment:structuredClone(env), workspace:structuredClone(workspace), capability_report:{facts:[]}, unresolved_mcp_ids:[], unresolved_skill_ids:[]}; },
    async ListEnvironmentSkills(id){ record('ListEnvironmentSkills',[id]); return {skills:[]}; },
    async ListVerifiers(id){ record('ListVerifiers',[id]); if(state.failVerifiers) throw new Error('verifiers unavailable'); return [{verifier_id:'verifier-a', name:'Verifier A', kind:'test', executable:'go', args:['test','./...'], enabled:true}]; },
    async ListProcesses(id){ record('ListProcesses',[id]); return [{id:'proc-a', state:'running', pid:2222, listening_ports:[8080]}]; },
    async ListRuns(id){ record('ListRuns',[id]); return [{id:'run-a', state:'succeeded', executable:'go', args:['test'], stdout:'RUN_A_OUTPUT', stderr:'', exit_code:0}]; },
    async ListGlobalMemory(){ record('ListGlobalMemory',[]); return [{key:'sentinel', value:'visible-after-explicit-load'}]; },
    async GetDesktopPreferences(){ record('GetDesktopPreferences',[]); return {launch_at_login_supported:false, launch_at_login:false}; },
    async SetLaunchAtLogin(value){ record('SetLaunchAtLogin',[value]); return {launch_at_login_supported:false, launch_at_login:false}; },
  };
  window.go = {desktop:{Adapter:adapter}};
  window.runtime = {EventsOn(){}, EventsEmit(){}};
  window.confirm = () => true;
  window.__fakeADM = {calls, state, browserErrors};
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
    check(document.getElementById('workspaceCount').textContent==='1', 'successful snapshot renders count');
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
    const detailButton=document.querySelector('#environmentList button[data-action="inspect-environment"]');
    detailButton.focus(); detailButton.click();
    await waitFor(() => !document.getElementById('environmentDetailPanel').hidden, 'Environment detail open');
    location.hash='#/skills'; await sleep(60);
    check(visibleRoute()==='environments' && location.hash==='#/environments', 'Environment detail guards route change');
    document.getElementById('closeEnvironmentDetail').click(); await sleep();
    check(document.activeElement===detailButton, 'closing Environment detail restores opener focus');

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

    const forbiddenAutoCalls=['ProbeMCPHealth','RefreshSkillSource','RunVerifier','StopProcess','CancelRun','WriteGlobalMemory','WriteEnvironmentMemory'];
    check(!window.__fakeADM.calls.some(c=>forbiddenAutoCalls.includes(c.name)), 'navigation/refresh makes no implicit probe/mutation call');

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
    '--virtual-time-budget=7000', '--dump-dom', pathToFileURL(fixture).href,
  ];
  const execution = spawnSync(browser, args, {encoding:'utf8', timeout:20000, maxBuffer:20*1024*1024});
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
