// Persisted intent (running / stopped) + edge ID, kept on disk so a
// user-initiated pause survives quit + relaunch.
//
// File location:
//   macOS:   ~/Library/Application Support/Liaison/state.json
//   Linux:   ~/.config/Liaison/state.json (via XDG_DATA_HOME)
//   Windows: %APPDATA%\Liaison\state.json
//
// Read failures (missing file, corrupt JSON, perm denied) are not
// fatal — we fall back to the default and let the user re-establish
// intent through the popup.

#![allow(dead_code)]

use std::path::PathBuf;

use serde::{Deserialize, Serialize};

use crate::edge_supervisor::{liaison_data_dir, IntendedState};

const STATE_FILENAME: &str = "state.json";

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PersistedState {
    #[serde(default = "default_intended")]
    pub intended: String,
    #[serde(default)]
    pub edge_id: Option<u64>,
    /// Liaison server base URL the user is logged into. Defaults to
    /// the public SaaS endpoint; users with their own deployment can
    /// switch via the popup's "更改服务器" entry.
    #[serde(default = "default_base_url")]
    pub base_url: String,
}

fn default_intended() -> String {
    "running".to_string()
}

pub fn default_base_url() -> String {
    "https://liaison.cloud".to_string()
}

impl Default for PersistedState {
    fn default() -> Self {
        Self {
            intended: default_intended(),
            edge_id: None,
            base_url: default_base_url(),
        }
    }
}

impl PersistedState {
    pub fn intended_as_enum(&self) -> IntendedState {
        match self.intended.as_str() {
            "stopped" => IntendedState::Stopped,
            _ => IntendedState::Running,
        }
    }

    pub fn set_intended(&mut self, state: IntendedState) {
        self.intended = match state {
            IntendedState::Running => "running",
            IntendedState::Stopped => "stopped",
        }
        .to_string();
    }
}

fn state_path() -> Option<PathBuf> {
    liaison_data_dir().ok().map(|d| d.join(STATE_FILENAME))
}

pub fn load() -> PersistedState {
    let Some(path) = state_path() else {
        return PersistedState::default();
    };
    match std::fs::read(&path) {
        Ok(bytes) => serde_json::from_slice(&bytes).unwrap_or_else(|e| {
            eprintln!(
                "[state] failed to parse {}: {e}; using defaults",
                path.display()
            );
            PersistedState::default()
        }),
        Err(e) if e.kind() == std::io::ErrorKind::NotFound => PersistedState::default(),
        Err(e) => {
            eprintln!(
                "[state] failed to read {}: {e}; using defaults",
                path.display()
            );
            PersistedState::default()
        }
    }
}

pub fn save(s: &PersistedState) {
    let Some(path) = state_path() else {
        return;
    };
    if let Some(parent) = path.parent() {
        if let Err(e) = std::fs::create_dir_all(parent) {
            eprintln!("[state] mkdir {} failed: {e}", parent.display());
            return;
        }
    }
    match serde_json::to_vec_pretty(s) {
        Ok(bytes) => {
            if let Err(e) = std::fs::write(&path, bytes) {
                eprintln!("[state] write {} failed: {e}", path.display());
            }
        }
        Err(e) => eprintln!("[state] serialize failed: {e}"),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn default_round_trip() {
        let s = PersistedState::default();
        let json = serde_json::to_string(&s).unwrap();
        let back: PersistedState = serde_json::from_str(&json).unwrap();
        assert_eq!(back.intended, "running");
        assert!(back.edge_id.is_none());
        assert_eq!(back.base_url, "https://liaison.cloud");
    }

    #[test]
    fn legacy_state_without_base_url_loads_with_default() {
        // state.json files written by earlier builds didn't have a
        // base_url field — they should still parse, defaulting to
        // the public SaaS endpoint so the user isn't logged out.
        let json = r#"{"intended":"running","edge_id":42}"#;
        let s: PersistedState = serde_json::from_str(json).unwrap();
        assert_eq!(s.base_url, "https://liaison.cloud");
        assert_eq!(s.edge_id, Some(42));
    }

    #[test]
    fn intended_enum_conversion() {
        let mut s = PersistedState::default();
        assert_eq!(s.intended_as_enum(), IntendedState::Running);

        s.set_intended(IntendedState::Stopped);
        assert_eq!(s.intended, "stopped");
        assert_eq!(s.intended_as_enum(), IntendedState::Stopped);

        s.set_intended(IntendedState::Running);
        assert_eq!(s.intended, "running");
    }

    #[test]
    fn unknown_intended_falls_back_to_running() {
        let s = PersistedState {
            intended: "weird".into(),
            edge_id: None,
            base_url: default_base_url(),
        };
        assert_eq!(s.intended_as_enum(), IntendedState::Running);
    }

    #[test]
    fn missing_field_uses_default() {
        // Past versions only wrote `intended`. Confirm we still load.
        let json = r#"{"intended":"stopped"}"#;
        let s: PersistedState = serde_json::from_str(json).unwrap();
        assert_eq!(s.intended, "stopped");
        assert!(s.edge_id.is_none());
    }
}
