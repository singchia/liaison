// Supervises the bundled liaison-edge child process.
//
// Responsibilities:
//   - Translate IntendedState (Running/Stopped) into actual subprocess
//     spawn / kill, with exponential backoff on unexpected exit.
//   - Broadcast TrayState transitions so the menubar UI can react.
//   - Render an EdgeConfig into the YAML schema liaison-edge expects
//     (see dist/edge/liaison-edge.yaml.template).
//
// Driving the supervisor to TrayState::Online is the *poller's* job
// (T9, via api_client::list_edges) — this module only knows whether
// the subprocess is alive, not whether it's actually connected.

#![allow(dead_code)]

use std::path::{Path, PathBuf};
use std::sync::Arc;
use std::time::{Duration, Instant};

use serde::Serialize;
use thiserror::Error;
use tokio::process::Command;
use tokio::sync::{broadcast, watch};

const BACKOFF_INITIAL: Duration = Duration::from_secs(2);
const BACKOFF_MAX: Duration = Duration::from_secs(30);
const HEALTHY_RUNTIME: Duration = Duration::from_secs(30);

#[derive(Debug, Error)]
pub enum SupervisorError {
    #[error("io: {0}")]
    Io(#[from] std::io::Error),
    #[error("yaml: {0}")]
    Yaml(#[from] serde_yml::Error),
    #[error("home directory not found")]
    NoHome,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize)]
pub enum IntendedState {
    Running,
    Stopped,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(tag = "kind", content = "message")]
pub enum TrayState {
    LoggedOut,
    Paused,
    Connecting,
    Online,
    Error(String),
}

#[derive(Debug, Clone)]
pub struct EdgeConfig {
    pub manager_addr: String,
    pub access_key: String,
    pub secret_key: String,
    pub tls_enable: bool,
    pub insecure_skip_verify: bool,
    pub log_file: PathBuf,
}

impl EdgeConfig {
    pub fn to_yaml_string(&self) -> Result<String, serde_yml::Error> {
        use serde_yml::{Mapping, Value};

        fn m() -> Mapping {
            Mapping::new()
        }

        let mut tls = m();
        tls.insert("enable".into(), Value::Bool(self.tls_enable));
        tls.insert(
            "insecure_skip_verify".into(),
            Value::Bool(self.insecure_skip_verify),
        );

        let mut dial = m();
        dial.insert(
            "addrs".into(),
            Value::Sequence(vec![Value::String(self.manager_addr.clone())]),
        );
        dial.insert("network".into(), Value::String("tcp".into()));
        dial.insert("tls".into(), Value::Mapping(tls));

        let mut auth = m();
        auth.insert(
            "access_key".into(),
            Value::String(self.access_key.clone()),
        );
        auth.insert(
            "secret_key".into(),
            Value::String(self.secret_key.clone()),
        );

        let mut manager = m();
        manager.insert("dial".into(), Value::Mapping(dial));
        manager.insert("auth".into(), Value::Mapping(auth));

        let mut log = m();
        log.insert("level".into(), Value::String("info".into()));
        log.insert(
            "file".into(),
            Value::String(self.log_file.to_string_lossy().into_owned()),
        );
        log.insert("maxsize".into(), Value::Number(100.into()));
        log.insert("maxrolls".into(), Value::Number(10.into()));

        let mut root = m();
        root.insert("manager".into(), Value::Mapping(manager));
        root.insert("log".into(), Value::Mapping(log));

        serde_yml::to_string(&Value::Mapping(root))
    }
}

pub fn write_edge_config(path: &Path, cfg: &EdgeConfig) -> Result<(), SupervisorError> {
    if let Some(parent) = path.parent() {
        std::fs::create_dir_all(parent)?;
    }
    let body = cfg.to_yaml_string()?;
    std::fs::write(path, body)?;
    #[cfg(unix)]
    {
        use std::os::unix::fs::PermissionsExt;
        let mut perms = std::fs::metadata(path)?.permissions();
        perms.set_mode(0o600);
        std::fs::set_permissions(path, perms)?;
    }
    Ok(())
}

pub fn liaison_data_dir() -> Result<PathBuf, SupervisorError> {
    let base = dirs::data_dir().ok_or(SupervisorError::NoHome)?;
    Ok(base.join("Liaison"))
}

pub fn liaison_log_dir() -> Result<PathBuf, SupervisorError> {
    #[cfg(target_os = "macos")]
    {
        let home = dirs::home_dir().ok_or(SupervisorError::NoHome)?;
        Ok(home.join("Library").join("Logs").join("Liaison"))
    }
    #[cfg(not(target_os = "macos"))]
    {
        let base = dirs::data_local_dir().ok_or(SupervisorError::NoHome)?;
        Ok(base.join("Liaison").join("logs"))
    }
}

#[derive(Clone)]
pub struct SupervisorHandle {
    intended_tx: watch::Sender<IntendedState>,
    state_tx: broadcast::Sender<TrayState>,
    last_state: Arc<std::sync::Mutex<TrayState>>,
}

impl SupervisorHandle {
    pub fn set_intended(&self, state: IntendedState) {
        let _ = self.intended_tx.send(state);
    }

    pub fn intended(&self) -> IntendedState {
        *self.intended_tx.borrow()
    }

    pub fn subscribe(&self) -> broadcast::Receiver<TrayState> {
        self.state_tx.subscribe()
    }

    pub fn current(&self) -> TrayState {
        self.last_state.lock().expect("poisoned").clone()
    }

    /// Force a tray state transition + broadcast — used by cmd_pause/
    /// cmd_resume so the popup sees the new state on the very next
    /// refresh() instead of racing the supervisor's run loop.
    pub fn force_state(&self, s: TrayState) {
        crate::debug_log(format!("handle.force_state {:?}", s));
        *self.last_state.lock().expect("poisoned") = s.clone();
        let _ = self.state_tx.send(s);
    }

    /// Promote Connecting -> Online when the API status poller sees the
    /// edge as online. Never overrides Paused / Error / LoggedOut.
    pub fn report_tunnel_online(&self) {
        let mut last = self.last_state.lock().expect("poisoned");
        if matches!(*last, TrayState::Connecting) {
            *last = TrayState::Online;
            let _ = self.state_tx.send(TrayState::Online);
        }
    }

    /// Demote Online -> Connecting when the poller sees the edge offline.
    pub fn report_tunnel_offline(&self) {
        let mut last = self.last_state.lock().expect("poisoned");
        if matches!(*last, TrayState::Online) {
            *last = TrayState::Connecting;
            let _ = self.state_tx.send(TrayState::Connecting);
        }
    }
}

pub struct EdgeSupervisor {
    binary_path: PathBuf,
    config_path: PathBuf,
    /// Where stdout/stderr from the edge subprocess get appended.
    /// Critical on Windows: with CREATE_NO_WINDOW the child has no
    /// console attached, so inheriting our parent stdio gives it a
    /// closed handle that some Go logging paths blow up on. Explicit
    /// file redirection both fixes that and gives us a place to read
    /// boot-time crashes from.
    supervisor_log: PathBuf,
    intended_tx: watch::Sender<IntendedState>,
    intended_rx: watch::Receiver<IntendedState>,
    state_tx: broadcast::Sender<TrayState>,
    last_state: Arc<std::sync::Mutex<TrayState>>,
}

impl EdgeSupervisor {
    pub fn new(
        binary_path: PathBuf,
        config_path: PathBuf,
        supervisor_log: PathBuf,
    ) -> Self {
        Self::with_intent(binary_path, config_path, supervisor_log, IntendedState::Running)
    }

    /// Construct with a specific initial intent so callers can seed
    /// the watch channel from persisted state without first calling
    /// set_intended() — which is technically a no-op if the value
    /// matches but still bumps the watch's version, causing the
    /// run-loop to fire a spurious "intent changed mid-run" on its
    /// first iteration and pointlessly kill+respawn the edge child.
    pub fn with_intent(
        binary_path: PathBuf,
        config_path: PathBuf,
        supervisor_log: PathBuf,
        initial_intent: IntendedState,
    ) -> Self {
        let (intended_tx, intended_rx) = watch::channel(initial_intent);
        let (state_tx, _) = broadcast::channel(16);
        Self {
            binary_path,
            config_path,
            supervisor_log,
            intended_tx,
            intended_rx,
            state_tx,
            last_state: Arc::new(std::sync::Mutex::new(TrayState::LoggedOut)),
        }
    }

    pub fn handle(&self) -> SupervisorHandle {
        SupervisorHandle {
            intended_tx: self.intended_tx.clone(),
            state_tx: self.state_tx.clone(),
            last_state: self.last_state.clone(),
        }
    }

    fn emit(&self, s: TrayState) {
        crate::debug_log(format!("supervisor: emit {:?}", s));
        *self.last_state.lock().expect("poisoned") = s.clone();
        let n = self.state_tx.send(s).unwrap_or(0);
        crate::debug_log(format!("supervisor: broadcast receivers={}", n));
    }

    pub async fn run(self) {
        let mut intended_rx = self.intended_rx.clone();
        let mut backoff = BACKOFF_INITIAL;

        loop {
            if *intended_rx.borrow() == IntendedState::Stopped {
                self.emit(TrayState::Paused);
                if intended_rx.changed().await.is_err() {
                    return;
                }
                continue;
            }

            self.emit(TrayState::Connecting);

            // Open / append the supervisor log to capture edge's
            // stdout + stderr. Best-effort: if we can't open it we
            // still spawn (the child gets the inherited handle, and
            // on Windows-with-no-console that may bite, but it's
            // strictly worse to block edge from running over a
            // missing log file).
            if let Some(parent) = self.supervisor_log.parent() {
                let _ = std::fs::create_dir_all(parent);
            }
            let log_handle = std::fs::OpenOptions::new()
                .create(true)
                .append(true)
                .open(&self.supervisor_log);

            let mut cmd = Command::new(&self.binary_path);
            // liaison-edge uses Go's flag package which only takes
            // single-dash short flags; --config is reported as
            // "flag provided but not defined: -config" and the
            // process exits with status 2.
            cmd.arg("-c")
                .arg(&self.config_path)
                .kill_on_drop(true);

            if let Ok(file) = log_handle {
                if let Ok(stdout_clone) = file.try_clone() {
                    cmd.stdout(std::process::Stdio::from(stdout_clone));
                }
                cmd.stderr(std::process::Stdio::from(file));
            }

            // Suppress the brief console window Windows flashes when a
            // CLI binary is launched from a GUI process. CREATE_NO_WINDOW
            // = 0x08000000. Has no effect on macOS / Linux. Combined
            // with the explicit stdio redirection above, this gives
            // edge valid file handles for stdout/stderr instead of the
            // closed handles it would inherit from the (windowless)
            // GUI process.
            #[cfg(target_os = "windows")]
            {
                use std::os::windows::process::CommandExt;
                const CREATE_NO_WINDOW: u32 = 0x08000000;
                cmd.creation_flags(CREATE_NO_WINDOW);
            }
            let spawned = cmd.spawn();

            match spawned {
                Err(e) => {
                    self.emit(TrayState::Error(format!(
                        "spawn failed: {e}\npath: {}",
                        self.binary_path.display()
                    )));
                    sleep_or_change(&mut intended_rx, backoff).await;
                    backoff = (backoff * 2).min(BACKOFF_MAX);
                }
                Ok(mut child) => {
                    let started = Instant::now();
                    tokio::select! {
                        res = child.wait() => {
                            let was_healthy = started.elapsed() > HEALTHY_RUNTIME;
                            let msg = match res {
                                Ok(status) if status.success() => "edge exited cleanly".to_string(),
                                Ok(status) => format!("edge exited: {status}"),
                                Err(e) => format!("edge wait failed: {e}"),
                            };
                            self.emit(TrayState::Error(msg));
                            sleep_or_change(&mut intended_rx, backoff).await;
                            backoff = if was_healthy {
                                BACKOFF_INITIAL
                            } else {
                                (backoff * 2).min(BACKOFF_MAX)
                            };
                        }
                        _ = intended_rx.changed() => {
                            let new_intent = *intended_rx.borrow();
                            crate::debug_log(format!(
                                "supervisor: intent changed mid-run -> {:?}, killing child",
                                new_intent
                            ));
                            let kill_res = child.start_kill();
                            let wait_res = child.wait().await;
                            crate::debug_log(format!(
                                "supervisor: child killed: start_kill={:?}, wait_status={:?}",
                                kill_res.is_ok(),
                                wait_res.is_ok()
                            ));
                            // Loop top will re-check intent and emit Paused.
                        }
                    }
                }
            }
        }
    }
}

async fn sleep_or_change(rx: &mut watch::Receiver<IntendedState>, delay: Duration) {
    tokio::select! {
        _ = tokio::time::sleep(delay) => {},
        _ = rx.changed() => {},
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::time::Duration;

    fn sample_cfg() -> EdgeConfig {
        EdgeConfig {
            manager_addr: "manager.liaison.cloud:30012".into(),
            access_key: "ak_test".into(),
            secret_key: "sk_test".into(),
            tls_enable: true,
            insecure_skip_verify: true,
            log_file: PathBuf::from("/tmp/liaison-edge.log"),
        }
    }

    #[test]
    fn yaml_contains_required_keys() {
        let yaml = sample_cfg().to_yaml_string().unwrap();
        // Round-trip through serde_yml to assert structure rather than
        // string-matching, which is fragile against indentation tweaks.
        let v: serde_yml::Value = serde_yml::from_str(&yaml).unwrap();
        let manager = v.get("manager").expect("manager");
        let dial = manager.get("dial").expect("dial");
        let addrs = dial.get("addrs").and_then(|x| x.as_sequence()).expect("addrs");
        assert_eq!(
            addrs[0].as_str().unwrap(),
            "manager.liaison.cloud:30012"
        );
        assert_eq!(dial.get("network").unwrap().as_str().unwrap(), "tcp");
        let tls = dial.get("tls").expect("tls");
        assert_eq!(tls.get("enable").unwrap().as_bool().unwrap(), true);
        assert_eq!(
            tls.get("insecure_skip_verify").unwrap().as_bool().unwrap(),
            true
        );
        let auth = manager.get("auth").expect("auth");
        assert_eq!(auth.get("access_key").unwrap().as_str().unwrap(), "ak_test");
        assert_eq!(auth.get("secret_key").unwrap().as_str().unwrap(), "sk_test");
        let log = v.get("log").expect("log");
        assert_eq!(log.get("level").unwrap().as_str().unwrap(), "info");
        assert_eq!(
            log.get("file").unwrap().as_str().unwrap(),
            "/tmp/liaison-edge.log"
        );
    }

    #[test]
    fn write_edge_config_creates_dir_and_locks_perms() {
        let dir = std::env::temp_dir().join(format!(
            "liaison-test-{}",
            std::time::SystemTime::now()
                .duration_since(std::time::UNIX_EPOCH)
                .unwrap()
                .as_nanos()
        ));
        let path = dir.join("nested").join("liaison-edge.yaml");
        write_edge_config(&path, &sample_cfg()).unwrap();
        let body = std::fs::read_to_string(&path).unwrap();
        assert!(body.contains("ak_test"));
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let perms = std::fs::metadata(&path).unwrap().permissions();
            assert_eq!(perms.mode() & 0o777, 0o600);
        }
        std::fs::remove_dir_all(&dir).ok();
    }

    #[tokio::test]
    async fn supervisor_emits_paused_when_intended_stopped_at_start() {
        let sup =
            EdgeSupervisor::new(
            PathBuf::from("/nonexistent"),
            PathBuf::from("/nonexistent"),
            std::env::temp_dir().join("liaison-supervisor-test.log"),
        );
        let handle = sup.handle();
        let mut rx = handle.subscribe();
        handle.set_intended(IntendedState::Stopped);
        let join = tokio::spawn(sup.run());

        let evt = tokio::time::timeout(Duration::from_secs(1), rx.recv())
            .await
            .expect("timeout")
            .expect("recv");
        assert_eq!(evt, TrayState::Paused);

        join.abort();
        let _ = join.await;
    }

    #[tokio::test]
    async fn supervisor_reports_error_when_binary_missing() {
        let sup = EdgeSupervisor::new(
            PathBuf::from("/nonexistent/liaison-edge"),
            PathBuf::from("/nonexistent/liaison-edge.yaml"),
            std::env::temp_dir().join("liaison-supervisor-test.log"),
        );
        let handle = sup.handle();
        let mut rx = handle.subscribe();
        let join = tokio::spawn(sup.run());

        let mut saw_connecting = false;
        let mut saw_error = false;
        for _ in 0..4 {
            match tokio::time::timeout(Duration::from_secs(1), rx.recv()).await {
                Ok(Ok(TrayState::Connecting)) => saw_connecting = true,
                Ok(Ok(TrayState::Error(_))) => {
                    saw_error = true;
                    break;
                }
                Ok(Ok(_)) => continue,
                _ => break,
            }
        }
        assert!(saw_connecting, "should have emitted Connecting");
        assert!(saw_error, "should have emitted Error for missing binary");

        join.abort();
        let _ = join.await;
    }

    fn handle_seeded_with(state: TrayState) -> SupervisorHandle {
        let sup = EdgeSupervisor::new(
            PathBuf::from("/x"),
            PathBuf::from("/x"),
            std::env::temp_dir().join("liaison-supervisor-test.log"),
        );
        let h = sup.handle();
        *h.last_state.lock().unwrap() = state;
        h
    }

    #[test]
    fn report_tunnel_online_promotes_connecting() {
        let h = handle_seeded_with(TrayState::Connecting);
        h.report_tunnel_online();
        assert_eq!(h.current(), TrayState::Online);
    }

    #[test]
    fn report_tunnel_online_does_not_override_paused() {
        let h = handle_seeded_with(TrayState::Paused);
        h.report_tunnel_online();
        assert_eq!(h.current(), TrayState::Paused);
    }

    #[test]
    fn report_tunnel_online_does_not_override_error() {
        let h = handle_seeded_with(TrayState::Error("boom".into()));
        h.report_tunnel_online();
        assert_eq!(h.current(), TrayState::Error("boom".into()));
    }

    #[test]
    fn report_tunnel_offline_demotes_online() {
        let h = handle_seeded_with(TrayState::Online);
        h.report_tunnel_offline();
        assert_eq!(h.current(), TrayState::Connecting);
    }

    #[test]
    fn report_tunnel_offline_does_not_override_paused() {
        let h = handle_seeded_with(TrayState::Paused);
        h.report_tunnel_offline();
        assert_eq!(h.current(), TrayState::Paused);
    }
}
