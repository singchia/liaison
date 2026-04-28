// Detect a pre-existing liaison-edge CLI installation so the GUI can
// warn the user before it auto-creates a second connector.
//
// Path layouts mirror the install scripts the CLI uses:
//   macOS / Linux: install.sh -> /usr/local/bin/liaison-edge
//                                ~/.config/liaison-edge/liaison-edge.yaml
//                                LaunchAgent / systemd unit
//   Windows:       install.ps1 -> C:\Program Files\Liaison\bin\liaison-edge.exe
//                                C:\Program Files\Liaison\conf\liaison-edge.yaml
//                                Scheduled Task "LiaisonEdge"
//                                HKLM:\...\Run\LiaisonEdge

#![allow(dead_code)]

use std::path::PathBuf;

use serde::Serialize;

#[derive(Debug, Clone, Serialize)]
pub struct CliHit {
    pub kind: HitKind,
    pub detail: String,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize)]
#[serde(rename_all = "snake_case")]
pub enum HitKind {
    Binary,
    Config,
    Autostart,
}

pub fn detect_existing_cli() -> Vec<CliHit> {
    let mut hits = Vec::new();

    for path in candidate_binaries() {
        if path.exists() {
            hits.push(CliHit {
                kind: HitKind::Binary,
                detail: path.display().to_string(),
            });
        }
    }

    for path in candidate_configs() {
        if path.exists() {
            hits.push(CliHit {
                kind: HitKind::Config,
                detail: path.display().to_string(),
            });
        }
    }

    for hit in detect_autostart() {
        hits.push(hit);
    }

    hits
}

fn candidate_binaries() -> Vec<PathBuf> {
    if cfg!(target_os = "windows") {
        vec![
            PathBuf::from(r"C:\Program Files\Liaison\bin\liaison-edge.exe"),
            PathBuf::from(r"C:\Program Files (x86)\Liaison\bin\liaison-edge.exe"),
        ]
    } else if cfg!(target_os = "macos") {
        let mut v = vec![
            PathBuf::from("/usr/local/bin/liaison-edge"),
            PathBuf::from("/opt/homebrew/bin/liaison-edge"),
            PathBuf::from("/opt/liaison-edge/bin/liaison-edge"),
        ];
        if let Some(home) = dirs::home_dir() {
            v.push(home.join(".local/bin/liaison-edge"));
        }
        v
    } else {
        let mut v = vec![
            PathBuf::from("/usr/local/bin/liaison-edge"),
            PathBuf::from("/usr/bin/liaison-edge"),
            PathBuf::from("/opt/liaison-edge/bin/liaison-edge"),
        ];
        if let Some(home) = dirs::home_dir() {
            v.push(home.join(".local/bin/liaison-edge"));
        }
        v
    }
}

fn candidate_configs() -> Vec<PathBuf> {
    if cfg!(target_os = "windows") {
        vec![
            PathBuf::from(r"C:\Program Files\Liaison\conf\liaison-edge.yaml"),
            PathBuf::from(r"C:\Program Files (x86)\Liaison\conf\liaison-edge.yaml"),
        ]
    } else {
        let mut v = vec![PathBuf::from("/etc/liaison-edge/liaison-edge.yaml")];
        if let Some(home) = dirs::home_dir() {
            v.push(home.join(".config/liaison-edge/liaison-edge.yaml"));
            v.push(home.join(".liaison-edge/liaison-edge.yaml"));
        }
        v
    }
}

fn detect_autostart() -> Vec<CliHit> {
    let mut hits = Vec::new();

    if cfg!(target_os = "windows") {
        // Scheduled task XML lives under %WINDIR%\System32\Tasks\<name>.
        // Reading the registry would need extra deps; the on-disk task
        // file is just as authoritative and works without elevation.
        if let Ok(windir) = std::env::var("SystemRoot") {
            let task = PathBuf::from(&windir)
                .join("System32")
                .join("Tasks")
                .join("LiaisonEdge");
            if task.exists() {
                hits.push(CliHit {
                    kind: HitKind::Autostart,
                    detail: format!("Scheduled Task: LiaisonEdge ({})", task.display()),
                });
            }
        }
    } else if cfg!(target_os = "macos") {
        if let Some(home) = dirs::home_dir() {
            for name in [
                "cloud.liaison.edge.plist",
                "com.liaison.edge.plist",
                "io.liaison.edge.plist",
            ] {
                let p = home.join("Library/LaunchAgents").join(name);
                if p.exists() {
                    hits.push(CliHit {
                        kind: HitKind::Autostart,
                        detail: format!("LaunchAgent: {}", p.display()),
                    });
                }
            }
        }
        for name in [
            "cloud.liaison.edge.plist",
            "com.liaison.edge.plist",
        ] {
            let p = PathBuf::from("/Library/LaunchDaemons").join(name);
            if p.exists() {
                hits.push(CliHit {
                    kind: HitKind::Autostart,
                    detail: format!("LaunchDaemon: {}", p.display()),
                });
            }
        }
    } else {
        for p in [
            "/etc/systemd/system/liaison-edge.service",
            "/lib/systemd/system/liaison-edge.service",
        ] {
            let p = PathBuf::from(p);
            if p.exists() {
                hits.push(CliHit {
                    kind: HitKind::Autostart,
                    detail: format!("systemd unit: {}", p.display()),
                });
            }
        }
        if let Some(home) = dirs::home_dir() {
            let p = home.join(".config/systemd/user/liaison-edge.service");
            if p.exists() {
                hits.push(CliHit {
                    kind: HitKind::Autostart,
                    detail: format!("user systemd unit: {}", p.display()),
                });
            }
        }
    }

    hits
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn candidate_lists_are_non_empty_per_platform() {
        // Lists must include at least one path so the detector actually
        // checks something on every supported OS.
        assert!(!candidate_binaries().is_empty(), "binary candidates empty");
        assert!(!candidate_configs().is_empty(), "config candidates empty");
    }

    #[test]
    fn detect_returns_vec_without_panicking() {
        // Smoke test: must complete on whatever OS runs this.
        let _ = detect_existing_cli();
    }

    #[test]
    fn hit_kind_serialises_as_snake_case() {
        let json = serde_json::to_string(&HitKind::Binary).unwrap();
        assert_eq!(json, "\"binary\"");
        let json = serde_json::to_string(&HitKind::Autostart).unwrap();
        assert_eq!(json, "\"autostart\"");
    }
}
